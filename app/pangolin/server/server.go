package server

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"pangolin/app/pangolin/server/controller"
	"syscall"
)

const DefaultPort = 10086

var FatalChannel = make(chan error, 10)

type Server struct {
}

func (server *Server) InitServer(ctl *controller.Controller) {
	// 最初只需要进行端口监听，可以根据flag，如果flag没有的话默认10086开始向下递增找到第一个可用的端口
	port, custom := server.parsePort()

	// 判断port
	if custom {

	}
	router := server.initGinServer(ctl)
	srv := endless.NewServer(fmt.Sprintf(":%d", port), router)
	srv.BeforeBegin = func(addr string) {
		log.Printf("Listening on host: %s. Actual pid is %d", addr, syscall.Getpid())
	}
	srv.RegisterOnShutdown(func() {
		<-FatalChannel
		os.Exit(-1)
	})
	srv.ListenAndServe()
}

func (server *Server) parsePort() (int, bool) {
	port := -1
	flag.IntVar(&port, "port", DefaultPort, "listen port")
	flag.Parse()
	if port == DefaultPort {
		return port, false
	}
	return port, true
}

func (server *Server) initGinServer(controller *controller.Controller) *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode("debug")
	r := gin.New()
	// trace
	// 记录接口访问日志的中间件
	// 处理 panic 异常的中间件
	// 处理错误的中间件
	// ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	ajaxGroup := r.Group("/api")
	//数据库管理api
	ajaxGroup.GET("/management/db", controller.GetCurrentConfig)
	ajaxGroup.POST("/management/db", controller.ConnectDB)
	// tunnel config相关
	ajaxGroup.GET("/tunnel_configs", controller.ListAllTunnel)
	ajaxGroup.POST("/tunnel_configs", controller.CreateTunnelConfig)
	// tunnel相关
	ajaxGroup.POST("/tunnel_configs/:id/@add", controller.CreateTunnel)
	ajaxGroup.POST("/tunnel_configs/:id/@statistic", controller.GetTunnelStatistic)
	return r
}
