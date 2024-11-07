// components/AuthProvider.tsx
'use client'

import { useEffect } from 'react';
import { useAuthStore } from '@/lib/auth';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { token, setAuth } = useAuthStore();

  useEffect(() => {
    // Check token validity on mount
    if (token) {
      // Optionally verify token with backend
      // If invalid, call setAuth(null, null)
    }
  }, []);

  return <>{children}</>;
}