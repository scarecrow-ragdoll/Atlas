import { execFileSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import fs from 'node:fs';
import path from 'node:path';

const root = process.cwd();
const configPath = path.join(root, 'tools/coverage/coverage.config.json');
const config = JSON.parse(fs.readFileSync(configPath, 'utf8'));
const testPostgresHost = process.env.TEST_POSTGRES_HOST ?? 'localhost';
const testPostgresPort = process.env.TEST_POSTGRES_PORT ?? '17501';
const testPostgresUser = process.env.TEST_POSTGRES_USER ?? 'app';
const testPostgresPassword = process.env.TEST_POSTGRES_PASSWORD ?? 'secret';
const testPostgresDB = process.env.TEST_POSTGRES_DB ?? 'monorepo_test';
const testRedisHost = process.env.TEST_REDIS_HOST ?? 'localhost';
const testRedisPort = process.env.TEST_REDIS_PORT ?? '17502';
const testRedisPassword = process.env.TEST_REDIS_PASSWORD ?? '';
const testDatabaseDSN =
  process.env.API_TEST_DATABASE_DSN ??
  `postgres://${testPostgresUser}:${testPostgresPassword}@${testPostgresHost}:${testPostgresPort}/${testPostgresDB}?sslmode=disable`;
const resourcePrefix =
  process.env.TEST_RESOURCE_PREFIX ??
  `mt-test-${createHash('sha1').update(root).digest('hex').slice(0, 8)}`;
const scopedTestEnv = {
  TEST_COMPOSE_PROJECT: process.env.TEST_COMPOSE_PROJECT ?? resourcePrefix,
  TEST_POSTGRES_CONTAINER_NAME:
    process.env.TEST_POSTGRES_CONTAINER_NAME ?? `${resourcePrefix}-postgres`,
  TEST_REDIS_CONTAINER_NAME: process.env.TEST_REDIS_CONTAINER_NAME ?? `${resourcePrefix}-redis`,
  TEST_POSTGRES_VOLUME: process.env.TEST_POSTGRES_VOLUME ?? `${resourcePrefix}-pg-test-data`,
};

function fail(message) {
  console.error(`[Coverage][gate] ${message}`);
  process.exitCode = 1;
}

function run(command, args, options = {}) {
  console.log(`[Coverage][run] ${command} ${args.join(' ')}`);
  execFileSync(command, args, {
    cwd: options.cwd || root,
    stdio: 'inherit',
    env: { ...process.env, ...options.env },
  });
}

function coverageGateEnv() {
  return {
    COVERAGE_GATE: '1',
    API_TEST_DATABASE_DSN: testDatabaseDSN,
    TEST_POSTGRES_HOST: testPostgresHost,
    TEST_POSTGRES_PORT: testPostgresPort,
    TEST_POSTGRES_USER: testPostgresUser,
    TEST_POSTGRES_PASSWORD: testPostgresPassword,
    TEST_POSTGRES_DB: testPostgresDB,
    TEST_REDIS_HOST: testRedisHost,
    TEST_REDIS_PORT: testRedisPort,
    TEST_REDIS_PASSWORD: testRedisPassword,
    ...scopedTestEnv,
  };
}

function assertSafeTestTarget() {
  if (testPostgresDB !== 'monorepo_test') {
    throw new Error(`[Coverage][run] unsafe target: test database target is ${testPostgresDB}`);
  }
  if (testPostgresPort === '7501') {
    throw new Error('[Coverage][run] unsafe target: test postgres port 7501 is the dev port');
  }
  if (testRedisPort === '7502') {
    throw new Error('[Coverage][run] unsafe target: test redis port 7502 is the dev port');
  }

  let dsn;
  try {
    dsn = new URL(testDatabaseDSN);
  } catch (error) {
    throw new Error(`[Coverage][run] unsafe target: malformed API_TEST_DATABASE_DSN (${error.message})`);
  }

  if (dsn.protocol !== 'postgres:' && dsn.protocol !== 'postgresql:') {
    throw new Error(`[Coverage][run] unsafe target: unsupported database DSN protocol ${dsn.protocol}`);
  }

  const dsnDB = dsn.pathname.replace(/^\/+/, '');
  if (dsnDB !== 'monorepo_test') {
    throw new Error(`[Coverage][run] unsafe target: database in DSN is ${dsnDB || '<empty>'}`);
  }
  if (dsnDB === 'monorepo_dev') {
    throw new Error('[Coverage][run] unsafe target: database in DSN is monorepo_dev');
  }
  if (dsn.port === '7501') {
    throw new Error('[Coverage][run] unsafe target: DSN postgres port 7501 is the dev port');
  }
}

function ensureDirFor(file) {
  fs.mkdirSync(path.dirname(path.join(root, file)), { recursive: true });
}

function isAllowlistedGoFile(file) {
  const normalizedFile = file.replaceAll('\\', '/');
  return config.allowlist
    .filter((item) => item.path.endsWith('.go'))
    .some((item) => {
      const normalizedPath = item.path.replaceAll('\\', '/');
      return normalizedFile.endsWith(normalizedPath) || normalizedFile.includes(normalizedPath);
    });
}

function parseGoFilteredTotal(profile) {
  const lines = fs.readFileSync(path.join(root, profile), 'utf8').trim().split('\n');
  let statements = 0;
  let covered = 0;

  for (const line of lines) {
    if (line === 'mode: set' || line === 'mode: count' || line === 'mode: atomic') {
      continue;
    }
    const match = line.match(/^(.+):\d+\.\d+,\d+\.\d+\s+(\d+)\s+(\d+)$/);
    if (!match) {
      throw new Error(`Cannot parse Go coverage line in ${profile}: ${line}`);
    }
    const [, file, statementCountRaw, hitCountRaw] = match;
    if (isAllowlistedGoFile(file)) {
      continue;
    }
    const statementCount = Number(statementCountRaw);
    const hitCount = Number(hitCountRaw);
    statements += statementCount;
    if (hitCount > 0) {
      covered += statementCount;
    }
  }

  if (statements === 0) {
    throw new Error(`No non-allowlisted Go statements found in ${profile}`);
  }

  return Number(((covered / statements) * 100).toFixed(1));
}

function readJson(file) {
  return JSON.parse(fs.readFileSync(path.join(root, file), 'utf8'));
}

run('node', ['tools/coverage/preflight.mjs']);
assertSafeTestTarget();
console.log(
  [
    `[Coverage][run] compose scope: project=${scopedTestEnv.TEST_COMPOSE_PROJECT}`,
    `postgres=${scopedTestEnv.TEST_POSTGRES_CONTAINER_NAME}`,
    `redis=${scopedTestEnv.TEST_REDIS_CONTAINER_NAME}`,
    `volume=${scopedTestEnv.TEST_POSTGRES_VOLUME}`,
  ].join(' '),
);
run(
  'docker',
  ['compose', '-f', 'docker/docker-compose.test.yml', 'up', '-d', '--wait', 'postgres', 'redis'],
  { env: scopedTestEnv },
);

for (const project of config.goProjects) {
  ensureDirFor(project.profile);
  run('go', ['test', `-coverprofile=${path.join(root, project.profile)}`, project.packages], {
    cwd: path.join(root, project.cwd),
    env: coverageGateEnv(),
  });
  const total = parseGoFilteredTotal(project.profile);
  if (total !== config.thresholds.goStatements) {
    fail(`${project.name}: Go statement coverage ${total}% != ${config.thresholds.goStatements}%`);
  }
}

run('bunx', ['nx', 'run', 'web-admin:test-coverage']);
run('bunx', ['nx', 'run', 'web:test-coverage']);
run('bunx', ['nx', 'run', 'nx-go:test']);
run('bunx', ['nx', 'run', 'codegen:validate']);

for (const item of config.typescriptCoverageSummaries) {
  const summary = readJson(item.summary).total;
  const checks = [
    ['statements', summary.statements.pct, config.thresholds.typescriptStatements],
    ['branches', summary.branches.pct, config.thresholds.typescriptBranches],
    ['functions', summary.functions.pct, config.thresholds.typescriptFunctions],
    ['lines', summary.lines.pct, config.thresholds.typescriptLines],
  ];
  for (const [metric, actual, expected] of checks) {
    if (actual !== expected) {
      fail(`${item.name}: ${metric} coverage ${actual}% != ${expected}%`);
    }
  }
}

if (process.exitCode) {
  process.exit(process.exitCode);
}

console.log('[Coverage][gate] all thresholds passed');
