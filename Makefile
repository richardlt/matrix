# `test` pipes go test through tee, and tee's exit status would otherwise hide a
# failing test from make.
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -c

GO_JUNIT_REPORT_VERSION := v2.1.0
ERRCHECK_VERSION := v1.20.0

# Every target is a command rather than a file it produces. Without this, `make build`
# does nothing once the build/ directory exists, because make considers the target
# already up to date.
.PHONY: reset-all clean-all install-all build-all check-all clean install build \
	build-armv6 check-armv6 package debpacker \
	check fmt check-fmt vet errcheck test test-with-report

reset-all:
	(cd gamepad && make reset)
	(cd emulator && make reset)

clean-all: clean
	(cd gamepad && make clean)
	(cd emulator && make clean)

install-all: install
	(cd gamepad && make install)
	(cd emulator && make install)

# The web apps are built first because the Go binary embeds their output.
build-all:
	(cd gamepad && make build)
	(cd emulator && make build)
	$(MAKE) build

check-all: check
	(cd gamepad && make check)
	(cd emulator && make check)

clean:
	rm -rf matrix-package
	rm -f matrix.tar.gz
	rm -f *.log
	rm -rf build
	rm -f *.xml
	rm -rf target

install:
	go mod tidy

build:
	go build -o build/matrix-local .

# --- Release build --------------------------------------------------------------------
#
# One artifact: a static linux/ARMv6 binary.
#
#   CGO_ENABLED=1  the device component reaches USB controllers through
#                  github.com/karalabe/hid, which is only compiled in under cgo. A
#                  cgo-less binary still builds and runs, but hid.Supported() reports
#                  false and no controller is ever found.
#   GOARM=6        ARMv6 is the baseline; the result also runs on later ARMv7 hardware.
#   CC             an ARMv6 musl toolchain. The C half of the build matters as much as
#                  the Go half: hid compiles a vendored libusb, and Debian's
#                  arm-linux-gnueabihf targets ARMv7 by default, which would leave an
#                  ARMv7 C payload inside an otherwise ARMv6 binary.
#   -static        with libc linked in, the artifact does not depend on glibc or musl at
#                  runtime.
#   netgo          use the pure Go resolver rather than the cgo one, which is what makes
#                  a static link safe.
#
# Install the toolchain first; it is not an apt package. See the README.
ARMV6_CC ?= armv6-linux-musleabihf-gcc

build-armv6:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 CC=$(ARMV6_CC) \
		go build -trimpath -tags netgo -ldflags '-extldflags "-static"' \
		-o build/matrix-linux-armv6 .

# Compile-check the release target without the cross toolchain. This proves the code
# builds for ARMv6; it does not produce a shippable binary, since cgo is off and the
# controller support is therefore missing.
check-armv6:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o /dev/null .

# The web apps are embedded in the binary, so only the data the user can replace is
# shipped alongside it.
package:
	rm -rf matrix-package
	mkdir -p matrix-package
	cp build/matrix-* matrix-package/
	cp -R themes matrix-package/
	cp -R fonts matrix-package/
	cp -R images matrix-package/
	cp -R animations matrix-package/
	tar czf matrix.tar.gz matrix-package

debpacker:
	rm -rf target
	docker run -it \
	-v $(PWD):/tmp/workspace \
	-w /tmp/workspace richardleterrier/debpacker:v0.0.2 debpacker make

# --- Checks ---------------------------------------------------------------------------

check: check-fmt vet errcheck

fmt:
	gofmt -w .

check-fmt:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "these files are not gofmt'd:"; echo "$$unformatted"; exit 1; \
	fi

vet:
	go vet ./...

# errcheck reports return values of type error that are discarded. Calls that are
# deliberately ignored are written as `_ = f()`, which errcheck accepts, so a report here
# means an error was dropped by accident.
errcheck:
	go run github.com/kisielk/errcheck@$(ERRCHECK_VERSION) -ignoretests ./...

test:
	go test -race github.com/richardlt/matrix/... -v | tee report.out

test-with-report: test
	go run github.com/jstemmer/go-junit-report/v2@$(GO_JUNIT_REPORT_VERSION) \
		< report.out > report.xml
