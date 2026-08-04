import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// Injected into every stylesheet so that components can use the shared variables
// without importing them explicitly, mirroring the old vue-cli `sass.data` option.
const variables = fileURLToPath(new URL('./src/scss/variables.scss', import.meta.url));

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
    css: {
        preprocessorOptions: {
            scss: {
                additionalData: `@use "${variables}" as *;\n`,
            },
        },
    },
    server: {
        // In production the API is served from the same origin under /api/, so
        // the dev server forwards that path to a locally running backend.
        proxy: {
            '/api': 'http://localhost:8118',
        },
    },
    build: {
        sourcemap: false,
    },
});
