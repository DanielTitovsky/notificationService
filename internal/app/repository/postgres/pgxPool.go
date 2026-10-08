package app_repository_postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, config Config) Pool {
	connectinString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DataBase,
	)

	poolConfig, err := pgxpool.ParseConfig(connectinString)

	if err != nil {
		log.Fatalf("Parse config error: %v", err)
	}

	poolConfig.MaxConns = int32(config.MaxConns)
	poolConfig.MinConns = int32(config.MinConns)
	poolConfig.MaxConnLifetime = config.MaxConnLifeTime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	err = pool.Ping(ctx)

	if err != nil {
		log.Fatalf("Failed created postgres pool: %v", err)
	}

	return Pool{
		Pool: pool,
	}
}
