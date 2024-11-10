'use client';

import { useEffect } from 'react';
import { useAuthStore } from '../../lib/auth/authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { setAuth } = useAuthStore();

  useEffect(() => {
    // Get the refresh_token from cookies
    const token = document.cookie.replace(
      /(?:(?:^|.*;\s*)refresh_token\s*=\s*([^;]*).*$)|^.*$/,
      '$1'
    );

    if (token) {
      fetch(`${API_BASE_URL}/auth/verify`, {
        headers: { Authorization: `Bearer ${token}` },
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.user_id) {
            setAuth(data.refresh_token, data.user_id, data.expires_at );
          } else {
            document.cookie = 'refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
            setAuth(null, null,null);
          }
        })
        .catch(() => {
          document.cookie = 'refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
          setAuth(null, null,null);
        });
    } else {
      setAuth(null, null,null);
    }
  }, [setAuth]);

  return <>{children}</>;
}
