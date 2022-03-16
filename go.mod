module github.com/vdbulcke/oidc-server-demo

go 1.18

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/oauth2-proxy/mockoidc v0.0.0-20220308204021-b9169deeb282
	github.com/spf13/cobra v1.4.0
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

replace github.com/oauth2-proxy/mockoidc => github.com/vdbulcke/mockoidc v0.1.0
