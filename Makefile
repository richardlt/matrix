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
	build-web build-armv6 build-armv7 check-armv6 package deb \
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

build-web:
	(cd gamepad && make build)
	(cd emulator && make build)

# The web apps are built first because the Go binary embeds their output.
build-all: build-web
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
	rm -rf alpine-card

install:
	go mod tidy

build:
	go build -o build/matrix-local .

# --- Release builds -------------------------------------------------------------------
#
# Two artifacts, one per generation of board:
#
#   ARMv6, as a tarball   for boards too small to carry a full distribution, installed by
#                         writing an SD card. See scripts/prepare-alpine-card.sh.
#   ARMv7, as a .deb      installed with dpkg and run as a systemd service.
#
# ARMv6 also runs on ARMv7 hardware, so the split is about the packaging rather than the
# instruction set: where a full distribution fits, a service that starts at boot and
# survives a crash is worth more than the last few percent of the CPU.
#
# Both are built the same way:
#
#   CGO_ENABLED=1  the device component reaches USB controllers through
#                  github.com/karalabe/hid, which is only compiled in under cgo. A
#                  cgo-less binary still builds and runs, but hid.Supported() reports
#                  false and no controller is ever found.
#   CC             a musl toolchain. The C half of the build matters as much as the Go
#                  half: hid compiles a vendored libusb, so the C flags below pin the
#                  architecture rather than trusting the compiler's default. Debian's
#                  arm-linux-gnueabihf, for one, targets ARMv7, which would otherwise
#                  leave an ARMv7 C payload inside an ARMv6 binary.
#   CGO_CFLAGS     -O2 -g are Go's own defaults, repeated because setting this replaces
#                  them rather than adding to them.
#   -static        with libc linked in, the artifact does not depend on glibc or musl at
#                  runtime, which is what lets one binary serve Alpine and Raspberry Pi OS.
#   netgo          use the pure Go resolver rather than the cgo one, which is what makes
#                  a static link safe.
#
# Both targets use the same cross compiler: armhf is one triple, and the architecture
# comes from the flags. Install the toolchain first; it is not an apt package. See
# docs/development.md.
ARMV6_CC ?= armv6-linux-musleabihf-gcc
ARMV6_CFLAGS ?= -O2 -g -march=armv6 -mfpu=vfp
ARMV7_CC ?= $(ARMV6_CC)
ARMV7_CFLAGS ?= -O2 -g -march=armv7-a -mfpu=vfpv3-d16

GO_RELEASE_FLAGS := -trimpath -tags netgo -ldflags '-extldflags "-static"'

build-armv6:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 \
		CC=$(ARMV6_CC) CGO_CFLAGS="$(ARMV6_CFLAGS)" \
		go build $(GO_RELEASE_FLAGS) -o build/matrix-linux-armv6 .

build-armv7:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
		CC=$(ARMV7_CC) CGO_CFLAGS="$(ARMV7_CFLAGS)" \
		go build $(GO_RELEASE_FLAGS) -o build/matrix-linux-armv7 .

# Compile-check the release target without the cross toolchain. This proves the code
# builds for ARM; it does not produce a shippable binary, since cgo is off and the
# controller support is therefore missing.
check-armv6:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o /dev/null .

# --- Release packaging ----------------------------------------------------------------

# The web apps are embedded in the binary, so what ships beside it is only the data a user
# may want to replace, plus the card-staging script the Alpine install is built around.
#
# They are rebuilt as a recipe step rather than a prerequisite so the order is guaranteed
# even under make -j. It matters: the embedded directory carries a .gitkeep, so a binary
# built before them compiles happily and then serves nothing.
package:
	$(MAKE) build-web
	$(MAKE) build-armv6
	rm -rf matrix-package
	mkdir -p matrix-package
	cp build/matrix-linux-armv6 matrix-package/
	cp -R themes fonts images animations matrix-package/
	cp scripts/prepare-alpine-card.sh matrix-package/
	tar czf matrix.tar.gz matrix-package

# The .deb, assembled with dpkg-deb rather than a packaging framework: the payload is one
# static binary, four data directories and a unit file, which is less than any of them
# would cost to configure.
#
# --root-owner-group is what makes the result reproducible from an unprivileged build.
# Without it every file in the package is owned by whoever ran make.
#
# /etc/default/matrix is listed as a conffile so dpkg keeps an edited component list
# across an upgrade instead of overwriting it.
DEB_VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo 0.0.0)
DEB_ARCH := armhf
DEB_ROOT := target/deb

deb:
	$(MAKE) build-web
	$(MAKE) build-armv7
	rm -rf $(DEB_ROOT)
	mkdir -p $(DEB_ROOT)/DEBIAN $(DEB_ROOT)/usr/bin $(DEB_ROOT)/etc/default \
		$(DEB_ROOT)/var/lib/matrix $(DEB_ROOT)/lib/systemd/system
	sed -e 's/@VERSION@/$(DEB_VERSION)/' -e 's/@ARCH@/$(DEB_ARCH)/' \
		packaging/deb/control > $(DEB_ROOT)/DEBIAN/control
	install -m 755 packaging/deb/postinst packaging/deb/prerm packaging/deb/postrm \
		$(DEB_ROOT)/DEBIAN/
	install -m 755 build/matrix-linux-armv7 $(DEB_ROOT)/usr/bin/matrix
	install -m 644 packaging/deb/matrix.service $(DEB_ROOT)/lib/systemd/system/
	install -m 644 packaging/deb/default $(DEB_ROOT)/etc/default/matrix
	echo /etc/default/matrix > $(DEB_ROOT)/DEBIAN/conffiles
	cp -R themes fonts images animations $(DEB_ROOT)/var/lib/matrix/
	dpkg-deb --build --root-owner-group $(DEB_ROOT) target/matrix_$(DEB_VERSION)_$(DEB_ARCH).deb

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
