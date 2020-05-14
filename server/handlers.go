package server

import (
	"context"
	"errors"
	pb "github.com/chenyu116/generator-mobile/proto"
	"github.com/chenyu116/generator-mobile/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

func jsonError(err string) gin.H {
	return gin.H{"error": errors.New(err)}
}

func projects(c *gin.Context) {
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("grpc connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.GeneratorProjects(context.Background(), &pb.GeneratorProjectsRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}

func projectFeatured(c *gin.Context) {
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.GeneratorProjectFeatured(context.Background(), &pb.GeneratorProjectFeaturedRequest{
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Projects)
}

func projectFeatures(c *gin.Context) {
	var req RequestProjectId
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	if req.ProjectId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("miss project_id"))
		return
	}
	dbServerConn, ok := utils.DbServerGrpcConn.Get().(*grpc.ClientConn)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError("GRPC connection lost"))
		return
	}
	defer utils.DbServerGrpcConn.Put(dbServerConn)
	client := pb.NewApiClient(dbServerConn)
	reply, err := client.GeneratorProjectFeaturesByProjectId(context.Background(), &pb.GeneratorProjectFeaturesByProjectIdRequest{
		ProjectId: req.ProjectId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsonError(err.Error()))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, reply.Features)
}
