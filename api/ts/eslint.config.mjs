import js from '@eslint/js';
import prettier from 'eslint-config-prettier';
import parser from '@typescript-eslint/parser';

export default [
  js.configs.recommended,
  prettier,
  {
    files: ['src/**/*.ts'],
    ignores: ['dist/**', 'src/gen/**'],
    languageOptions: {
      parser,
      parserOptions: {
        project: ['./tsconfig.json'],
      },
    },
  },
];
