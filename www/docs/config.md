# Server Configuration 

You can find a complete example of the client configuration in [example/config.yaml](https://github.com/vdbulcke/oidc-server-demo/blob/main/example/config.yaml).

## Client Authentication Settings


### Client ID and Secret

!!! important 
    Mandatory (either in config file or ENV variables) **unless** using pkce flow. In which case the `client_secret` is not required.

    See section [PKCE](https://vdbulcke.github.io/oidc-server-demo/config/#pkce).

The OIDC client credentials can be passed either in the main `config.yaml` config file, or as environment variables. 

#### Config File

```yaml
## Client Credentials: (Mandatory)
### NOTE: those client_id/client_secret can be passed
###       as environment variables with: 
###
###    export OIDC_CLIENT_ID=my_client_id
###    export OIDC_CLIENT_SECRET=my_client_id
###
client_id: my_client_id
client_secret: my_client_secret
```

#### Env Variables

```bash
export OIDC_CLIENT_ID=my_client_id
export OIDC_CLIENT_SECRET=my_client_secret
```

### Client Authentication Method

!!! important 
    Only `client_secret_post` is currently supported, where  ClientID/ClientSecret are passed in the POST body as `application/x-www-form-urlencoded` parameters.




## PKCE



!!! info
    More information about pkce can be found [https://www.oauth.com/oauth2-servers/pkce/](https://www.oauth.com/oauth2-servers/pkce/).

#### Pkce Challenge Method

!!! note 
    Optional Settings


```yaml

### Supported challenge method (Optional)
###
### Default: 
###  - S256
###  - plain
###
# pkce_challenge_methods:
# - plain
# - S256
```

## Scopes


!!! note 
    Optional Settings

You can update the list of scopes requested using the `supported_scopes` setting:

```yaml
## Supported Scropes (optional)
### List of supported scropes by the OIDC server
### Default to same default as https://github.com/oauth2-proxy/mockoidc
# supported_scopes: 
# - openid
# - profile
# - roles
```

!!! warning
    The oidc mock server will validate that the scopes requested are present in this list

## Authorization Server 

### Issuer 

!!! note 
    Optional Settings


You can specify the `issuer_base_url` setting that will be used for construct the Issuer by adding the base path `/oidc`.

```yaml
## Issuer Base Url (Optional)
### Set the base url for the OIDC server
### Issuer is generated using
###    issuer_base_url + '/oidc'
###
### Default: http://127.0.0.1:5557
# issuer_base_url: http://oidc.example.com:8080
```

!!! tips
    Use this when the OIDC server is access via a reverse proxy/LB or when using port mapping (e.g. via docker)


### Discovery Endpoint (and well-known configuration)

The discovery endpoint is build using the `issuer_base_url` (or the default `http://127.0.0.1:5557`) by adding the base path `/oidc`, so by default: 

* discovery endpoint: `http://127.0.0.1:5557/oidc/.well-known/openid-configuration`

* authorization_endpoint: `http://127.0.0.1:5557/oidc/authorize`
* issuer: `http://127.0.0.1:5557/oidc`
* token_endpoint: `http://127.0.0.1:5557/oidc/token`
* userinfo_endpoint: `http://127.0.0.1:5557/oidc/userinfo`
* jwks_uri: `http://127.0.0.1:5557/oidc/.well-known/jwks.json`


