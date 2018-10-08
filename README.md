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

## Development setup

1. Requirements.
* [Go](https://golang.org/dl/) (version 1.10+)
* [Node.js](https://nodejs.org/en/download/) (with npm, version 8+)
* [Bower](https://bower.io/) (latest)
* [Polymer CLI](https://www.polymer-project.org/2.0/docs/tools/polymer-cli) (latest)

2. Install JS projects dependencies.
```sh
$ make install
```

3. Run it.
```sh
$ go run main.go start --log-level info core gamepad emulator demo
$ (cd $GOPATH/src/github.com/richardlt/matrix/emulator && npm start)
$ (cd $GOPATH/src/github.com/richardlt/matrix/gamepad && npm start)
```

4. Open emulator at http://localhost:3001 and/or gamepad at http://localhost:4002.