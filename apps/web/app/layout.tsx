// FILE: apps/web/app/layout.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the public Next root layout.
//   SCOPE: Defines metadata, document shell, global styles, and provider wrapping; excludes page data fetching and route proxy behavior.
//   DEPENDS: apps/web/app/globals.css, apps/web/src/app/providers.tsx, next.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   metadata - Public web page metadata.
//   default - Next root layout wrapping page children.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added public Next root layout.
// END_CHANGE_SUMMARY

import type { Metadata } from 'next';
import { Providers } from '../src/app/providers';
import './globals.css';

export const metadata: Metadata = {
  title: 'Monorepo Template',
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
