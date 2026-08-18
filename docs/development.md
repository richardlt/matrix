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

CI uses `arm-unknown-linux-musleabihf` from
[musl-cross](https://github.com/musl-cross/musl-cross/releases), pinned to a release and a
checksum in [`.github/actions/release-artifacts`](../.github/actions/release-artifacts).
Unpacking one of those tarballs gives you the same compiler CI builds with:

```sh
make package ARMV6_CC=arm-unknown-linux-musleabihf-gcc
```

Check the result targets what you expect with:

```sh
readelf -A build/matrix-linux-armv6 | grep Tag_CPU_arch
```

`make check-armv6` compiles for ARMv6 without the toolchain. It is a compile check only:
cgo is off, so the result has no controller support and is not shippable.

The package contents are under [`packaging/deb/`](../packaging/deb);
`/etc/default/matrix` is a conffile, so a user's edits to the component list survive an
upgrade.

The Makefile explains why the builds are put together the way they are.

## Versions

The version is taken from the tags and stamped into the binary, so it reports the same
thing the `.deb` around it claims:

```sh
make print-version        # 0.2.3+18.g8736069, what a build here would carry
./build/matrix-local --version
```

A tagged commit gives the tag alone, `0.2.3`, which is what a release is. Anywhere else
`git describe`'s suffix says how far past the tag the build sits and which commit it is,
with `.dirty` on the end if the tree had uncommitted changes. Pass `VERSION=1.2.3` to
override, and note that a checkout with no tag in reach — a shallow clone, for one — has
nothing to describe and falls back to `0.0.0`.

A binary built with a bare `go build` reports `dev`: the version is a linker flag, and
only the Makefile passes it.

## CI and releases

[`ci.yml`](../.github/workflows/ci.yml) runs on every push: the Go checks, a TypeScript
check of both web apps, the tests, and the two release artifacts. The artifacts are
attached to the run, so any commit can be written to a card without a cross toolchain on
your own machine, and each carries the describe version of the commit it came from.

Nothing runs on a pull request, so a fork cannot start a workflow here.

Tagging is the release:

```sh
git tag 1.2.3
git push origin 1.2.3
```

[`release.yml`](../.github/workflows/release.yml) then builds the same two artifacts, which
this time report the tag, and publishes them as a GitHub release. Nothing is built by hand
and nothing is uploaded by hand.

Both workflows share [`.github/actions/release-artifacts`](../.github/actions/release-artifacts),
which installs the toolchain, builds the two artifacts and checks the architecture, the
static linking and the stamped version of each before anything is published.
