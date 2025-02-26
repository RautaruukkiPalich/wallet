package transaction

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"wallet/internal/entity"
)

type Repository struct{}

func New() *Repository {
	return &Repository{}
}

func (r Repository) Insert(ctx context.Context, tx pgx.Tx, tr *entity.Transaction) error {
	stmt, args, err := sq.
		Insert("transactions").
		Columns(
			"wallet_uuid",
			"idempotency_key",
			"operation",
			"amount",
			"created_at",
			"updated_at",
		).
		Values(
			tr.WalletUUID,
			tr.IdempotencyKey,
			tr.Operation,
			tr.Amount,
			tr.CreatedAt,
			tr.UpdatedAt,
		).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if err := tx.QueryRow(ctx, stmt, args...).Scan(&tr.ID); err != nil {
		return err
	}

	return nil
}
