package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrdersModel struct {
	BaseModel
}

type Order struct {
	ID       string `json:"id"`
	Distance uint32 `json:"distance"`
	Status   string `json:"status"`
}

// Place order
func (om *OrdersModel) Place(origin, destination [2]string, distance uint32, status string) (id string, err error) {
	id = uuid.New().String()
	_, err = om.DB.Exec("INSERT INTO orders (id, origin_lat, origin_long, destination_lat, destination_long, distance, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, origin[0], origin[1], destination[0], destination[1], distance, status)

	return
}

// Take order
func (om *OrdersModel) Take(id, status string) (err error) {
	var order Order

	// Begin transaction
	tx, err := om.DB.Beginx()
	if err != nil {
		return
	}

	// Add 5s timeout for getting order
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get and lock the order, return an error if order not exists
	err = tx.GetContext(ctx, &order, "SELECT status FROM orders WHERE id=? FOR UPDATE", id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}

		return
	}

	// Check if order is unassigned
	if order.Status != "UNASSIGNED" {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}

		return errors.New("Order is not unassigned")
	}

	// Take the order
	_, err = tx.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}

		return
	}

	return tx.Commit()
}

// List orders
func (om *OrdersModel) List(page, limit int) ([]Order, error) {
	orders := []Order{}

	// Add 5s timeout for getting order
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := om.DB.SelectContext(ctx, &orders, "SELECT id, distance, status FROM orders INNER JOIN (SELECT id FROM orders LIMIT ?, ?) t USING (id)", (page-1)*limit, limit)

	return orders, err
}
