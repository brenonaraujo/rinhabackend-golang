package service

import (
	domain "brenonaraujo/rinhabackend-q12024/domain/entities"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"context"
)

var cachedCustomers = make(map[int]*domain.Customer)

func GetCustomer(customerId int) (*domain.Customer, error) {
	if customer, ok := cachedCustomers[customerId]; ok {
		return customer, nil
	}

	db := database.GetDBPool()
	row := db.QueryRow(context.Background(), "SELECT id, limite FROM clientes WHERE id = $1", customerId)

	var customer domain.Customer
	err := row.Scan(&customer.Id, &customer.AccountLimit)
	if err != nil {
		return nil, err
	}

	cachedCustomers[customer.Id] = &customer
	return &customer, nil
}
