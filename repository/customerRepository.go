package repository

type CustomerRepository interface {
	DeductBalance(customerId string, amount float64) error
	AddBalance(customerId string, amount float64) error
	CreateTransaction(customerId string, amount float64, transactionType string, description string) error
	GetBalance(customerId string) (float64, error)
}
