package transaction

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"wallet/internal/entity"
	"wallet/internal/utils/metrics"
)

type consumer interface {
	Consume(context.Context) ([]byte, error)
}

type publisher interface {
	Publish(context.Context, []byte) error
}

type Repository struct {
	consumer  consumer
	publisher publisher
}

func New(consumer consumer, publisher publisher) *Repository {
	return &Repository{
		consumer:  consumer,
		publisher: publisher,
	}
}

const (
	insertTransactionFn  = "insert transaction"
	isExistTransactionFn = "is exist transaction"
)

func (r Repository) Insert(ctx context.Context, tx pgx.Tx, tr *entity.Transaction) error {
	stmt, args, err := sq.
		Insert("transactions").
		Columns(
			"wallet_uuid",
			"idempotency_key",
			"operation",
			"amount",
			"status",
			"created_at",
			"updated_at",
		).
		Values(
			tr.WalletUUID,
			tr.IdempotencyKey,
			tr.Operation,
			tr.Amount,
			tr.Status,
			tr.CreatedAt,
			tr.UpdatedAt,
		).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if err := metrics.Tx().QueryRow(insertTransactionFn, ctx, tx, stmt, args...).Scan(&tr.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrDuplicateTransaction
			}
		}
		return err
	}

	return nil
}

func (r Repository) Exists(ctx context.Context, tx pgx.Tx, tr *entity.Transaction) (bool, error) {
	stmt, args, err := sq.
		Select("COUNT(*)").
		From("transactions").
		Where(sq.Eq{
			"wallet_uuid":     tr.WalletUUID,
			"idempotency_key": tr.IdempotencyKey,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}
	var count int64
	if err := metrics.Tx().QueryRow(isExistTransactionFn, ctx, tx, stmt, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r Repository) Publish(ctx context.Context, tr *entity.Transaction) error {
	data, err := tr.Marshall()
	if err != nil {
		return err
	}

	return r.publisher.Publish(ctx, data)
}

func (r Repository) Consume(ctx context.Context) (*entity.Transaction, error) {
	data, err := r.consumer.Consume(ctx)
	if err != nil {
		return nil, err
	}

	tr := new(entity.Transaction)
	if err := tr.Unmarshall(data); err != nil {
		return nil, err
	}

	return tr, nil
}
