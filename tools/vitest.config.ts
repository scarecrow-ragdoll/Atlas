import { resolve } from 'node:path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    globals: true,
    include: ['tools/**/*.test.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary'],
      reportsDirectory: resolve(__dirname, '../dist/coverage/tools'),
      include: ['tools/nx-go/src/**/*.ts', 'tools/codegen/**/*.ts', 'tools/ci/src/**/*.ts'],
      exclude: ['tools/**/*.test.ts', 'tools/**/schema.json', 'tools/**/package.json'],
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
});
