package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
	JSONFormat
	calls sync.Map
	clock timer
	err string
	
}

func (l *logger) CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start the watch to measure latency
		start := time.Now()

		// Get the X-Request-ID header from the request
		requestID := r.Header.Get("X-Request-ID")

		// Add the X-Request-ID to the context
		ctx := context.WithValue(r.Context(), requestIdKey, requestID)

		// Add the current time to the context
		ctx = context.WithValue(ctx, timeKey, time.Now())

		// Create a log entry with the request ID and current time
		logEntry := fmt.Sprintf("Request ID: %s, Time: %s", requestID, time.Now())

		// Do something with the log entry, such as logging it or passing it to downstream services

		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))

		// Measure the latency
		latency := time.Since(start)
		fmt.Printf("Latency: %s\n", latency)
	})
}


		// Add the X-Request-ID to the context
		ctx = context.WithValue(ctx, requestIdKey, requestID)

		// Create a log entry with the request ID and current time
		logEntry := fmt.Sprintf("Request ID: %s, Time: %s", requestID, time.Now())

		// Do something with the log entry, such as logging it or passing it to downstream services

		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
		requestID := r.Header.Get("X-Request-ID")

		// Do something with the X-Request-ID, such as logging or passing it to downstream services

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	}
}
