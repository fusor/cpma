## cpma [![Build Status](https://travis-ci.com/fusor/cpma.svg?branch=master)](https://travis-ci.com/fusor/cpma) [![Maintainability](https://api.codeclimate.com/v1/badges/aac7d46fd7899042ce52/maintainability)](https://codeclimate.com/github/fusor/cpma/maintainability)
Control Plane Migration Assistance:  Intended to help migration cluster
configuration of a OCP 3.x cluster to OCP 4.x

## Build

Requires go >= v1.11

This project is Go Modules for managing dependencies. This means it can be
cloned and compiled outside of `GOPATH`.

```console
$ git checkout https://github.com/fusor/cpma.git
$ cd cpma
$ make
$ ./bin/cpma
```

## Usage

Flags:
```
  -i, --allow-insecure-host   allow insecure ssh host key
  -c, --cluster-name string   OCP3 cluster kubeconfig name
      --config string         config file (Default searches ./cpma.yaml, $HOME/cpma.yml)
  -d, --debug                 show debug ouput
  -h, --help                  help for cpma
  -n, --hostname string       OCP3 cluster hostname
  -o, --output-dir string     set the directory to store extracted configuration. (default ".")
  -k, --ssh-keyfile string    OCP3 ssh keyfile path
  -l, --ssh-login string      OCP3 ssh login
  -p, --ssh-port string       OCP3 ssh port
  -v, --verbose               verbose output
```

You can find an example config in `examples/`. If a config is not provided CPMA will prompt for configuration information and offer to save inputs to a new configuration file.

Example:

```console
$ ./bin/cpma --config /path/to/config/.yml --verbose --debug
```

## CPMA Image
CPMA is also available in as an image, quay.io/ocpmigrate/cpma

Example Usage:
```
docker run -it --rm -v ${PWD}:/mnt:z -v $HOME/.kube:/.kube:z -v $HOME/.ssh:/.ssh:z -u ${UID} \
quay.io/ocpmigrate/cpma:latest
```

In these examples `${PWD}` is mounted in the working directory of the image (`/mnt`). This means that paths provided to --config and --output-dir will need to specified be relative to your present working directory.

To make it a little more intuitive it can also be run via an alias, for example:
```
$ alias cpma="docker run -it --rm -v ${PWD}:/mnt:z -v $HOME/.kube:/.kube:z -v $HOME/.ssh:/.ssh:z \
-u ${UID} quay.io/ocpmigrate/cpma:latest"

$ cpma --help
Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x

Usage:
  cpma [flags]

Flags:
      --config string       config file (Default searches ./cpma.yaml, $HOME/cpma.yml)
      --console-logs        output log to console
      --debug               show debug ouput
  -h, --help                help for cpma
      --insecure-key        allow insecure host key
  -k, --key string          OCP3 ssh key path
  -l, --login string        OCP3 ssh login
  -o, --output-dir string   set the directory to store extracted configuration. (default ".")
  -p, --port string         OCP3 ssh port
  -s, --source string       OCP3 cluster hostname
```

## IO

The data file structure looks like the following tree structure example. The
cluster endpoints subfolders contain the configuration files retrieved and to
process. The manifests directory contains the generated CRDs.

```
data
├── manifests
├── master-0.example.com
|   └── etc
|       └── origin
|           ├── master
|               ├── htpasswd
|               └── master-config.yaml
└── node-1.example.com
    └── etc
        └── origin
            └── node
                └── node-config.yaml
```

The configuration files are retrieved from local disk (`outputDir/<Hostname>/`),
If a file is not available it's retrieved from `<Hostname>` and stored on local disk.

To trigger a total or partial network file fetch, remove any prior data from
`<Hostname>` sub directory.

## Unit tests

In order to add new unit test bundle create `*_test.go` file in package you
want to test(ex: `foo.go`, `foo_test.go`).  To execute tests run `make test`.

https://golang.org/pkg/testing/
