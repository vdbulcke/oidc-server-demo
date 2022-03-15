package oidcserver

import (
	"crypto/rand"
	"crypto/rsa"
	"net"

	"github.com/oauth2-proxy/mockoidc"
	"go.uber.org/zap"
)

type OIDCServer struct {
	logger *zap.Logger

	m *mockoidc.MockOIDC

	ln net.Listener

	config *OIDCServerConfig
}

func NewOIDCServer(l *zap.Logger, c *OIDCServerConfig) (*OIDCServer, error) {

	// set Supported Scopes
	if len(c.SupportedScopes) != 0 {
		mockoidc.ScopesSupported = c.SupportedScopes

	}

	// Create a fresh RSA Private Key for token signing
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		l.Error("generating RSA key", zap.Error(err))
		return nil, err
	}

	// Create an unstarted MockOIDC server
	m, err := mockoidc.NewServer(rsaKey)
	if err != nil {
		l.Error("generating MockOIDC", zap.Error(err))
		return nil, err
	}

	// Create the net.Listener on the exact IP:Port you want
	ln, err := net.Listen("tcp", c.GetListenAddress())
	if err != nil {
		l.Error("generating http listener", zap.Error(err))
		return nil, err
	}

	// setting mock oidc config
	m.ClientID = c.ClientID
	m.ClientSecret = c.ClientSecret
	m.IssuerBaseUrl = c.IssuerBaseUrl
	if len(c.PKCEChallengeMethodsSupported) != 0 {
		m.CodeChallengeMethodsSupported = c.PKCEChallengeMethodsSupported
	}

	return &OIDCServer{
		logger: l,
		config: c,
		m:      m,
		ln:     ln,
	}, nil
}
