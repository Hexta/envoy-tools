## envoy-tools completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	envoy-tools completion fish | source

To load completions for every new session, execute once:

	envoy-tools completion fish > ~/.config/fish/completions/envoy-tools.fish

You will need to start a new shell for this setup to take effect.


```
envoy-tools completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -o, --output Format   output format (json, yaml, text, jq) (default yaml)
```

### SEE ALSO

* [envoy-tools completion](envoy-tools_completion.md)	 - Generate the autocompletion script for the specified shell

