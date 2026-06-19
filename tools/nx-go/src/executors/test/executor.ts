import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface TestExecutorOptions {
  coverage: boolean;
  short: boolean;
  packages?: string[];
}

export default async function runExecutor(
  options: TestExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);

  const args: string[] = ['go', 'test'];

  if (options.short) args.push('-short');
  if (options.coverage) args.push('-coverprofile=coverage.out');

  const pkgs = options.packages?.length ? options.packages.join(' ') : './...';
  args.push(pkgs);

  console.log(`Running Go tests in ${cwd}...`);

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
