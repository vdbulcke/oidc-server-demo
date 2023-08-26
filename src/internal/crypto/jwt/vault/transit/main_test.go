package transit

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	vault "github.com/hashicorp/vault/api"
	"github.com/vdbulcke/oidc-server-demo/src/logger"
)

func TestJWTsigningES512(t *testing.T) {
	ctx := context.Background()

	config := vault.DefaultConfig()
	config.Address = "http://127.0.0.1:8200"
	client, err := vault.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("root-token")
	logger := logger.GetZapLogger(false)
	key, err := NewVaultTransitKey(ctx, logger, client, "transit", "rsa")
	if err != nil {
		t.Fatal(err)
	}

	standardClaims := &jwt.RegisteredClaims{
		Audience: []string{"audience"},

		ID: "someid",

		Issuer: "issuer",

		Subject: "123456789",
	}

	token := jwt.NewWithClaims(NewTransitSigningMethodRS256(), standardClaims)

	// setting kid
	token.Header["kid"] = fmt.Sprintf("transit:%s:v:%d", key.Name, key.Version)

	// get string signature
	sig, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("sig: ", sig)

}
