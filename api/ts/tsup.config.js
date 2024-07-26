import { defineConfig } from 'tsup';

export default defineConfig({
  clean: true, // Clean output directory before bundling
  dts: true, // Generate .d.ts file
  entry: ['src/**/*.ts', 'src/**/*.tsx'],
  format: ['esm', 'cjs'],
  sourcemap: true, // Generate sourcemaps
  splitting: true, // Enable code splitting into chunks
  treeshake: true, // Remove dead code
});
