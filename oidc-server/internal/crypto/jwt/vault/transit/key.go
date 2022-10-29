package transit

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/savaki/jq"
	"go.uber.org/zap"
)

type TransitPublicKey struct {
	// pub key for JWKS
	PublicKey crypto.PublicKey
	// Version
	Version int

	// Name
	Name string
}

func NewTransitPublicKey(pub crypto.PublicKey, v int, name string) *TransitPublicKey {
	return &TransitPublicKey{
		PublicKey: pub,
		Version:   v,
		Name:      name,
	}
}

func (k *TransitPublicKey) KeyID() string {
	return fmt.Sprintf("transit:%s:v:%d", k.Name, k.Version)
}

type VaultTransitKey struct {
	// transit backend mount
	MountPath string
	// transit Key Name
	Name string

	// 'key' type
	Type string

	// Version
	Version int

	// Set sig version
	SigVersion int

	// List of public keys
	PublicKeys []*TransitPublicKey

	// the vault api client
	client *vault.Client
	// context for vault client
	ctx context.Context

	// logger
	logger *zap.Logger
}

func NewVaultTransitKey(ctx context.Context, l *zap.Logger, client *vault.Client, mount string, name string) (*VaultTransitKey, error) {

	// instantiate Transit vault key
	k := &VaultTransitKey{
		Name:      name,
		ctx:       ctx,
		client:    client,
		MountPath: mount,
		logger:    l,
	}

	// fetch key info
	err := k.SyncKeyInfo()
	if err != nil {
		return nil, err
	}

	return k, nil

}

// LatestKeyID return latest kid for adding in jwt header
func (k *VaultTransitKey) LatestKeyID() string {
	return fmt.Sprintf("transit:%s:v:%d", k.Name, k.Version)
}

func (k *VaultTransitKey) SetSignVersionFromKeyID(kid string) error {
	kidParts := strings.Split(kid, ":")
	if len(kidParts) != 4 {
		k.logger.Error("error parsing kid", zap.String("kid", kid))
		return fmt.Errorf("error parsing kid: %s", kid)
	}

	v, err := strconv.Atoi(kidParts[3])
	if err != nil {
		k.logger.Error("error parsing kid version", zap.String("kid", kid), zap.Error(err))
		return err
	}

	k.SigVersion = v
	return nil
}

// SyncKeyInfo read transit key info
func (k *VaultTransitKey) SyncKeyInfo() error {

	// transit key api path
	keyPath := fmt.Sprintf("%s/keys/%s", k.MountPath, k.Name)
	// read transit key
	keyInfo, err := k.client.Logical().ReadWithContext(k.ctx, keyPath)
	if err != nil {
		return err
	}

	k.logger.Debug("transit read key response", zap.Any("resp", keyInfo))

	if keyInfo == nil {
		k.logger.Error("No response for transit key read", zap.String("key", keyPath))
		return fmt.Errorf("error reading key %s", keyPath)
	}

	// parse key type
	keyType, ok := keyInfo.Data["type"].(string)
	if !ok {
		k.logger.Debug("Key type not found in transit read response", zap.Any("resp", keyInfo))
		return fmt.Errorf("key type not found for %s", keyPath)
	}

	keyVersionJson, ok := keyInfo.Data["latest_version"].(json.Number)
	if !ok {
		k.logger.Debug("Key latest_version not found in transit read response", zap.Any("resp", keyInfo))
		return fmt.Errorf("key latest_version not found for %s", keyPath)
	}

	keyVersion, err := keyVersionJson.Int64()
	if err != nil {
		return err
	}

	minVersionJson, ok := keyInfo.Data["min_decryption_version"].(json.Number)
	if !ok {
		k.logger.Debug("Key min_decryption_version not found in transit read response", zap.Any("resp", keyInfo))
		return fmt.Errorf("key min_decryption_version not found for %s", keyPath)
	}

	minVersion, err := minVersionJson.Int64()
	if err != nil {
		return err
	}

	pubKeys := []*TransitPublicKey{}

	// for each pub keys within range min version to latest_version
	for i := int(minVersion); i <= int(keyVersion); i++ {

		pub, err := k.GetPublicKeyFromTransitResponse(keyInfo, i)
		if err != nil {
			k.logger.Error("error parsing pub key from response", zap.Error(err))
			return err
		}

		pubKeys = append(pubKeys, NewTransitPublicKey(pub, i, k.Name))

	}

	k.Type = keyType
	k.Version = int(keyVersion)
	k.SigVersion = int(keyVersion)
	k.PublicKeys = pubKeys

	return nil

}

// JWTSign signs JWT signingString with the alg
func (k *VaultTransitKey) JWTSign(signingString, alg string) (string, error) {

	// Get properties for this JWT signing alg
	jwtKeyProp, err := k.Validate(alg)
	if err != nil {
		k.logger.Error("Error verifying with Transit key", zap.Error(err))
		return "", err
	}

	// get transit API properties
	sigAlg := jwtKeyProp.APISignatureAlgorithm
	hashAlg := jwtKeyProp.APIHashAlgorithm
	marshallingAlg := "jws"
	prehashed := false

	// Sign byte payload with transit API
	transitSig, err := k.Sign([]byte(signingString), sigAlg, hashAlg, marshallingAlg, prehashed)
	if err != nil {
		k.logger.Error("Error signing with Transit key", zap.Error(err))
		return "", err
	}

	// check format
	if !strings.HasPrefix(transitSig, "vault:v") {
		return "", fmt.Errorf("invalid signature expecting prefix 'vault:v' but got %s", transitSig)
	}

	// Vault transit signature are prefixed with
	// 'vault:vX:' indicating the version of key used for this signature
	// split on ':' and return last part
	sigParts := strings.Split(transitSig, ":")
	return sigParts[2], nil

}

// JWTVerify verifies the 'signature' matches the 'signingString' for the 'alg'
func (k *VaultTransitKey) JWTVerify(signingString, signature, alg string) error {

	// Get properties for this JWT signing alg
	jwtKeyProp, err := k.Validate(alg)
	if err != nil {
		k.logger.Error("Error verifying with Transit key", zap.Error(err))
		return err
	}

	// get transit API properties
	sigAlg := jwtKeyProp.APISignatureAlgorithm
	hashAlg := jwtKeyProp.APIHashAlgorithm
	marshallingAlg := "jws"
	prehashed := false

	// verify signing string bytes against the signature using transit API
	ok, err := k.Verify([]byte(signingString), signature, sigAlg, hashAlg, marshallingAlg, prehashed)
	if err != nil {
		k.logger.Error("Error verifying with Transit key", zap.Error(err))
		return err
	}

	if !ok {
		k.logger.Error("invalid transit signature", zap.String("signingString", signingString), zap.String("signature", signature))
		return fmt.Errorf("invalid signature %s for input %s", signingString, signature)
	}

	return nil

}

// Sign byte payload, and returns "signature" output of transit sign api
func (k *VaultTransitKey) Sign(inputBytes []byte, apiSigAlg string, apiHashAlg string, marshallingAlg string, prehashed bool) (string, error) {

	args := map[string]interface{}{
		// transit required input to base64 encoded
		"input":                base64.StdEncoding.EncodeToString(inputBytes),
		"signature_algorithm":  apiSigAlg,
		"marshaling_algorithm": marshallingAlg,
		"prehashed":            prehashed,
		"key_version":          k.Version,
	}

	// sign with transit API
	signingPath := fmt.Sprintf("%s/sign/%s/%s", k.MountPath, k.Name, apiHashAlg)
	transitResp, err := k.client.Logical().WriteWithContext(k.ctx, signingPath, args)
	if err != nil {
		return "", err
	}

	if transitResp == nil {
		k.logger.Error("No response for transit signing ", zap.String("key", signingPath))
		return "", fmt.Errorf("error signing key %s", signingPath)
	}

	sig, ok := transitResp.Data["signature"].(string)
	if !ok {
		return "", fmt.Errorf("unable to get 'signature' from transit response")
	}

	return sig, nil
}

// verify byte payload,  and signature (without the "vault:v1")
//
//	returns true if signature is valid for byte payload
func (k *VaultTransitKey) Verify(inputBytes []byte, signature string, apiSigAlg string, apiHashAlg string, marshallingAlg string, prehashed bool) (bool, error) {

	args := map[string]interface{}{
		// transit required input to base64 encoded
		"input":                base64.StdEncoding.EncodeToString(inputBytes),
		"signature":            fmt.Sprintf("vault:v%d:%s", k.SigVersion, signature),
		"signature_algorithm":  apiSigAlg,
		"marshaling_algorithm": marshallingAlg,
		"prehashed":            prehashed,
	}

	// sign with transit API
	signingPath := fmt.Sprintf("%s/verify/%s/%s", k.MountPath, k.Name, apiHashAlg)
	transitResp, err := k.client.Logical().WriteWithContext(k.ctx, signingPath, args)
	if err != nil {
		return false, err
	}

	if transitResp == nil {
		k.logger.Error("No response for transit verifying ", zap.String("key", signingPath))
		return false, fmt.Errorf("error verifying key %s", signingPath)
	}

	sigValid, ok := transitResp.Data["valid"].(bool)
	if !ok {
		return false, fmt.Errorf("unable to get 'valid' from transit response")
	}

	return sigValid, nil
}

func (k *VaultTransitKey) SetSigKeyVersion(v int) {
	k.SigVersion = v
}

// GetPublicKeyFromTransitResponse return parsed public key from the keyInfo transit read API response
func (k *VaultTransitKey) GetPublicKeyFromTransitResponse(keyInfo *vault.Secret, version int) (crypto.PublicKey, error) {

	// Build jq query
	jqQuery := fmt.Sprintf(".keys.%d.public_key", version)

	// extract and parse public key PEM
	op, err := jq.Parse(jqQuery)
	if err != nil {
		k.logger.Debug("jq query", zap.String("query", jqQuery))
		return nil, err
	}
	data, err := json.Marshal(keyInfo.Data)
	if err != nil {
		return nil, err
	}
	value, err := op.Apply(data)
	if err != nil {
		return nil, err
	}
	key, err := strconv.Unquote(string(value))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(key))

	if block == nil {
		return nil, fmt.Errorf("error Pem Decoding pub key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

// Validate this transit key supports the jwt alg
func (k *VaultTransitKey) Validate(alg string) (*JWTKeyTransit, error) {

	// Get properties for this JWT signing alg
	jwtKeyProp := SupportedKeyTypeMap[alg]
	if jwtKeyProp == nil {
		k.logger.Error("Unsupported algorithm", zap.String("alg", alg))
		return nil, fmt.Errorf("Unsupported algorithm: %s", alg)
	}

	// check if Transit key supports this JWT alg
	if !stringInSlice(k.Type, jwtKeyProp.SupportedKeyType) {
		k.logger.Error("Unsupported transit key type for this alg", zap.String("alg", alg), zap.String("type", k.Type))
		return nil, fmt.Errorf("Unsupported transit key type %s , for this alg %s", alg, k.Type)
	}

	return jwtKeyProp, nil

}
