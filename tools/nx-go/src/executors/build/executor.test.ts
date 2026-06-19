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

describe('nx-go build executor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('runs go build with output path and main package', async () => {
    execSyncMock.mockReturnValue(Buffer.from('ok'));

    const result = await runExecutor({ outputPath: 'dist/apps/api', main: 'cmd/server' }, context);

    expect(result).toEqual({ success: true });
    expect(execSyncMock).toHaveBeenCalledWith(
      'go build -o /repo/dist/apps/api ./cmd/server',
      expect.objectContaining({
        cwd: '/repo/apps/api',
        stdio: 'inherit',
        env: expect.objectContaining(process.env),
      }),
    );
  });

  it('returns false when go build fails', async () => {
    execSyncMock.mockImplementation(() => {
      throw new Error('failed');
    });

    await expect(
      runExecutor({ outputPath: 'dist/apps/api', main: 'cmd/server' }, context),
    ).resolves.toEqual({ success: false });
  });
});
