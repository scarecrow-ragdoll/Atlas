import { ExecutorContext } from '@nx/devkit';
import { execSync } from 'child_process';
import * as path from 'path';

interface BuildExecutorOptions {
  outputPath: string;
  main: string;
}

export default async function runExecutor(
  options: BuildExecutorOptions,
  context: ExecutorContext,
): Promise<{ success: boolean }> {
  const projectRoot = context.projectsConfigurations!.projects[context.projectName!].root;
  const cwd = path.join(context.root, projectRoot);
  const outputPath = path.join(context.root, options.outputPath);

  console.log(`Building Go application in ${cwd}...`);

  try {
    execSync(`go build -o ${outputPath} ./${options.main}`, {
      cwd,
      stdio: 'inherit',
      env: { ...process.env },
    });
    return { success: true };
  } catch {
    return { success: false };
  }
}
