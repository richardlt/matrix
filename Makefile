VERSION := $(if ${VERSION},${VERSION},snapshot)
TARGET_LDFLAGS = -ldflags "-X 'main.VERSION=$(VERSION)'"
GO_BUILD = go build -buildvcs=false

reset-all:
	(cd gamepad && make reset)
	(cd emulator && make reset)

clean-all: clean
	(cd gamepad && make clean)
	(cd emulator && make clean)

install-all:
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
	rm -f *.xml
	rm -rf target

build:
	docker run --rm \
		-e VERSION \
		-v $(PWD):/tmp/workspace -w /tmp/workspace \
		--entrypoint /usr/bin/make \
		ghcr.io/goreleaser/goreleaser-cross:v1.21 \
		build-linux-arm-7 build-linux-amd64 build-darwin-amd64 build-windows-amd64

build-linux-arm-7:
	CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 $(GO_BUILD) $(TARGET_LDFLAGS) -o build/matrix-linux-arm-7 .

build-linux-amd64: 
	CC=x86_64-linux-gnu-gcc CXX=x86_64-linux-gnu-g++ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO_BUILD) $(TARGET_LDFLAGS) -o build/matrix-linux-amd64 .

build-darwin-amd64:
	CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GO_BUILD) $(TARGET_LDFLAGS) -o build/matrix-darwin-amd64 .

build-windows-amd64:
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(GO_BUILD) $(TARGET_LDFLAGS) -o build/matrix-windows-4.0-amd64.exe .

package:	
	rm -rf matrix-package
	mkdir -p matrix-package/gamepad/public
	mkdir -p matrix-package/emulator/public
	cp build/matrix-* matrix-package/
	cp -R themes matrix-package/
	cp -R fonts matrix-package/
	cp -R images matrix-package/
	cp -R animations matrix-package/
	cp -R gamepad/public/. matrix-package/gamepad/public/
	cp -R emulator/public/. matrix-package/emulator/public/
	zip -r matrix.zip matrix-package

debpacker:
	rm -rf target
	docker run --rm \
		-v $(PWD):/tmp/workspace -w /tmp/workspace \
		richardleterrier/debpacker:v0.0.2 \
		debpacker make

test: 	
	go test -race github.com/richardlt/matrix/... -v | tee report.out

test-with-report: test	
	go install github.com/jstemmer/go-junit-report@latest
	cat report.out | go-junit-report > report.xml