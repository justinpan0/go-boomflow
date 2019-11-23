package rest

import (
	"time"

	"github.com/zimengpan/go-boomflow/models"
	"github.com/zimengpan/go-boomflow/service"
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
	Hash                  string  `json:"hash"`
	MakerAddress          string  `json:"makerAddress"`
	TakerAddress          string  `json:"takerAddress"`
	FeeRecipientAddress   string  `json:"feeRecipientAddress"`
	SenderAddress         string  `json:"senderAddress"`
	MakerAssetAmount      float64 `json:"makerAssetAmount"`
	TakerAssetAmount      float64 `json:"takerAssetAmount"`
	MakerFee              float64 `json:"makerFee"`
	TakerFee              float64 `json:"takerFee"`
	ExpirationTimeSeconds float64 `json:"expirationTimeSeconds"`
	Salt                  float64 `json:"salt"`
	MakerAssetData        string  `json:"makerAssetData"`
	TakerAssetData        string  `json:"takerAssetData"`
	MakerFeeAssetData     string  `json:"makerFeeAssetData"`
	TakerFeeAssetData     string  `json:"takerFeeAssetData"`
	Signature             string  `json:"signature"`
}

type orderVo struct {
	Id                    string `json:"Id"`
	CreatedAt             string `json:"CreatedAt"`
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
	Side                  string `json:"Side"`
	ProductId             string `json:"ProductId"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	Signature             string `json:"signature"`
	Status                string `json:"Status"`
	Settled               bool   `json:"Settled"`
}

type ProductVo struct {
	Id             string `json:"id"`
	BaseCurrency   string `json:"baseCurrency"`
	QuoteCurrency  string `json:"quoteCurrency"`
	BaseAssetData  string `json:"BaseAssetData"`
	QuoteAssetData string `json:"QuoteAssetData"`
}

type orderBookVo struct {
	Sequence string           `json:"sequence"`
	Asks     [][3]interface{} `json:"asks"`
	Bids     [][3]interface{} `json:"bids"`
}

func newProductVo(product *models.Product) *ProductVo {
	base, _ := service.GetAssetByCurrency(product.BaseCurrency)
	quote, _ := service.GetAssetByCurrency(product.QuoteCurrency)
	return &ProductVo{
		Id:             product.Id,
		BaseCurrency:   product.BaseCurrency,
		QuoteCurrency:  product.QuoteCurrency,
		BaseAssetData:  base.AssetData,
		QuoteAssetData: quote.AssetData,
	}
}

func newOrderVo(order *models.Order) *orderVo {
	return &orderVo{
		CreatedAt:             order.CreatedAt.Format(time.RFC3339),
		MakerAddress:          order.MakerAddress,
		TakerAddress:          order.TakerAddress,
		FeeRecipientAddress:   order.FeeRecipientAddress,
		SenderAddress:         order.SenderAddress,
		MakerAssetAmount:      order.MakerAssetAmount.String(),
		TakerAssetAmount:      order.TakerAssetAmount.String(),
		MakerFee:              order.MakerFee.String(),
		TakerFee:              order.TakerFee.String(),
		ExpirationTimeSeconds: order.ExpirationTimeSeconds.String(),
		Salt:                  order.Salt.String(),
		Side:                  order.Side.String(),
		ProductId:             order.ProductId,
		MakerFeeAssetData:     order.MakerFeeAssetData,
		TakerFeeAssetData:     order.TakerFeeAssetData,
		Signature:             order.Signature,
		Status:                order.Status.String(),
		Settled:               order.Settled,
	}
}
