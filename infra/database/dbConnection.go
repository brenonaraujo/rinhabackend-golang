package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectDB(dbUrl string) error {
	var err error
	dbPool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	dbPool.Config().MaxConnLifetime = 0
	dbPool.Config().MaxConns = 200
	dbPool.Config().MinConns = 100

	return nil
}

func GetDBPool() *pgxpool.Pool {
	return dbPool
}
