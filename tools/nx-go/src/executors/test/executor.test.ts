import type { ExecutorContext } from '@nx/devkit';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import runExecutor from './executor';

const execSyncMock = vi.hoisted(() => vi.fn());

vi.mock('child_process', () => ({
  execSync: execSyncMock,
}));

const context = {
  root: '/repo',
  projectName: 'api',
  projectsConfigurations: {
    projects: {
      api: { root: 'apps/api' },
    },
  },
} as ExecutorContext;

describe('nx-go test executor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('runs go test with coverage, short mode, and explicit packages', async () => {
    execSyncMock.mockReturnValue(Buffer.from('ok'));

    const result = await runExecutor(
      { coverage: true, short: true, packages: ['./internal/service', './internal/handler'] },
      context,
    );

    expect(result).toEqual({ success: true });
    expect(execSyncMock).toHaveBeenCalledWith(
      'go test -short -coverprofile=coverage.out ./internal/service ./internal/handler',
      expect.objectContaining({
        cwd: '/repo/apps/api',
        stdio: 'inherit',
        env: expect.objectContaining(process.env),
      }),
    );
  });

  it('runs go test for all packages without optional flags', async () => {
    execSyncMock.mockReturnValue(Buffer.from('ok'));

    const result = await runExecutor({ coverage: false, short: false }, context);

    expect(result).toEqual({ success: true });
    expect(execSyncMock).toHaveBeenCalledWith(
      'go test ./...',
      expect.objectContaining({ cwd: '/repo/apps/api', stdio: 'inherit' }),
    );
  });

  it('returns false when go test fails', async () => {
    execSyncMock.mockImplementation(() => {
      throw new Error('failed');
    });

    await expect(runExecutor({ coverage: false, short: false }, context)).resolves.toEqual({
      success: false,
    });
  });
});
