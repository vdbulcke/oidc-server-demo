package oidcserver

import (
	"os"
	"os/signal"

	"go.uber.org/zap"
)

func (s *OIDCServer) StartServer(debug bool) error {

	// add middleware
	s.logger.Debug("Adding InfiniteMockUserMiddleware")
	err := s.m.AddMiddleware(s.InfiniteMockUserMiddleware)
	if err != nil {
		s.logger.Error("adding middleware", zap.String("middleware", "InfiniteMockUserMiddleware"), zap.Error(err))
		return err
	}

	if debug {
		s.logger.Debug("Adding DebugLoggerMiddleware")
		err := s.m.AddMiddleware(s.DebugLoggerMiddleware)
		if err != nil {
			s.logger.Error("adding middleware", zap.String("middleware", "DebugLoggerMiddleware"), zap.Error(err))
			return err
		}
	}

	// tlsConfig can be nil if you want HTTP
	err = s.m.Start(s.ln, nil)
	if err != nil {
		s.logger.Error("starting mockoidc server", zap.Error(err))
		return err
	}
	//nolint
	defer s.m.Shutdown()

	// print config
	suggar := s.logger.Sugar()
	suggar.Info("starting server")
	suggar.Infow("server config",
		"Issuer", s.m.Issuer(),
	)

	// trap sigterm or interupt and gracefully shutdown the server
	close := make(chan os.Signal, 1)
	signal.Notify(close, os.Interrupt)

	sig := <-close
	s.logger.Info("Got signal", zap.Any("sig", sig))

	return nil
}
