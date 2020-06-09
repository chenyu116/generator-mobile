package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/chenyu116/generator-mobile/config"
	"github.com/chenyu116/generator-mobile/logger"
	"github.com/chenyu116/log"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var RabbitMQ *rabbitMQClient

type rabbitMQClient struct {
	rmqcfg           config.RabbitmqConfig
	conn             *amqp.Connection
	ch               *amqp.Channel
	amqpCfg          amqp.Config
	ProcessorMessage func(msg amqp.Delivery) error
	ctx              context.Context
	CtxCancel        context.CancelFunc
	deliveryTag      uint64
	deliveryMu       sync.Mutex
	confirmMap       map[string]chan bool
	confirmMapMu     sync.Mutex
	confirm          chan uint64
	rtn              chan string
	publish          chan amqp.Confirmation
}

func NewRabbitMQClient(cfg config.RabbitmqConfig) *rabbitMQClient {
	if cfg.VHost == "" {
		cfg.VHost = "/"
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &rabbitMQClient{
		rmqcfg: cfg,
		amqpCfg: amqp.Config{
			Heartbeat: time.Second * 2,
			Vhost:     cfg.VHost,
		},
		ctx:        ctx,
		CtxCancel:  cancel,
		confirmMap: make(map[string]chan bool),
	}
}

func (p *rabbitMQClient) Channel() *amqp.Channel {
	return p.ch
}

func (p *rabbitMQClient) Publish(exchange, routeKey string,
	msg amqp.Publishing, confirmChan chan bool) (err error) {
	if p.isClosed() {
		err = errors.New("消息系统未连接")
		return
	}
	if confirmChan != nil && msg.ReplyTo == "" {
		err = errors.New("要求回执但未填写回执编号")
		return
	}
	err = p.ch.Publish(
		exchange,
		routeKey,
		false,
		false,
		msg)
	if err != nil {
		return
	}
	select {
	case c := <-p.publish:
		if !c.Ack {
			err = errors.New("投递服务器失败")
			return
		}
		if confirmChan != nil {
			clientChan := make(chan bool, 1)
			p.confirmMapMu.Lock()
			p.confirmMap[msg.ReplyTo] = clientChan
			p.confirmMapMu.Unlock()

			timeout := time.Second * 10
			if msg.Expiration != "" {
				tt, _ := time.ParseDuration(msg.Expiration + "ms")
				if tt > 0 {
					timeout = tt
				}
			}
			tc := time.After(timeout)
			select {
			case <-clientChan:
				confirmChan <- true
			case <-tc:
				err = errors.New("设备无响应")
				p.confirmMapMu.Lock()
				delete(p.confirmMap, msg.ReplyTo)
				p.confirmMapMu.Unlock()
			}
		}
		return
	}
}

func (p *rabbitMQClient) Init() (rabbitMQClient *rabbitMQClient, err error) {
	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s/", p.rmqcfg.Username, p.rmqcfg.Password, p.rmqcfg.HostPort)
	p.conn, err = amqp.DialConfig(amqpUrl, p.amqpCfg)
	if err != nil {
		return p, err
	}
	p.ch, err = p.conn.Channel()
	if err != nil {
		return nil, err
	}
	p.publish = make(chan amqp.Confirmation, 1)
	_ = p.ch.Confirm(false)
	p.ch.NotifyPublish(p.publish)

	err = p.ch.Qos(p.rmqcfg.Prefetch, 0, false)
	if err != nil {
		return nil, err
	}
	messageQueueName := p.rmqcfg.QueuePrefix + "-" + time.Now().String()
	_, err = p.ch.QueueDeclare(
		messageQueueName, // name
		true,             // durable
		true,             // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

	for _, v := range p.rmqcfg.Exchanges {
		err = p.ch.ExchangeDeclare(v["name"], v["kind"], true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		err = p.ch.QueueBind(
			messageQueueName, // queue name
			messageQueueName, // routing key
			v["name"],        // exchange
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	go p.consumerMessage(messageQueueName)

	return p, err
}

func (p *rabbitMQClient) consumerMessage(queueName string) {
	defer func() {
		log.Debugf("consumerMessage stopped")
		p.CtxCancel()
	}()
	messages, err := p.ch.Consume(
		queueName, // queue
		queueName, // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.ZapLogger.Error("[Rabbitmq]Consume", zap.Error(err))
		return
	}

	go func() {
		defer log.Debugf("listener off")
		for {
			select {
			case <-p.ctx.Done():
				_ = p.ch.Close()
				_ = p.conn.Close()
				p.conn = nil
				return
			default:
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()
	for d := range messages {
		go func() {
			err = nil
			if d.ReplyTo != "" {
				p.confirmMapMu.Lock()
				if c, ok := p.confirmMap[d.ReplyTo]; ok {
					c <- true
					delete(p.confirmMap, d.ReplyTo)
				}
				p.confirmMapMu.Unlock()
			} else {
				if p.ProcessorMessage != nil && strings.Compare(d.AppId, queueName) != 0 {
					err = p.ProcessorMessage(d)
				}
			}
			if err == nil {
				log.Debugf("d.Ack %d %s", d.DeliveryTag, d.Body)
				_ = d.Ack(false)
			}
		}()
	}
}

func (p *rabbitMQClient) Recovery() {
	recoveryTicker := time.NewTicker(time.Second * 3)
	reconnecting := false
	for range recoveryTicker.C {
		if !p.isClosed() || reconnecting {
			continue
		}
		reconnecting = true
		p.CtxCancel()
		_, err := p.Init()
		if err != nil {
			logger.ZapLogger.Error("[Rabbitmq]Recovery", zap.Error(err))
		} else {
			p.ctx, p.CtxCancel = context.WithCancel(context.Background())
		}
		reconnecting = false
	}
}

func (p *rabbitMQClient) isClosed() bool {
	if p.conn == nil || p.conn.IsClosed() {
		return true
	}
	return false
}
