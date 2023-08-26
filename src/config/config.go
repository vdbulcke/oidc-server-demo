package oidcserver

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

var (
	TransitDefaultMount      = "transit"
	VaultDefaultSyncDuration = "5m"
)

type OIDCServerConfig struct {
	ClientID     string `yaml:"client_id"  validate:"required"`
	ClientSecret string `yaml:"client_secret" `

	SupportedScopes               []string `yaml:"supported_scopes" `
	PKCEChallengeMethodsSupported []string `yaml:"pkce_challenge_methods" `
	IssuerBaseUrl                 string   `yaml:"issuer_base_url" `

	VaultCryptoBackend *VaultCryptoBackendConfig `yaml:"vault_crypto_backend" validate:"omitempty"`

	MockUser YAMLUser `yaml:"mock_user" validate:"required"`
	// Folder where to find mocked user if not defined the user in mock_user will be returned
	MockUserFolder string `yaml:"mock_user_folder"`

	// Listen Address
	ListenAddress string
	// Listen Port
	ListenPort int

	// internal
	AccessLog bool
	Debug     bool
}

type VaultCryptoBackendConfig struct {
	VaultAddress string `yaml:"address"  validate:"required"`
	VaultToken   string `yaml:"token"  validate:"required"`

	TransitKeyName string `yaml:"transit_key"  validate:"required"`
	TransitMount   string `yaml:"transit_mount" validate:"required"`
	JWTSigningAlg  string `yaml:"jwt_signing_alg"  validate:"required,oneof=RS256 RS384 RS512 ES256 ES384 ES512"`

	SyncPeriodDuration string `yaml:"sync_duration" validate:"required"`
}

// ValidateConfig validate config
func ValidateConfig(config *OIDCServerConfig) bool {

	validate := validator.New()
	errs := validate.Struct(config)

	if errs == nil {
		return true
	}

	for _, e := range errs.(validator.ValidationErrors) {
		fmt.Println(e)
	}

	return false

}

// ParseConfig Parse config file
func ParseConfig(configFile string) (*OIDCServerConfig, error) {

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := OIDCServerConfig{}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	// override properties with env variable if declared
	parseEnv(&config)

	// handle default value
	if config.VaultCryptoBackend != nil {

		if config.VaultCryptoBackend.TransitMount == "" {
			config.VaultCryptoBackend.TransitMount = TransitDefaultMount
		}

		if config.VaultCryptoBackend.SyncPeriodDuration == "" {
			config.VaultCryptoBackend.SyncPeriodDuration = VaultDefaultSyncDuration
		}

		_, err := time.ParseDuration(config.VaultCryptoBackend.SyncPeriodDuration)
		if err != nil {
			return nil, err
		}

	}

	// return Parse config struct
	return &config, nil

}

// parseEnv Parse config file
func parseEnv(config *OIDCServerConfig) {

	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")

	if clientID != "" {
		config.ClientID = clientID
	}

	if clientSecret != "" {
		config.ClientSecret = clientSecret
	}

	if config.VaultCryptoBackend != nil {
		VAULT_ADDR := os.Getenv("VAULT_ADDR")
		VAULT_TOKEN := os.Getenv("VAULT_TOKEN")
		if VAULT_ADDR != "" {
			config.VaultCryptoBackend.VaultAddress = VAULT_ADDR
		}
		if VAULT_TOKEN != "" {
			config.VaultCryptoBackend.VaultToken = VAULT_TOKEN
		}

	}

}

// ListenAddress returns http listener address
func (c *OIDCServerConfig) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
}
