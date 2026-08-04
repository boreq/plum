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
            // Components are named after the thing they render, e.g. Table.
            'vue/multi-word-component-names': 'off',

            // Table and Summary shadow HTML element names. They are only ever
            // used through explicit `components` registration, where Vue
            // resolves the component before the native element.
            'vue/no-reserved-component-names': 'off',
        },
    },
);
