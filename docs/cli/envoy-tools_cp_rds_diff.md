## envoy-tools cp rds diff

Compare Envoy CDS configuration from two Envoy instances

```
envoy-tools cp rds diff IP:PORT IP:PORT [flags]
```

### Examples

```
# Diff all virtual hosts
$ envoy-tools cp rds diff 127.0.0.1:18000 127.0.0.1:18001

# Diff specific virtual hosts
$ envoy-tools cp rds diff 127.0.0.1:18000 127.0.0.1:18001 -r virtual-host-1 -r virtual-host-2

```

### Options

```
  -h, --help                       help for diff
  -i, --indent int                 Indentation level (default 4)
      --route-config-name string   Route config name (default "default")
  -s, --stats                      Display stats only
  -r, --virtualhost strings        Virtual host name
```

### Options inherited from parent commands

```
      --max-grpc-message-size int   Max size of gRPC message (default 104857600)
      --node-id string              Node id used in discovery requests
  -o, --output Format               output format (json, yaml, text, jq) (default yaml)
```

### SEE ALSO

* [envoy-tools cp rds](envoy-tools_cp_rds.md)	 - RDS tools

