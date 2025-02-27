package wallet

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strconv"
	"time"
	"wallet/internal/entity"
	"wallet/internal/utils/metrics"
)

type cache interface {
	SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Repository struct {
	cache cache
}

func New(cache cache) *Repository {
	return &Repository{
		cache: cache,
	}
}

const (
	insertWalletFn    = "insert wallet"
	updateWalletFn    = "update wallet"
	getWalletByUUIDFn = "get wallet by uuid"
)

func (r Repository) Insert(ctx context.Context, tx pgx.Tx, w *entity.Wallet) error {
	stmt, args, err := sq.
		Insert("wallets").
		Columns(
			"amount",
			"version",
			"created_at",
			"updated_at",
		).
		Values(
			w.Amount,
			w.Version,
			w.CreatedAt,
			w.UpdatedAt,
		).
		Suffix("RETURNING \"uuid\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if err := metrics.Tx().QueryRow(insertWalletFn, ctx, tx, stmt, args...).Scan(&w.UUID); err != nil {
		return err
	}

	return nil
}

func (r Repository) Update(ctx context.Context, tx pgx.Tx, w *entity.Wallet) error {
	stmt, args, err := sq.Update("wallets").
		Set("amount", w.Amount).
		Set("version", w.Version+1).
		Set("updated_at", w.UpdatedAt).
		Where(sq.Eq{
			"uuid":    w.UUID,
			"version": w.Version,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	res, err := metrics.Tx().Exec(updateWalletFn, ctx, tx, stmt, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func (r Repository) GetByUUID(ctx context.Context, tx pgx.Tx, uid uuid.UUID) (*entity.Wallet, error) {
	stmt, args, err := sq.Select(
		"uuid",
		"amount",
		"version",
		"created_at",
		"updated_at",
	).
		From("wallets").
		Where(sq.Eq{"uuid": uid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	w := new(entity.Wallet)

	if err := metrics.Tx().QueryRow(getWalletByUUIDFn, ctx, tx, stmt, args...).
		Scan(
			&w.UUID,
			&w.Amount,
			&w.Version,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	return w, nil
}

func (r Repository) SetBalance(ctx context.Context, uid uuid.UUID, balance int64) error {
	return r.cache.SetWithTTL(ctx, uid.String(), strconv.FormatInt(balance, 10), time.Second*5)
}
func (r Repository) GetBalance(ctx context.Context, uid uuid.UUID) (int64, error) {
	res, err := r.cache.Get(ctx, uid.String())
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(res, 10, 64)
}
