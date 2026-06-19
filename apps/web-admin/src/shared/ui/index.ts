// FILE: apps/web-admin/src/shared/ui/index.ts
// VERSION: 1.0.2
// START_MODULE_CONTRACT
//   PURPOSE: Expose the approved public UI-kit surface for web-admin pages.
//   SCOPE: Re-exports shadcn primitives and admin layout compositions; excludes implementation-only utility subpaths from page imports.
//   DEPENDS: apps/web-admin/src/shared/ui/primitives, apps/web-admin/src/shared/ui/layout.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: BARREL
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Alert - Approved alert primitive export.
//   AdminAppShell - Approved global sidebar app shell composition export.
//   AdminEmptyState - Approved admin empty-state composition export.
//   AdminBreadcrumbItem - Approved admin shell breadcrumb item type export.
//   AdminIcon - Approved admin shell icon type export.
//   AdminNavigationChild - Approved admin shell child navigation type export.
//   AdminNavigationGroup - Approved admin shell navigation group type export.
//   AdminNavigationItem - Approved admin shell navigation item type export.
//   AdminPageHeader - Approved admin page-header composition export.
//   AdminPageShell - Approved admin page-shell composition export.
//   AdminProjectItem - Approved admin shell reference item type export.
//   AdminSection - Approved admin section composition export.
//   AdminShellHeader - Approved sidebar shell header composition export.
//   AdminTeamItem - Approved admin shell team item type export.
//   AdminToolbar - Approved admin toolbar composition export.
//   AdminUser - Approved admin shell user type export.
//   AlertDescription - Approved alert description primitive export.
//   AlertTitle - Approved alert title primitive export.
//   Avatar - Approved avatar primitive export.
//   AvatarFallback - Approved avatar fallback primitive export.
//   AvatarImage - Approved avatar image primitive export.
//   Badge - Approved badge primitive export.
//   Breadcrumb - Approved breadcrumb nav primitive export.
//   BreadcrumbEllipsis - Approved breadcrumb ellipsis primitive export.
//   BreadcrumbItem - Approved breadcrumb item primitive export.
//   BreadcrumbLink - Approved breadcrumb link primitive export.
//   BreadcrumbList - Approved breadcrumb list primitive export.
//   BreadcrumbPage - Approved breadcrumb page primitive export.
//   BreadcrumbSeparator - Approved breadcrumb separator primitive export.
//   Button - Approved button primitive export.
//   Card - Approved card primitive export.
//   CardAction - Approved card action primitive export.
//   CardContent - Approved card content primitive export.
//   CardDescription - Approved card description primitive export.
//   CardFooter - Approved card footer primitive export.
//   CardHeader - Approved card header primitive export.
//   CardTitle - Approved card title primitive export.
//   Checkbox - Approved checkbox primitive export.
//   Collapsible - Approved collapsible root primitive export.
//   CollapsibleContent - Approved collapsible content primitive export.
//   CollapsibleTrigger - Approved collapsible trigger primitive export.
//   Dialog - Approved dialog root primitive export.
//   DialogClose - Approved dialog close primitive export.
//   DialogContent - Approved dialog content primitive export.
//   DialogDescription - Approved dialog description primitive export.
//   DialogFooter - Approved dialog footer primitive export.
//   DialogHeader - Approved dialog header primitive export.
//   DialogOverlay - Approved dialog overlay primitive export.
//   DialogPortal - Approved dialog portal primitive export.
//   DialogTitle - Approved dialog title primitive export.
//   DialogTrigger - Approved dialog trigger primitive export.
//   DropdownMenu - Approved dropdown menu root primitive export.
//   DropdownMenuCheckboxItem - Approved dropdown menu checkbox item primitive export.
//   DropdownMenuContent - Approved dropdown menu content primitive export.
//   DropdownMenuGroup - Approved dropdown menu group primitive export.
//   DropdownMenuItem - Approved dropdown menu item primitive export.
//   DropdownMenuLabel - Approved dropdown menu label primitive export.
//   DropdownMenuPortal - Approved dropdown menu portal primitive export.
//   DropdownMenuRadioGroup - Approved dropdown menu radio group primitive export.
//   DropdownMenuRadioItem - Approved dropdown menu radio item primitive export.
//   DropdownMenuSeparator - Approved dropdown menu separator primitive export.
//   DropdownMenuShortcut - Approved dropdown menu shortcut primitive export.
//   DropdownMenuSub - Approved dropdown menu submenu root primitive export.
//   DropdownMenuSubContent - Approved dropdown menu submenu content primitive export.
//   DropdownMenuSubTrigger - Approved dropdown menu submenu trigger primitive export.
//   DropdownMenuTrigger - Approved dropdown menu trigger primitive export.
//   Input - Approved input primitive export.
//   Label - Approved label primitive export.
//   Select - Approved select root primitive export.
//   SelectContent - Approved select content primitive export.
//   SelectGroup - Approved select group primitive export.
//   SelectItem - Approved select item primitive export.
//   SelectLabel - Approved select label primitive export.
//   SelectScrollDownButton - Approved select scroll-down button primitive export.
//   SelectScrollUpButton - Approved select scroll-up button primitive export.
//   SelectSeparator - Approved select separator primitive export.
//   SelectTrigger - Approved select trigger primitive export.
//   SelectValue - Approved select value primitive export.
//   Separator - Approved separator primitive export.
//   Sheet - Approved sheet root primitive export.
//   SheetClose - Approved sheet close primitive export.
//   SheetContent - Approved sheet content primitive export.
//   SheetDescription - Approved sheet description primitive export.
//   SheetFooter - Approved sheet footer primitive export.
//   SheetHeader - Approved sheet header primitive export.
//   SheetTitle - Approved sheet title primitive export.
//   SheetTrigger - Approved sheet trigger primitive export.
//   Sidebar - Approved sidebar primitive export.
//   SidebarContent - Approved sidebar content primitive export.
//   SidebarFooter - Approved sidebar footer primitive export.
//   SidebarGroup - Approved sidebar group primitive export.
//   SidebarGroupAction - Approved sidebar group action primitive export.
//   SidebarGroupContent - Approved sidebar group content primitive export.
//   SidebarGroupLabel - Approved sidebar group label primitive export.
//   SidebarHeader - Approved sidebar header primitive export.
//   SidebarInput - Approved sidebar input primitive export.
//   SidebarInset - Approved sidebar inset primitive export.
//   SidebarMenu - Approved sidebar menu primitive export.
//   SidebarMenuAction - Approved sidebar menu action primitive export.
//   SidebarMenuBadge - Approved sidebar menu badge primitive export.
//   SidebarMenuButton - Approved sidebar menu button primitive export.
//   SidebarMenuItem - Approved sidebar menu item primitive export.
//   SidebarMenuSkeleton - Approved sidebar menu skeleton primitive export.
//   SidebarMenuSub - Approved sidebar submenu primitive export.
//   SidebarMenuSubButton - Approved sidebar submenu button primitive export.
//   SidebarMenuSubItem - Approved sidebar submenu item primitive export.
//   SidebarProvider - Approved sidebar provider primitive export.
//   SidebarRail - Approved sidebar rail primitive export.
//   SidebarSeparator - Approved sidebar separator primitive export.
//   SidebarTrigger - Approved sidebar trigger primitive export.
//   Skeleton - Approved skeleton primitive export.
//   Switch - Approved switch primitive export.
//   Table - Approved table primitive export.
//   TableBody - Approved table body primitive export.
//   TableCaption - Approved table caption primitive export.
//   TableCell - Approved table cell primitive export.
//   TableFooter - Approved table footer primitive export.
//   TableHead - Approved table head primitive export.
//   TableHeader - Approved table header primitive export.
//   TableRow - Approved table row primitive export.
//   Tabs - Approved tabs root primitive export.
//   TabsContent - Approved tabs content primitive export.
//   TabsList - Approved tabs list primitive export.
//   TabsTrigger - Approved tabs trigger primitive export.
//   Textarea - Approved textarea primitive export.
//   ThemeToggle - Approved admin theme toggle composition export.
//   Tooltip - Approved tooltip root primitive export.
//   TooltipContent - Approved tooltip content primitive export.
//   TooltipProvider - Approved tooltip provider primitive export.
//   TooltipTrigger - Approved tooltip trigger primitive export.
//   useSidebar - Approved sidebar context hook export.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.3 - Exported sidebar-07 foundation primitives through the curated UI-kit barrel.
// END_CHANGE_SUMMARY

export { Alert, AlertDescription, AlertTitle } from './primitives/alert';
export { AdminAppShell } from './layout/admin-app-shell';
export { AdminEmptyState } from './layout/admin-empty-state';
export { AdminPageHeader } from './layout/admin-page-header';
export { AdminPageShell } from './layout/admin-page-shell';
export { AdminSection } from './layout/admin-section';
export { AdminShellHeader } from './layout/admin-shell-header';
export { AdminToolbar } from './layout/admin-toolbar';
export type {
  AdminBreadcrumbItem,
  AdminIcon,
  AdminNavigationChild,
  AdminNavigationGroup,
  AdminNavigationItem,
  AdminProjectItem,
  AdminTeamItem,
  AdminUser,
} from './layout/admin-shell-types';
export { ThemeToggle } from './layout/theme-toggle';
export { Avatar, AvatarFallback, AvatarImage } from './primitives/avatar';
export { Badge } from './primitives/badge';
export {
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from './primitives/breadcrumb';
export { Button } from './primitives/button';
export {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from './primitives/card';
export { Checkbox } from './primitives/checkbox';
export { Collapsible, CollapsibleContent, CollapsibleTrigger } from './primitives/collapsible';
export {
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
} from './primitives/dialog';
export {
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
} from './primitives/dropdown-menu';
export { Input } from './primitives/input';
export { Label } from './primitives/label';
export {
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
} from './primitives/select';
export { Separator } from './primitives/separator';
export {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from './primitives/sheet';
export {
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
  useSidebar,
} from './primitives/sidebar';
export { Skeleton } from './primitives/skeleton';
export { Switch } from './primitives/switch';
export {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from './primitives/table';
export { Tabs, TabsContent, TabsList, TabsTrigger } from './primitives/tabs';
export { Textarea } from './primitives/textarea';
export { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './primitives/tooltip';
