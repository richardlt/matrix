# <img src="https://raw.githubusercontent.com/richardlt/matrix/master/logo.png" width="30"/>&#8239;Matrix

[![Go Report Card](https://goreportcard.com/badge/github.com/richardlt/matrix)](https://goreportcard.com/report/github.com/richardlt/matrix)

Video game console operating system that displays on a 16*9 RGB LED matrix.

<p align="center">
  <br/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/demo.gif" height="75%"/><br/>
</p>

## Existing softwares

<p align="center">
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/demo.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/yumyum.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/clock.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/zigzag.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/draw.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/device.gif" height="75%"/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/docs/blocks.gif" height="75%"/>
</p>

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
$ go run main.go all
$ (cd $GOPATH/src/github.com/richardlt/matrix/emulator && npm start)
$ (cd $GOPATH/src/github.com/richardlt/matrix/gamepad && npm start)
```

4. Open emulator at http://localhost:3001 and/or gamepad at http://localhost:4002.