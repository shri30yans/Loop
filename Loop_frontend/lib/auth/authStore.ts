import { create, StateCreator } from "zustand";
import { persist, PersistOptions } from "zustand/middleware";
import { AuthError, AUTH_ERROR_CODES, NetworkError } from "@/utils/errors";

// Define the partial state that will be persisted
type PersistedAuthState = {
  userId: string | null;
  token: string | null;
  expiresAt: string | null;
};

export interface AuthState {
  userId: string | null;
  token: string | null;
  expiresAt: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: AuthError | NetworkError | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  register: (
    username: string,
    email: string,
    password: string,
  ) => Promise<void>;
  checkAuth: () => Promise<void>;
}

type AuthPersist = (
  config: StateCreator<AuthState>,
  options: PersistOptions<AuthState, PersistedAuthState>,
) => StateCreator<AuthState>;

export const useAuth = create<AuthState>()(
  (persist as AuthPersist)(
    (set, get): AuthState => ({
      userId: null,
      token: null,
      expiresAt: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(
            `${process.env.NEXT_PUBLIC_BASE_API_URL}/auth/login`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ email, password }),
            },
          );

          if (!response.ok) {
            const errorData = await response.json();
            
            // Map backend error codes to our frontend error types
            switch (response.status) {
              case 401:
                throw new AuthError(
                  errorData.message || "Invalid credentials",
                  AUTH_ERROR_CODES.INVALID_CREDENTIALS
                );
              case 404:
                throw new AuthError(
                  errorData.message || "User not found",
                  AUTH_ERROR_CODES.USER_NOT_FOUND
                );
              case 503:
                throw new AuthError(
                  "Server is currently unavailable",
                  AUTH_ERROR_CODES.BACKEND_DOWN
                );
              default:
                throw new AuthError(
                  errorData.message || "Login failed",
                  AUTH_ERROR_CODES.NETWORK_ERROR
                );
            }
          }

          const data = await response.json();
          set({
            userId: data.user_id,
            token: data.access_token,
            expiresAt: data.expires_at,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            error: error instanceof Error ? error : new AuthError("Login failed", AUTH_ERROR_CODES.NETWORK_ERROR),
            isLoading: false,
          });
          throw error;
        }
      },

      logout: () => {
        set({
          userId: null,
          token: null,
          expiresAt: null,
          isAuthenticated: false,
        });
      },

      register: async (username: string, email: string, password: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(
            `${process.env.NEXT_PUBLIC_BASE_API_URL}/auth/register`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ username, email, password }),
            },
          );

          if (!response.ok) {
            const errorData = await response.json();
            
            // Map backend error codes to our frontend error types
            switch (response.status) {
              case 409:
                throw new AuthError(
                  errorData.message || "Email already exists",
                  AUTH_ERROR_CODES.EMAIL_EXISTS
                );
              case 503:
                throw new AuthError(
                  "Server is currently unavailable",
                  AUTH_ERROR_CODES.BACKEND_DOWN
                );
              default:
                throw new AuthError(
                  errorData.message || "Registration failed",
                  AUTH_ERROR_CODES.NETWORK_ERROR
                );
            }
          }

          const data = await response.json();
          set({
            userId: data.user_id,
            token: data.access_token,
            expiresAt: data.expires_at,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          const authError = error instanceof Error ? 
            error instanceof AuthError ? 
              error : 
              new AuthError(error.message, AUTH_ERROR_CODES.NETWORK_ERROR) 
            : new AuthError("Registration failed", AUTH_ERROR_CODES.NETWORK_ERROR);
          
          set({
            error: authError,
            isLoading: false,
          });
          throw error;
        }
      },

      checkAuth: async () => {
        const token = get().token;
        if (!token) {
          set({ isAuthenticated: false });
          return;
        }

        set({ isLoading: true });
        try {
          const response = await fetch(
            `${process.env.NEXT_PUBLIC_BASE_API_URL}/auth/verify`,
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            },
          );

          if (!response.ok) {
            throw new AuthError(
              "Authentication failed",
              AUTH_ERROR_CODES.INVALID_CREDENTIALS
            );
          }

          set({
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            userId: null,
            token: null,
            expiresAt: null,
            isAuthenticated: false,
            isLoading: false,
            error: error instanceof AuthError ? error : new AuthError(
              "Authentication failed",
              AUTH_ERROR_CODES.NETWORK_ERROR
            ),
          });
        }
      },
    }),
    {
      name: "auth-storage",
      partialize: (state: AuthState): PersistedAuthState => ({
        userId: state.userId,
        token: state.token,
        expiresAt: state.expiresAt,
      }),
    },
  ),
);
