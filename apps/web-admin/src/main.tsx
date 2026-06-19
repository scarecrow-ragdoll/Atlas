// FILE: apps/web-admin/src/main.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Bootstrap the web-admin Vite React application.
//   SCOPE: Mounts the app into the Vite HTML root with providers; excludes route definitions and business behavior.
//   DEPENDS: react, react-dom/client, apps/web-admin/src/App.tsx, apps/web-admin/src/app/providers.tsx, apps/web-admin/src/styles.css.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: NONE
// END_MODULE_CONTRACT
// START_MODULE_MAP
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Vite React bootstrap for web-admin.
// END_CHANGE_SUMMARY

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { Providers } from './app/providers';
import './styles.css';

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found');
}

createRoot(rootElement).render(
  <StrictMode>
    <Providers>
      <App />
    </Providers>
  </StrictMode>,
);
