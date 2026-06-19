import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface LintExecutorOptions {
  fix: boolean;
  config?: string;
}

export default async function runExecutor(
  options: LintExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  const args: string[] = ['golangci-lint', 'run'];
  if (options.fix) args.push('--fix');
  if (options.config) args.push(`--config=${options.config}`);

  console.log(`Linting Go code in ${cwd}...`);

  try {
    execSync(args.join(' '), {
      cwd,
      stdio: 'inherit',
      env: { ...process.env },
    });
    return { success: true };
  } catch {
    return { success: false };
  }
}
