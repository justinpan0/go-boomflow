package match

import (
	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/conf"
	"github.com/zimengpan/go-boomflow/service"
)

func StartEngine() {
	gbeConfig := conf.GetConfig()

	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		orderReader := NewKafkaOrderReader(product.Id, gbeConfig.Kafka.Brokers)
		//snapshotStore := NewRedisSnapshotStore(product.Id)
		//logStore := NewKafkaLogStore(product.Id, gbeConfig.Kafka.Brokers)
		//matchEngine := NewEngine(product, orderReader, logStore, snapshotStore)
		matchEngine := NewEngine(product, orderReader)

		matchEngine.Start()
	}

	log.Info("match engine ok")
}
