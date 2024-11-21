## envoy-tools cp rds show

Show Envoy RDS configuration

```
envoy-tools cp rds show IP:PORT [virtual host name]... [flags]
```

### Examples

```
# Show all route configs
$ envoy-tools cp rds show 127.0.0.1:18000

# Show specific route configs
$ envoy-tools cp rds show 127.0.0.1:18000 route-config-1 route-config-2

```

### Options

```
  -h, --help   help for show
```

### Options inherited from parent commands

```
      --max-grpc-message-size int   Max size of gRPC message (default 104857600)
      --node-id string              Node id used in discovery requests
  -o, --output Format               output format (json, yaml, text, jq) (default yaml)
```

### SEE ALSO

* [envoy-tools cp rds](envoy-tools_cp_rds.md)	 - RDS tools

