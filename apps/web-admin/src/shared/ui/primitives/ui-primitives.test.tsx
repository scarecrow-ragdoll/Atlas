// FILE: apps/web-admin/src/shared/ui/primitives/ui-primitives.test.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify the public web-admin UI primitive export surface.
//   SCOPE: Covers rendering and importability for primitives used by admin pages and /ui-kit; excludes visual pixel assertions.
//   DEPENDS: apps/web-admin/src/shared/ui/index.ts, @testing-library/react, vitest.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   primitive exports test - Prove primitives are available through @shared/ui and render basic accessible output.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Added inactive sidebar data-active regression coverage.
// END_CHANGE_SUMMARY

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Avatar,
  AvatarFallback,
  AvatarImage,
  Badge,
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
  Button,
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Checkbox,
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
  DialogTrigger,
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
  Input,
  Label,
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
  Separator,
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInput,
  SidebarInset,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSkeleton,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
  Skeleton,
  Switch,
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
  ThemeToggle,
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
  useSidebar,
} from '@shared/ui';

const publicPrimitiveExports = {
  Alert,
  AlertDescription,
  AlertTitle,
  Avatar,
  AvatarFallback,
  AvatarImage,
  Badge,
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
  Button,
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Checkbox,
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
  DialogTrigger,
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
  Input,
  Label,
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
  Separator,
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInput,
  SidebarInset,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSkeleton,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
  Skeleton,
  Switch,
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
  ThemeToggle,
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
  useSidebar,
};

const originalResizeObserver = globalThis.ResizeObserver;
const originalScrollIntoView = Element.prototype.scrollIntoView;
const originalMatchMedia = window.matchMedia;
const originalInnerWidth = window.innerWidth;

type PrimitiveMediaListener = () => void;

class TestResizeObserver {
  observe() {
    return undefined;
  }

  unobserve() {
    return undefined;
  }

  disconnect() {
    return undefined;
  }
}

function setPrimitiveViewportWidth(width: number) {
  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: width,
  });
}

function installPrimitiveMatchMedia() {
  const listeners = new Set<PrimitiveMediaListener>();
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query.includes('767') ? window.innerWidth < 768 : false,
    media: query,
    onchange: null,
    addEventListener: (_event: string, listener: PrimitiveMediaListener) => listeners.add(listener),
    removeEventListener: (_event: string, listener: PrimitiveMediaListener) =>
      listeners.delete(listener),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
  return listeners;
}

beforeAll(() => {
  globalThis.ResizeObserver = TestResizeObserver as typeof ResizeObserver;
  Element.prototype.scrollIntoView = vi.fn();
  setPrimitiveViewportWidth(1024);
  installPrimitiveMatchMedia();
});

afterEach(() => {
  cleanup();
  document.body.removeAttribute('data-scroll-locked');
  document.body.style.removeProperty('pointer-events');
  document.querySelectorAll('[data-radix-focus-guard]').forEach((element) => element.remove());
  setPrimitiveViewportWidth(1024);
  installPrimitiveMatchMedia();
  document.cookie = 'sidebar_state=; path=/; max-age=0';
});

afterAll(() => {
  if (originalResizeObserver) {
    globalThis.ResizeObserver = originalResizeObserver;
  } else {
    delete (globalThis as Partial<typeof globalThis>).ResizeObserver;
  }

  if (originalScrollIntoView) {
    Element.prototype.scrollIntoView = originalScrollIntoView;
  } else {
    delete (Element.prototype as Partial<Element>).scrollIntoView;
  }

  if (originalMatchMedia) {
    window.matchMedia = originalMatchMedia;
  } else {
    delete (window as Partial<Window>).matchMedia;
  }

  setPrimitiveViewportWidth(originalInnerWidth);
});

function SidebarContextProbe() {
  useSidebar();
  return <span>Inside sidebar context</span>;
}

function SidebarDirectSetOpenProbe() {
  const { setOpen } = useSidebar();
  return <button onClick={() => setOpen(false)}>Close sidebar directly</button>;
}

function SidebarOutsideProviderProbe() {
  useSidebar();
  return null;
}

describe('web-admin UI primitive exports', () => {
  it('exports every approved primitive through the public @shared/ui barrel', () => {
    expect(Object.keys(publicPrimitiveExports)).toEqual([
      'Alert',
      'AlertDescription',
      'AlertTitle',
      'Avatar',
      'AvatarFallback',
      'AvatarImage',
      'Badge',
      'Breadcrumb',
      'BreadcrumbEllipsis',
      'BreadcrumbItem',
      'BreadcrumbLink',
      'BreadcrumbList',
      'BreadcrumbPage',
      'BreadcrumbSeparator',
      'Button',
      'Card',
      'CardAction',
      'CardContent',
      'CardDescription',
      'CardFooter',
      'CardHeader',
      'CardTitle',
      'Checkbox',
      'Collapsible',
      'CollapsibleContent',
      'CollapsibleTrigger',
      'Dialog',
      'DialogClose',
      'DialogContent',
      'DialogDescription',
      'DialogFooter',
      'DialogHeader',
      'DialogOverlay',
      'DialogPortal',
      'DialogTitle',
      'DialogTrigger',
      'DropdownMenu',
      'DropdownMenuCheckboxItem',
      'DropdownMenuContent',
      'DropdownMenuGroup',
      'DropdownMenuItem',
      'DropdownMenuLabel',
      'DropdownMenuPortal',
      'DropdownMenuRadioGroup',
      'DropdownMenuRadioItem',
      'DropdownMenuSeparator',
      'DropdownMenuShortcut',
      'DropdownMenuSub',
      'DropdownMenuSubContent',
      'DropdownMenuSubTrigger',
      'DropdownMenuTrigger',
      'Input',
      'Label',
      'Select',
      'SelectContent',
      'SelectGroup',
      'SelectItem',
      'SelectLabel',
      'SelectScrollDownButton',
      'SelectScrollUpButton',
      'SelectSeparator',
      'SelectTrigger',
      'SelectValue',
      'Separator',
      'Sheet',
      'SheetClose',
      'SheetContent',
      'SheetDescription',
      'SheetFooter',
      'SheetHeader',
      'SheetTitle',
      'SheetTrigger',
      'Sidebar',
      'SidebarContent',
      'SidebarFooter',
      'SidebarGroup',
      'SidebarGroupAction',
      'SidebarGroupContent',
      'SidebarGroupLabel',
      'SidebarHeader',
      'SidebarInput',
      'SidebarInset',
      'SidebarMenu',
      'SidebarMenuAction',
      'SidebarMenuBadge',
      'SidebarMenuButton',
      'SidebarMenuItem',
      'SidebarMenuSkeleton',
      'SidebarMenuSub',
      'SidebarMenuSubButton',
      'SidebarMenuSubItem',
      'SidebarProvider',
      'SidebarRail',
      'SidebarSeparator',
      'SidebarTrigger',
      'Skeleton',
      'Switch',
      'Table',
      'TableBody',
      'TableCaption',
      'TableCell',
      'TableFooter',
      'TableHead',
      'TableHeader',
      'TableRow',
      'Tabs',
      'TabsContent',
      'TabsList',
      'TabsTrigger',
      'Textarea',
      'ThemeToggle',
      'Tooltip',
      'TooltipContent',
      'TooltipProvider',
      'TooltipTrigger',
      'useSidebar',
    ]);

    expect(Object.values(publicPrimitiveExports).every(Boolean)).toBe(true);
  });

  it('renders the primitive set through the public @shared/ui barrel', () => {
    const triggerClick = vi.fn();

    render(
      <TooltipProvider>
        <Alert>
          <AlertTitle>Saved</AlertTitle>
          <AlertDescription>Changes are ready.</AlertDescription>
        </Alert>
        <Avatar>
          <AvatarImage alt="Developer avatar" src="/missing-avatar.png" />
          <AvatarFallback>DV</AvatarFallback>
        </Avatar>
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink href="/">Overview</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>Users</BreadcrumbPage>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <BreadcrumbEllipsis />
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <Badge>Active</Badge>
        <Badge variant="secondary">Secondary</Badge>
        <Badge asChild variant="outline">
          <a href="/status">Status link</a>
        </Badge>
        <Button>Save</Button>
        <Card>
          <CardHeader>
            <CardTitle>Card title</CardTitle>
            <CardAction>Card action</CardAction>
            <CardDescription>Card description</CardDescription>
          </CardHeader>
          <CardContent>Card body</CardContent>
          <CardFooter>Card footer</CardFooter>
        </Card>
        <Label htmlFor="name">Name</Label>
        <Input id="name" defaultValue="Ada" />
        <Textarea aria-label="Notes" defaultValue="Reference notes" />
        <Checkbox aria-label="Enabled" defaultChecked />
        <Collapsible open>
          <CollapsibleTrigger>Toggle section</CollapsibleTrigger>
          <CollapsibleContent>Expanded section</CollapsibleContent>
        </Collapsible>
        <Switch aria-label="Published" defaultChecked />
        <Separator />
        <SidebarProvider defaultOpen>
          <Sidebar>
            <SidebarHeader>Shell header</SidebarHeader>
            <SidebarContent className="custom-sidebar-content">
              <SidebarGroup>
                <SidebarGroupLabel asChild>
                  <a href="/platform">Platform</a>
                </SidebarGroupLabel>
                <SidebarGroupContent className="custom-sidebar-group-content">
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <SidebarMenuButton isActive tooltip="Overview">
                        <span>Overview</span>
                      </SidebarMenuButton>
                      <SidebarMenuBadge>New</SidebarMenuBadge>
                      <SidebarMenuAction asChild>
                        <a href="/overview-actions">More</a>
                      </SidebarMenuAction>
                    </SidebarMenuItem>
                  </SidebarMenu>
                </SidebarGroupContent>
                <SidebarGroupAction asChild>
                  <a href="/group-action">Group action</a>
                </SidebarGroupAction>
              </SidebarGroup>
              <SidebarMenuSkeleton />
              <SidebarMenuSub>
                <SidebarMenuSubItem>
                  <SidebarMenuSubButton>Child item</SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
              <SidebarInput aria-label="Filter navigation" />
              <SidebarSeparator />
            </SidebarContent>
            <SidebarFooter>Shell footer</SidebarFooter>
            <SidebarRail />
          </Sidebar>
          <SidebarInset>
            <SidebarDirectSetOpenProbe />
            <SidebarTrigger aria-label="Toggle sidebar" onClick={triggerClick} />
            <p>Sidebar content slot</p>
          </SidebarInset>
        </SidebarProvider>
        <Skeleton data-testid="loading-row" />
        <Tabs defaultValue="overview">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
          </TabsList>
          <TabsContent value="overview">Overview content</TabsContent>
        </Tabs>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>ada@example.com</TableCell>
            </TableRow>
          </TableBody>
          <TableFooter>
            <TableRow>
              <TableCell>Total users</TableCell>
            </TableRow>
          </TableFooter>
          <TableCaption>Users table caption</TableCaption>
        </Table>
      </TooltipProvider>,
    );

    expect(screen.getByText('DV')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Overview' })).toHaveAttribute('href', '/');
    expect(screen.getByText('Users')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByText('Active')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Status link' })).toHaveAttribute('href', '/status');
    expect(screen.getByText('Card action')).toBeInTheDocument();
    expect(screen.getByText('Card footer')).toBeInTheDocument();
    expect(screen.getByLabelText('Name')).toHaveValue('Ada');
    expect(screen.getByText('Expanded section')).toBeInTheDocument();
    expect(screen.getByText('Platform')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'More' })).toHaveAttribute('href', '/overview-actions');
    expect(screen.getByRole('link', { name: 'Group action' })).toHaveAttribute(
      'href',
      '/group-action',
    );
    const sidebarWrapper = document.querySelector('[data-slot="sidebar-wrapper"]');
    expect(sidebarWrapper).toHaveAttribute('data-state', 'expanded');
    fireEvent.click(screen.getByRole('button', { name: 'Close sidebar directly' }));
    expect(sidebarWrapper).toHaveAttribute('data-state', 'collapsed');
    fireEvent.click(screen.getByRole('button', { name: 'Toggle sidebar' }));
    expect(triggerClick).toHaveBeenCalledTimes(1);
    expect(sidebarWrapper).toHaveAttribute('data-state', 'expanded');
    expect(screen.getByText('Sidebar content slot')).toBeInTheDocument();
    expect(screen.getByText('Overview content')).toBeInTheDocument();
    expect(screen.getByText('ada@example.com')).toBeInTheDocument();
    expect(screen.getByText('Total users')).toBeInTheDocument();
    expect(screen.getByText('Users table caption')).toBeInTheDocument();
  });

  it('renders sheet wrappers in an open mobile-navigation state', () => {
    render(
      <Sheet open>
        <SheetTrigger>Open sheet</SheetTrigger>
        <SheetContent aria-describedby={undefined}>
          <SheetHeader>
            <SheetTitle>Mobile navigation</SheetTitle>
            <SheetDescription>Sidebar sheet content.</SheetDescription>
          </SheetHeader>
          <SheetFooter>
            <SheetClose>Close sheet</SheetClose>
          </SheetFooter>
        </SheetContent>
      </Sheet>,
    );

    expect(screen.getByRole('dialog', { name: 'Mobile navigation' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Close sheet' })).toBeInTheDocument();
  });

  it('supports controlled sidebar provider state and keyboard toggling', () => {
    const onOpenChange = vi.fn();
    render(
      <SidebarProvider open onOpenChange={onOpenChange}>
        <Sidebar>
          <SidebarContent>Controlled sidebar</SidebarContent>
        </Sidebar>
        <SidebarInset>
          <SidebarTrigger aria-label="Toggle controlled sidebar" />
        </SidebarInset>
      </SidebarProvider>,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Toggle controlled sidebar' }));
    expect(onOpenChange).toHaveBeenCalledWith(false);

    fireEvent.keyDown(window, { key: 'b', ctrlKey: true });
    expect(onOpenChange).toHaveBeenCalledWith(false);
  });

  it('opens navigation through the sidebar mobile sheet branch', async () => {
    setPrimitiveViewportWidth(500);
    installPrimitiveMatchMedia();

    render(
      <TooltipProvider>
        <SidebarProvider>
          <Sidebar aria-label="Mobile admin navigation" role="navigation">
            <SidebarHeader>Mobile shell</SidebarHeader>
            <SidebarContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild tooltip="Mobile Users">
                    <a href="/users">Mobile Users</a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarContent>
          </Sidebar>
          <SidebarInset>
            <SidebarTrigger aria-label="Open mobile sidebar" />
            <span>Mobile route content</span>
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>,
    );

    expect(screen.getByText('Mobile route content')).toBeInTheDocument();

    await waitFor(() =>
      expect(
        screen.queryByRole('navigation', { name: 'Mobile admin navigation' }),
      ).not.toBeInTheDocument(),
    );

    fireEvent.click(screen.getByRole('button', { name: 'Open mobile sidebar' }));

    await expect(screen.findByRole('dialog', { name: 'Sidebar' })).resolves.toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Mobile Users' })).toBeVisible();
  });

  it('covers sidebar context, fixed sidebar, and desktop variant helpers', async () => {
    setPrimitiveViewportWidth(1280);
    installPrimitiveMatchMedia();

    render(
      <TooltipProvider>
        <SidebarProvider defaultOpen={false}>
          <SidebarContextProbe />
          <Sidebar
            aria-label="Variant navigation"
            collapsible="icon"
            role="navigation"
            side="right"
            variant="floating"
          >
            <SidebarHeader>Variant header</SidebarHeader>
            <SidebarSeparator />
            <SidebarContent>
              <SidebarInput aria-label="Search navigation" />
              <SidebarGroup>
                <SidebarGroupLabel>Variant group</SidebarGroupLabel>
                <SidebarGroupAction aria-label="Add item">+</SidebarGroupAction>
                <SidebarGroupContent>
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <SidebarMenuButton isActive size="lg" tooltip="Reports" variant="outline">
                        Reports
                      </SidebarMenuButton>
                      <SidebarMenuAction aria-label="Reports actions" showOnHover />
                      <SidebarMenuBadge>3</SidebarMenuBadge>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuSkeleton showIcon />
                    </SidebarMenuItem>
                    <SidebarMenuSub>
                      <SidebarMenuSubItem>
                        <SidebarMenuSubButton isActive size="sm">
                          Child report
                        </SidebarMenuSubButton>
                      </SidebarMenuSubItem>
                    </SidebarMenuSub>
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            </SidebarContent>
            <SidebarFooter>Variant footer</SidebarFooter>
            <SidebarRail />
          </Sidebar>
          <SidebarInset>Variant content</SidebarInset>
        </SidebarProvider>
      </TooltipProvider>,
    );

    expect(screen.getByText('Inside sidebar context')).toBeInTheDocument();
    expect(screen.getByRole('navigation', { name: 'Variant navigation' })).toHaveAttribute(
      'data-side',
      'right',
    );
    expect(screen.getByRole('button', { name: 'Reports' })).toHaveAttribute('data-active', 'true');
    expect(screen.getByRole('button', { name: 'Add item' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Reports actions' })).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('Child report')).toHaveAttribute('data-active', 'true');
    expect(screen.getByLabelText('Search navigation')).toBeInTheDocument();

    const reportsButton = screen.getByRole('button', { name: 'Reports' });
    fireEvent.pointerMove(reportsButton, { pointerType: 'mouse' });
    fireEvent.pointerEnter(reportsButton, { pointerType: 'mouse' });
    fireEvent.mouseEnter(reportsButton);
    await waitFor(() => expect(screen.getByRole('tooltip')).toHaveTextContent('Reports'));
  });

  it('omits sidebar active-state data attributes for inactive menu items', () => {
    setPrimitiveViewportWidth(1280);
    installPrimitiveMatchMedia();

    render(
      <SidebarProvider defaultOpen>
        <Sidebar aria-label="Active state navigation" role="navigation">
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton>Overview</SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton isActive>Users</SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuSub>
                <SidebarMenuSubItem>
                  <SidebarMenuSubButton>Directory</SidebarMenuSubButton>
                </SidebarMenuSubItem>
                <SidebarMenuSubItem>
                  <SidebarMenuSubButton isActive>User detail</SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );

    expect(screen.getByRole('button', { name: 'Overview' })).not.toHaveAttribute('data-active');
    expect(screen.getByRole('button', { name: 'Users' })).toHaveAttribute('data-active', 'true');
    expect(screen.getByText('Directory')).not.toHaveAttribute('data-active');
    expect(screen.getByText('User detail')).toHaveAttribute('data-active', 'true');
  });

  it('covers non-collapsible sidebar rendering and provider misuse errors', () => {
    setPrimitiveViewportWidth(1280);
    installPrimitiveMatchMedia();

    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    expect(() => render(<SidebarOutsideProviderProbe />)).toThrow(
      'useSidebar must be used within a SidebarProvider.',
    );
    consoleError.mockRestore();

    render(
      <SidebarProvider>
        <Sidebar aria-label="Pinned navigation" collapsible="none" role="navigation">
          <SidebarContent>Pinned sidebar</SidebarContent>
        </Sidebar>
        <SidebarInset>Pinned content</SidebarInset>
      </SidebarProvider>,
    );

    expect(screen.getByRole('navigation', { name: 'Pinned navigation' })).toHaveTextContent(
      'Pinned sidebar',
    );
    expect(screen.getByText('Pinned content')).toBeInTheDocument();
  });

  it('renders dialog wrappers in open and close-button-free states', () => {
    render(
      <Dialog open>
        <DialogTrigger>Open dialog</DialogTrigger>
        <DialogContent aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Dialog title</DialogTitle>
            <DialogDescription>Dialog description</DialogDescription>
          </DialogHeader>
          <DialogFooter showCloseButton>
            <Button>Confirm dialog</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>,
    );

    expect(screen.getByRole('dialog', { name: 'Dialog title' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Confirm dialog' })).toBeInTheDocument();
  });

  it('renders dialog wrappers without generated close controls', () => {
    render(
      <Dialog open>
        <DialogContent aria-describedby={undefined} showCloseButton={false}>
          <DialogTitle>Dialog without close</DialogTitle>
          <DialogDescription>Dialog body without close button.</DialogDescription>
          <DialogFooter>
            <DialogClose>Manual close</DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>,
    );

    expect(screen.getByRole('dialog', { name: 'Dialog without close' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Manual close' })).toBeInTheDocument();
  });

  it('renders dropdown menu wrappers and submenu content', () => {
    render(
      <DropdownMenu open>
        <DropdownMenuTrigger>Open menu</DropdownMenuTrigger>
        <DropdownMenuPortal>
          <div data-testid="dropdown-portal">Dropdown portal</div>
        </DropdownMenuPortal>
        <DropdownMenuContent>
          <DropdownMenuGroup>
            <DropdownMenuLabel inset>Menu group</DropdownMenuLabel>
            <DropdownMenuItem inset>Default item</DropdownMenuItem>
            <DropdownMenuItem variant="destructive">Delete item</DropdownMenuItem>
            <DropdownMenuCheckboxItem checked>Checked item</DropdownMenuCheckboxItem>
            <DropdownMenuRadioGroup value="one">
              <DropdownMenuRadioItem value="one">Radio item</DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
            <DropdownMenuSeparator />
            <DropdownMenuShortcut>Cmd K</DropdownMenuShortcut>
            <DropdownMenuSub open>
              <DropdownMenuSubTrigger inset>More actions</DropdownMenuSubTrigger>
              <DropdownMenuSubContent>Nested action</DropdownMenuSubContent>
            </DropdownMenuSub>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>,
    );

    expect(screen.getByTestId('dropdown-portal')).toHaveTextContent('Dropdown portal');
    expect(screen.getByText('Default item')).toBeInTheDocument();
    expect(screen.getByText('Delete item')).toBeInTheDocument();
    expect(screen.getByText('Nested action')).toBeInTheDocument();
  });

  it('renders select wrappers with default and popper content', () => {
    render(
      <div>
        <Select open value="admin">
          <SelectTrigger aria-label="Role">
            <SelectValue placeholder="Select role" />
          </SelectTrigger>
          <SelectContent position="popper">
            <SelectGroup>
              <SelectLabel>Roles</SelectLabel>
              <SelectItem value="admin">Admin</SelectItem>
              <SelectSeparator />
              <SelectItem value="viewer">Viewer</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select open value="default">
          <SelectTrigger aria-label="Default role" size="sm">
            <SelectValue placeholder="Select default role" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="default">Default</SelectItem>
          </SelectContent>
        </Select>
      </div>,
    );

    expect(screen.getByText('Roles')).toBeInTheDocument();
    expect(screen.getAllByText('Default')).toHaveLength(2);
  });

  it('renders tooltip wrappers through the provider', () => {
    render(
      <TooltipProvider>
        <Tooltip open>
          <TooltipTrigger>Hover target</TooltipTrigger>
          <TooltipContent>Tooltip copy</TooltipContent>
        </Tooltip>
      </TooltipProvider>,
    );

    expect(screen.getByRole('tooltip')).toHaveTextContent('Tooltip copy');
  });
});
