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
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range,x-token,x-valid-code,x-refresh")
			c.Header("Content-Type", "text/plain; charset=utf-8")
			c.Header("Content-Length", "0")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
	v1 := r.Group("/v1")

	v1.GET("/projects", projects)
	v1.GET("/project/initialized", projectInitialized)
	v1.GET("/project/features", projectFeatures)
	v1.GET("/project/init", projectInit)
	v1.GET("/features", features)
	v1.GET("/feature", feature)
	v1.PUT("/install", install)
	v1.POST("/edit", edit)
	v1.POST("/upload", upload)
	v1.PUT("/build", build)

	log.Infof("Server Started! Addr: \"%s\"", cf.Serve.HostPort)

	err = http.Serve(providerListener, r)
	if err != nil {
		log.Fatal("Server Start err:%s", err.Error())
	}
}
