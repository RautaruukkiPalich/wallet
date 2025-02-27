package response

import (
	"errors"
	"net/http"
	"wallet/internal/dto"
	"wallet/internal/entity"
	"wallet/internal/presenter"
	walletRepository "wallet/internal/repository/wallet"
	"wallet/internal/service"
)

type Builder struct {
	response Response
}

func Resp() *Builder {
	return &Builder{response: NewResponse()}
}

func (b *Builder) WithCode(code int) *Builder {
	b.response.Code = code
	return b
}

func (b *Builder) WithPayload(payload any) *Builder {
	b.response.Payload = payload
	return b
}

func (b *Builder) WithHeader(key string, value string) *Builder {
	b.response.Headers.Set(key, value)
	return b
}

func (b *Builder) WithHeaders(headers http.Header) *Builder {
	for key, value := range headers {
		for _, v := range value {
			b.response.Headers.Add(key, v)
		}
	}
	return b
}

func (b *Builder) WithMapHeaders(headers map[string]string) *Builder {
	for key, value := range headers {
		b.response.Headers.Set(key, value)
	}
	return b
}

func (b *Builder) WithError(err error) *Builder {
	b.response.Payload = dto.ErrorResponse{Message: err.Error()}
	return b
}

func (b *Builder) HandleError(err error) *Builder {
	if errors.Is(err, service.ErrInvalidUUID) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, presenter.ErrInvalidUUID) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, walletRepository.ErrWalletNotFound) {
		return b.WithCode(http.StatusNotFound).WithError(err)
	}

	if errors.Is(err, entity.ErrWalletUUIDIsEmpty) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, entity.ErrNotEnoughFunds) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, entity.ErrInvalidOperationType) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, entity.ErrInvalidOperationUUID) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, entity.ErrInvalidStatus) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}
	if errors.Is(err, entity.ErrAmountIsOrBelowZero) {
		return b.WithCode(http.StatusBadRequest).WithError(err)
	}

	return b.WithCode(http.StatusInternalServerError)
}

func (b *Builder) Build() *Response {
	if b.response.Code == 0 {
		b.response.Code = http.StatusOK
	}
	return &b.response
}
