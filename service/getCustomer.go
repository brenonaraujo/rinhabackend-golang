package service

import (
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
)

type Customer struct {
	Id           int `db:"id"`
	AccountLimit int `db:"limite"`
	Balance      int `db:"valor"`
}

var cachedCustomers = make(map[int]*Customer)

func GetCustomer(customerId int) (*Customer, error) {
	if customer, ok := cachedCustomers[customerId]; ok {
		return customer, nil
	}

	db := database.GetDBPool()
	row := db.QueryRow(context.Background(), "SELECT id, limite FROM clientes WHERE id = $1", customerId)

	var customer Customer
	err := row.Scan(&customer.Id, &customer.AccountLimit)
	if err != nil {
		return nil, err
	}

	cachedCustomers[customer.Id] = &customer
	return &customer, nil
}
