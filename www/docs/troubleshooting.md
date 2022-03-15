# Troubleshooting

This page regroups some of the most common errors.

## Closing Local Server 

Use `CTL+C` to close the local http server 

```
{"level":"info","ts":1647367899.411411,"caller":"oidc-server/server.go:40","msg":"starting server"}
{"level":"info","ts":1647367899.4114451,"caller":"oidc-server/server.go:41","msg":"server config","Issuer":"http://127.0.0.1:5557/oidc"}
^C{"level":"info","ts":1647367900.9016469,"caller":"oidc-server/server.go:50","msg":"Got signal","sig":"interrupt"}

```

## Configuration Validation Errors

### Missing Mandatory Settings

* example missing client_id:
```
oidc-server start --config example/wrong.yaml 
Key: 'OIDCServerConfig.ClientID' Error:Field validation for 'ClientID' failed on the 'required' tag
{"level":"error","ts":1647367972.9460135,"caller":"cmd/start.go:76","msg":"validating config","error":"Validation Error","stacktrace":"github.com/vdbulcke/oidc-server-demo/cmd.runServer\n\t/home/runner/work/oidc-server-demo/oidc-server-demo/cmd/start.go:76\ngithub.com/spf13/cobra.(*Command).execute\n\t/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.4.0/command.go:860\ngithub.com/spf13/cobra.(*Command).ExecuteC\n\t/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.4.0/command.go:974\ngithub.com/spf13/cobra.(*Command).Execute\n\t/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.4.0/command.go:902\ngithub.com/vdbulcke/oidc-server-demo/cmd.Execute\n\t/home/runner/work/oidc-server-demo/oidc-server-demo/cmd/main.go:36\nmain.main\n\t/home/runner/work/oidc-server-demo/oidc-server-demo/main.go:8\nruntime.main\n\t/opt/hostedtoolcache/go/1.17.7/x64/src/runtime/proc.go:255"}

```

