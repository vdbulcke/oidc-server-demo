package oidcserver

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/oauth2-proxy/mockoidc"
	"gopkg.in/yaml.v2"
)

// source: https://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

type YAMLUser struct {
	Subject        string                      `yaml:"sub,omitempty" validate:"required"`
	IDTokenClaims  jwt.MapClaims               `yaml:"id_token_claims,omitempty"  validate:"required"`
	UserinfoClaims map[interface{}]interface{} `yaml:"userinfo_claims,omitempty"  validate:"required"`
}

func (u *YAMLUser) ID() string {
	return u.Subject
}

func (u *YAMLUser) Userinfo(scope []string) ([]byte, error) {

	return json.Marshal(convert(u.UserinfoClaims))
}

func (u *YAMLUser) Claims(scope []string, claims *mockoidc.IDTokenClaims) (jwt.Claims, error) {

	u.IDTokenClaims["aud"] = claims.Audience
	u.IDTokenClaims["exp"] = claims.ExpiresAt
	u.IDTokenClaims["jti"] = claims.Id
	u.IDTokenClaims["iat"] = claims.IssuedAt
	u.IDTokenClaims["iss"] = claims.Issuer
	u.IDTokenClaims["nbf"] = claims.NotBefore
	u.IDTokenClaims["sub"] = claims.Subject
	u.IDTokenClaims["nonce"] = claims.Nonce

	return u.IDTokenClaims, nil
}

func NewYAMLUser(filename string) (*YAMLUser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	user := YAMLUser{}

	err = yaml.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
