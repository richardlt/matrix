import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    // Built into the Go binary by the gamepad package's go:embed directive.
    // Vite owns this subdirectory and empties it on each build. The embed root is the
    // parent, which also holds a checked-in .gitkeep so the go:embed pattern still
    // matches on a fresh clone.
    outDir: 'public/app',
    emptyOutDir: true,
  },
  server: {
    port: 4002,
    proxy: {
      // Player commands go up and display frames come down over this socket.
      '/websocket': {
        target: 'http://localhost:4000',
        ws: true,
      },
    },
  },
});
