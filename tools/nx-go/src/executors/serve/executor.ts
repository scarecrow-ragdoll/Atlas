import { ExecutorContext } from '@nx/devkit';
import { spawn } from 'child_process';
import * as path from 'path';

interface ServeExecutorOptions {
  port: number;
  configPath?: string;
}

export default async function runExecutor(
  options: ServeExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  console.log(`Serving Go application on port ${options.port}...`);

  const airConfig = options.configPath || 'air.toml';

  return new Promise((resolve) => {
    const child = spawn('air', ['-c', airConfig], {
      cwd,
      stdio: 'inherit',
      env: { ...process.env, API_PORT: String(options.port) },
    });

    const signalHandler = (signal: NodeJS.Signals) => {
      child.kill(signal);
    };
    process.on('SIGINT', signalHandler);
    process.on('SIGTERM', signalHandler);

    child.on('close', (code) => {
      process.removeListener('SIGINT', signalHandler);
      process.removeListener('SIGTERM', signalHandler);
      resolve({ success: code === 0 || code === null });
    });
  });
}
