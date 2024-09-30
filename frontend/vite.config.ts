import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '', // https://stackoverflow.com/questions/69744253/vite-build-always-using-static-paths
})
