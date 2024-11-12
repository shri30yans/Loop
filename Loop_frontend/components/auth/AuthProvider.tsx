// AuthProvider.tsx
'use client'
import { useEffect, ReactNode } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '../../lib/auth/authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

interface AuthProviderProps {
  children: ReactNode;
}

const protectedPaths = ['/projects', '/create'];

export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { refresh_token, logout } = useAuthStore();

  useEffect(() => {
    async function verifyAuth() {
      // Redirect if there's no refresh token on protected paths
      if (!refresh_token && protectedPaths.includes(pathname)) {
        router.push('/auth/login');
        return;
      }

      // If refresh token exists, verify it only once per session
      if (refresh_token && protectedPaths.includes(pathname)) {
        try {
          const response = await fetch(`${API_BASE_URL}/auth/verify`, {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${refresh_token}`,
              'Content-Type': 'application/json'
            },
            credentials: 'include',
            mode: 'cors' 
          });

          if (!response.ok) {
            throw new Error('Verification failed');
          }

          const data = await response.json();
          if (!data.session_id) {
            throw new Error('Invalid session');
          }
        } catch (error) {
          console.error('Verification error:', error);
          logout();
          router.push('/auth/login');
        }
      }
    }

    verifyAuth();
  }, [refresh_token, pathname, router, logout]);

  return <>{children}</>;
}