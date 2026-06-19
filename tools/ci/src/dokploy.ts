export type FetchResponse = {
  ok: boolean;
  status: number;
  text(): Promise<string>;
};

export type FetchLike = (
  url: string,
  init: {
    method: string;
    headers: Record<string, string>;
    body?: string;
  },
) => Promise<FetchResponse>;

export type DeployDokployComposeInput = {
  baseUrl: string;
  apiKey: string;
  composeId: string;
  imageEnv: Record<string, string>;
  fetchImpl?: FetchLike;
};

function normalizeBaseUrl(baseUrl: string): string {
  return baseUrl.replace(/\/+$/, '');
}

function mergeEnv(existingEnv: string, updates: Record<string, string>): string {
  const seen = new Set<string>();
  const mergedLines = existingEnv
    .split('\n')
    .filter((line) => line.trim() !== '')
    .map((line) => {
      const match = line.match(/^([A-Za-z_][A-Za-z0-9_]*)=(.*)$/);
      if (!match) {
        return line;
      }

      const key = match[1];
      if (!(key in updates)) {
        return line;
      }

      seen.add(key);
      return `${key}=${updates[key]}`;
    });

  for (const [key, value] of Object.entries(updates)) {
    if (!seen.has(key)) {
      mergedLines.push(`${key}=${value}`);
    }
  }

  return mergedLines.join('\n');
}

function defaultFetch(url: string, init: Parameters<FetchLike>[1]): Promise<FetchResponse> {
  return fetch(url, init) as Promise<FetchResponse>;
}

async function callDokploy(
  path: string,
  input: DeployDokployComposeInput,
  options: { method: string; body?: Record<string, unknown> },
): Promise<string> {
  const fetchImpl = input.fetchImpl || defaultFetch;
  const response = await fetchImpl(`${normalizeBaseUrl(input.baseUrl)}/api/${path}`, {
    method: options.method,
    headers: {
      'Content-Type': 'application/json',
      'x-api-key': input.apiKey,
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  if (!response.ok) {
    throw new Error(`Dokploy ${path} failed with status ${response.status}: ${text}`);
  }

  return text;
}

async function getComposeEnv(input: DeployDokployComposeInput): Promise<string> {
  const text = await callDokploy(
    `compose.one?composeId=${encodeURIComponent(input.composeId)}`,
    input,
    {
      method: 'GET',
    },
  );
  const data = JSON.parse(text) as { env?: string | null };

  return data.env ?? '';
}

async function updateComposeEnv(input: DeployDokployComposeInput, env: string): Promise<void> {
  await callDokploy('compose.update', input, {
    method: 'POST',
    body: {
      composeId: input.composeId,
      env,
    },
  });
}

async function deployCompose(input: DeployDokployComposeInput): Promise<void> {
  await callDokploy('compose.deploy', input, {
    method: 'POST',
    body: {
      composeId: input.composeId,
    },
  });
}

export async function deployDokployCompose(input: DeployDokployComposeInput): Promise<void> {
  const currentEnv = await getComposeEnv(input);
  await updateComposeEnv(input, mergeEnv(currentEnv, input.imageEnv));
  await deployCompose(input);
}
