reset:
	(cd gamepad && make reset)
	(cd emulator && make reset)
	(cd device && make reset)

install:
	(cd gamepad && make install)
	(cd emulator && make install)
	(cd device && make install)