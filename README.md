# envoy-tools

A toolset for debugging Envoy and Envoy Control Plane configuration:
* `cp diff` displays a diff between configuration produced by two control planes

## Getting Started

### Installation

```shell
go install github.com/Hexta/envoy-tools@latest
```

## Usage

### Examples

* Compare CDS configuration returned by Envoy Control Plane `envoy-cp-first` and `envoy-cp-second`
    ```shell
    envoy-tools cp cds diff --node-id node envoy-cp-first:18000 envoy-cp-second:18000
    ```
* Compare RDS configuration returned by Envoy Control Plane `envoy-cp-first` and `envoy-cp-second`
    ```shell
    envoy-tools cp rds diff --node-id node envoy-cp-first:18000 envoy-cp-second:18000
    ```

## Contributing

Please see our [contributing guidelines](CONTRIBUTING.md).
