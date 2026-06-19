// FILE: apps/web-admin/src/shared/ui/layout/admin-shell.test.tsx
// VERSION: 1.1.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the adapted sidebar-07 web-admin app shell.
//   SCOPE: Covers shell rendering, active navigation, disabled placeholders, breadcrumbs, theme placement, sidebar trigger, authenticated user menu logout, and content slot; excludes page data behavior.
//   DEPENDS: apps/web-admin/src/shared/ui/layout/admin-app-shell.tsx, @testing-library/react, lucide-react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   admin shell tests - Prove sidebar-07 shell composition renders template-native chrome and route content.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.1.0 - Added authenticated admin logout menu coverage.
// END_CHANGE_SUMMARY

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { HomeIcon, SettingsIcon, UsersIcon } from 'lucide-react';
import { MemoryRouter } from 'react-router';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import { AdminAppShell } from './admin-app-shell';
import type {
  AdminBreadcrumbItem,
  AdminNavigationGroup,
  AdminProjectItem,
  AdminTeamItem,
  AdminUser,
} from './admin-shell-types';

const navigation: AdminNavigationGroup[] = [
  {
    id: 'platform',
    label: 'Platform',
    items: [
      { id: 'overview', label: 'Overview', href: '/', icon: HomeIcon },
      { id: 'users', label: 'Users', href: '/users', icon: UsersIcon, isActive: true },
      {
        id: 'reports',
        label: 'Reports',
        href: '/reports',
        icon: SettingsIcon,
        isActive: true,
        children: [
          { id: 'activity-report', label: 'Activity report', href: '/reports/activity' },
          {
            id: 'draft-report',
            label: 'Draft report',
            href: '#draft-report',
            disabled: true,
          },
        ],
      },
      {
        id: 'disabled-admin',
        label: 'Disabled admin',
        href: '#disabled-admin',
        icon: SettingsIcon,
        disabled: true,
      },
    ],
  },
];

const referenceItems: AdminProjectItem[] = [
  {
    id: 'live-reference',
    name: 'Live reference',
    href: '#live-reference',
    icon: HomeIcon,
  },
  {
    id: 'system-settings',
    name: 'System/Settings',
    href: '#system-settings',
    icon: SettingsIcon,
    disabled: true,
  },
];

const originalMatchMedia = window.matchMedia;
const originalInnerWidth = window.innerWidth;

function installShellMatchMedia(isMobile = false) {
  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: isMobile ? 500 : 1024,
  });
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query.includes('767') ? isMobile : false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
}

beforeAll(() => {
  installShellMatchMedia();
});

type RenderShellOptions = {
  breadcrumbs?: AdminBreadcrumbItem[];
  isLogoutPending?: boolean;
  navigation?: AdminNavigationGroup[];
  onLogout?: () => void | Promise<void>;
  pathname?: string;
  referenceItems?: AdminProjectItem[];
  teams?: AdminTeamItem[];
  user?: AdminUser;
};

function renderShell({
  breadcrumbs = [{ label: 'Users', href: '/users' }, { label: 'User detail' }],
  isLogoutPending,
  navigation: shellNavigation = navigation,
  onLogout,
  pathname = '/users/1',
  referenceItems: shellReferenceItems = referenceItems,
  teams = [
    { id: 'template', name: 'Monorepo Template', plan: 'Admin shell', icon: HomeIcon },
    { id: 'secondary', name: 'Secondary Workspace', plan: 'Branch preview', icon: SettingsIcon },
  ],
  user = {
    name: 'Developer',
    email: 'developer@example.local',
    initials: 'DV',
    avatarUrl: '/avatar.png',
  },
}: RenderShellOptions = {}) {
  return render(
    <MemoryRouter>
      <AdminAppShell
        breadcrumbs={breadcrumbs}
        isLogoutPending={isLogoutPending}
        navigation={shellNavigation}
        onLogout={onLogout}
        pathname={pathname}
        referenceItems={shellReferenceItems}
        teams={teams}
        user={user}
      >
        <section>Route content</section>
      </AdminAppShell>
    </MemoryRouter>,
  );
}

afterEach(() => {
  cleanup();
  document.documentElement.classList.remove('dark');
  document.cookie = 'sidebar_state=; path=/; max-age=0';
  window.localStorage.clear();
  installShellMatchMedia();
});

afterAll(() => {
  if (originalMatchMedia) {
    window.matchMedia = originalMatchMedia;
  } else {
    delete (window as Partial<Window>).matchMedia;
  }

  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: originalInnerWidth,
  });
});

describe('AdminAppShell', () => {
  it('renders sidebar navigation, breadcrumbs, global controls, and content', () => {
    renderShell();

    expect(screen.getByRole('main')).toHaveTextContent('Route content');
    expect(screen.getByRole('navigation', { name: 'Admin navigation' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Users' })).toHaveAttribute('href', '/users');
    expect(screen.getByRole('link', { name: 'Users breadcrumb' })).toHaveAttribute(
      'href',
      '/users',
    );
    expect(screen.getByText('User detail')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Live reference' })).toHaveAttribute(
      'href',
      '/#live-reference',
    );
    expect(screen.getByText('System/Settings')).toBeInTheDocument();
    expect(screen.getByText('Coming soon')).toBeInTheDocument();
    expect(screen.getByText('Activity report')).toBeInTheDocument();
    expect(screen.getByText('Draft report')).toBeInTheDocument();
    expect(screen.getByText('Disabled admin')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Switch to dark theme' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Toggle sidebar' })).toBeInTheDocument();
    expect(screen.getByText('Monorepo Template')).toBeInTheDocument();
    expect(screen.getByText('Developer')).toBeInTheDocument();
  });

  it('updates local team menu state and renders the user menu', async () => {
    renderShell({ onLogout: vi.fn() });

    fireEvent.pointerDown(screen.getByRole('button', { name: /Monorepo Template Admin shell/i }), {
      button: 0,
      ctrlKey: false,
      pointerType: 'mouse',
    });
    fireEvent.keyDown(screen.getByRole('button', { name: /Monorepo Template Admin shell/i }), {
      key: 'Enter',
    });
    fireEvent.click(await screen.findByRole('menuitem', { name: /Secondary Workspace/ }));
    await waitFor(() =>
      expect(
        screen.getByRole('button', { name: /Secondary Workspace Branch preview/i }),
      ).toBeInTheDocument(),
    );

    fireEvent.pointerDown(
      screen.getByRole('button', { name: /Developer developer@example\.local/i }),
      {
        button: 0,
        ctrlKey: false,
        pointerType: 'mouse',
      },
    );
    fireEvent.keyDown(screen.getByRole('button', { name: /Developer developer@example\.local/i }), {
      key: 'Enter',
    });
    expect(await screen.findByRole('menuitem', { name: 'Logout' })).toBeInTheDocument();
    expect(screen.queryByText('Account placeholder')).not.toBeInTheDocument();
    expect(screen.queryByText('Settings placeholder')).not.toBeInTheDocument();
  });

  it('renders the authenticated admin and calls logout from the user menu', async () => {
    const onLogout = vi.fn();

    renderShell({
      breadcrumbs: [{ label: 'Overview' }],
      navigation: [],
      onLogout,
      pathname: '/',
      referenceItems: [],
      teams: [],
      user: { name: 'Owner Admin', email: 'owner@example.test', initials: 'OA' },
    });

    const userButton = screen.getByRole('button', { name: /Owner Admin owner@example\.test/i });
    fireEvent.pointerDown(userButton, { button: 0, ctrlKey: false, pointerType: 'mouse' });
    fireEvent.keyDown(userButton, { key: 'Enter' });
    fireEvent.click(await screen.findByRole('menuitem', { name: 'Logout' }));

    expect(onLogout).toHaveBeenCalledTimes(1);
    expect(screen.getByText('owner@example.test')).toBeInTheDocument();
  });

  it('disables logout while logout is pending', async () => {
    renderShell({
      breadcrumbs: [{ label: 'Overview' }],
      isLogoutPending: true,
      navigation: [],
      onLogout: vi.fn(),
      pathname: '/',
      referenceItems: [],
      teams: [],
      user: { name: 'Owner Admin', email: 'owner@example.test', initials: 'OA' },
    });

    const userButton = screen.getByRole('button', { name: /Owner Admin owner@example\.test/i });
    fireEvent.pointerDown(userButton, { button: 0, ctrlKey: false, pointerType: 'mouse' });
    fireEvent.keyDown(userButton, { key: 'Enter' });

    expect(await screen.findByRole('menuitem', { name: 'Logging out...' })).toHaveAttribute(
      'data-disabled',
    );
  });

  it('renders content when no teams are available', () => {
    renderShell({ teams: [] });

    expect(screen.getByRole('main')).toHaveTextContent('Route content');
    expect(screen.queryByText('Monorepo Template')).not.toBeInTheDocument();
  });

  it('closes the mobile sidebar sheet after route navigation', async () => {
    installShellMatchMedia(true);
    renderShell();

    fireEvent.click(screen.getByRole('button', { name: 'Toggle sidebar' }));
    const mobileSidebar = await screen.findByRole('dialog', { name: 'Sidebar' });
    fireEvent.click(within(mobileSidebar).getByRole('link', { name: 'Users' }));

    await waitFor(() =>
      expect(screen.queryByRole('dialog', { name: 'Sidebar' })).not.toBeInTheDocument(),
    );
  });
});
