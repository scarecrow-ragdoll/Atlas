import { describe, expect, it } from 'vitest';
import buildExecutor from './executors/build/executor';
import lintExecutor from './executors/lint/executor';
import serveExecutor from './executors/serve/executor';
import testExecutor from './executors/test/executor';
import {
  buildExecutor as exportedBuildExecutor,
  lintExecutor as exportedLintExecutor,
  serveExecutor as exportedServeExecutor,
  testExecutor as exportedTestExecutor,
} from './index';

describe('nx-go public exports', () => {
  it('re-exports all executor entrypoints', () => {
    expect(exportedBuildExecutor).toBe(buildExecutor);
    expect(exportedLintExecutor).toBe(lintExecutor);
    expect(exportedServeExecutor).toBe(serveExecutor);
    expect(exportedTestExecutor).toBe(testExecutor);
  });
});
