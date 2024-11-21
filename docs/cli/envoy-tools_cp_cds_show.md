## envoy-tools cp cds show

Show Envoy CDS configuration

```
envoy-tools cp cds show IP:PORT [cluster name]... [flags]
```

### Examples

```
# Show all clusters
$ envoy-tools cp cds show 127.0.0.1:18000

# Show specific clusters
$ envoy-tools cp cds show 127.0.0.1:18000 cluster1 cluster2

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

* [envoy-tools cp cds](envoy-tools_cp_cds.md)	 - CDS tools

