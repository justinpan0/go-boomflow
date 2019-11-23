package rest

import (
	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/conf"
)

func StartServer() {
	gbeConfig := conf.GetConfig()

	httpServer := NewHttpServer(gbeConfig.RestServer.Addr)
	go httpServer.Start()

	log.Info("rest server ok")
}
