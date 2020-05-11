package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/chenyu116/generator-mobile/config"
	"github.com/chenyu116/generator-mobile/logger"
	"go.uber.org/zap"
)

// Server Server
type Server struct {
	waitGroup *sync.WaitGroup
}

// NewServer NewServer
func NewServer() *Server {
	s := &Server{
		waitGroup: &sync.WaitGroup{},
	}
	return s
}

// Start Start server
func (s *Server) Start() {
	fileService := http.FileServer(http.Dir("./dist"))
	mux := http.NewServeMux()
	mux.Handle("/", fileService)

	logger.InitLogger(true, "debug")
	cf := config.GetConfig()


	_, err := net.Listen("tcp", cf.Websocket.HostPort)
	if err != nil {
		logger.ZapLogger.Fatal("tcp.Listen", zap.Error(err))
	}

	logger.ZapLogger.Debug("Websocket Serving", zap.String("hostPort", cf.Websocket.HostPort), zap.Bool("TLS",
		cf.Websocket.Tls.Enable))

	fmt.Println("OK")
}
