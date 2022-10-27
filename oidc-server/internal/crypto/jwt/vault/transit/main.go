package transit

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTKeyTransit metadata for Vault transit key properties
type JWTKeyTransit struct {
	// vault transit 'signature_algorithm'
	APISignatureAlgorithm string
	// vault transit 'hash_algorithm'
	APIHashAlgorithm string

	// Vault supported key type
	SupportedKeyType []string
}

var (

	// Maps JWT Alg to Vault Transit properties
	SupportedKeyTypeMap = map[string]*JWTKeyTransit{
		"RS256": {
			APISignatureAlgorithm: "pkcs1v15",
			APIHashAlgorithm:      "sha2-256",
			SupportedKeyType: []string{
				"rsa-2048",
				"rsa-3072",
				"rsa-4096",
			},
		},
		"RS384": {
			APISignatureAlgorithm: "pkcs1v15",
			APIHashAlgorithm:      "sha2-384",
			SupportedKeyType: []string{
				"rsa-2048",
				"rsa-3072",
				"rsa-4096",
			},
		},
		"RS512": {
			APISignatureAlgorithm: "pkcs1v15",
			APIHashAlgorithm:      "sha2-512",
			SupportedKeyType: []string{
				"rsa-2048",
				"rsa-3072",
				"rsa-4096",
			},
		},
		"ES256": {
			APISignatureAlgorithm: "pss",
			APIHashAlgorithm:      "sha2-256",
			SupportedKeyType:      []string{"ecdsa-p256"},
		},
		"ES384": {
			APISignatureAlgorithm: "pss",
			APIHashAlgorithm:      "sha2-384",
			SupportedKeyType:      []string{"ecdsa-p384"},
		},
		"ES512": {
			APISignatureAlgorithm: "pss",
			APIHashAlgorithm:      "sha2-512",
			SupportedKeyType:      []string{"ecdsa-p521"},
		},
	}
)

func validateKey(key interface{}) (*VaultTransitKey, error) {
	var transitKey *VaultTransitKey
	switch k := key.(type) {
	case *VaultTransitKey:
		transitKey = k
	default:
		return nil, jwt.ErrInvalidKeyType
	}

	return transitKey, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
