package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS384')
type TransitSigningMethodRS384 struct{}

var (
	algRS384 = "RS384"
)

// NewTransitSigningMethodRS384 creates a new  TransitSigningMethodRS384
func NewTransitSigningMethodRS384() *TransitSigningMethodRS384 {
	return &TransitSigningMethodRS384{}
}

func (m *TransitSigningMethodRS384) Alg() string {
	return algRS384
}

func (m *TransitSigningMethodRS384) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algRS384)
}

func (m *TransitSigningMethodRS384) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algRS384)
}
