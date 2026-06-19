// FILE: apps/web/next.config.js
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure the public Next web application.
//   SCOPE: Enables standalone output for Docker/runtime deployment; excludes app route behavior and CI orchestration.
//   DEPENDS: next.
//   LINKS: M-WEB / V-M-WEB / M-CI-CD.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Next.js standalone build configuration for public web.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Next standalone configuration for public web.
// END_CHANGE_SUMMARY

const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
};

export default nextConfig;
