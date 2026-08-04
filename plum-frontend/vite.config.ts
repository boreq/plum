import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

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
        proxy: {
            '/api': 'http://localhost:8118',
        },
    },
    build: {
        sourcemap: false,
    },
});
