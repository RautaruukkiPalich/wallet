package metrics

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
	GetStatusCode() int
}

type CustomResponseWriter struct {
	rw         http.ResponseWriter
	statusCode int
}

func NewCustomResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &CustomResponseWriter{w, http.StatusOK}
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.rw.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.rw.Write(b)
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *CustomResponseWriter) GetStatusCode() int {
	return w.statusCode
}
