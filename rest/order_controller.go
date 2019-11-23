package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/conf"
	"github.com/zimengpan/go-boomflow/matching"
	"github.com/zimengpan/go-boomflow/models"
	"github.com/zimengpan/go-boomflow/service"
	"github.com/zimengpan/go-boomflow/utils"
)

var productId2Writer sync.Map
var assetPari2ProductId sync.Map

func getWriter(assetA string, assetB string) *kafka.Writer {
	key := assetA + assetB
	if assetA > assetB {
		key := assetB + assetA
	}

	productId, found := assetPari2ProductId.Load(key)
	if !found {
		product, err := GetProductByAssetPair(makerAssetData, takerAssetData)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New(fmt.Sprintf("product not found: %v", productId))
		}
		productId = product.Id
		assetPari2ProductId.Store(key, productId)
	}

	writer, found := productId2Writer.Load(productId)
	if found {
		return writer.(*kafka.Writer)
	}

	gbeConfig := conf.GetConfig()

	newWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      gbeConfig.Kafka.Brokers,
		Topic:        matching.TopicOrderPrefix + productId,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 5 * time.Millisecond,
	})
	productId2Writer.Store(productId, newWriter)
	return newWriter
}

func submitOrder(order *models.Order) {
	buf, err := json.Marshal(order)
	if err != nil {
		log.Error(err)
		return
	}

	err = getWriter(order.MakerAssetData, order.TakerAssetData).WriteMessages(context.Background(), kafka.Message{Value: buf})
	if err != nil {
		log.Error(err)
	}
}

// POST /orders
func PlaceOrder(ctx *gin.Context) {
	var req placeOrderRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newMessageVo(err))
		return
	}

	//TODO: Validate order, signature and balances
	makerAddress := req.MakerAddress
	takerAddress := req.TakerAddress
	if takerAddress != "0x0000000000000000000000000000000000000000" {
		ctx.JSON(http.StatusBadRequest, newMessageVo(fmt.Errorf("Taker Address is not Zero")))
		return
	}
	feeRecipientAddress := req.FeeRecipientAddress
	senderAddress := req.SenderAddress
	makerAssetAmount := req.MakerAssetAmount
	takerAssetAmount := req.TakerAssetAmount
	makerFee := req.MakerFee
	takerFee := req.TakerFee
	expirationTimeSeconds := req.ExpirationTimeSeconds
	salt := req.Salt
	makerAssetData := req.MakerAssetData
	takerAssetData := req.TakerAssetData
	makerFeeAssetData := req.MakerFeeAssetData
	takerFeeAssetData := req.TakerFeeAssetData
	signature := req.Signature

	// Place Order to SQL DB
	order, err := service.PlaceOrder(
		makerAddress,
		takerAddress,
		feeRecipientAddress,
		senderAddress,
		makerAssetAmount,
		takerAssetAmount,
		makerFee,
		takerFee,
		expirationTimeSeconds,
		salt,
		makerAssetData,
		takerAssetData,
		makerFeeAssetData,
		takerFeeAssetData,
		signature)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	submitOrder(order)

	ctx.JSON(http.StatusOK, order)
}

// 撤销指定id的订单
// DELETE /orders/1
// DELETE /orders/client:1
func CancelOrder(ctx *gin.Context) {
	rawOrderId := ctx.Param("orderId")

	var order *models.Order
	var err error
	if strings.HasPrefix(rawOrderId, "client:") {
		clientOid := strings.Split(rawOrderId, ":")[1]
		order, err = service.GetOrderByClientOid(GetCurrentUser(ctx).Id, clientOid)
	} else {
		orderId, _ := utils.AToInt64(rawOrderId)
		order, err = service.GetOrderById(orderId)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	if order == nil || order.UserId != GetCurrentUser(ctx).Id {
		ctx.JSON(http.StatusNotFound, newMessageVo(errors.New("order not found")))
		return
	}

	order.Status = models.OrderStatusCancelling
	submitOrder(order)

	ctx.JSON(http.StatusOK, nil)
}

// 批量撤单
// DELETE /orders/?productId=BTC-USDT&side=[buy,sell]
func CancelOrders(ctx *gin.Context) {
	productId := ctx.Query("productId")

	var side *models.Side
	var err error
	rawSide := ctx.Query("side")
	if len(rawSide) > 0 {
		side, err = models.NewSideFromString(rawSide)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newMessageVo(err))
			return
		}
	}

	orders, err := service.GetOrdersByUserId(GetCurrentUser(ctx).Id,
		[]models.OrderStatus{models.OrderStatusOpen, models.OrderStatusNew}, side, productId, 0, 0, 10000)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	for _, order := range orders {
		order.Status = models.OrderStatusCancelling
		submitOrder(order)
	}

	ctx.JSON(http.StatusOK, nil)
}

// GET /orders
func GetOrders(ctx *gin.Context) {
	makerAssetData := ctx.Query("makerAssetData")
	takerAssetData := ctx.Query("takerAssetData")

	var err error

	var statuses []models.OrderStatus
	statusValues := ctx.QueryArray("status")
	for _, statusValue := range statusValues {
		status, err := models.NewOrderStatusFromString(statusValue)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newMessageVo(err))
			return
		}
		statuses = append(statuses, *status)
	}

	before, _ := strconv.ParseInt(ctx.Query("before"), 10, 64)
	after, _ := strconv.ParseInt(ctx.Query("after"), 10, 64)

	orders, err := service.GetOrdersByUserId(makerAddress, statuses, makerAssetData, takerAssetData, before, after)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMessageVo(err))
		return
	}

	orderVos := []*orderVo{}
	for _, order := range orders {
		orderVos = append(orderVos, newOrderVo(order))
	}

	var newBefore, newAfter int64 = 0, 0
	if len(orders) > 0 {
		newBefore = orders[0].Id
		newAfter = orders[len(orders)-1].Id
	}
	ctx.Header("gbe-before", strconv.FormatInt(newBefore, 10))
	ctx.Header("gbe-after", strconv.FormatInt(newAfter, 10))

	ctx.JSON(http.StatusOK, orderVos)
}