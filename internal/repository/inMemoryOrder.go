package repository

import (
	"applicationDesignTest/internal/models"
	"context"
	"sync"
)

var orders []models.Order

type Order struct {
	orders []models.Order
	sync.Mutex
}

func NewOrder() *Order {
	return &Order{
		orders: orders,
	}
}

func (o *Order) CreateOrder(ctx context.Context, order models.Order) (id int, err error) {
	o.Lock()
	defer o.Unlock()
	if len(o.orders) > 0 {
		id = o.orders[len(o.orders)-1].ID + 1
	} else {
		id++
	}

	order.ID = id
	o.orders = append(o.orders, order)

	return
}

func (o *Order) DeleteOrder(ctx context.Context, id int) error {
	o.Lock()
	defer o.Unlock()
	var newOrders []models.Order
	for _, existingOrder := range o.orders {
		if existingOrder.ID != id {
			newOrders = append(newOrders, existingOrder)
		}
	}
	o.orders = newOrders

	return nil
}
