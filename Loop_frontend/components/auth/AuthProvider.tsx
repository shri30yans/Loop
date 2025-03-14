// AuthProvider.tsx
'use client'
import { useEffect, ReactNode } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '../../lib/auth/authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

interface AuthProviderProps {
  children: ReactNode;
}

// Define path patterns with regex
const protectedPathPatterns = [
  /^\/projects(\/.*)?$/, // Matches /projects and all sub-paths
  /^\/create(\/.*)?$/,   // Matches /create and all sub-paths
  /^\/account(\/.*)?$/,     // Matches /user and all sub-paths
  /^\/edit\/\d+(\/.*)?$/ // Matches /edit/123 and similar paths
];

// Helper function to check if path needs protection
const isProtectedPath = (path: string): boolean => {
  return protectedPathPatterns.some(pattern => pattern.test(path));
};

export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { access_token, logout } = useAuthStore();

  useEffect(() => {
    async function verifyAuth() {
      // Redirect if there's no access token on protected paths
      if (!access_token && isProtectedPath(pathname)) {
        router.push('/auth/login');
        return;
      }

      // If access token exists, verify it only once per session
      if (access_token && isProtectedPath(pathname)) {
        try {
          console.log("SENT TO VERIFY!!!")
          const response = await fetch(`${API_BASE_URL}/auth/verify`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${access_token}`,
          'Content-Type': 'application/json'
        },
        credentials: 'include'
          });

          if (!response.ok) {
            throw new Error('Verification failed');
          }

          // const data = await response.json();
          // if (!data.session_id) {
          //   throw new Error('Invalid session');
          // }
        } catch (error) {
          console.error('Verification error:', error);
          logout();
          router.push('/auth/login');
        }
      }
    }

    verifyAuth();
  }, [access_token, pathname, router, logout]);

  return <>{children}</>;
}