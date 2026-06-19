import { EventEmitter } from 'node:events';
import type { ExecutorContext } from '@nx/devkit';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import runExecutor from './executor';

const spawnMock = vi.hoisted(() => vi.fn());

vi.mock('child_process', () => ({
  spawn: spawnMock,
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

function childProcess() {
  const child = new EventEmitter() as EventEmitter & { kill: ReturnType<typeof vi.fn> };
  child.kill = vi.fn();
  return child;
}

describe('nx-go serve executor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    process.removeAllListeners('SIGINT');
    process.removeAllListeners('SIGTERM');
  });

  it('spawns air with config path and API_PORT', async () => {
    const child = childProcess();
    spawnMock.mockReturnValue(child);

    const promise = runExecutor({ port: 8080, configPath: 'custom.air.toml' }, context);
    child.emit('close', 0);

    await expect(promise).resolves.toEqual({ success: true });
    expect(spawnMock).toHaveBeenCalledWith(
      'air',
      ['-c', 'custom.air.toml'],
      expect.objectContaining({
        cwd: '/repo/apps/api',
        stdio: 'inherit',
        env: expect.objectContaining({ API_PORT: '8080' }),
      }),
    );
  });

  it('uses the default air config when none is supplied', async () => {
    const child = childProcess();
    spawnMock.mockReturnValue(child);

    const promise = runExecutor({ port: 8080 }, context);
    child.emit('close', null);

    await expect(promise).resolves.toEqual({ success: true });
    expect(spawnMock).toHaveBeenCalledWith(
      'air',
      ['-c', 'air.toml'],
      expect.objectContaining({ cwd: '/repo/apps/api' }),
    );
  });

  it('returns false for non-zero close code', async () => {
    const child = childProcess();
    spawnMock.mockReturnValue(child);

    const promise = runExecutor({ port: 8080 }, context);
    child.emit('close', 1);

    await expect(promise).resolves.toEqual({ success: false });
  });

  it('forwards termination signals and cleans listeners after close', async () => {
    const child = childProcess();
    spawnMock.mockReturnValue(child);

    const promise = runExecutor({ port: 8080 }, context);
    process.emit('SIGTERM', 'SIGTERM');
    child.emit('close', 0);

    await promise;
    expect(child.kill).toHaveBeenCalledWith('SIGTERM');
    expect(process.listenerCount('SIGTERM')).toBe(0);
    expect(process.listenerCount('SIGINT')).toBe(0);
  });
});
