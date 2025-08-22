import react from '@vitejs/plugin-react'
import * as path from "path";
import {defineConfig} from 'vite'
import {createHtmlPlugin} from 'vite-plugin-html'

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        react(),
        createHtmlPlugin({
            minify: true
        })
    ],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src")
        }
    },
    build: {
        outDir: './static',
        emptyOutDir: true,
    },
})
