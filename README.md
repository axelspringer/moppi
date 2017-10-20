[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# kombinat api

A golang app created with `yo golang`.

## Publish a Package

Each package has its own directory, with subdirectories for package revisions. Their package folder contains various files describing how to install, uninstall and configure it.

### Folder Structure

> Universes are best kept in Gits

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

Make an install request with the `example/data.json`.

```json
curl -d "@install.json" -X POST http://localhost:8080/install
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
[MIT](/LICENSE)
