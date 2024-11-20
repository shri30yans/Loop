import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

interface AuthState {
  refresh_token: string | null;
  user_id: string | null;
  expires_at: Date | null;
  setAuth: (refresh_token: string | null, user_id: string | null, expires_at: Date | null) => void;
  logout: () => Promise<void>;
}

// Utility function to manage cookies
const setCookie = (name: string, value: string | null, maxAgeInSeconds?: number) => {
  if (typeof window !== 'undefined') {
    const cookieValue =
      value !== null
        ? `${name}=${value}; path=/; ${maxAgeInSeconds ? `max-age=${maxAgeInSeconds};` : ''} Secure; SameSite=Strict`
        : `${name}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC; Secure; SameSite=Strict`;
    document.cookie = cookieValue;
  }
};

const getCookie = (name: string): string | null => {
  if (typeof window === 'undefined') return null;
  const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : null;
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      refresh_token: getCookie('refresh_token'),
      user_id: getCookie('user_id'),
      expires_at: (() => {
        const expiresAtStr = getCookie('expires_at');
        return expiresAtStr ? new Date(expiresAtStr) : null;
      })(),
      setAuth: (refresh_token, user_id, expires_at) => {
        const expiresAtDate = expires_at instanceof Date ? expires_at : expires_at ? new Date(expires_at) : null;
        if (expiresAtDate === null || isNaN(expiresAtDate.getTime())) {
          console.error('Invalid expires_at:', expires_at);
          return;
        }

        const maxAge = 60 * 60 * 24 * 7; 
        setCookie('refresh_token', refresh_token, maxAge);
        setCookie('user_id', user_id, maxAge);
        setCookie('expires_at', expiresAtDate.toISOString(), maxAge);
        set({ refresh_token, user_id, expires_at: expiresAtDate });
      },
      logout: async () => {
        try {
          const cookiesToClear = ['refresh_token', 'user_id', 'expires_at'];
          cookiesToClear.forEach((cookie) => setCookie(cookie, null));
          set({ refresh_token: null, user_id: null, expires_at: null });
        } catch (error) {
          console.error('Logout failed:', error);
          throw error;
        }
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        if (state?.expires_at) {
          try {
            state.expires_at = new Date(state.expires_at as unknown as string);
            if (isNaN(state.expires_at.getTime())) {
              state.expires_at = null;
            }
          } catch {
            state.expires_at = null;
          }
        }
        console.log('Hydrated state:', state);
      },
    }
  )
);
