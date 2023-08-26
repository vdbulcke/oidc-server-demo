package oidcserver

import (
	"encoding/json"
	"os"

	"github.com/golang-jwt/jwt/v4"
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
	Subject               string                      `yaml:"sub,omitempty" validate:"required"`
	IDTokenClaims         map[interface{}]interface{} `yaml:"id_token_claims,omitempty"  validate:"required"`
	UserAccessTokenClaims map[interface{}]interface{} `yaml:"access_token_claims,omitempty" `
	UserinfoClaims        map[interface{}]interface{} `yaml:"userinfo_claims,omitempty"  validate:"required"`
}

func (u *YAMLUser) ID() string {
	return u.Subject
}

func (u *YAMLUser) Userinfo(scope []string) ([]byte, error) {

	return json.Marshal(convert(u.UserinfoClaims))
}

func (u *YAMLUser) Claims(scope []string, claims *mockoidc.IDTokenClaims) (jwt.Claims, error) {

	userClaims := jwt.MapClaims{}

	// merge standard claims into User Access Token Claims
	userClaims["aud"] = claims.Audience
	userClaims["exp"] = claims.ExpiresAt
	userClaims["jti"] = claims.ID
	userClaims["iat"] = claims.IssuedAt
	userClaims["iss"] = claims.Issuer
	userClaims["nbf"] = claims.NotBefore
	userClaims["sub"] = claims.Subject
	userClaims["nonce"] = claims.Nonce

	// convert to map[string]interface{}
	for k, v := range u.IDTokenClaims {
		userClaims[k.(string)] = convert(v)
	}

	// return u.UserAccessTokenClaims, nil
	return userClaims, nil

}

// AccessTokenClaims just return standard claims
func (u *YAMLUser) AccessTokenClaims(claims *jwt.RegisteredClaims) (jwt.Claims, error) {

	if u.UserAccessTokenClaims != nil {

		userClaims := jwt.MapClaims{}

		// merge standard claims into User Access Token Claims
		userClaims["aud"] = claims.Audience
		userClaims["exp"] = claims.ExpiresAt
		userClaims["jti"] = claims.ID
		userClaims["iat"] = claims.IssuedAt
		userClaims["iss"] = claims.Issuer
		userClaims["nbf"] = claims.NotBefore
		userClaims["sub"] = claims.Subject

		for k, v := range u.UserAccessTokenClaims {
			userClaims[k.(string)] = convert(v)
		}

		// return u.UserAccessTokenClaims, nil
		return userClaims, nil
	}

	return claims, nil
}

func NewYAMLUser(filename string) (*YAMLUser, error) {

	data, err := os.ReadFile(filename)
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
