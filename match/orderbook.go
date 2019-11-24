package match

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
	"github.com/siddontang/go-log/log"
	"github.com/zimengpan/go-boomflow/models"
)

const (
	orderIdWindowCap = 10000
)

type orderBook struct {
	// one product corresponds to one order book
	product *models.Product

	// depths: asks & bids
	depths map[models.Side]*depth

	// strictly continuously increasing transaction ID, used for the primary key ID of trade
	tradeSeq int64

	// strictly continuously increasing log SEQ, used to write matching log
	logSeq int64

	// to prevent the order from being submitted to the order book repeatedly,
	// a sliding window de duplication strategy is adopted.
	orderIdWindow Window
}

type priceOrderIdKey struct {
	price   decimal.Decimal
	orderId int64
}

func NewOrderBook(product *models.Product) *orderBook {
	asks := &depth{
		queue:  treemap.NewWith(priceOrderIdKeyAscComparator),
		orders: map[int64]*BookOrder{},
	}
	bids := &depth{
		queue:  treemap.NewWith(priceOrderIdKeyDescComparator),
		orders: map[int64]*BookOrder{},
	}

	orderBook := &orderBook{
		product:       product,
		depths:        map[models.Side]*depth{models.SideBuy: bids, models.SideSell: asks},
		orderIdWindow: newWindow(0, orderIdWindowCap),
	}
	return orderBook
}

func (o *orderBook) ApplyOrder(order *models.Order) (logs []Log) {
	// prevent orders from being submitted repeatedly to the matching engine
	log.Info("Order: ", order.Signature)
	err := o.orderIdWindow.put(order.Salt)
	if err != nil {
		log.Error(err)
		return logs
	}

	takerOrder := newBookOrder(order)
	log.Info("Order Price", takerOrder.Price)
	if takerOrder.Size.GreaterThan(decimal.Zero) {
		// If taker has an uncompleted size, put taker in orderBook
		o.depths[takerOrder.Side].add(*takerOrder)

		openLog := newOpenLog(o.nextLogSeq(), o.product.Id, takerOrder)
		logs = append(logs, openLog)
	} else {
		var remainingSize = takerOrder.Size
		var reason = models.DoneReasonFilled

		doneLog := newDoneLog(o.nextLogSeq(), o.product.Id, takerOrder, remainingSize, reason)
		logs = append(logs, doneLog)
	}
	return logs
}

/*func (o *orderBook) CancelOrder(order *models.Order) (logs []Log) {
	_ = o.orderIdWindow.put(order.Id)

	bookOrder, found := o.depths[order.Side].orders[order.Id]
	if !found {
		return logs
	}

	// 将order的size全部decr，等于remove操作
	remainingSize := bookOrder.Size
	err := o.depths[order.Side].decrSize(order.Id, bookOrder.Size)
	if err != nil {
		panic(err)
	}

	doneLog := newDoneLog(o.nextLogSeq(), o.product.Id, bookOrder, remainingSize, models.DoneReasonCancelled)
	return append(logs, doneLog)
}*/

/*func (o *orderBook) Snapshot() orderBookSnapshot {
	snapshot := orderBookSnapshot{
		Orders:        make([]BookOrder, len(o.depths[models.SideSell].orders)+len(o.depths[models.SideBuy].orders)),
		LogSeq:        o.logSeq,
		TradeSeq:      o.tradeSeq,
		OrderIdWindow: o.orderIdWindow,
	}

	i := 0
	for _, order := range o.depths[models.SideSell].orders {
		snapshot.Orders[i] = *order
		i++
	}
	for _, order := range o.depths[models.SideBuy].orders {
		snapshot.Orders[i] = *order
		i++
	}

	return snapshot
}

func (o *orderBook) Restore(snapshot *orderBookSnapshot) {
	o.logSeq = snapshot.LogSeq
	o.tradeSeq = snapshot.TradeSeq
	o.orderIdWindow = snapshot.OrderIdWindow
	if o.orderIdWindow.Cap == 0 {
		o.orderIdWindow = newWindow(0, orderIdWindowCap)
	}

	for _, order := range snapshot.Orders {
		o.depths[order.Side].add(order)
	}
}*/

func (o *orderBook) nextLogSeq() int64 {
	o.logSeq++
	return o.logSeq
}

func (o *orderBook) nextTradeSeq() int64 {
	o.tradeSeq++
	return o.tradeSeq
}

type depth struct {
	// all orders
	orders map[int64]*BookOrder

	// price first, time first order queue for order match
	// priceOrderIdKey -> orderId
	queue *treemap.Map
}

func (d *depth) add(order BookOrder) {
	d.orders[order.OrderId] = &order
	d.queue.Put(&priceOrderIdKey{order.Price, order.OrderId}, order.OrderId)
}

func (d *depth) decrSize(orderId int64, size decimal.Decimal) error {
	order, found := d.orders[orderId]
	if !found {
		return errors.New(fmt.Sprintf("order %v not found on book", orderId))
	}

	if order.Size.LessThan(size) {
		return errors.New(fmt.Sprintf("order %v Size %v less than %v", orderId, order.Size, size))
	}

	order.Size = order.Size.Sub(size)
	if order.Size.IsZero() {
		delete(d.orders, orderId)
		d.queue.Remove(&priceOrderIdKey{order.Price, order.OrderId})
	}

	return nil
}

type BookOrder struct {
	OrderId int64
	Size    decimal.Decimal
	Funds   decimal.Decimal
	Price   decimal.Decimal
	Side    models.Side
}

func newBookOrder(order *models.Order) *BookOrder {
	return &BookOrder{
		OrderId: order.Id,
		Size:    order.TakerAssetAmount,
		Funds:   order.MakerAssetAmount,
		Price:   (order.TakerAssetAmount).Div(order.MakerAssetAmount),
		Side:    order.Side,
	}
}

func priceOrderIdKeyAscComparator(a, b interface{}) int {
	aAsserted := a.(*priceOrderIdKey)
	bAsserted := b.(*priceOrderIdKey)

	x := aAsserted.price.Cmp(bAsserted.price)
	if x != 0 {
		return x
	}

	y := aAsserted.orderId - bAsserted.orderId
	if y == 0 {
		return 0
	} else if y > 0 {
		return 1
	} else {
		return -1
	}
}

func priceOrderIdKeyDescComparator(a, b interface{}) int {
	aAsserted := a.(*priceOrderIdKey)
	bAsserted := b.(*priceOrderIdKey)

	x := aAsserted.price.Cmp(bAsserted.price)
	if x != 0 {
		return -x
	}

	y := aAsserted.orderId - bAsserted.orderId
	if y == 0 {
		return 0
	} else if y > 0 {
		return 1
	} else {
		return -1
	}
}
