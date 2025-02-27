package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Code    int
	Payload any
	Headers http.Header
}

func NewResponse() Response {
	return Response{
		Code:    http.StatusOK, // Дефолтный код ответа
		Payload: nil,           // Дефолтный payload
		Headers: http.Header{}, // Дефолтные заголовки
	}
}

func (b *Response) Write(w http.ResponseWriter) {

	var payload []byte
	var err error

	if b.Payload != nil {
		payload, err = b.preparePayload()
		if err != nil {
			log.Println("json marshal failed: ", err.Error())
		}
	}

	for k, v := range b.Headers {
		w.Header().Set(k, strings.Join(v, ","))
	}

	w.WriteHeader(b.Code)

	_, err = w.Write(payload)
	if err != nil {
		log.Println("write response failed: ", err.Error())
	}

}

func (b *Response) preparePayload() ([]byte, error) {
	jsoned, err := json.Marshal(b.Payload)
	if err != nil {
		return []byte{}, err
	}

	b.Headers.Set("Content-Type", "application/json")
	b.Headers.Set("Content-Length", strconv.Itoa(len(jsoned)))

	return jsoned, nil
}
