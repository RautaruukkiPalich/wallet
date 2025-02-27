package metrics

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxStruct struct {
}

func Tx() *TxStruct {
	return &TxStruct{}
}

func (t TxStruct) QueryRow(queryName string, ctx context.Context, tx querier, sql string, args ...any) pgx.Row {
	start := time.Now()
	defer func() {
		ObserveHistogramTimeQueryCounter(
			queryName, time.Since(start),
		)
	}()
	return tx.QueryRow(ctx, sql, args...)
}

func (t TxStruct) Query(queryName string, ctx context.Context, tx querier, sql string, args ...any) (pgx.Rows, error) {
	start := time.Now()
	defer func() {
		ObserveHistogramTimeQueryCounter(
			queryName, time.Since(start))
	}()
	return tx.Query(ctx, sql, args...)
}

func (t TxStruct) Exec(queryName string, ctx context.Context, tx executor, sql string, args ...any) (pgconn.CommandTag, error) {
	start := time.Now()
	defer func() {
		ObserveHistogramTimeQueryCounter(
			queryName, time.Since(start))
	}()
	return tx.Exec(ctx, sql, args...)
}

var histogramTimeQueryCounter = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "hist",
		Name:      "query_time",
		Help:      "sql query time",
	},
	[]string{
		"query",
	},
)

func ObserveHistogramTimeQueryCounter(query string, dur time.Duration) {
	histogramTimeQueryCounter.WithLabelValues(
		query,
	).Observe(dur.Seconds())
}
