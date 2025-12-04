import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // cacheComponents: true,
  reactCompiler: true,
  images: {
    dangerouslyAllowLocalIP: true,
    qualities: [75, 100],
    remotePatterns: [
      {
        protocol: "http",
        hostname: "localhost",
        port: "9001",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "dev-api.nexttemplate.com",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "api.nexttemplate.com",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "cdn.nexttemplate.com",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "dev-cdn.nexttemplate.com",
        port: "",
        pathname: "/**",
      },

    ],
  },
  env: {
    TZ: "Asia/Jakarta",
  }
};

export default nextConfig;
