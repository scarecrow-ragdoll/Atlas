'use client';

// FILE: apps/web/src/shared/ui/theme-toggle.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the approved public web light/dark theme toggle.
//   SCOPE: Reads and persists the public web theme preference, toggles the root dark class, and renders an accessible icon button; excludes route data behavior.
//   DEPENDS: react, lucide-react, apps/web/src/shared/ui/primitives/button.tsx.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ThemeToggle - Public web icon button for switching between light and dark themes.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added persisted public web light/dark theme switching.
// END_CHANGE_SUMMARY

import { useEffect, useState } from 'react';
import { MoonIcon, SunIcon } from 'lucide-react';
import { Button } from './primitives/button';

const storageKey = 'web-theme';

type Theme = 'light' | 'dark';

function applyTheme(theme: Theme) {
  document.documentElement.classList.toggle('dark', theme === 'dark');
}

// START_CONTRACT: ThemeToggle
//   PURPOSE: Render a persisted accessible public web theme switch button.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - icon button whose accessible name describes the next theme }
//   SIDE_EFFECTS: Reads/writes localStorage and toggles documentElement.dark.
//   LINKS: M-WEB / V-M-WEB.
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
    <Button
      aria-label={label}
      className="theme-toggle"
      onClick={toggleTheme}
      size="icon"
      type="button"
      variant="outline"
    >
      <Icon aria-hidden="true" />
    </Button>
  );
}
