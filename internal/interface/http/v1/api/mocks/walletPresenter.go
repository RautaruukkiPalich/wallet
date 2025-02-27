// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "wallet/internal/dto"

	mock "github.com/stretchr/testify/mock"
)

// WalletPresenter is an autogenerated mock type for the walletPresenter type
type WalletPresenter struct {
	mock.Mock
}

// GetBalance provides a mock function with given fields: _a0, _a1
func (_m *WalletPresenter) GetBalance(_a0 context.Context, _a1 string) (*dto.GetBalanceResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 *dto.GetBalanceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*dto.GetBalanceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *dto.GetBalanceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.GetBalanceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWallet provides a mock function with given fields: ctx
func (_m *WalletPresenter) NewWallet(ctx context.Context) (*dto.WalletResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for NewWallet")
	}

	var r0 *dto.WalletResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*dto.WalletResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *dto.WalletResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Transaction provides a mock function with given fields: _a0, _a1
func (_m *WalletPresenter) Transaction(_a0 context.Context, _a1 *dto.PostOperationRequest) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Transaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.PostOperationRequest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWalletPresenter creates a new instance of WalletPresenter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWalletPresenter(t interface {
	mock.TestingT
	Cleanup(func())
}) *WalletPresenter {
	mock := &WalletPresenter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
