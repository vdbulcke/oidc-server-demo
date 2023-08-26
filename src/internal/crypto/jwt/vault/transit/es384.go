package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
type TransitSigningMethodES384 struct{}

var (
	algES384 = "ES384"
)

// NewTransitSigningMethodES384 creates a new  TransitSigningMethodRS256
func NewTransitSigningMethodES384() *TransitSigningMethodES384 {
	return &TransitSigningMethodES384{}
}

func (m *TransitSigningMethodES384) Alg() string {
	return algES384
}

func (m *TransitSigningMethodES384) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algES384)
}

func (m *TransitSigningMethodES384) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algES384)
}
