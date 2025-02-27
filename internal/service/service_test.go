package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sync"
	"testing"
	"time"
	"wallet/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wallet/internal/service/mocks"
)

func TestService_NewWallet(t *testing.T) {
	ctx := context.Background()

	// Создаем моки
	walletRepoMock := &mocks.WalletRepo{}
	storeMock := &mocks.Store{}
	txMock := &mocks.MockTx{}

	// Настраиваем ожидания
	txMock.On("Begin").Return(txMock, nil)
	txMock.On("Commit").Return(nil)
	txMock.On("Rollback").Return(nil)

	storeMock.
		On("WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(pgx.Tx) error)
			_ = fn(txMock)
		}).
		Return(nil)
	walletRepoMock.
		On("Insert", ctx, mock.AnythingOfType("*mocks.MockTx"), mock.AnythingOfType("*entity.Wallet")).
		Return(nil)

	// Создаем сервис с моками
	service := &Service{
		walletRepo: walletRepoMock,
		store:      storeMock,
	}

	// Вызываем метод
	wallet, err := service.NewWallet(ctx)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, int64(0), wallet.Amount) // Проверяем, что кошелек создан с нулевым балансом

	// Проверяем, что моки были вызваны
	storeMock.AssertCalled(t, "WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error"))
	walletRepoMock.AssertCalled(t, "Insert", ctx, mock.AnythingOfType("*mocks.MockTx"), mock.AnythingOfType("*entity.Wallet"))
}

func TestService_GetBalance(t *testing.T) {
	ctx := context.Background()
	walletUUID := uuid.New()

	// Создаем моки
	walletCacheMock := &mocks.WalletCache{}
	storeMock := &mocks.Store{}
	walletRepoMock := &mocks.WalletRepo{}

	// Настраиваем ожидания
	walletCacheMock.On("GetBalance", ctx, walletUUID).Return(int64(100), nil)

	// Создаем сервис с моками
	service := &Service{
		walletCache: walletCacheMock,
		store:       storeMock,
		walletRepo:  walletRepoMock,
		mu:          &sync.RWMutex{},
	}

	// Вызываем метод
	balance, err := service.GetBalance(ctx, walletUUID)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, int64(100), balance)

	// Проверяем, что моки были вызваны
	walletCacheMock.AssertCalled(t, "GetBalance", ctx, walletUUID)
}

func TestService_NewTransaction(t *testing.T) {
	ctx := context.Background()
	walletUUID := uuid.New()

	// Создаем моки
	walletRepoMock := &mocks.WalletRepo{}
	storeMock := &mocks.Store{}
	transactionBrokerMock := &mocks.TransactionBroker{}
	txMock := &mocks.MockTx{}

	// Настраиваем ожидания
	txMock.On("Begin").Return(txMock, nil)
	txMock.On("Commit").Return(nil)
	txMock.On("Rollback").Return(nil)

	// Настраиваем ожидания
	storeMock.
		On("WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(pgx.Tx) error)
			_ = fn(txMock)
		}).
		Return(nil)
	walletRepoMock.
		On("GetByUUID", ctx, mock.AnythingOfType("*mocks.MockTx"), walletUUID).
		Return(&entity.Wallet{UUID: walletUUID, Amount: 200}, nil)
	transactionBrokerMock.
		On("Publish", ctx, mock.AnythingOfType("*entity.Transaction")).
		Return(nil)

	// Создаем сервис с моками
	service := &Service{
		walletRepo:        walletRepoMock,
		store:             storeMock,
		transactionBroker: transactionBrokerMock,
		mu:                &sync.RWMutex{},
	}

	// Создаем транзакцию
	transaction := &entity.Transaction{
		WalletUUID: walletUUID,
		Operation:  entity.Withdraw,
		Amount:     100,
	}

	// Вызываем метод
	err := service.NewTransaction(ctx, transaction)

	// Проверяем результаты
	assert.NoError(t, err)

	// Проверяем, что моки были вызваны
	storeMock.AssertCalled(t, "WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error"))
	walletRepoMock.AssertCalled(t, "GetByUUID", ctx, mock.AnythingOfType("*mocks.MockTx"), walletUUID)
	transactionBrokerMock.AssertCalled(t, "Publish", ctx, mock.AnythingOfType("*entity.Transaction"))
}

func TestService_consumeTransactions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем моки
	transactionBrokerMock := &mocks.TransactionBroker{}
	walletRepoMock := &mocks.WalletRepo{}
	transactionRepoMock := &mocks.TransactionRepo{}
	storeMock := &mocks.Store{}
	walletCacheMock := &mocks.WalletCache{}

	// Настраиваем ожидания
	transactionBrokerMock.On("Consume", ctx).Return(&entity.Transaction{
		WalletUUID:     uuid.New(),
		Operation:      entity.Deposit,
		Amount:         100,
		Status:         entity.New,
		IdempotencyKey: uuid.New(),
	}, nil)
	storeMock.
		On("WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error")).
		Return(nil)
	transactionRepoMock.
		On("Exists", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Transaction")).
		Return(false, nil)
	walletRepoMock.
		On("GetByUUID", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Wallet{Amount: 0}, nil)
	walletRepoMock.
		On("Update", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Wallet")).
		Return(nil)
	transactionRepoMock.
		On("Insert", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Transaction")).
		Return(nil)
	walletCacheMock.
		On("SetBalance", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("int64")).
		Return(nil)

	// Создаем сервис с моками
	service := &Service{
		transactionBroker: transactionBrokerMock,
		walletRepo:        walletRepoMock,
		transactionRepo:   transactionRepoMock,
		store:             storeMock,
		walletCache:       walletCacheMock,
		mu:                &sync.RWMutex{},
		walletMutex:       sync.Map{},
	}

	// Запускаем consumeTransactions в отдельной горутине
	go service.consumeTransactions(ctx)

	// Даем время для обработки
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что моки были вызваны
	transactionBrokerMock.AssertCalled(t, "Consume", ctx)
	storeMock.AssertCalled(t, "WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error"))
}

func TestService_consumeTransactions2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем моки
	transactionBrokerMock := &mocks.TransactionBroker{}
	walletRepoMock := &mocks.WalletRepo{}
	transactionRepoMock := &mocks.TransactionRepo{}
	storeMock := &mocks.Store{}
	walletCacheMock := &mocks.WalletCache{}

	// Настраиваем ожидания
	transactionBrokerMock.On("Consume", ctx).Return(&entity.Transaction{
		WalletUUID:     uuid.New(),
		Operation:      entity.Deposit,
		Amount:         100,
		Status:         entity.New,
		IdempotencyKey: uuid.New(),
	}, nil)
	storeMock.
		On("WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error")).
		Return(nil)
	transactionRepoMock.
		On("Exists", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Transaction")).
		Return(false, nil)
	walletRepoMock.
		On("GetByUUID", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Wallet{Amount: 0}, nil)
	walletRepoMock.
		On("Update", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Wallet")).
		Return(nil)
	transactionRepoMock.
		On("Insert", ctx, mock.AnythingOfType("pgx.Tx"), mock.AnythingOfType("*entity.Transaction")).
		Return(nil)
	walletCacheMock.
		On("SetBalance", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("int64")).
		Return(nil)

	// Создаем сервис с моками
	service := &Service{
		transactionBroker: transactionBrokerMock,
		walletRepo:        walletRepoMock,
		transactionRepo:   transactionRepoMock,
		store:             storeMock,
		walletCache:       walletCacheMock,
		mu:                &sync.RWMutex{},
		walletMutex:       sync.Map{},
	}

	// Запускаем consumeTransactions в отдельной горутине
	go service.consumeTransactions(ctx)

	// Даем время для обработки
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что моки были вызваны
	transactionBrokerMock.AssertCalled(t, "Consume", ctx)
	storeMock.AssertCalled(t, "WithTransact", ctx, mock.AnythingOfType("func(pgx.Tx) error"))
}
