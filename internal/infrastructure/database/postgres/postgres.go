package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, cfg DBConfig) (*Store, error) {
	psqlURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Pass, cfg.Database)

	if cfg.Port != "" {
		psqlURL += fmt.Sprintf(" port=%s", cfg.Port)
	}

	pool, err := pgxpool.New(ctx, psqlURL)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrConnectToDB, err)
	}
	if pool.Ping(ctx) != nil {
		return nil, ErrPingToDB
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) WithTransact(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		log.Println("err creating transaction", err)
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = fn(tx)
	if err != nil {
		log.Println("err executing transaction: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("err committing transaction", err)
		return err
	}

	return nil
}

//func (s *Store) WithSerializableTransact(ctx context.Context, fn func(pgx.Tx) error) error {
//	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
//		IsoLevel: pgx.Serializable,
//	})
//
//	if err != nil {
//		fmt.Println("err creating transaction", err)
//		return err
//	}
//	defer func() {
//		_ = tx.Rollback(ctx)
//	}()
//
//	err = fn(tx)
//	if err != nil {
//		fmt.Println("err executing transaction: ", err)
//		return err
//	}
//
//	if err := tx.Commit(ctx); err != nil {
//		fmt.Println("err committing transaction", err)
//		return err
//	}
//
//	return nil
//}
