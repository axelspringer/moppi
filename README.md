[![Build Status](https://travis-ci.org/axelspringer/moppi.svg?branch=master)](https://travis-ci.org/axelspringer/moppi)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# :dizzy: Moppi

> :construction_worker: This is work in progress and many things are subject to change

> A golang app created with `yo golang`.

This is a universe for Mesos, Marathon and Chronos in many different KVs. We currently use [etcd](https://coreos.com/etcd/), but we use [libkv](https://github.com/docker/libkv) and there do support many more KVs.

## Docs

You can find a full documentation of the API in `/docs` or [here](https://axelspringer.github.io/moppi/).

## Config

We support config files, Environment variables and config parameters for `moppi`.

This is an example (`moppi.yml`) config

```yaml
verbose: true
etcd:
  prefix: "/moppi"
  endpoint: "localhost:2379"
chronos: "https://localhost:8181/"
marathon: "https://localhost:8080
```

### `--help` 

Displays the available options for `moppi`.

### `init`

Does initializes a KV for moppi.

## Publish a Package

Each package has its own directory, with subdirectories for package revisions. Their package folder contains various files describing how to install, uninstall and configure it.

### Folder Structure

> Your universe is best kept in a Git to have revisions

A universe has the following folder structure.

```
└── packages/F/foo
    ├── 0
    │   ├── install.json
    │   ├── uninstall.json
    │   ├── marathon.json
    │   ├── package.json
    ├── 1
    │   ├── install.json
    │   ├── uninstall.json
    │   ├── marathon.json
    │   ├── package.json
    └── ...
```

All meta information about a universe must also be stored within.

```
└── meta
    ├── version.json
    └── ...
```

### `package.json`

Contains general information about a package when published.

### `install.json`

Contains all information necessary to install this package.

### `uninstall.json`

Contains all information necessary to uninstall this package.

### `marathon.json`

> we currently support Marathon up to `1.4.x`, but working hard to move to `1.5.x`

When using [Marathon](https://mesosphere.github.io) to run long standing task on Mesos this config contains the description of such. Please, consult the [Marathon Docs](https://mesosphere.github.io/marathon/docs) as how to write such config.

### `chronos.json`

When using [Chronos](https://github.com/mesos/chronos) to run scheduled task, you can provide a config for it within your package. It is then applied in the installation process.


## Examples

You find the example data in `examples/`. etcd import is a s simple as `PUT`ing the relevant JSON in the KV.

```bash
curl http://127.0.0.1:2379/v2/keys/moppi/universes/dev/packages/example/1/uninstall -XPUT -d value="$(cat marathon.json)"
```

## Getting Started

Install neat tools and dependencies.

```
make deps && make restore
```

Build the app.

```
make build
```


## License
[Apache-2.0](/LICENSE)
