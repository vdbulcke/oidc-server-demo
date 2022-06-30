# Mock User Configuration


## Default (Infinite) Mock User

!!! note
    The _Default Mock User_ is used if no other mock user is defined or when the Mock User Queue is empty.


You configure the default mock user using the following settings in the main `config.yaml`: 

```yaml

## Mock User (Mandatory)
## 
mock_user:
  ## Sub  (Mandatory)
  ###  the user's subject
  sub: bob@acme.com

  ## ID Token Claims (mandatory)
  ### Arbitrary key/values claims to 
  ### add in the id_token
  ### 
  ### Empty set to '{}'
  id_token_claims: 
    ## example adding amr values
    amr:
      - eid
      - urn:be:fedict:iam:fas:Level500
    
    ## dummy claims
    foo: bar

  ## Access Token Claims (Optional)
  ### Arbitrary key/values claims to 
  ### add in the access_token
  ### 
  access_token_claims: 
    amr:
      - eid
      - urn:be:fedict:iam:fas:Level500
    
    foo: 
      - hello: world
      - bar: baz



  ## Userinfo Claims (optional)
  ### Arbitrary key/values claims to 
  ### add in the userinfo response
  ### 
  ### Empty set to '{}'
  userinfo_claims: 

    ## example 
    fedid: "73691e9e7beee3becdf78fc9394d28fe548fe249"
    surname: Jane
    givenName: Doe
    
```


!!! note
    Access Token Claims are introduced in release `v0.4.0`, but are not mandatory for backward compatibility

## Multiple Mock Users

Optionally if you need additional different mock users, you can configure the following property in the main `config.yaml`: 

```yaml

##
## Additional Mock Users
##  since version v0.3.0
### Users loaded in the UserQueue
mock_user_folder: ./example/users
```

The `mock_user_folder` property is a path to a directory containing the additional mock users definition (one per YAML file). 

### Mock User Files

Each files represents a mock user definition:

```yaml

## Sub  (Mandatory)
###  the user's subject
sub: john.doe@acme.com

## ID Token Claims (mandatory)
### Arbitrary key/values claims to
### add in the id_token
###
### Empty set to '{}'
id_token_claims:
  ## example adding amr values
  amr:
    - eid
    - urn:be:fedict:iam:fas:Level500

  ## dummy claims
  foo: bar

## Access Token Claims (Optional)
### Arbitrary key/values claims to 
### add in the access_token
### 
access_token_claims: 
  custom: claims


## Userinfo Claims (optional)
### Arbitrary key/values claims to
### add in the userinfo response
###
### Empty set to '{}'
userinfo_claims:

  ## example
  fedid: "F8nZR6mFFlvyGd85CE5Qu5iFD4qaAGivWHdW1evt"
  surname: John
  givenName: Doe
```


### Mock User Queue

The `oidc-server` will load the mock user definition files found in the `mock_user_folder` and put them **in that order** in the User Queue. 

For each new Authorize call, the `oidc-server` will Pop the first user from the Queue and use it for the id_token and userinfo response, until the queue is empty. In which case, it will continue returning the _default mock user_.


For example, if we have the following files in `mock_user_folder`: 

* `01-john_doe.yaml`
* `02-babs_jensen.yaml` 

And we make the following requests: 

1. Authorize request returns the profile of  `01-john_doe.yaml`
1. Authorize request returns the profile of `02-babs_jensen.yaml`
1. Authorize request returns the profile of  _default mock user_ (`mock_user:` in the main `config.yaml`)
1. Authorize request returns the profile of  _default mock user_ 
1. Authorize request returns the profile of  _default mock user_ 
1. ...

