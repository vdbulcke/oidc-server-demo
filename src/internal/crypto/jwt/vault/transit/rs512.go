package transit

// implement jwt SigningMethod interface:
// Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
// Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
// Alg() string                                                   // returns the alg identifier for this method (example: 'HS512')
type TransitSigningMethodRS512 struct{}

var (
	algRS512 = "RS512"
)

// NewTransitSigningMethodRS512 creates a new  TransitSigningMethodRS512
func NewTransitSigningMethodRS512() *TransitSigningMethodRS512 {
	return &TransitSigningMethodRS512{}
}

func (m *TransitSigningMethodRS512) Alg() string {
	return algRS512
}

func (m *TransitSigningMethodRS512) Verify(signingString, signature string, key interface{}) error {
	// get key format
	k, err := validateKey(key)
	if err != nil {
		return err
	}

	// Verify JWT signature with Transit key
	return k.JWTVerify(signingString, signature, algRS512)
}

func (m *TransitSigningMethodRS512) Sign(signingString string, key interface{}) (string, error) {

	// get key format
	k, err := validateKey(key)
	if err != nil {
		return "", err
	}
	// Sign JWT  with Transit key
	return k.JWTSign(signingString, algRS512)
}
