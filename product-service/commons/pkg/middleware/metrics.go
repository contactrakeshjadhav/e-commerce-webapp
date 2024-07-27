package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bsm/openmetrics"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/urfave/negroni"
)

// HTTP metrics middleware
func Metrics(requestCount openmetrics.CounterFamily, responseTime openmetrics.HistogramFamily, log logger.Logger, serviceName string) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID, err := model.GetReqIDFromCTX(r.Context())
			//TODO ERROR if the request correlation id is not found
			if err != nil {
				// TODO -- ERROR if there is no request ID
				reqID = "UNKNOWN"
			}

			start := time.Now()
			// DEBUG log entry into request
			log.Debugf("Request Entry, %v, %v, %v, %v, %v", reqID, serviceName, r.Method, r.URL.Path, start.Format(time.RFC3339))
			// TODO - it would be nice for debug to be able to print out the request body

			rw := negroni.NewResponseWriter(w)
			// Serve the request
			next.ServeHTTP(rw, r)

			httpRespCode := strconv.Itoa(rw.Status())
			user, ok := model.GetUserFromCTX(r.Context())
			if !ok {
				// Do not ERROR here because it could be gate unauthenticate request
				//		  later, myabe we will improve this
				user.Username = "UNAUTHENTICATED"
			}

			//increment number of http requests
			requestCount.With(httpRespCode, reqID, user.Username, serviceName, r.Method, r.URL.Path).Add(1)

			//calculate response time and observe it for Histogram
			dur := time.Since(start)
			respTimeSeconds := float64(dur.Microseconds()) / (1000.00 * 1000.00)
			responseTime.With(httpRespCode, reqID, user.Username, serviceName, r.Method, r.URL.Path).Observe(respTimeSeconds)

			log.Infof("Request Metrics, %v, %v, %v, %v, %v, %v, %f", reqID, user.Username, serviceName, r.Method, r.URL.Path, httpRespCode, dur.Seconds())
			//TODO - INFO log exit of request with response time, User, request ID, serviceName
			//
		})
	}
	return
}
