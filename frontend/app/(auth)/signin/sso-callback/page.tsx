'use client';
import { useEffect, useState } from 'react';
import { useClerk } from '@clerk/nextjs';
import { useRouter } from 'next/navigation';
import { FaSpinner } from 'react-icons/fa';

export default function SSOCallback() {
  const { handleRedirectCallback } = useClerk();
  const router = useRouter();
  const [error, setError] = useState('');

  useEffect(() => {
    const doCallback = async () => {
      try {
        await handleRedirectCallback({ signUpForceRedirectUrl: '/home' });
      } catch (err) {
        setError('OAuth callback error. Please try again.');
        console.error('OAuth callback error:', err);
      }
    };
    doCallback();
  }, [handleRedirectCallback, router]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      {error ? (
        <p className="text-red-600">{error}</p>
      ) : (
        <div className="flex flex-col items-center gap-2">
          <FaSpinner className="animate-spin text-2xl text-blue-600" />
          <p className="text-lg font-medium">Signing you in...</p>
        </div>
      )}
    </div>
  );
}
