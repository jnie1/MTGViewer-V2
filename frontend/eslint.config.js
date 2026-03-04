import pluginVue from 'eslint-plugin-vue';
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript';
import pluginVitest from '@vitest/eslint-plugin';
import prettierConfig from '@vue/eslint-config-prettier';

export default defineConfigWithVueTs([
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },

  {
    name: 'app/files-to-ignore',
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
  },

  ...pluginVue.configs['flat/recommended-error'],
  {
    name: 'vue/rule-overrides',
    rules: {
      'vue/v-bind-style': ['error', 'shorthand', { sameNameShorthand: 'always' }],
      'vue/component-name-in-template-casing': ['error', 'kebab-case'],
      'vue/component-api-style': ['error', ['script-setup']],
      'vue/require-default-prop': ['off'],
    },
  },
  vueTsConfigs.strict,

  {
    ...pluginVitest.configs.recommended,
    files: ['src/**/__tests__/*'],
  },

  prettierConfig,
]);
