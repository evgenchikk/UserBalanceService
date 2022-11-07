package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func New(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Println("Could not connect to database:", err)
		return nil, err
	}

	return conn, nil
}
