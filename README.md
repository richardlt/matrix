# <img src="https://raw.githubusercontent.com/richardlt/matrix/master/logo.png" width="30"/>&#8239;Matrix

Video game console operating system that displays on a 16*9 RGB LED matrix. WIP.

<p align="center">
  <br/>
  <img src="https://raw.githubusercontent.com/richardlt/matrix/master/demo.gif" height="75%"/><br/>
</p>

## Development setup

1. Requirements
* [Go](https://golang.org/dl/) (version 1.9+)
* [Node.js](https://nodejs.org/en/download/) (with npm, version 8+)
* [Bower](https://bower.io/) (latest)

2. Install JS projects dependencies
```sh
$ make install
```

3. Run it
```sh
$ go run main.go all
$ (cd $GOPATH/src/github.com/richardlt/matrix/emulator && npm start)
$ (cd $GOPATH/src/github.com/richardlt/matrix/gamepad && npm start)
```