## oidc-server completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	oidc-server completion fish | source

To load completions for every new session, execute once:

	oidc-server completion fish > ~/.config/fish/completions/oidc-server.fish

You will need to start a new shell for this setup to take effect.


```
oidc-server completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   debug mode enabled
```

### SEE ALSO

* [oidc-server completion](oidc-server_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 15-Mar-2022