package metrics

import (
	"net/http"
	"time"
)

func MW(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				TotalGetRequests.Inc()
			case http.MethodPost:
				TotalPostRequests.Inc()
			default:
			}

			start := time.Now()
			rw := NewCustomResponseWriter(w)

			next.ServeHTTP(rw, r)

			ObserveHistogramCodeResponseVec(rw.GetStatusCode(), time.Since(start))
			ObserveHistogramMethodResponseVec(r.Method, r.RequestURI, time.Since(start))
		},
	)
}
