import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

const protectedPaths = ['/projects'];

export async function middleware(request: NextRequest) {
  // Retrieve the token from cookies or the Authorization header
  const token = request.cookies.get('refresh_token') || request.headers.get('Authorization');

  // Check if the request is for a protected path
  if (protectedPaths.some((path) => request.nextUrl.pathname.startsWith(path))) {
    // If there's no token, redirect to login
    if (!token) {
      return NextResponse.redirect(new URL('/auth/login', request.url));
    }

    // Verify the token by calling the backend verify API
    const response = await fetch(`${API_BASE_URL}/auth/verify`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    });

    // If token verification fails, redirect to login
    const data = await response.json();
    if (!response.ok || !data.user) {
      return NextResponse.redirect(new URL('/auth/login', request.url));
    }

    // If token is valid, proceed with the request
    return NextResponse.next();
  }

  // If the request is not for a protected route, proceed normally
  return NextResponse.next();
}

export const config = {
  matcher: protectedPaths, // Apply this middleware to these routes
};
