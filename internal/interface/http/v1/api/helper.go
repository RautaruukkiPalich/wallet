package api

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

func getFromBody(r *http.Request, dest any) error {
	defer func() {
		_ = r.Body.Close()
	}()
	return jsoniter.NewDecoder(r.Body).Decode(dest)
}
