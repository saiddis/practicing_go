package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PG struct {
	DB *pgxpool.Pool
}

var (
	pgInstance *PG
	pgOnce     sync.Once
)

func New(dbName string, opts ...Option) (*PG, error) {
	options := options{
		user:     "postgres",
		host:     "localhost",
		password: "root",
		port:     5432,
		params:   map[string]string{},
	}
	for _, opt := range opts {
		opt(&options)
	}

	var url string
	if options.url != "" {
		url = options.url
		return connectToDB(url)
	}

	if options.port <= 0 {
		return nil, errors.New("port parameter must be more than 0")
	}
	url = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		options.user,
		options.password,
		options.host,
		options.port,
		dbName)

	if options.params != nil {
		paramAdder := addParams(url)
		if v, ok := options.params["sslmode"]; ok {
			url = paramAdder("sslmode=" + v)
		}
		if v, ok := options.params["TimeZone"]; ok {
			url = paramAdder("TimeZone=" + v)
		}
	}

	return connectToDB(url)
}

func (pg *PG) Close() {
	pg.DB.Close()
}

func addParams(url string) func(string) string {
	url += "?"
	return func(param string) string {
		url += param + "&"
		return url[:len(url)-1]
	}
}

func connectToDB(url string) (*PG, error) {
	var err error
	pgOnce.Do(func() {
		config := PgxPoolConfig(url)
		connPool, err := pgxpool.NewWithConfig(context.TODO(), config)
		if err != nil {
			err = fmt.Errorf("Can't connect to database: %v", err)
		}

		connection, err := connPool.Acquire(context.TODO())
		if err != nil {
			err = fmt.Errorf("Error acquiring connnection to the database pool: %v", err)
		}

		err = connection.Ping(context.TODO())
		if err != nil {
			err = fmt.Errorf("Couldn't ping database: %v", err)
		}
		pgInstance = &PG{DB: connPool}
	})
	return pgInstance, err
}
