package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
type TransitSigningMethodRS256 struct{}

var (
	algRS256 = "RS256"
)

// NewTransitSigningMethodRS256 creates a new  TransitSigningMethodRS256
func NewTransitSigningMethodRS256() *TransitSigningMethodRS256 {
	return &TransitSigningMethodRS256{}
}

func (m *TransitSigningMethodRS256) Alg() string {
	return algRS256
}

func (m *TransitSigningMethodRS256) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algRS256)
}

func (m *TransitSigningMethodRS256) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algRS256)
}
