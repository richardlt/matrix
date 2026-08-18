# Installing matrix

Matrix runs on a Raspberry Pi. `core` and `device` are the two components a console needs:
`device` is the display, the player and a software all at once, since it pushes frames to
the Arduino over serial and reads USB controllers through HID. Everything else — the games
and the two web apps — is optional and can also run on your desk against a remote core,
with `--core-uri`.

[Flash the Arduino](./arduino.md) first, or the panel stays dark whatever else you do.

## Pi 2 or later — the .deb

Install Raspberry Pi OS yourself, then download the package from the
[releases](https://github.com/richardlt/matrix/releases) and install it:

```sh
sudo dpkg -i matrix_<version>_armhf.deb
systemctl status matrix
matrix --version
```

That puts a static ARMv7 binary at `/usr/bin/matrix`, the data it renders from under
`/var/lib/matrix`, and a systemd service that starts at boot and restarts on failure.

Choose what it runs in `/etc/default/matrix`:

```sh
MATRIX_COMPONENTS="core device gamepad emulator demo zigzag yumyum clock draw blocks getout animate light"
MATRIX_GAMEPAD_PORT=80
MATRIX_LOG_LEVEL=info
```

then `sudo systemctl restart matrix`. The file survives an upgrade with your edits intact.

On a **64-bit** Raspberry Pi OS, `dpkg` refuses a 32-bit package until armhf is added as a
foreign architecture:

```sh
sudo dpkg --add-architecture armhf
```

The binary itself is statically linked and runs either way.

`sudo apt remove matrix` stops and removes it; `apt purge` also deletes `/var/lib/matrix`,
including any theme or image you put there.

## Anything else — the tarball

The tarball holds the ARMv6 build, which also runs on later ARMv7 hardware:

```sh
tar xzf matrix.tar.gz
cd matrix-package
./matrix-linux-armv6 start --log-level info --gamepad-port 80 \
    core device gamepad emulator demo zigzag yumyum clock draw blocks getout
```

Run it from the directory holding `themes/`, `fonts/`, `images/` and `animations/`: core
reads them from the working directory and tolerates them missing without saying so.

On anything that is not a Raspberry Pi, build from source with `make build` — see
[development.md](./development.md).
