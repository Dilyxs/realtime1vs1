package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLRequestType int

const (
	Query SQLRequestType = iota
	Exec
)

type SQLMessage struct {
	Query   string
	Args    any
	OutChan chan SQLResult
	ScanFn  func(rows pgx.Rows) (any, error)
	SQLType SQLRequestType
}

type SQLResult struct {
	SQLRequestType SQLRequestType
	Err            error
	Pgtag          pgconn.CommandTag
	Results        any
}

type PoolManager struct {
	Pool *pgxpool.Pool
	Chan chan SQLMessage
	Wg   *sync.WaitGroup
}

func GetPool(connString string, wg *sync.WaitGroup) *PoolManager {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("cannot start pool with config: %v", err)
	}
	config.MaxConns = 10
	config.HealthCheckPeriod = 30 * time.Second
	workerPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("[FATAL] Cannot create connection pool: %v", err)
	}
	return &PoolManager{
		Pool: workerPool,
		Chan: make(chan SQLMessage, 100),
		Wg:   wg,
	}
}

func SpawnSQLWorkers(pool *PoolManager, numWorkers int) {
	for range numWorkers {
		go Worker(pool)
	}
}

func Worker(pool *PoolManager) {
	defer pool.Wg.Done()
	for request := range pool.Chan {
		switch request.SQLType {
		case Query:
			rows, err := pool.Pool.Query(context.Background(), request.Query, request.Args.([]any)...)
			if err != nil {
				fmt.Printf("problem !: %v", err)
				request.OutChan <- SQLResult{
					SQLRequestType: request.SQLType,
					Err:            err,
					Pgtag:          pgconn.CommandTag{},
					Results:        nil,
				}
				continue
			}
			scannedrows, err := request.ScanFn(rows)
			request.OutChan <- SQLResult{
				SQLRequestType: request.SQLType,
				Err:            err,
				Pgtag:          pgconn.CommandTag{},
				Results:        scannedrows,
			}
		case Exec:
			tag, err := pool.Pool.Exec(context.Background(), request.Query, request.Args.([]any)...)
			request.OutChan <- SQLResult{
				Err:   err,
				Pgtag: tag,
			}
		}
	}
}
