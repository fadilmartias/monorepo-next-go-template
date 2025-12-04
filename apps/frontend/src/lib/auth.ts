"use server";

import { createApiServer } from "./api-server";
import { cookies } from "next/headers";
import {jwtDecode} from "jwt-decode";

type User = {
    id: string;
    name: string;
    phone: string;
    email: string;
    role: string;
}

export async function logoutAction(): Promise<void | Error> {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.post("/v1/auth/logout");
    if (response.data.success) {
      const cookieStore = await cookies();
      cookieStore.delete("access_token_" + process.env.NEXT_PUBLIC_APP_ENV);
      cookieStore.delete("refresh_token_" + process.env.NEXT_PUBLIC_APP_ENV);
    }
  } catch (error) {
    console.error("Error logging out:", error);
    throw error;
  }
}

export async function isLoggedIn(): Promise<boolean> {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get("access_token_" + process.env.NEXT_PUBLIC_APP_ENV)?.value || "";
  return !!accessToken;
}

export async function getUser(): Promise<User | Error> {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get("access_token_" + process.env.NEXT_PUBLIC_APP_ENV)?.value || "";
  if (!accessToken) return new Error("Unauthorized");
  try {
    const apiServer = await createApiServer();
    const res = await apiServer.get("/v1/auth/me");
    return res.data.data as User;
  } catch (error) {
    console.error("Error getting user:", error);
    throw error;
  }
}

type TokenPayload = {
    id:    string,
    name:  string,
    phone: string,
    email: string,
    role:  string,
}

export async function getUserFromToken() {
  const cookieStore = await cookies()
  
  const accessToken = cookieStore.get("access_token_" + process.env.NEXT_PUBLIC_APP_ENV)?.value || "";
  
  if (!accessToken) return null

  try {
    const payload = jwtDecode<TokenPayload>(accessToken)
    
    return { id: payload.id, name: payload.name, phone: payload.phone, email: payload.email, role: payload.role }
  } catch (e) {
    return null // token rusak atau format salah
  }
}
