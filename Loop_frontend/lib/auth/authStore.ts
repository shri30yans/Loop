import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

interface AuthState {
  refresh_token: string | null;
  user_id: any | null;
  expires_at: any | null;
  setAuth: (refresh_token: string | null, user_id: any | null, expires_at: any | null) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      refresh_token: document.cookie.replace(/(?:(?:^|.*;\s*)refresh_token\s*=\s*([^;]*).*$)|^.*$/, '$1') || null,
      user_id: null,
      expires_at: null,
      setAuth: (refresh_token, user_id, expires_at) => {
        // Set cookie explicitly
        document.cookie = `refresh_token=${refresh_token}; user_id=${user_id}; expires_at=${expires_at}`; // 7 days
        set({ refresh_token, user_id, expires_at });
      },
      logout: () => {
        // Clear cookies and reset state on logout
        document.cookie = 'refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
        set({ user_id: null, refresh_token: null, expires_at: null });
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        console.log('Hydrated state:', state);
      },
    }
  )
);
