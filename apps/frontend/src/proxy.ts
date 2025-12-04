import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { getUserFromToken } from "./lib/auth";
import { createApiServer } from "./lib/api-server";
import { disableConsoleInProduction } from "@/lib/disable-console";
disableConsoleInProduction();
export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("access_token_" + process.env.NEXT_PUBLIC_APP_ENV)?.value;
  const refreshToken = request.cookies.get("refresh_token_" + process.env.NEXT_PUBLIC_APP_ENV)?.value;
  let visitorId = request.cookies.get("visitor_id")?.value;
  const domain = "nexttemplate.com";

  function getCookieDomain(host: string) {
    if (process.env.NEXT_PUBLIC_APP_ENV === "local") return undefined;
    return `.${host}`;
  }
  
  // Mulai response
  const response = NextResponse.next();
  response.headers.set("X-Pathname", pathname);

  // Inject Visitor ID
  if (!visitorId) {
    visitorId = crypto.randomUUID();
    response.cookies.set("visitor_id", visitorId, {
      httpOnly: false,
      sameSite: "lax",
      path: "/",
      domain: getCookieDomain(domain),
    });
  }

  // ✅ Skip refresh kalau sedang di halaman auth
  const isAuthPage = [
    "/login",
    "/register",
    "/forgot-password",
    "/reset-password",
  ].some((p) => pathname.startsWith(p));

  // 1️⃣ Refresh token jika accessToken hilang dan ada refreshToken
  if (!accessToken && refreshToken && !isAuthPage) {
    try {
      const apiServer = await createApiServer();
      const refreshRes = await apiServer.post(`/v1/auth/refresh-access-token`);
      const data = refreshRes.data;

      // Set cookie baru
      response.cookies.set("access_token_" + process.env.NEXT_PUBLIC_APP_ENV, data.data.access_token, {
        httpOnly: true,
        secure: true,
        path: "/",
        maxAge: 60 * 60 * 1,
        sameSite: "none",
        domain: getCookieDomain(domain),
      });
      response.cookies.set("refresh_token_" + process.env.NEXT_PUBLIC_APP_ENV, data.data.refresh_token, {
        httpOnly: true,
        secure: true,
        path: "/",
        maxAge: 60 * 60 * 24,
        sameSite: "none",
        domain: getCookieDomain(domain),
      });

      return response; // lanjut ke halaman admin
    } catch (error: any) {
      console.error("Refresh token failed:", error.message);
      // Hapus cookie biar tidak loop
      response.cookies.set("access_token_" + process.env.NEXT_PUBLIC_APP_ENV, "", { maxAge: 0, domain: getCookieDomain(domain) });
      response.cookies.set("refresh_token_" + process.env.NEXT_PUBLIC_APP_ENV, "", { maxAge: 0, domain: getCookieDomain(domain) });

      return NextResponse.redirect(new URL("/login?session=expired", request.url));
    }
  }

  // Ambil user dari token
  const user = await getUserFromToken();

  const userMustLogin = [
    "/orders",
  ].some((p) => pathname.startsWith(p));

  // 2️⃣ Proteksi route user
  if (userMustLogin && !user) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // 2️⃣ Proteksi route admin
  if (pathname.startsWith("/admin")) {
    if (!user) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    if (user.role !== "admin") {
      return NextResponse.redirect(new URL("/", request.url));
    }
  }

  // 3️⃣ Proteksi halaman auth (redirect jika sudah login)
  if (isAuthPage && user) {
    return NextResponse.redirect(
      new URL(user.role === "admin" ? "/admin/dashboard" : "/", request.url)
    );
  }

  return response;
}

export const config = {
  matcher: [
    "/((?!api|_next|favicon.ico).*)", // header
    "/admin/:path*", // admin
    "/login",
    "/register",
    "/forgot-password",
    "/reset-password",
  ],
};
