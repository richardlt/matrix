# Installing on an ARMv6 board, with Alpine

For a Raspberry Pi Zero or a Pi 1. The card is prepared in one go, and the console then
needs no keyboard, no screen and no network: it boots and plays.

## What you need

- An SD card, 2 GB is enough, and a card reader.
- `matrix.tar.gz` from the [releases](https://github.com/richardlt/matrix/releases), or
  built with `make package` (see [development.md](./development.md)).
- The Arduino [already flashed](./arduino.md).
- `curl`, `tar` and `sha256sum` on the machine preparing the card.

## 1. Stage the card contents

```sh
tar xzf matrix.tar.gz
./matrix-package/prepare-alpine-card.sh
```

The script downloads Alpine, verifies its checksum, and builds the directory tree that
belongs on the card in `./alpine-card`. It never formats, partitions or writes to a block
device.

Two options are worth knowing about:

| | |
| - | - |
| `--components "core device clock"` | start fewer components. The default is `core`, `device` and every game |
| `--debug` | log at debug level to `matrix.log` on the card. Off by default, because nothing rotates that file |

`--help` lists the rest.

## 2. Write the card

Format it as a single FAT32 partition, marked bootable, then copy the **contents** of
`alpine-card/` to its root — not the directory itself. Alpine's Pi release is a tarball
extracted onto a filesystem rather than an image written to a disk, which is why no `dd`
and no root are needed.

## 3. Play

Put the card in the board and power it on. Matrix starts by itself, and pulling the power
is a fine way to switch it off.

If the panel stays dark, re-stage the card with `--debug` and read `matrix.log` from it
afterwards. The previous boot's log is kept as `matrix.log.prev`.
