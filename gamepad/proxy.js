const browserSync = require('browser-sync');
const proxy = require('http-proxy-middleware');

browserSync.init({
  port: 4002,
  ui: false,
  files: ["src/**/*.*", "index.html", "*.js"],
  proxy: {
    target: "localhost:4001",
    middleware: [proxy("/socket.io/**", {
      target: 'http://localhost:4000',
      ws: true,
      changeOrigin: true
    })]
  },
  open: false
});