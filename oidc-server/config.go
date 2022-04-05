package oidcserver

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

type OIDCServerConfig struct {
	ClientID     string `yaml:"client_id"  validate:"required"`
	ClientSecret string `yaml:"client_secret" `

	SupportedScopes               []string `yaml:"supported_scopes" `
	PKCEChallengeMethodsSupported []string `yaml:"pkce_challenge_methods" `
	IssuerBaseUrl                 string   `yaml:"issuer_base_url" `

	MockUser YAMLUser `yaml:"mock_user" validate:"required"`
	// Folder where to find mocked user if not defined the user in mock_user will be returned
	MockUserFolder string `yaml:"mock_user_folder"`

	// Listen Address
	ListenAddress string
	// Listen Port
	ListenPort int
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

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
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

}

// ListenAddress returns http listener address
func (c *OIDCServerConfig) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
}
