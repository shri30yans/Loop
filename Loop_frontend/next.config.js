// Import required modules
const crypto = require('crypto');

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Essential configuration for Next.js
  reactStrictMode: true,

  // Required for NextUI components
  experimental: {
    optimizePackageImports: ['@nextui-org/react']
  },

  // Required Node.js polyfills for browser environment
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        crypto: require.resolve('crypto-browserify'),
        stream: require.resolve('stream-browserify'),
        buffer: require.resolve('buffer'),
      };
    }
    return config;
  }
}

module.exports = nextConfig
