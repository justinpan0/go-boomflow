package rest

import (
	"time"

	"github.com/zimengpan/go-boomflow/models"
	"github.com/zimengpan/go-boomflow/utils"
)

type messageVo struct {
	Message string `json:"message"`
}

func newMessageVo(error error) *messageVo {
	return &messageVo{
		Message: error.Error(),
	}
}

type placeOrderRequest struct {
	Hash                  string `json:"hash"`
	MakerAddress          string `json:"makerAddress"`
	TakerAddress          string `json:"takerAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	SenderAddress         string `json:"senderAddress"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerFee              string `json:"takerFee"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
	MakerAssetData        string `json:"makerAssetData"`
	TakerAssetData        string `json:"takerAssetData"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	Signature             string `json:"signature"`
}

type orderVo struct {
	Id                    string `json:"Id"`
	CreatedAt             string `json:"CreatedAt"`
	UpdatedAt             string `json:"UpdatedAt"`
	makerAddress          string `json:"makerAddress"`
	takerAddress          string `json:"takerAddress"`
	feeRecipientAddress   string `json:"feeRecipientAddress"`
	senderAddress         string `json:"senderAddress"`
	makerAssetAmount      string `json:"makerAssetAmount"`
	takerAssetAmount      string `json:"takerAssetAmount"`
	makerFee              string `json:"makerFee"`
	takerFee              string `json:"takerFee"`
	expirationTimeSeconds string `json:"expirationTimeSeconds"`
	salt                  string `json:"salt"`
	makerAssetData        string `json:"makerAssetData"`
	takerAssetData        string `json:"takerAssetData"`
	makerFeeAssetData     string `json:"makerFeeAssetData"`
	takerFeeAssetData     string `json:"takerFeeAssetData"`
	signature             string `json:"signature"`
	Status                string `json:"Status"`
	Settled               string `json:"Settled"`
}

type ProductVo struct {
	Id             string `json:"id"`
	BaseCurrency   string `json:"baseCurrency"`
	QuoteCurrency  string `json:"quoteCurrency"`
	BaseMinSize    string `json:"baseMinSize"`
	BaseMaxSize    string `json:"baseMaxSize"`
	QuoteIncrement string `json:"quoteIncrement"`
	BaseScale      int32  `json:"baseScale"`
	QuoteScale     int32  `json:"quoteScale"`
}

type orderBookVo struct {
	Sequence string           `json:"sequence"`
	Asks     [][3]interface{} `json:"asks"`
	Bids     [][3]interface{} `json:"bids"`
}

func newProductVo(product *models.Product) *ProductVo {
	return &ProductVo{
		Id:             product.Id,
		BaseCurrency:   product.BaseCurrency,
		QuoteCurrency:  product.QuoteCurrency,
		BaseMinSize:    product.BaseMinSize.String(),
		BaseMaxSize:    product.BaseMaxSize.String(),
		QuoteIncrement: utils.F64ToA(product.QuoteIncrement),
		BaseScale:      product.BaseScale,
		QuoteScale:     product.QuoteScale,
	}
}

func newOrderVo(order *models.Order) *orderVo {
	return &orderVo{
		CreatedAt:             order.CreatedAt,
		MakerAddress:          order.MakerAddress,
		TakerAddress:          order.TakerAddress,
		FeeRecipientAddress:   order.FeeRecipientAddress,
		SenderAddress:         order.SenderAddress,
		MakerAssetAmount:      order.MakerAssetAmount,
		TakerAssetAmount:      order.TakerAssetAmount,
		MakerFee:              order.MakerFee,
		TakerFee:              order.TakerFee,
		ExpirationTimeSeconds: order.ExpirationTimeSeconds,
		Salt:                  order.Salt,
		MakerAssetData:        order.MakerAssetData,
		TakerAssetData:        order.TakerAssetData,
		MakerFeeAssetData:     order.MakerFeeAssetData,
		TakerFeeAssetData:     order.TakerFeeAssetData,
		Signature:             order.Signature,
		Status:                order.Status,
		Settled:               order.Settled
	}
}
