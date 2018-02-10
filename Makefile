reset:
	(cd gamepad && make reset)
	(cd emulator && make reset)

clean:
	rm -rf matrix-package
	rm -f matrix
	rm -f matrix.zip
	rm *.log
	(cd gamepad && make clean)
	(cd emulator && make clean)

install:
	(cd gamepad && make install)
	(cd emulator && make install)

build:
	go build -a -o matrix .
	(cd gamepad && make build)
	(cd emulator && make build)

package:	
	rm -rf matrix-package
	mkdir -p matrix-package/gamepad/public
	mkdir -p matrix-package/emulator/public
	mkdir -p matrix-package/device
	cp matrix matrix-package/
	cp -R themes matrix-package/
	cp -R fonts matrix-package/
	cp -R images matrix-package/
	cp -R gamepad/build/default/ matrix-package/gamepad/public/
	cp -R emulator/client/public/ matrix-package/emulator/public/
	cp -R device/build/* matrix-package/device/
	cp -R device/node_modules matrix-package/device/
	cd matrix-package && zip -ro ../matrix.zip *