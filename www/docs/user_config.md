# Mock User Configuration

You configure a mock user using the following settings: 


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


