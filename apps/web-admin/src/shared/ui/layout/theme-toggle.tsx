'use client';

// FILE: apps/web-admin/src/shared/ui/layout/theme-toggle.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the approved web-admin light/dark theme toggle.
//   SCOPE: Reads and persists the admin theme preference, toggles the root dark class, and renders an accessible icon button; excludes route-specific placement.
//   DEPENDS: react, lucide-react, apps/web-admin/src/shared/ui/primitives/button.tsx.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ThemeToggle - Admin icon button for switching between light and dark themes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added persisted admin light/dark theme switching.
// END_CHANGE_SUMMARY

import { useEffect, useState } from 'react';
import { MoonIcon, SunIcon } from 'lucide-react';
import { Button } from '../primitives/button';

const storageKey = 'web-admin-theme';

type Theme = 'light' | 'dark';

function applyTheme(theme: Theme) {
  document.documentElement.classList.toggle('dark', theme === 'dark');
}

// START_CONTRACT: ThemeToggle
//   PURPOSE: Render a persisted accessible admin theme switch button.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - icon button whose accessible name describes the next theme }
//   SIDE_EFFECTS: Reads/writes localStorage and toggles documentElement.dark.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: ThemeToggle
export function ThemeToggle() {
  const [theme, setTheme] = useState<Theme>('light');

  useEffect(() => {
    const savedTheme = window.localStorage.getItem(storageKey) === 'dark' ? 'dark' : 'light';
    setTheme(savedTheme);
    applyTheme(savedTheme);
  }, []);

  function toggleTheme() {
    setTheme((currentTheme) => {
      const nextTheme = currentTheme === 'dark' ? 'light' : 'dark';
      applyTheme(nextTheme);
      window.localStorage.setItem(storageKey, nextTheme);
      return nextTheme;
    });
  }

  const isDark = theme === 'dark';
  const label = isDark ? 'Switch to light theme' : 'Switch to dark theme';
  const Icon = isDark ? SunIcon : MoonIcon;

  return (
    <Button aria-label={label} onClick={toggleTheme} size="icon" type="button" variant="outline">
      <Icon aria-hidden="true" />
    </Button>
  );
}
