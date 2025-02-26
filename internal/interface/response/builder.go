package response

import (
	"net/http"
	"wallet/internal/dto"
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

// TODO: error handling
func (b *Builder) HandleError(err error) *Builder {
	return b.WithCode(http.StatusBadRequest).WithError(err)
}

func (b *Builder) Build() *Response {
	if b.response.Code == 0 {
		b.response.Code = http.StatusOK
	}
	return &b.response
}
