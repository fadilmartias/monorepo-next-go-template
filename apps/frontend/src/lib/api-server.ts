"use server";

import { cookies } from "next/headers";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { createSignature } from "@/utils/signature";
import type { APIResponse, APIErrorResponse } from "@/types/api";

export async function createApiServer() {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");

  async function serverFetch<T>(
    method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE",
    path: string,
    body?: any,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> {
    // --- Logic params/data harus sama dengan interceptor ---
    const isGetOrDelete = method === "GET" || method === "DELETE";

    // Untuk GET/DELETE → body dianggap sebagai query params
    const params = isGetOrDelete ? body : undefined;

    const data = isGetOrDelete ? undefined : body;

    // Build full URL untuk GET dengan params

    // Signature harus pakai params untuk GET, dan pakai body untuk POST
    const signatureHeaders = createSignature({
      method,
      url: path, // path tanpa query
      body: isGetOrDelete ? params : data,
      clientSecret: process.env.NEXT_PUBLIC_API_CLIENT_KEY || "",
      isMultipart: false,
    });

    try {
      const res = await axios.request<APIResponse<T>>({
        baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
        method,
        url: path,
        params: isGetOrDelete ? params : undefined, // untuk axios
        data: data,
        withCredentials: true,
        ...config,
        headers: {
          "Content-Type": "application/json",
          Cookie: cookieHeader,
          ...signatureHeaders,
        },
      });

      res.data.rc = res.status;

      return res;
    } catch (err) {
      const error = err as AxiosError<APIErrorResponse>;

      if (error.response?.data) {
        error.response.data.rc = error.response.status;
        throw error.response.data;
      }
      throw {
        rc: error.response?.status || 500,
        success: false,
        message: error.message || "Unknown error",
      };
    }
  }

  return {
    get: <T = any>(path: string, config?: AxiosRequestConfig) =>
      serverFetch<T>("GET", path, config?.params, config),
    post: <T = any>(path: string, data?: any, config?: AxiosRequestConfig) =>
      serverFetch<T>("POST", path, data, config),
    put: <T = any>(path: string, data?: any, config?: AxiosRequestConfig) =>
      serverFetch<T>("PUT", path, data, config),
    patch: <T = any>(path: string, data?: any, config?: AxiosRequestConfig) =>
      serverFetch<T>("PATCH", path, data, config),
    delete: <T = any>(path: string, config?: AxiosRequestConfig) =>
      serverFetch<T>("DELETE", path, config?.params, config),
  };
}
