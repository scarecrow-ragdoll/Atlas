// FILE: apps/web/src/shared/ui/primitives/alert.tsx
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide shadcn-compatible alert primitives for the public web app.
//   SCOPE: Owns public web alert rendering and variants; excludes page-specific error translation.
//   DEPENDS: react, class-variance-authority, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Alert - Public web alert primitive.
//   AlertTitle - Public web alert title primitive.
//   AlertDescription - Public web alert body primitive.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added independent public web alert primitives.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

import { cn } from '../lib/utils';

const alertVariants = cva('web-alert', {
  variants: {
    variant: {
      default: 'web-alert--default',
      destructive: 'web-alert--destructive',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
});

function Alert({
  className,
  variant = 'default',
  ...props
}: React.ComponentProps<'div'> & VariantProps<typeof alertVariants>) {
  return (
    <div
      data-slot="alert"
      data-variant={variant}
      role="alert"
      className={cn(alertVariants({ className, variant }))}
      {...props}
    />
  );
}

function AlertTitle({ className, ...props }: React.ComponentProps<'p'>) {
  return <p data-slot="alert-title" className={cn('web-alert-title', className)} {...props} />;
}

function AlertDescription({ className, ...props }: React.ComponentProps<'p'>) {
  return (
    <p
      data-slot="alert-description"
      className={cn('web-alert-description', className)}
      {...props}
    />
  );
}

export { Alert, AlertDescription, AlertTitle };
