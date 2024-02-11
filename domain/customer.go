package domain

type Customer struct {
	Id           int `db:"id"`
	AccountLimit int `db:"limite"`
	Balance      int `db:"valor"`
}
