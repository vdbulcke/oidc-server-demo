package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
type TransitSigningMethodES512 struct{}

var (
	algES512 = "ES512"
)

// NewTransitSigningMethodES512 creates a new  TransitSigningMethodRS256
func NewTransitSigningMethodES512() *TransitSigningMethodES512 {
	return &TransitSigningMethodES512{}
}

func (m *TransitSigningMethodES512) Alg() string {
	return algES512
}

func (m *TransitSigningMethodES512) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algES512)
}

func (m *TransitSigningMethodES512) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algES512)
}
