package middleware

import (
	"context"
	"net/http"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/google/uuid"
)

// create a new Request ID and add it to the HTTP headers and the request context
// This should only be used by the gate service
// The Request ID is a unique ID that is passed from service to service and remains the same
// through the lifetime of an individual request
func WithReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().String()
		r.Header.Add(model.ReqIDKey.String(), reqId)
		ctx := context.WithValue(r.Context(), model.ReqIDKey, reqId)
		req := r.Clone(ctx)
		next.ServeHTTP(w, req)
	})
}

// milddleware to prepare the logger context with the incoming request ID
// All log records will include the unique request ID
func WithRequestIDLogger(handler func(rw http.ResponseWriter, r *http.Request, log logger.Logger), log logger.Logger) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqID, err := model.GetReqIDFromRequest(r)
		if err != nil {
			// we send a warning that the request doesn't contain the require ID
			// we need to set the request ID in our internal services call to throw an error
			log.
				Warnf(err.Error())
			handler(rw, r, log)
			return
		}
		loggerWithReqID := log.WithStr(model.ReqIDKey.String(), reqID)
		handler(rw, r, loggerWithReqID)
	}
}
