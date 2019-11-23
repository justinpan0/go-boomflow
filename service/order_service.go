package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
	"github.com/zimengpan/go-boomflow/models"
)

func PlaceOrder(
	makerAddress string,
	takerAddress string,
	feeRecipientAddress string,
	senderAddress string,
	makerAssetAmount decimal.Decimal,
	takerAssetAmount decimal.Decimal,
	makerFee decimal.Decimal,
	takerFee decimal.Decimal,
	expirationTimeSeconds decimal.Decimal,
	salt decimal.Decimal,
	makerAssetData string,
	takerAssetData string,
	makerFeeAssetData string,
	takerFeeAssetData string,
	signature string
) (*models.Order, error) {
	product, err := GetProductByAssetPair(makerAssetData, takerAssetData)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New(fmt.Sprintf("product not found: %v", productId))
	}

	order := &models.Order{
		MakerAddress: 			makerAddress,
		TakerAddress: 			takerAddress,
		FeeRecipientAddress: 	feeRecipientAddress,
		SenderAddress: 			senderAddress,
		MakerAssetAmount: 		makerAssetAmount,
		TakerAssetAmount: 		takerAssetAmount,
		MakerFee: 				makerFee,
		TakerFee: 				takerFee,
		ExpirationTimeSeconds: 	expirationTimeSeconds,
		Salt: 					salt,
		MakerAssetData: 		makerAssetData,
		TakerAssetData: 		takerAssetData,
		MakerFeeAssetData: 		makerFeeAssetData,
		TakerFeeAssetData: 		takerFeeAssetData,
		Signature: 				signature,
		Status:    				models.OrderStatusNew,
	}
	return order
	// tx
	/*
		var holdCurrency string
		var holdSize decimal.Decimal

		db, err := mysql.SharedStore().BeginTx()
		if err != nil {
			return nil, err
		}
		defer func() { _ = db.Rollback() }()

		err = HoldBalance(db, userId, holdCurrency, holdSize, models.BillTypeTrade)
		if err != nil {
			return nil, err
		}

		err = db.AddOrder(order)
		if err != nil {
			return nil, err
		}

		return order, db.CommitTx()*/
}

func UpdateOrderStatus(orderId int64, oldStatus, newStatus models.OrderStatus) (bool, error) {
	return mysql.SharedStore().UpdateOrderStatus(orderId, oldStatus, newStatus)
}

func ExecuteFill(orderId int64) error {
	// tx
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	order, err := db.GetOrderByIdForUpdate(orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found: %v", orderId)
	}
	if order.Status == models.OrderStatusFilled || order.Status == models.OrderStatusCancelled {
		return fmt.Errorf("order status invalid: %v %v", orderId, order.Status)
	}

	product, err := GetProductById(order.ProductId)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product not found: %v", order.ProductId)
	}

	fills, err := mysql.SharedStore().GetUnsettledFillsByOrderId(orderId)
	if err != nil {
		return err
	}
	if len(fills) == 0 {
		return nil
	}

	var bills []*models.Bill
	for _, fill := range fills {
		fill.Settled = true

		notes := fmt.Sprintf("%v-%v", fill.OrderId, fill.Id)

		if !fill.Done {
			executedValue := fill.Size.Mul(fill.Price)
			order.ExecutedValue = order.ExecutedValue.Add(executedValue)
			order.FilledSize = order.FilledSize.Add(fill.Size)

			if order.Side == models.SideBuy {
				// 买单，incr base
				bill, err := AddDelayBill(db, order.UserId, product.BaseCurrency, fill.Size, decimal.Zero,
					models.BillTypeTrade, notes)
				if err != nil {
					return err
				}
				bills = append(bills, bill)

				// 买单，decr quote
				bill, err = AddDelayBill(db, order.UserId, product.QuoteCurrency, decimal.Zero, executedValue.Neg(),
					models.BillTypeTrade, notes)
				if err != nil {
					return err
				}
				bills = append(bills, bill)

			} else {
				// 卖单，decr base
				bill, err := AddDelayBill(db, order.UserId, product.BaseCurrency, decimal.Zero, fill.Size.Neg(),
					models.BillTypeTrade, notes)
				if err != nil {
					return err
				}
				bills = append(bills, bill)

				// 卖单，incr quote
				bill, err = AddDelayBill(db, order.UserId, product.QuoteCurrency, executedValue, decimal.Zero,
					models.BillTypeTrade, notes)
				if err != nil {
					return err
				}
				bills = append(bills, bill)
			}

		} else {
			if fill.DoneReason == models.DoneReasonCancelled {
				order.Status = models.OrderStatusCancelled
			} else if fill.DoneReason == models.DoneReasonFilled {
				order.Status = models.OrderStatusFilled
			} else {
				log.Fatalf("unknown done reason: %v", fill.DoneReason)
			}

			if order.Side == models.SideBuy {
				// 如果是是买单，需要解冻剩余的funds
				remainingFunds := order.Funds.Sub(order.ExecutedValue)
				if remainingFunds.GreaterThan(decimal.Zero) {
					bill, err := AddDelayBill(db, order.UserId, product.QuoteCurrency, remainingFunds, remainingFunds.Neg(),
						models.BillTypeTrade, notes)
					if err != nil {
						return err
					}
					bills = append(bills, bill)
				}

			} else {
				// 如果是卖单，解冻剩余的size
				remainingSize := order.Size.Sub(order.FilledSize)
				if remainingSize.GreaterThan(decimal.Zero) {
					bill, err := AddDelayBill(db, order.UserId, product.BaseCurrency, remainingSize, remainingSize.Neg(),
						models.BillTypeTrade, notes)
					if err != nil {
						return err
					}
					bills = append(bills, bill)
				}
			}

			break
		}
	}

	err = db.UpdateOrder(order)
	if err != nil {
		return err
	}

	for _, fill := range fills {
		err = db.UpdateFill(fill)
		if err != nil {
			return err
		}
	}

	return db.CommitTx()
}

func GetOrderById(orderId int64) (*models.Order, error) {
	return mysql.SharedStore().GetOrderById(orderId)
}

func GetOrderByClientOid(userId int64, clientOid string) (*models.Order, error) {
	return mysql.SharedStore().GetOrderByClientOid(userId, clientOid)
}

func GetOrdersByUserId(userId int64, statuses []models.OrderStatus, side *models.Side, productId string,
	beforeId, afterId int64, limit int) ([]*models.Order, error) {
	return mysql.SharedStore().GetOrdersByUserId(userId, statuses, side, productId, beforeId, afterId, limit)
}
