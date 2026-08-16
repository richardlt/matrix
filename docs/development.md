# Development

## Requirements

- [Go](https://golang.org/dl/) 1.25+
- [Node.js](https://nodejs.org/en/download/) 20.19+ or 22.12+, required by Vite

```sh
make install-all
```

## Running it

```sh
go run main.go start --log-level info core gamepad emulator demo  # add any other software
(cd emulator && npm start)
(cd gamepad && npm start)
```

The emulator is at http://localhost:3001 and the gamepad at http://localhost:4002. Both
web apps are served from the Go binary in a release build, where they are embedded, but in
development Vite serves them so you get hot reload.

## Checks

```sh
make check       # gofmt, go vet and errcheck
make check-all   # the same, plus a TypeScript check of both web apps
make test        # unit tests with the race detector
```

Vite transpiles without type checking, so `make check-all` is the only thing that rejects
a TypeScript error.

`errcheck` reports discarded `error` returns. A call whose error is deliberately ignored
is written `_ = f()`, which errcheck accepts, so a report means one was dropped by
accident.

## Protobuf

The gRPC bindings are generated from the `.proto` files under `sdk-go/`:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
(cd sdk-go && make proto)
```

It needs [protoc](https://github.com/protocolbuffers/protobuf/releases) as well. The
generated files land next to their `.proto`, and the result is reproducible, so a diff
after regenerating means a real change.

## Release builds

Two artifacts, one per generation of board:

| | | |
| - | - | - |
| `make package` | ARMv6 tarball | Pi Zero and Pi 1, written to an SD card |
| `make deb` | ARMv7 `.deb` | Raspberry Pi OS on a Pi 2 or later |

Both need an ARM musl cross toolchain, which is not an apt package. The same one covers
both targets:

```sh
make package ARMV6_CC=<your arm musl gcc>
make deb ARMV7_CC=<your arm musl gcc>
```

Check the result targets what you expect with:

```sh
readelf -A build/matrix-linux-armv6 | grep Tag_CPU_arch
```

`make check-armv6` compiles for ARMv6 without the toolchain. It is a compile check only:
cgo is off, so the result has no controller support and is not shippable.

`make deb` takes its version from the latest git tag; override with
`make deb DEB_VERSION=1.2.3`. The package contents are under
[`packaging/deb/`](../packaging/deb); `/etc/default/matrix` is a conffile, so a user's
edits to the component list survive an upgrade.

The Makefile explains why the builds are put together the way they are.
