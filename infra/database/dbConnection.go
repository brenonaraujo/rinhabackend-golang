package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectDB(dbUrl string) error {
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	dbPool.Config().MaxConnLifetime = 0 // Means no limit
	dbPool.Config().MaxConns = 50       // Adjust based on your needs
	dbPool.Config().MinConns = 5        // Adjust based on your needs

	return nil
}

func GetDBPool() *pgxpool.Pool {
	return dbPool
}
