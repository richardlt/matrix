import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    // Built into the Go binary by the emulator package's go:embed directive.
    // Vite owns this subdirectory and empties it on each build. The embed root is the
    // parent, which also holds a checked-in .gitkeep so the go:embed pattern still
    // matches on a fresh clone.
    outDir: 'public/app',
    emptyOutDir: true,
  },
  server: {
    port: 3001,
    proxy: {
      // Frames reach the browser over this socket; core is behind the emulator server.
      '/websocket': {
        target: 'http://localhost:3000',
        ws: true,
      },
    },
  },
});
