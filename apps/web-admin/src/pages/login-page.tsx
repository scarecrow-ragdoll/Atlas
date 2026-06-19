// FILE: apps/web-admin/src/pages/login-page.tsx
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Render the public web-admin login route.
//   SCOPE: Owns login form state, error presentation, safe return-to navigation, and already-authenticated redirects; excludes backend session storage and protected shell rendering.
//   DEPENDS: react, react-router, apps/web-admin/src/entities/admin-auth/provider.tsx, apps/web-admin/src/entities/admin-auth/model.ts, apps/web-admin/src/shared/ui.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   default - Public login page for backend cookie-backed admin sessions.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Cleared one-shot logout redirect state on the public login route.
// END_CHANGE_SUMMARY

import { type FormEvent, useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { resolveSafeReturnTo } from '@entities/admin-auth/model';
import { useAdminAuth } from '@entities/admin-auth/provider';
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Button,
  Card,
  CardContent,
  CardHeader,
  Input,
  Label,
} from '@shared/ui';

const NETWORK_ERROR_MESSAGE = 'Unable to sign in. Try again after the API is available.';

function networkErrorMessage() {
  return NETWORK_ERROR_MESSAGE;
}

// START_CONTRACT: LoginPage
//   PURPOSE: Authenticate admins through the backend cookie session and navigate to a safe return path.
//   INPUTS: none.
//   OUTPUTS: { JSX.Element - public login route }
//   SIDE_EFFECTS: Calls AuthProvider.login and navigates within the same app after success.
//   LINKS: M-WEB-ADMIN / V-M-WEB-ADMIN.
// END_CONTRACT: LoginPage
export default function LoginPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { admin, clearCompletedLogout, hasCompletedLogout, login } = useAdminAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const returnTo = resolveSafeReturnTo(new URLSearchParams(location.search).get('from'));

  useEffect(() => {
    if (admin) {
      navigate(returnTo, { replace: true });
      return;
    }
    if (hasCompletedLogout) {
      clearCompletedLogout();
    }
  }, [admin, clearCompletedLogout, hasCompletedLogout, navigate, returnTo]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const result = await login({ email, password });
      if (!result.ok) {
        setError(
          result.error.field
            ? `${result.error.field}: ${result.error.message}`
            : result.error.message,
        );
        return;
      }
      navigate(returnTo, { replace: true });
    } catch {
      setError(networkErrorMessage());
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="grid min-h-svh place-items-center bg-background px-4 py-8">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <h1 className="text-xl font-semibold leading-none">Admin sign in</h1>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="admin-email">Email</Label>
              <Input
                autoComplete="email"
                id="admin-email"
                onChange={(event) => setEmail(event.target.value)}
                required
                type="email"
                value={email}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="admin-password">Password</Label>
              <Input
                autoComplete="current-password"
                id="admin-password"
                onChange={(event) => setPassword(event.target.value)}
                required
                type="password"
                value={password}
              />
            </div>
            {error ? (
              <Alert variant="destructive">
                <AlertTitle>Sign in failed</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : null}
            <Button className="w-full" disabled={isSubmitting} type="submit">
              {isSubmitting ? 'Signing in...' : 'Sign in'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
