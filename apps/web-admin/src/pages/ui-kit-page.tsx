// FILE: apps/web-admin/src/pages/ui-kit-page.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Render the web-admin UI-kit reference route.
//   SCOPE: Demonstrates approved primitives and admin compositions using local static data; excludes API calls and product-specific workflows.
//   DEPENDS: react-router, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Visible admin UI-kit showcase route.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added shell foundation primitives to the reference page.
// END_CHANGE_SUMMARY

import { Link } from 'react-router';
import {
  AdminEmptyState,
  AdminPageHeader,
  AdminPageShell,
  AdminSection,
  AdminToolbar,
  Alert,
  AlertDescription,
  AlertTitle,
  Avatar,
  AvatarFallback,
  Badge,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Checkbox,
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Separator,
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
  SidebarProvider,
  SidebarTrigger,
  Skeleton,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@shared/ui';

const demoRows = [
  { id: 'usr_ada', name: 'Ada Lovelace', email: 'ada@example.com', status: 'Active' },
  { id: 'usr_grace', name: 'Grace Hopper', email: 'grace@example.com', status: 'Pending' },
];

// START_CONTRACT: UiKitPage
//   PURPOSE: Render a local-only reference of approved admin UI components and compositions.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - visible /ui-kit route using static demo data }
//   SIDE_EFFECTS: none.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: UiKitPage
export default function UiKitPage() {
  return (
    <TooltipProvider>
      <AdminPageShell>
        <AdminPageHeader
          title="UI Kit"
          description="Approved web-admin components and page compositions."
          actions={
            <Button asChild variant="outline">
              <Link to="/">Home</Link>
            </Button>
          }
        />

        <AdminSection title="Foundation" description="Theme tokens used by admin screens.">
          <div className="grid gap-4 lg:grid-cols-2">
            <div className="space-y-3">
              <h3 className="text-base font-semibold">Typography scale</h3>
              <div className="space-y-1">
                <p className="text-2xl font-semibold">Page title</p>
                <p className="text-sm text-muted-foreground">Muted helper text</p>
                <p className="text-xs uppercase tracking-normal text-muted-foreground">
                  Section label
                </p>
              </div>
            </div>
            <div className="space-y-3">
              <h3 className="text-base font-semibold">Spacing examples</h3>
              <div className="flex items-center gap-2">
                <span className="h-4 w-4 rounded-sm bg-primary" />
                <span className="h-6 w-6 rounded-sm bg-primary" />
                <span className="h-8 w-8 rounded-sm bg-primary" />
              </div>
              <h3 className="text-base font-semibold">Radius examples</h3>
              <div className="flex items-center gap-2">
                <span className="h-10 w-16 rounded-sm border bg-card" />
                <span className="h-10 w-16 rounded-md border bg-card" />
                <span className="h-10 w-16 rounded-lg border bg-card" />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-3 lg:col-span-2">
              {[
                'bg-background text-foreground',
                'bg-primary text-primary-foreground',
                'bg-muted text-muted-foreground',
              ].map((className) => (
                <div className={`rounded-md border p-4 text-sm ${className}`} key={className}>
                  {className}
                </div>
              ))}
            </div>
          </div>
        </AdminSection>

        <AdminSection title="Actions" description="Approved command surfaces.">
          <div className="flex flex-wrap gap-2">
            <Button>Primary</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="destructive">Destructive</Button>
            <Button disabled>Disabled</Button>
            <Button aria-busy="true" disabled>
              Creating...
            </Button>
          </div>
        </AdminSection>

        <AdminSection title="Forms" description="Form primitives for admin CRUD flows.">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="kit-name">Name</Label>
              <Input id="kit-name" defaultValue="Ada Lovelace" />
            </div>
            <div className="space-y-2">
              <Label>Status</Label>
              <Select defaultValue="active">
                <SelectTrigger aria-label="Status">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="kit-notes">Notes</Label>
              <Textarea id="kit-notes" defaultValue="Reference form copy." />
            </div>
            <label className="flex items-center gap-2 text-sm">
              <Checkbox defaultChecked />
              Receive status updates
            </label>
            <label className="flex items-center gap-2 text-sm">
              <Switch defaultChecked />
              Published
            </label>
          </div>
        </AdminSection>

        <AdminSection title="Feedback" description="Status, loading, and empty states.">
          <div className="grid gap-4 lg:grid-cols-2">
            <Alert>
              <AlertTitle>Validation message</AlertTitle>
              <AlertDescription>Email is already used by another user.</AlertDescription>
            </Alert>
            <Card>
              <CardHeader>
                <CardTitle>Loading skeleton</CardTitle>
                <CardDescription>Use for pending list or detail content.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                <Skeleton className="h-4 w-2/3" />
                <Skeleton className="h-4 w-1/2" />
              </CardContent>
            </Card>
            <AdminEmptyState
              title="No records"
              description="Create a record to populate this table."
            />
            <div className="flex items-center gap-2">
              <Badge>Active</Badge>
              <Badge variant="secondary">Draft</Badge>
              <Badge variant="destructive">Blocked</Badge>
            </div>
          </div>
        </AdminSection>

        <AdminSection title="Data" description="Table pattern for admin collections.">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {demoRows.map((row) => (
                <TableRow key={row.id}>
                  <TableCell>{row.name}</TableCell>
                  <TableCell>{row.email}</TableCell>
                  <TableCell>
                    <Badge variant={row.status === 'Active' ? 'default' : 'secondary'}>
                      {row.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button size="sm" variant="outline">
                          Actions
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>Open</DropdownMenuItem>
                        <DropdownMenuItem>Archive</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </AdminSection>

        <AdminSection
          title="Overlays And Navigation"
          description="Use overlays for bounded secondary tasks."
        >
          <Tabs defaultValue="dialog">
            <TabsList>
              <TabsTrigger value="dialog">Dialog</TabsTrigger>
              <TabsTrigger value="tooltip">Tooltip</TabsTrigger>
            </TabsList>
            <TabsContent value="dialog" className="pt-4">
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline">Open dialog</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm action</DialogTitle>
                    <DialogDescription>This is the approved modal composition.</DialogDescription>
                  </DialogHeader>
                </DialogContent>
              </Dialog>
            </TabsContent>
            <TabsContent value="tooltip" className="pt-4">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline">Hover target</Button>
                </TooltipTrigger>
                <TooltipContent>Tooltip content</TooltipContent>
              </Tooltip>
            </TabsContent>
          </Tabs>
          <Separator className="my-4" />
          <p className="text-sm text-muted-foreground">Separators divide related admin panels.</p>
        </AdminSection>

        <AdminSection
          title="Shell Foundation"
          description="Responsive shell primitives used by the admin navigation frame."
        >
          <div className="grid gap-4 lg:grid-cols-2">
            <div className="space-y-3 rounded-md border p-4">
              <h3 className="text-base font-semibold">SidebarProvider</h3>
              <SidebarProvider className="min-h-0 rounded-md border bg-muted/30">
                <div className="flex min-h-16 w-full items-center gap-2 p-3">
                  <SidebarTrigger aria-label="Toggle UI-kit sidebar" />
                  <span className="text-sm text-muted-foreground">
                    Persistent navigation state wrapper
                  </span>
                </div>
              </SidebarProvider>
            </div>

            <div className="space-y-3 rounded-md border p-4">
              <h3 className="text-base font-semibold">Breadcrumb</h3>
              <Breadcrumb aria-label="Shell breadcrumb">
                <BreadcrumbList>
                  <BreadcrumbItem>
                    <BreadcrumbPage>Current route</BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </div>

            <div className="space-y-3 rounded-md border p-4">
              <h3 className="text-base font-semibold">Avatar</h3>
              <div className="flex items-center gap-3">
                <Avatar>
                  <AvatarFallback>AV</AvatarFallback>
                </Avatar>
                <span className="text-sm text-muted-foreground">User menu identity slot</span>
              </div>
            </div>

            <div className="space-y-3 rounded-md border p-4">
              <h3 className="text-base font-semibold">Collapsible</h3>
              <Collapsible defaultOpen>
                <CollapsibleTrigger asChild>
                  <Button variant="outline">Collapsible</Button>
                </CollapsibleTrigger>
                <CollapsibleContent className="pt-3 text-sm text-muted-foreground">
                  Nested navigation branch content
                </CollapsibleContent>
              </Collapsible>
            </div>

            <div className="space-y-3 rounded-md border p-4 lg:col-span-2">
              <h3 className="text-base font-semibold">Sheet</h3>
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="outline">Sheet</Button>
                </SheetTrigger>
                <SheetContent>
                  <SheetHeader>
                    <SheetTitle>Mobile shell panel</SheetTitle>
                    <SheetDescription>Responsive navigation content lives here.</SheetDescription>
                  </SheetHeader>
                </SheetContent>
              </Sheet>
            </div>
          </div>
        </AdminSection>

        <AdminSection title="Admin Compositions" description="Recommended route-building blocks.">
          <div className="grid gap-3 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>AdminPageShell</CardTitle>
                <CardDescription>Wraps every admin route.</CardDescription>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminPageHeader</CardTitle>
                <CardDescription>Standardizes titles, descriptions, and actions.</CardDescription>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminToolbar</CardTitle>
                <CardDescription>Groups filters and route commands.</CardDescription>
              </CardHeader>
              <CardContent>
                <AdminToolbar>
                  <Input aria-label="Filter users" placeholder="Filter users" />
                  <Button variant="outline">Refresh</Button>
                </AdminToolbar>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>AdminSection</CardTitle>
                <CardDescription>Frames repeated route panels.</CardDescription>
              </CardHeader>
            </Card>
            <AdminEmptyState
              title="No filtered users"
              description="Clear filters to show the full directory."
              action={<Button variant="outline">Clear filters</Button>}
            />
          </div>
        </AdminSection>
      </AdminPageShell>
    </TooltipProvider>
  );
}
