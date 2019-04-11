reset-all:
	(cd gamepad && make reset)
	(cd emulator && make reset)

clean-all: clean
	(cd gamepad && make clean)
	(cd emulator && make clean)

install-all: install
	(cd gamepad && make install)
	(cd emulator && make install)

build-all: build
	(cd gamepad && make build)
	(cd emulator && make build)

clean:
	rm -rf matrix-package
	rm -f matrix.zip
	rm -f *.log
	rm -rf build
	rm -rf vendor

install:
	GO111MODULE=off go clean -modcache || true
	GOPROXY=https://gocenter.io GO111MODULE=on go mod tidy

build:
	go build -o build/matrix-local .

build-arm:
	docker run -it -e TARGETS="linux/arm-7" -e OUT=matrix -e EXT_GOPATH=/gopath \
	-v $(PWD):/gopath/src/github.com/richardlt/matrix \
	-v $(PWD)/build:/build karalabe/xgo-1.11 github.com/richardlt/matrix

build-windows:
	docker run -it -e TARGETS="windows/amd64" -e OUT=matrix -e EXT_GOPATH=/gopath \
	-v $(PWD):/gopath/src/github.com/richardlt/matrix \
	-v $(PWD)/build:/build karalabe/xgo-1.11 github.com/richardlt/matrix

package:	
	rm -rf matrix-package
	mkdir -p matrix-package/gamepad/public
	mkdir -p matrix-package/emulator/public
	cp build/matrix-* matrix-package/
	cp -R themes matrix-package/
	cp -R fonts matrix-package/
	cp -R images matrix-package/
	cp -R animations matrix-package/
	cp -R gamepad/build/default/. matrix-package/gamepad/public/
	cp -R emulator/client/public/. matrix-package/emulator/public/
	zip -r matrix.zip matrix-package

test: 	
	go test -race github.com/richardlt/matrix/... -v

test-with-report: 	
	GOPROXY=https://gocenter.io GO111MODULE=on go test -race github.com/richardlt/matrix/... -v | go-junit-report > report.xml