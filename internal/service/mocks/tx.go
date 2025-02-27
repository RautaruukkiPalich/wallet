package mocks

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

// Реализуем методы pgx.Tx (например, Begin, Commit, Rollback и т.д.)
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}
