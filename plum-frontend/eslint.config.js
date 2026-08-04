import pluginVue from 'eslint-plugin-vue';
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript';

export default defineConfigWithVueTs(
    {
        name: 'plum/files',
        files: ['**/*.ts', '**/*.vue'],
    },
    {
        name: 'plum/ignores',
        ignores: ['dist/**', 'node_modules/**'],
    },

    pluginVue.configs['flat/essential'],
    vueTsConfigs.recommended,

    {
        name: 'plum/rules',
        rules: {
            'vue/multi-word-component-names': 'off',
            'vue/no-reserved-component-names': 'off',
        },
    },
);
