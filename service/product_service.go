package service

import (
	"errors"
	"time"

	"github.com/zimengpan/go-boomflow/models"
	//"github.com/zimengpan/go-boomflow/models/mysql"
)

func GetProductByAssetPair(assetA string, assetB string) (*models.Product, error) {
	if assetA == "0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498"
		&& assetB == "0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063" {
		return &models.Product{
			"1",
			time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
			time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
			"A",
			"B",
			"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
			"0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063",
		}, nil
	}

	return nil, errors.New("GetProductByID")
	//return mysql.SharedStore().GetProductById(id)
}

func GetProducts() ([]*models.Product, error) {
	return []*models.Product{&models.Product{
		"1",
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		"A",
		"B",
		"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
		"0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063",
	}}, nil
	//return mysql.SharedStore().GetProducts()
}
