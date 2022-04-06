package oidcserver

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

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

	// load user in the queue
	if c.MockUserFolder != "" {
		// go over each YAML file in mock user dir
		err := filepath.Walk(c.MockUserFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				l.Error("Error accessing file", zap.String("path", path), zap.Error(err))
				return err
			}

			// parse yaml files only
			if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				user, err := ReadMockUser(path)
				if err != nil {
					log.Fatal(err)
				}

				l.Debug("adding mock user to queue", zap.String("user", user.ID()), zap.String("path", path))
				m.QueueUser(user)
			}

			// skip other files
			return nil
		})
		if err != nil {
			l.Error("error walking dir", zap.String("dir", c.MockUserFolder), zap.Error(err))
			return nil, err
		}
	}

	return &OIDCServer{
		logger: l,
		config: c,
		m:      m,
		ln:     ln,
	}, nil
}

func ReadMockUser(path string) (*YAMLUser, error) {
	var user YAMLUser
	//log.Printf("Reading mock user from filename: %s\n", path)
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
