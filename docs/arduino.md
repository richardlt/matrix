# Flashing the Arduino

The board drives the LED panel and talks to the `device` component over USB serial. The
firmware is [`device/firmware/firmware.ino`](../device/firmware/firmware.ino).

Check the three values at the top of the sketch first:

| | |
| - | - |
| `PIN` | the Arduino pin wired to the panel's data line |
| `NUMPIXELS` | number of LEDs, 144 for a 16x9 panel |
| `BAUD_RATE` | must equal `serialBaudRate` in [`device/matrix.go`](../device/matrix.go) |

[arduino-cli](https://github.com/arduino/arduino-cli/releases) is a single binary and
needs no installer. The commands below are the same on Linux and Windows; only the port
name differs.

```sh
arduino-cli config init
arduino-cli core update-index
arduino-cli core install arduino:avr        # compiler and avrdude
arduino-cli lib install "Adafruit NeoPixel" # the sketch needs this to compile

arduino-cli board list                      # find the port and the FQBN
arduino-cli compile --fqbn arduino:avr:nano:cpu=atmega328old device/firmware
arduino-cli upload -p /dev/ttyUSB0 --fqbn arduino:avr:nano:cpu=atmega328old device/firmware
```

Use `-p COM3` or similar on Windows. `atmega328old` covers the clones, which nearly all
carry the old bootloader; a genuine recent Nano is `arduino:avr:nano:cpu=atmega328`, and
an Uno is `arduino:avr:uno`. The [Arduino IDE](https://www.arduino.cc/en/software) does
the same job if you prefer it — add `Adafruit NeoPixel` through **Manage Libraries** and
set **Processor → ATmega328P (Old Bootloader)** on a clone.

A successful upload ends with `avrdude done. Thank you.`, and matrix then logs
`Serial /dev/ttyUSB0 connected` when it finds the board.
