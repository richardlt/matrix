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
	rm -f matrix
	rm -f matrix.zip
	rm -f *.log
	rm -rf build

build:
	go build -o matrix .

build-arm:
	docker run -it -e TARGETS="linux/arm-6 linux/arm-7" -e OUT=matrix -e EXT_GOPATH=/gopath \
	-v $(PWD):/gopath/src/github.com/richardlt/matrix \
	-v $(PWD)/build:/build karalabe/xgo-latest github.com/richardlt/matrix

package:	
	rm -rf matrix-package
	mkdir -p matrix-package/gamepad/public
	mkdir -p matrix-package/emulator/public
	cp matrix matrix-package/
	cp -R themes matrix-package/
	cp -R fonts matrix-package/
	cp -R images matrix-package/
	cp -R gamepad/build/default/. matrix-package/gamepad/public/
	cp -R emulator/client/public/. matrix-package/emulator/public/
	zip -r matrix.zip matrix-package

test: 	
	go test -race github.com/richardlt/matrix/... -v

test-with-report: 	
	go test -race github.com/richardlt/matrix/... -v | go-junit-report > report.xml