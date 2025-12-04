"use server";

import { cookies } from "next/headers";
import { createApiServer } from "@/lib/api-server";

export async function loginAction(formData: {
  credential: string;
  password: string;
}) {
  const apiServer = await createApiServer();
  try {
    const res = await apiServer.post("/v1/auth/login", formData);

    if (!res.data.success) {
      throw new Error(res.data.message);
    }

    const data = res.data.data;

    const expiresRefreshToken = new Date(Date.now() + 24 * 60 * 60 * 1000);
    const expiresAccessToken = new Date(Date.now() + 60 * 60 * 1000);
    // Set cookie di server (host-only di nexttemplate.com)
    const cookieStore = await cookies();
    cookieStore.set("access_token_" + process.env.NEXT_PUBLIC_APP_ENV, data.access_token, {
      httpOnly: true,
      secure: true,
      path: "/",
      expires: expiresAccessToken,
      sameSite: "none",
    });
    cookieStore.set("refresh_token_" + process.env.NEXT_PUBLIC_APP_ENV, data.refresh_token, {
      httpOnly: true,
      secure: true,
      path: "/",
      expires: expiresRefreshToken,
      sameSite: "none",
    });

    return data.user;
  } catch (error: any) {
    console.error("Error logging in:", error);
    throw new Error(error.message || "Terjadi kesalahan");
  }
}
