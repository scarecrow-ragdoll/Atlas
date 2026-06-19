// FILE: apps/web/src/shared/ui/primitives/card.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shadcn-compatible card primitives for the public web app.
//   SCOPE: Owns public web card container, header, title, description, and content rendering; excludes page-specific layout behavior.
//   DEPENDS: react, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Card - Public web card container primitive.
//   CardHeader - Public web card header primitive.
//   CardTitle - Public web card title primitive.
//   CardDescription - Public web card description primitive.
//   CardContent - Public web card content primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added independent public web card primitives.
// END_CHANGE_SUMMARY

import * as React from 'react';

import { cn } from '../lib/utils';

function Card({ className, ...props }: React.ComponentProps<'section'>) {
  return <section data-slot="card" className={cn('web-card', className)} {...props} />;
}

function CardHeader({ className, ...props }: React.ComponentProps<'div'>) {
  return <div data-slot="card-header" className={cn('web-card-header', className)} {...props} />;
}

function CardTitle({ className, ...props }: React.ComponentProps<'h1'>) {
  return <h1 data-slot="card-title" className={cn('web-card-title', className)} {...props} />;
}

function CardDescription({ className, ...props }: React.ComponentProps<'p'>) {
  return (
    <p data-slot="card-description" className={cn('web-card-description', className)} {...props} />
  );
}

function CardContent({ className, ...props }: React.ComponentProps<'div'>) {
  return <div data-slot="card-content" className={cn('web-card-content', className)} {...props} />;
}

export { Card, CardContent, CardDescription, CardHeader, CardTitle };
