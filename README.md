# oidc-server-demo

OIDC Server Demo is Mock OIDC server that can be used to test OIDC integration. 

Built on the shoulders of giants [https://github.com/oauth2-proxy/mockoidc](https://github.com/oauth2-proxy/mockoidc), or more specifically a fork of it ([https://github.com/vdbulcke/mockoidc](https://github.com/vdbulcke/mockoidc)) for specific features for supporting a full standalone Mock OIDC server.




## Features

* OIDC Flows: Authorization Code, Refresh Token, PKCE => from `https://github.com/oauth2-proxy/mockoidc` 
* Pushed Authorization Request
* Generic Mock Users: Inject Arbitrary structured claims in ID Token and/or userinfo.
* Support for Hashicorp Vault transit backend as a Software Security Module
## Documentation

Install, configuration documentation can be found at [https://vdbulcke.github.io/oidc-server-demo/](https://vdbulcke.github.io/oidc-server-demo/).

## tl;dr

Start the server with
```bash
oidc-server start --config example/config.yaml
```

Docker [container](https://github.com/vdbulcke/oidc-server/pkgs/container/oidc-server)
### Debugging

Using the `--debug` flags will log each request made to the the mock oidc server

```bash
oidc-server start --config example/config.yaml --debug
```
