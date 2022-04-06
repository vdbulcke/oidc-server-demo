# Home

`oidc-server` is a standalone Mock OIDC server built on top of [https://github.com/oauth2-proxy/mockoidc](https://github.com/oauth2-proxy/mockoidc).

## Features

* OIDC Authorization Code flow: from [mockoidc](https://github.com/oauth2-proxy/mockoidc)
* Provider Discovery (`./well-known/openid-configuration`): from [mockoidc](https://github.com/oauth2-proxy/mockoidc)
* Token Signature validation (from jwk provider endpoint): from [mockoidc](https://github.com/oauth2-proxy/mockoidc)
* Client Auth Method (`client_secret_post`): from [mockoidc](https://github.com/oauth2-proxy/mockoidc)
* Refresh Token Flow: : from [mockoidc](https://github.com/oauth2-proxy/mockoidc)
* Standalone Server: `oidc-server start` 
* Custom Mock Users: inject arbitrary Claims in `id_token` and/or `userinfo`
* Docker container (TODO)