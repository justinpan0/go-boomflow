package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/conf"
	"github.com/zimengpan/go-boomflow/matching"
	"github.com/zimengpan/go-boomflow/rest"
)

func main() {
	gbeConfig := conf.GetConfig()

	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	initValidator()

	matching.StartEngine()

	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}

	rest.StartServer()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/v1/order/{productID}", setOrder)
	router.HandleFunc("/v1/orders/{orderHash}", getOrderByHash)
	router.HandleFunc("/v1/orders", getOrders)
	router.HandleFunc("/v1/orderbook", getOrderbook)

	logger.Fatal(http.ListenAndServe(":9000", router))
}
