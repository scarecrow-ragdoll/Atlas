// FILE: apps/web/src/shared/ui/primitives/button.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn-compatible button primitive for the public web app.
//   SCOPE: Owns public web button variants and primitive rendering; excludes page-specific mutation behavior.
//   DEPENDS: react, class-variance-authority, apps/web/src/shared/ui/lib/utils.ts.
//   LINKS: M-WEB / V-M-WEB.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Button - Public web button primitive.
//   buttonVariants - Variant class generator for public web buttons.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Added the icon size variant for public theme toggle buttons.
// END_CHANGE_SUMMARY

import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

import { cn } from '../lib/utils';

const buttonVariants = cva('web-button', {
  variants: {
    size: {
      default: 'web-button--default-size',
      icon: 'web-button--icon',
      sm: 'web-button--sm',
    },
    variant: {
      default: 'web-button--default',
      ghost: 'web-button--ghost',
      outline: 'web-button--outline',
    },
  },
  defaultVariants: {
    size: 'default',
    variant: 'default',
  },
});

function Button({
  className,
  size = 'default',
  variant = 'default',
  ...props
}: React.ComponentProps<'button'> & VariantProps<typeof buttonVariants>) {
  return (
    <button
      data-slot="button"
      data-size={size}
      data-variant={variant}
      className={cn(buttonVariants({ className, size, variant }))}
      {...props}
    />
  );
}

export { Button, buttonVariants };
