[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# kombinat api

A golang app created with `yo golang`.

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
