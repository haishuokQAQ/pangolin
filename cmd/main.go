package main

import (
	"pangolin/app/pangolin/server"
	"pangolin/app/pangolin/server/controller"
	"pangolin/app/pangolin/server/service"
)

func main() {
	ctl := initController(initService())
	ajaxServer := server.Server{}
	ajaxServer.InitServer(ctl)

}

func initService() *service.Service {
	srv := &service.Service{}
	srv.InitResources()
	return srv
}

func initController(srv *service.Service) *controller.Controller {
	return &controller.Controller{Srv: srv}
}
