package service

import (
	"time"

	"github.com/zimengpan/go-boomflow/models"
	//"github.com/zimengpan/go-boomflow/models/mysql"
)

var mockAssetDB1 = map[string]*models.Asset{
	"A": &models.Asset{
		"A",
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
	},
	"B": &models.Asset{
		"B",
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		"0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063",
	},
}

var mockAssetDB2 = map[string]*models.Asset{
	"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498": &models.Asset{
		"A",
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498",
	},
	"0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063": &models.Asset{
		"B",
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
		"0x02571792000000000000000000000000371b13d97f4bf77d724e78c16b7dc74099f40e840000000000000000000000000000000000000000000000000000000000000063",
	},
}

func GetAssetByCurrency(currency string) (*models.Asset, error) {
	//TODO: Implementation
	return mockAssetDB1[currency], nil
}

func GetAssetByAssetData(assetData string) (*models.Asset, error) {
	//TODO: Implementation
	return mockAssetDB2[assetData], nil
}
