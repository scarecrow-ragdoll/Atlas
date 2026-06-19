// FILE: apps/web-admin/src/shared/ui/primitives/collapsible.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shadcn collapsible primitives for web-admin UI compositions.
//   SCOPE: Owns collapsible root, trigger, and content rendering; excludes page-specific disclosure state decisions.
//   DEPENDS: react, radix-ui Collapsible.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Collapsible - shadcn collapsible root primitive.
//   CollapsibleContent - Collapsible content primitive.
//   CollapsibleTrigger - Collapsible trigger primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added shadcn collapsible primitives under the web-admin UI kit.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { Collapsible as CollapsiblePrimitive } from 'radix-ui';

function Collapsible({ ...props }: React.ComponentProps<typeof CollapsiblePrimitive.Root>) {
  return <CollapsiblePrimitive.Root data-slot="collapsible" {...props} />;
}

function CollapsibleTrigger({
  ...props
}: React.ComponentProps<typeof CollapsiblePrimitive.CollapsibleTrigger>) {
  return <CollapsiblePrimitive.CollapsibleTrigger data-slot="collapsible-trigger" {...props} />;
}

function CollapsibleContent({
  ...props
}: React.ComponentProps<typeof CollapsiblePrimitive.CollapsibleContent>) {
  return <CollapsiblePrimitive.CollapsibleContent data-slot="collapsible-content" {...props} />;
}

export { Collapsible, CollapsibleContent, CollapsibleTrigger };
