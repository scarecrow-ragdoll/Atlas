module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'scope-enum': [
      2,
      'always',
      [
        'api',
        'bot',
        'web',
        'graphql',
        'codegen',
        'nx-go',
        'docker',
        'ci',
        'deps',
        'logger',
        'config',
        'docs',
      ],
    ],
    'scope-empty': [1, 'never'],
  },
};
