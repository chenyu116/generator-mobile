package utils

import (
	"context"
	"github.com/chenyu116/generator-mobile/config"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	DbServerGrpcConn sync.Pool
)

func InitPool() {
	cf := config.GetConfig()
	DbServerGrpcConn = sync.Pool{New: func() interface{} {
		ctxTimeout, ctxTimeoutCancel := context.WithTimeout(context.Background(), time.Second*5)
		dbServerConn, err := grpc.DialContext(ctxTimeout, cf.DbServer.HostPort, grpc.WithInsecure(),
			grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*12)))
		if err != nil {
			ctxTimeoutCancel()
			return nil
		}
		ctxTimeoutCancel()
		return dbServerConn
	}}
}
