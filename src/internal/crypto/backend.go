package crypto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/golang-jwt/jwt/v4"
	vault "github.com/hashicorp/vault/api"
	"github.com/vdbulcke/oidc-server-demo/src/internal/crypto/jwt/vault/transit"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
)

type VaultTransitCryptoBackend struct {
	Key              *transit.VaultTransitKey
	JWTSigningMethod jwt.SigningMethod
	logger           *zap.Logger
}

// NewVaultTransitCryptoBackend Create an new Vault Crypto Backend
func NewVaultTransitCryptoBackend(ctx context.Context, l *zap.Logger, client *vault.Client, alg, schedulerDuration, transitMount, transitKey string) (*VaultTransitCryptoBackend, error) {

	// create the transit key
	key, err := transit.NewVaultTransitKey(ctx, l, client, transitMount, transitKey)
	if err != nil {
		return nil, err
	}

	_, err = key.Validate(alg)
	if err != nil {
		return nil, err
	}

	// select the corresponding jwt.SigningMethod interface
	var method jwt.SigningMethod
	switch alg {
	case "RS256":
		method = transit.NewTransitSigningMethodRS256()
	case "RS384":
		method = transit.NewTransitSigningMethodRS384()
	case "RS512":
		method = transit.NewTransitSigningMethodRS512()
	case "ES256":
		method = transit.NewTransitSigningMethodES256()
	case "ES384":
		method = transit.NewTransitSigningMethodES384()
	case "ES512":
		method = transit.NewTransitSigningMethodES512()

	default:
		return nil, fmt.Errorf("not implemented method %s", alg)
	}

	// create the Crypto backend
	v := &VaultTransitCryptoBackend{
		Key:              key,
		JWTSigningMethod: method,
		logger:           l,
	}

	// start the scheduler
	err = v.StartScheduler(schedulerDuration)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *VaultTransitCryptoBackend) JWKS() ([]byte, error) {

	jwks := []jose.JSONWebKey{}
	// for each public keys
	for _, pub := range v.Key.PublicKeys {

		jwk := jose.JSONWebKey{
			Use:       "sig",
			Algorithm: v.JWTSigningMethod.Alg(),
			Key:       pub.PublicKey,
			KeyID:     pub.KeyID(),
		}

		jwks = append(jwks, jwk)

	}

	jwksSet := &jose.JSONWebKeySet{
		Keys: jwks,
	}

	return json.Marshal(jwksSet)

}

func (v *VaultTransitCryptoBackend) SignJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(v.JWTSigningMethod, claims)

	token.Header["kid"] = v.Key.LatestKeyID()

	return token.SignedString(v.Key)
}

func (v *VaultTransitCryptoBackend) VerifyJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		token.Method = v.JWTSigningMethod
		kid, ok := token.Header["kid"].(string)
		if ok {
			// set the version number from kid for verification
			err := v.Key.SetSignVersionFromKeyID(kid)
			if err != nil {
				return nil, err
			}

			return v.Key, nil
		}
		return nil, errors.New("token kid does not match or is not present")
	})
}

func (v *VaultTransitCryptoBackend) StartScheduler(duration string) error { // scheduler
	frequencyDuration, err := time.ParseDuration(duration)
	if err != nil {
		v.logger.Error("Error parsing  SchedulerPeriodDuration", zap.Error(err))
		return err
	}

	// start scheduler with job
	_, err = scheduler.Every(int(frequencyDuration.Seconds())).Seconds().Run(v.SyncJob)
	if err != nil {
		v.logger.Error("Error starting scheduler", zap.Error(err))
		return err
	}

	return nil
}

func (v *VaultTransitCryptoBackend) SyncJob() { // scheduler
	v.logger.Debug("running syncJob")

	err := v.Key.SyncKeyInfo()
	if err != nil {

		v.logger.Error("error synchronizing key info in scheduler", zap.Error(err))
	}

}
