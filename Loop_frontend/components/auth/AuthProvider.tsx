// AuthProvider.tsx
"use client";
import { useEffect, ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "../../lib/auth/authStore";

const API_URL = process.env.NEXT_PUBLIC_BASE_API_URL;

interface AuthProviderProps {
  children: ReactNode;
}

// Define path patterns with regex
const protectedPathPatterns = [
  /^\/projects(\/.*)?$/, // Matches /projects and all sub-paths
  /^\/create(\/.*)?$/, // Matches /create and all sub-paths
  /^\/account(\/.*)?$/, // Matches /user and all sub-paths
  /^\/edit\/\d+(\/.*)?$/, // Matches /edit/123 and similar paths
];

// Helper function to check if path needs protection
const isProtectedPath = (path: string): boolean => {
  return protectedPathPatterns.some((pattern) => pattern.test(path));
};

export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { token, logout } = useAuth();
  
  useEffect(() => {
    // Immediately redirect if no token on protected paths
    if (!token && isProtectedPath(pathname)) {
      router.push("/auth/login");
      return;
    }

    // Non-blocking token verification
    if (token && isProtectedPath(pathname)) {
      // Use AbortController to cancel pending requests if component unmounts
      const controller = new AbortController();
      
      fetch(`${API_URL}/auth/verify`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        credentials: "include",
        signal: controller.signal
      })
      .then(response => {
        if (!response.ok) {
          throw new Error("Verification failed");
        }
      })
      .catch(error => {
        if (error.name === 'AbortError') return;
        
        console.error("Verification error:", error);
        logout();
        router.push("/auth/login");
      });

      // Cleanup function to abort pending requests
      return () => controller.abort();
    }
  }, [token, pathname]);

  // Render children immediately while verification happens in background
  return <>{children}</>;
}
