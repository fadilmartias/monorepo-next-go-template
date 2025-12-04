import { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
    const baseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL || 'https://nexttemplate.com'
  return {
    rules: [
      {
        userAgent: "*",
        allow: "/",
        disallow: ["/admin*","/login","/register","/forgot-password","/reset-password","/payment*"], // semua route admin & turunannya
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
  };
}
