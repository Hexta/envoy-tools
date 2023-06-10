# envoy-tools

A toolset for debugging Envoy and Envoy Control Plane configuration:
* `cp diff` displays a diff between configuration produced by two control planes

## Getting Started

```shell
go install github.com/Hexta/envoy-tools@latest
```

## Usage

```shell
envoy-tools cp --node-id node diff -t cds envoy-cp-first:18000 envoy-cp-second:18000
```
