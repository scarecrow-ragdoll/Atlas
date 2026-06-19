import { describe, expect, it } from 'vitest';
import config from './codegen';

describe('codegen config', () => {
  it('uses shared GraphQL schema and web-admin operation documents', () => {
    expect(config.schema).toBe('../../libs/graphql/schema/**/*.graphql');
    expect(config.documents).toEqual([
      '../../apps/web-admin/src/features/**/api/**/*.graphql',
      '../../apps/web-admin/src/entities/**/api/**/*.graphql',
    ]);
  });

  it('generates web-admin types with expected plugins and scalar mappings', () => {
    const output = config.generates?.['../../apps/web-admin/src/shared/api/generated/types.ts'];

    expect(output).toBeDefined();
    expect(output).toMatchObject({
      plugins: ['typescript', 'typescript-operations'],
      config: {
        scalars: {
          DateTime: 'string',
          UUID: 'string',
        },
      },
    });
  });

  it('allows empty document globs so a new template can codegen before features exist', () => {
    expect(config.ignoreNoDocuments).toBe(true);
  });
});
