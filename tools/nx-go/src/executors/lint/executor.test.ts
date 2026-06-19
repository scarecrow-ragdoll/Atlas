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

describe('nx-go lint executor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('runs golangci-lint without optional flags', async () => {
    execSyncMock.mockReturnValue(Buffer.from('ok'));

    const result = await runExecutor({ fix: false }, context);

    expect(result).toEqual({ success: true });
    expect(execSyncMock).toHaveBeenCalledWith(
      'golangci-lint run',
      expect.objectContaining({ cwd: '/repo/apps/api', stdio: 'inherit' }),
    );
  });

  it('runs golangci-lint with fix and config flags', async () => {
    execSyncMock.mockReturnValue(Buffer.from('ok'));

    const result = await runExecutor({ fix: true, config: '.golangci.yml' }, context);

    expect(result).toEqual({ success: true });
    expect(execSyncMock).toHaveBeenCalledWith(
      'golangci-lint run --fix --config=.golangci.yml',
      expect.objectContaining({
        cwd: '/repo/apps/api',
        stdio: 'inherit',
        env: expect.objectContaining(process.env),
      }),
    );
  });

  it('returns false when golangci-lint fails', async () => {
    execSyncMock.mockImplementation(() => {
      throw new Error('failed');
    });

    await expect(runExecutor({ fix: false }, context)).resolves.toEqual({ success: false });
  });
});
