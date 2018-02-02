const IO = require('socket.io-client'),
  argv = require('minimist')(process.argv.slice(2));

const Matrix = require('./matrix/matrix'),
  Gamepad = require('./gamepad/gamepad');

const start = async () => {
  const matrix = new Matrix(144);
  await matrix.loadPorts();
  matrix.openPorts();

  const gamepad = new Gamepad();
  gamepad.loadDevices();
  gamepad.openDevices();

  const socket = IO(argv.uri ? argv.uri : 'http://localhost:5000');
  socket.on('connect', _ => { console.log('connect socket'); });
  socket.on('frame', data => { matrix.setFrame(data.pixels); });
  socket.on('disconnect', _ => { console.log('disconnect socket') });

  gamepad.on('event', event => { socket.emit('command', event); });
};

start();