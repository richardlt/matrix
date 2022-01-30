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
	rm -f *.xml

install:
	go mod tidy

build:
	go build -o build/matrix-local .

build-arm:
	docker run -it -e TARGETS="linux/arm-7" -e OUT=matrix -e EXT_GOPATH=/gopath \
	-v $(PWD):/gopath/src/github.com/richardlt/matrix \
	-v $(PWD)/build:/build richardleterrier/xgo:v1.13.1 github.com/richardlt/matrix

build-windows:
	docker run -it -e TARGETS="windows/amd64" -e OUT=matrix -e EXT_GOPATH=/gopath \
	-v $(PWD):/gopath/src/github.com/richardlt/matrix \
	-v $(PWD)/build:/build richardleterrier/xgo:v1.13.1 github.com/richardlt/matrix

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
	docker run -it \
	-v $(PWD):/tmp/workspace \
	-w /tmp/workspace richardleterrier/debpacker:v0.0.2 debpacker make

test: 	
	go test -race github.com/richardlt/matrix/... -v

test-with-report: 	
	go test -race github.com/richardlt/matrix/... -v | go-junit-report > report.xml