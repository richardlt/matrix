# <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/logo.png" width="30"/>&#8239;Matrix

[![Go Report Card](https://goreportcard.com/badge/github.com/richardlt/matrix)](https://goreportcard.com/report/github.com/richardlt/matrix)

Video game console operating system that displays on a 16*9 RGB LED matrix.

<p align="center">
  <br/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/gamepad.gif" width="400"/>
  <br/>
  <br/>
</p>

## Existing softwares

| | Name | Description | |
| - | - | - | - |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/demo.png" width="60"/> | Demo | A demo software that uses all drivers from the SDK. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/demo.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/yumyum.png" width="60"/> | Yumyum | Eat all the candies with your monster to win the game. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/yumyum.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/clock.png" width="60"/> | Clock | What time is it? | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/clock.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/zigzag.png" width="60"/> | Zigzag | Turn left then right, eat candies but not yourself. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/zigzag.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/draw.png" width="60"/> | Draw | For those who like pixel art. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/draw.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/device.png" width="60"/> | Device | The Device software allows you to change the luminosity of the LEDs. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/device.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/blocks.png" width="60"/> | Blocks | A puzzle game, score a maximum of points by clearing complete lines. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/blocks.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/getout.png" width="60"/> | Getout | A labyrinth game, try to get out if you can. | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/getout.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/rollup-dice.png" width="60"/> | Rollup dice | Random dice generator (https://github.com/gwenker/matrix-rollup-dice). | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/rollup-dice.gif" width="150"/> |
| <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/animate.png" width="60"/> | Animate | Player for animations generated with Glediator (http://www.solderlab.de/index.php/software/glediator). | <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/animate.gif" width="150"/> |

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

## Production setup

Matrix is designed to run on a Raspberry Pi (at least model 3), it is composed by multiple softwares (core, device, gamepad...). All softwares can run on the Raspberry Pi but you can also start a software on your desk that will communicate remotely with the Matrix's core (with flag --core-uri).

Here are the few steps to install your own Matrix:

1. Download Matrix latest release [here](https://github.com/richardlt/matrix/releases). If you want to install it on Raspbian or Debian there is a .deb file available that will create a service to start Matrix automatically at boot.

2. Extract/install and run Matrix package.
```sh
$ dpkg -i matrix.deb # for Raspbian/Debian users
$ service matrix status
```
```sh
$ unzip matrix.zip # for others
$ cd matrix-package && ./matrix-[REPLACE_DEPENDING_OS] start --log-level info --gamepad-port 80 core device gamepad emulator demo zigzag yumyum clock draw blocks getout # select the right executable depending on your os 
```

3. Install firmware on the Arduino from file in Matrix source code (inside folder at ./device/firmware/firmware.ino). Source code can be downloaded from [release](https://github.com/richardlt/matrix/releases).

## Development setup (linux/darwin)

1. Requirements.
* [Go](https://golang.org/dl/) (version 1.12+)
* [Node.js](https://nodejs.org/en/download/) (with npm, version 10+)

2. Install JS projects dependencies.
```sh
$ make install-all
```

3. Run it.
```sh
$ export GO111MODULE=on
$ export GOPROXY=https://gocenter.io
$ go run main.go start --log-level info core gamepad emulator demo # you can start all other softwares by adding their names
$ (cd emulator && npm start)
$ (cd gamepad && npm start)
```

4. Open emulator at http://localhost:3001 and/or gamepad at http://localhost:4002.