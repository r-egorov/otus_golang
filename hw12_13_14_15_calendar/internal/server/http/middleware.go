package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(fmt.Sprintf(`%s [%s] %s %s %s %d %d "%s"`,
			r.RemoteAddr,
			start.Format("01/Jan/2001:12:00:00 +0300"),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			time.Since(start)/time.Second,
			r.UserAgent(),
		))
	})
}
