// FILE: apps/web-admin/src/shared/ui/hooks/use-mobile.ts
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide the shadcn-compatible mobile breakpoint hook used by sidebar primitives.
//   SCOPE: Owns viewport media-query detection for sidebar responsive rendering; excludes layout composition.
//   DEPENDS: react, browser matchMedia.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   useIsMobile - Return whether the current viewport is below the sidebar mobile breakpoint.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added sidebar mobile breakpoint detection for the web-admin shell foundation.
// END_CHANGE_SUMMARY

import * as React from 'react';

const MOBILE_BREAKPOINT = 768;

// START_CONTRACT: useIsMobile
//   PURPOSE: Detect whether the current viewport should render mobile sidebar behavior.
//   INPUTS: none.
//   OUTPUTS: { boolean - true when window.innerWidth is below the sidebar mobile breakpoint }
//   SIDE_EFFECTS: Registers and removes a matchMedia change listener while mounted.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: useIsMobile
export function useIsMobile() {
  const [isMobile, setIsMobile] = React.useState<boolean | undefined>(undefined);

  React.useEffect(() => {
    const mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`);
    const onChange = () => {
      setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    };
    mql.addEventListener('change', onChange);
    setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    return () => mql.removeEventListener('change', onChange);
  }, []);

  return !!isMobile;
}
