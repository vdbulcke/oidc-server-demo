package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
type TransitSigningMethodES256 struct{}

var (
	algES256 = "ES256"
)

// NewTransitSigningMethodES256 creates a new  TransitSigningMethodRS256
func NewTransitSigningMethodES256() *TransitSigningMethodES256 {
	return &TransitSigningMethodES256{}
}

func (m *TransitSigningMethodES256) Alg() string {
	return algES256
}

func (m *TransitSigningMethodES256) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algES256)
}

func (m *TransitSigningMethodES256) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algES256)
}
