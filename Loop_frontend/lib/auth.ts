import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';
import type { User } from './types';


interface AuthState {
  token: string | null;
  user: User | null;
  setAuth: (token: string | null, user: User | null) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      setAuth: (token: any, user: any) => set({ token, user }),
      logout: () => set({ token: null, user: null }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => sessionStorage), // or localStorage
    }
  )
);