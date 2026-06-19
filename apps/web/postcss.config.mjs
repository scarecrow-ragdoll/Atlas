// FILE: apps/web/postcss.config.mjs
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Configure PostCSS for the public Next app Tailwind v4 pipeline.
//   SCOPE: Registers the Tailwind PostCSS plugin used by Next CSS compilation; excludes Vite web-admin styling.
//   DEPENDS: @tailwindcss/postcss, apps/web/app/globals.css.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: CONFIG
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - PostCSS plugin configuration for Tailwind v4 in apps/web.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Tailwind v4 PostCSS integration for public web.
// END_CHANGE_SUMMARY

const config = {
  plugins: {
    '@tailwindcss/postcss': {},
  },
};

export default config;
