## envoy-tools completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	envoy-tools completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
envoy-tools completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -o, --output Format   output format (json, yaml, text, jq) (default yaml)
```

### SEE ALSO

* [envoy-tools completion](envoy-tools_completion.md)	 - Generate the autocompletion script for the specified shell

