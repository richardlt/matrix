# <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/logo.png" width="30"/>&#8239;Matrix

[![Go Report Card](https://goreportcard.com/badge/github.com/richardlt/matrix)](https://goreportcard.com/report/github.com/richardlt/matrix)

Video game console operating system that displays on a 16*9 RGB LED matrix.

<p align="center">
  <br/>
  <img src="./docs/gamepad.gif" width="400"/>
  <br/>
  <br/>
</p>

## Existing softwares

| | Name | Description | |
| - | - | - | - |
| <img src="./docs/demo.png" width="60"/> | Demo | A demo software that uses all drivers from the SDK. | <img src="./docs/demo.gif" width="150"/> |
| <img src="./docs/yumyum.png" width="60"/> | Yumyum | Eat all the candies with your monster to win the game. | <img src="./docs/yumyum.gif" width="150"/> |
| <img src="./docs/clock.png" width="60"/> | Clock | What time is it? | <img src="./docs/clock.gif" width="150"/> |
| <img src="./docs/zigzag.png" width="60"/> | Zigzag | Turn left then right, eat candies but not yourself. | <img src="./docs/zigzag.gif" width="150"/> |
| <img src="./docs/draw.png" width="60"/> | Draw | For those who like pixel art. | <img src="./docs/draw.gif" width="150"/> |
| <img src="./docs/device.png" width="60"/> | Device | The Device software allows you to change the luminosity of the LEDs. | <img src="./docs/device.gif" width="150"/> |
| <img src="./docs/blocks.png" width="60"/> | Blocks | A puzzle game, score a maximum of points by clearing complete lines. | <img src="./docs/blocks.gif" width="150"/> |
| <img src="./docs/getout.png" width="60"/> | Getout | A labyrinth game, try to get out if you can. | <img src="./docs/getout.gif" width="150"/> |
| <img src="./docs/rollup-dice.png" width="60"/> | Rollup dice | Roll two dice, press A for a new throw. | <img src="./docs/rollup-dice.gif" width="150"/> |
| <img src="./docs/animate.png" width="60"/> | Animate | Player for animations generated with Glediator (http://www.solderlab.de/index.php/software/glediator). | <img src="./docs/animate.gif" width="150"/> |
| <img src="./docs/light.png" width="60"/> | Light | Simple software to generate mood light. | <img src="./docs/light.gif" width="150"/> |

## Matrix types

There are 3 main types that exists in Matrix's sdk:
- A **display** receives a live stream of frames from core.
- A **player** sends actions for a slot to core, an action is a button press/release event.  
- A **software** receives player's actions from core and use sdk's rendering features to generate frames in Matrix core. 

## Matrix components

| Name | Description |
| - | - |
| Core | The heart of the Matrix system that managed software's lifecycle. All softwares, players and displays are connected to core. |
| Device | The component that interacts with usb controllers and Arduino. |
| Gamepad | A web application that contains a virtual controller with display. |
| Emulator | A web application built for development purpose. It displays Matrix main screen and player's screens. |

## Install

Matrix runs on a Raspberry Pi. On a Pi 2 or later there is a `.deb` in the
[releases](https://github.com/richardlt/matrix/releases) that installs a service starting
at boot; ARMv6 boards get an SD card staged by a script. Either way the Arduino has to be
flashed first.

- [Installing matrix](./docs/install.md)
- [Installing on an ARMv6 board, with Alpine](./docs/install-alpine.md)
- [Flashing the Arduino](./docs/arduino.md)

All components can run on the Pi, but any software can equally run on your desk against a
remote core with `--core-uri`.

## Contributing

See [development.md](./docs/development.md) for the development setup, the checks and how
the release artifacts are built.
