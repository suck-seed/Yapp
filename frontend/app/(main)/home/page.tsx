'use client';

import { SignedIn, SignedOut, useUser } from '@clerk/nextjs';

export default function HomePage() {
  const { user } = useUser();

  return (
    <div>
      <SignedIn>
        <h1>Hello {user?.firstName}, welcome to Yapp!</h1>
      </SignedIn>

      <SignedOut>
        <h1>Hello, you are not signed in.</h1>
      </SignedOut>

      {/* CAPTCHA Widget - Clerk will inject the CAPTCHA */}
      <div id="clerk-captcha" className="mb-4"></div>
    </div>
  );
}

// This is a protected page that only signed-in users can access. But it not protection is not enforced yet.
// Protection will be enforced in the future.