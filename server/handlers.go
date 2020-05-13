package server

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/chenyu116/generator-mobile/proto"
	"github.com/chenyu116/generator-mobile/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

func projects(c *gin.Context) {
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, zap.Error(errors.New("grpc connection lost")))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.GeneratorProjects(context.Background(), &pb.GeneratorProjectsRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, zap.Error(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}

func projectFeatured(c *gin.Context) {
	id, _ := uuid.NewUUID()
	fmt.Println(id , len(id))
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, zap.Error(errors.New("grpc connection lost")))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.GeneratorProjectFeatured(context.Background(), &pb.GeneratorProjectFeaturedRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, zap.Error(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}
