## envoy-tools cp cds diff

Compare Envoy CDS configuration from two Envoy instances

```
envoy-tools cp cds diff IP:PORT IP:PORT [flags]
```

### Examples

```
# Diff all clusters
$ envoy-tools cp cds diff 127.0.0.1:18000 127.0.0.1:18001

# Diff specific clusters
$ envoy-tools cp cds diff 127.0.0.1:18000 127.0.0.1:18001 -c cluster-1 -c cluster-2

```

### Options

```
  -c, --cluster strings   Cluster name
  -h, --help              help for diff
  -i, --indent int        Indentation level (default 4)
  -s, --stats             Display stats only
```

### Options inherited from parent commands

```
      --max-grpc-message-size int   Max size of gRPC message (default 104857600)
      --node-id string              Node id used in discovery requests
```

### SEE ALSO

* [envoy-tools cp cds](envoy-tools_cp_cds.md)	 - CDS tools

