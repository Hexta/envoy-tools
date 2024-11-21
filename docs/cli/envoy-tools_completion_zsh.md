## envoy-tools completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(envoy-tools completion zsh)

To load completions for every new session, execute once:

#### Linux:

	envoy-tools completion zsh > "${fpath[1]}/_envoy-tools"

#### macOS:

	envoy-tools completion zsh > $(brew --prefix)/share/zsh/site-functions/_envoy-tools

You will need to start a new shell for this setup to take effect.


```
envoy-tools completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -o, --output Format   output format (json, yaml, text, jq) (default yaml)
```

### SEE ALSO

* [envoy-tools completion](envoy-tools_completion.md)	 - Generate the autocompletion script for the specified shell

