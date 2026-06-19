import { createHash } from 'node:crypto';
import { execFileSync } from 'node:child_process';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const e2eDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(e2eDir, '../../..');
const composeFile = resolve(repoRoot, 'docker/docker-compose.test.yml');
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';
const resourcePrefix =
  process.env.TEST_RESOURCE_PREFIX ??
  `mt-test-${createHash('sha1').update(repoRoot).digest('hex').slice(0, 8)}`;
const composeEnv = {
  ...process.env,
  TEST_COMPOSE_PROJECT: process.env.TEST_COMPOSE_PROJECT ?? resourcePrefix,
  TEST_POSTGRES_CONTAINER_NAME:
    process.env.TEST_POSTGRES_CONTAINER_NAME ?? `${resourcePrefix}-postgres`,
  TEST_REDIS_CONTAINER_NAME: process.env.TEST_REDIS_CONTAINER_NAME ?? `${resourcePrefix}-redis`,
  TEST_POSTGRES_VOLUME: process.env.TEST_POSTGRES_VOLUME ?? `${resourcePrefix}-pg-test-data`,
};

function run(command, args) {
  console.log(`[e2e:preflight] ${command} ${args.join(' ')}`);
  execFileSync(command, args, { cwd: repoRoot, stdio: 'inherit', env: composeEnv });
}

function assertSafeTestTarget() {
  if (testPostgresDB !== 'monorepo_test') {
    throw new Error(`[e2e:preflight] unsafe test database target: ${testPostgresDB}`);
  }
  if (testPostgresPort === '7501') {
    throw new Error('[e2e:preflight] unsafe test postgres port: 7501 is the dev port');
  }
  if (testRedisPort === '7502') {
    throw new Error('[e2e:preflight] unsafe test redis port: 7502 is the dev port');
  }
}

assertSafeTestTarget();
console.log(
  [
    `[e2e:preflight] compose scope: project=${composeEnv.TEST_COMPOSE_PROJECT}`,
    `postgres=${composeEnv.TEST_POSTGRES_CONTAINER_NAME}`,
    `redis=${composeEnv.TEST_REDIS_CONTAINER_NAME}`,
    `volume=${composeEnv.TEST_POSTGRES_VOLUME}`,
  ].join(' '),
);
run('docker', ['compose', '-f', composeFile, 'up', '-d', '--wait', 'postgres', 'redis']);
run('docker', ['compose', '-f', composeFile, 'ps', 'postgres', 'redis']);
console.log(`[e2e:preflight] test docker services ready: postgres:${testPostgresPort} redis:${testRedisPort}`);
