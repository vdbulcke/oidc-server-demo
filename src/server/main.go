package oidcserver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/oauth2-proxy/mockoidc"
	cfg "github.com/vdbulcke/oidc-server-demo/src/config"
	"github.com/vdbulcke/oidc-server-demo/src/internal/crypto"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type OIDCServer struct {
	logger *zap.Logger

	m *mockoidc.MockOIDC

	ln net.Listener

	config *cfg.OIDCServerConfig
}

func NewOIDCServer(l *zap.Logger, c *cfg.OIDCServerConfig) (*OIDCServer, error) {

	// set Supported Scopes
	if len(c.SupportedScopes) != 0 {
		mockoidc.ScopesSupported = c.SupportedScopes

	}

	// instantiate the mock oidc server
	// with the selected crypto backend
	var m *mockoidc.MockOIDC

	if c.VaultCryptoBackend == nil {

		// Create a fresh RSA Private Key for token signing
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			l.Error("generating RSA key", zap.Error(err))
			return nil, err
		}

		// Create an unstarted MockOIDC server
		m, err = mockoidc.NewServer(rsaKey)
		if err != nil {
			l.Error("generating MockOIDC", zap.Error(err))
			return nil, err
		}
	} else {

		// create vault client
		client, err := CreateVaultClient(c)
		if err != nil {
			l.Error("Error creating Vault client", zap.Error(err))
			return nil, err
		}

		ctx := context.Background()

		// create the Vault Crypto backend
		vaultCB, err := crypto.NewVaultTransitCryptoBackend(
			ctx, l, client,
			c.VaultCryptoBackend.JWTSigningAlg,
			c.VaultCryptoBackend.SyncPeriodDuration,
			c.VaultCryptoBackend.TransitMount,
			c.VaultCryptoBackend.TransitKeyName)
		if err != nil {
			l.Error("generating Crypto backend", zap.Error(err))
			return nil, err
		}

		m, err = mockoidc.NewServerWithCryptoBackend(vaultCB)
		if err != nil {
			l.Error("generating MockOIDC", zap.Error(err))
			return nil, err
		}

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

	m.SetAccessTokenTLL(c.AccessTokenTTL)
	m.SetRefreshTokenTLL(c.RefreshTokenTTL)
	if c.IssueNewRefreshTokenOnRefreshToken {
		m.EnableIssueNewRefreshTokenOnRefreshToken()
	}

	return &OIDCServer{
		logger: l,
		config: c,
		m:      m,
		ln:     ln,
	}, nil
}

func ReadMockUser(path string) (*cfg.YAMLUser, error) {
	var user cfg.YAMLUser
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

func CreateVaultClient(c *cfg.OIDCServerConfig) (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = c.VaultCryptoBackend.VaultAddress

	// tls
	tlsConfig := &vault.TLSConfig{
		// CACert:   c.CACertPEMPath,
		Insecure: false,
	}

	err := config.ConfigureTLS(tlsConfig)
	if err != nil {
		return nil, err
	}

	// create new client
	client, err := vault.NewClient(config)
	if err != nil {

		return nil, err
	}

	// Authentication Token
	client.SetToken(c.VaultCryptoBackend.VaultToken)

	return client, nil
}
