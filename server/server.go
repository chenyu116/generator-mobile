package server

import (
	"github.com/chenyu116/generator-mobile/config"
	"github.com/chenyu116/generator-mobile/logger"
	"github.com/chenyu116/generator-mobile/utils"
	"github.com/chenyu116/log"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"sync"
)

// Server Server
type Server struct {
	waitGroup *sync.WaitGroup
}

func NotFound(c *gin.Context) {
	log.Errorf("NotFound %s", c.Request.URL.Path)
	c.AbortWithStatus(http.StatusNotFound)
	return
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
	utils.InitPool()
	logger.InitLogger(true, "debug")
	cf := config.GetConfig()
	providerListener, err := net.Listen("tcp", cf.Serve.HostPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.NoMethod(NotFound)
	r.NoRoute(NotFound)

	v1 := r.Group("/v1")

	v1.GET("/projects", projects)
	v1.GET("/project/initialized", projectInitialized)
	v1.GET("/project/features", projectFeatures)
	v1.GET("/project/init", projectInit)
	v1.GET("/features", features)

	log.Infof("Server Started! Addr: \"%s\"", cf.Serve.HostPort)

	err = http.Serve(providerListener, r)
	if err != nil {
		log.Fatal("Server Start err:%s", err.Error())
	}
}
