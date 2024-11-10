import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

const protectedPaths = ['/projects'];

export async function middleware(request: NextRequest) {
  const token = request.cookies.get('refresh_token') || request.headers.get('Authorization');

  if (protectedPaths.some((path) => request.nextUrl.pathname.startsWith(path))) {
    if (!token) {
      return NextResponse.redirect(new URL('/auth/login', request.url));
    }

    const response = await fetch(`${API_BASE_URL}/auth/verify`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    });

    const data = await response.json();
    console.log(data)
    if (!response.ok || !data.user) {
      // If token is invalid or user is not found, remove the token and redirect to login
      return NextResponse.redirect(new URL('/auth/login', request.url));
    }

    // If the token is valid, proceed with the request
    return NextResponse.next();
  }

  // Proceed with non-protected routes
  return NextResponse.next();
}

export const config = {
  matcher: protectedPaths,
};
