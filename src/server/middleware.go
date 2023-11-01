package oidcserver

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// InfiniteMockUserMiddleware add new mock user to queue on authorize endpoint
func (s *OIDCServer) InfiniteMockUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if strings.HasSuffix(req.URL.Path, "authorize") {
			if len(s.m.UserQueue.Queue) == 0 {
				s.logger.Debug("adding mock user to queue")
				s.m.QueueUser(&s.config.MockUser)
			}
		}

		err := req.ParseForm()
		if err == nil {

			if strings.HasSuffix(req.URL.Path, "token") && req.FormValue("grant_type") == "client_credentials" {
				if len(s.m.UserQueue.Queue) == 0 {
					s.logger.Debug("adding mock user to queue")
					s.m.QueueUser(&s.config.MockUser)
				}
			}
		}

		// custom middleware logic here...
		next.ServeHTTP(rw, req)
	})
}

// DebugLoggerMiddleware logs every request
func (s *OIDCServer) DebugLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// defer req.Body.Close()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			s.logger.Error("reading body", zap.Error(err))
			// custom middleware logic here...
			next.ServeHTTP(rw, req)
		}

		// don't close body here
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		s.logger.Info("Request",
			zap.String("method", req.Method),
			zap.String("url", req.URL.Path),
			zap.Any("headers", req.Header),
			zap.Any("param", req.Form),
			zap.String("body", string(body)),
		)
		// custom middleware logic here...
		next.ServeHTTP(rw, req)
	})
}
