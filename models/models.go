package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// 用于表示一笔订单或者交易的方向：买，卖
type Side string

func NewSideFromString(s string) (*Side, error) {
	side := Side(s)
	switch side {
	case SideBuy:
	case SideSell:
	default:
		return nil, fmt.Errorf("invalid side: %v", s)
	}
	return &side, nil
}

func (s Side) Opposite() Side {
	if s == SideBuy {
		return SideSell
	}
	return SideBuy
}

func (s Side) String() string {
	return string(s)
}

// 用于表示订单状态
type OrderStatus string

func NewOrderStatusFromString(s string) (*OrderStatus, error) {
	status := OrderStatus(s)
	switch status {
	case OrderStatusNew:
	case OrderStatusOpen:
	case OrderStatusCancelling:
	case OrderStatusCancelled:
	case OrderStatusFilled:
	default:
		return nil, fmt.Errorf("invalid status: %v", s)
	}
	return &status, nil
}

func (t OrderStatus) String() string {
	return string(t)
}

// 用于表示一条fill完成的原因
type DoneReason string

type TransactionStatus string

const (
	SideBuy  = Side("buy")
	SideSell = Side("sell")

	// 初始状态
	OrderStatusNew = OrderStatus("new")
	// 已经加入orderBook
	OrderStatusOpen = OrderStatus("open")
	// 中间状态，请求取消订单
	OrderStatusCancelling = OrderStatus("cancelling")
	// 订单已经被取消，部分成交的订单也是cancelled
	OrderStatusCancelled = OrderStatus("cancelled")
	// 订单完全成交
	OrderStatusFilled = OrderStatus("filled")

	DoneReasonFilled    = DoneReason("filled")
	DoneReasonCancelled = DoneReason("cancelled")

	TransactionStatusPending   = TransactionStatus("pending")
	TransactionStatusCompleted = TransactionStatus("completed")
)

type Asset struct {
	Currency  string `gorm:"column:currency;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AssetData string
}

type Product struct {
	Id            string `gorm:"column:id;primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	BaseCurrency  string
	QuoteCurrency string
}

type Order struct {
	Id                    int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	MakerAddress          string
	TakerAddress          string
	FeeRecipientAddress   string
	SenderAddress         string
	MakerAssetAmount      decimal.Decimal `sql:"type:decimal(32,16);"`
	TakerAssetAmount      decimal.Decimal `sql:"type:decimal(32,16);"`
	MakerFee              decimal.Decimal
	TakerFee              decimal.Decimal
	ExpirationTimeSeconds decimal.Decimal `sql:"type:decimal(32,16);"`
	Salt                  decimal.Decimal
	Side                  Side
	ProductId             string
	MakerFeeAssetData     string
	TakerFeeAssetData     string
	Signature             string
	Status                OrderStatus
	Settled               bool
}

type Config struct {
	Id        int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
}
