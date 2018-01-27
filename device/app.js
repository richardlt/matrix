const Matrix = require("./matrix");

const start = async () => {
  const matrix = new Matrix(144);
  await matrix.loadPorts();
  matrix.openPorts();

  setInterval(_ => {
    let frame = [];
    for (let i = 0; i < 144; i++) {
      frame.push({
        r: Math.floor(Math.random() * 255),
        g: Math.floor(Math.random() * 255),
        b: Math.floor(Math.random() * 255),
        a: 1
      });
    }
    matrix.setFrame(frame);
  }, 10);
};

start();