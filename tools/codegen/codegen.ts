import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../../libs/graphql/schema/**/*.graphql',
  documents: [
    '../../apps/web-admin/src/features/**/api/**/*.graphql',
    '../../apps/web-admin/src/entities/**/api/**/*.graphql',
  ],
  ignoreNoDocuments: true,
  generates: {
    '../../apps/web-admin/src/shared/api/generated/types.ts': {
      plugins: ['typescript', 'typescript-operations'],
      config: {
        scalars: {
          DateTime: 'string',
          UUID: 'string',
        },
      },
    },
  },
};

export default config;
