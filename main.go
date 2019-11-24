package main

import (
	"net/http"

	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/match"
	"github.com/zimengpan/go-boomflow/rest"
)

func main() {
	//gbeConfig := conf.GetConfig()

	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	match.StartEngine()

	/*products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}*/

	rest.StartServer()

	select {}
}
