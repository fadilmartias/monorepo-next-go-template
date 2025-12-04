"use client";

import { APIErrorResponse, APIResponse } from "@/types/api";
import { createSignature } from "@/utils/signature";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { toast } from "sonner";

const apiInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Tenant-Id": process.env.NEXT_PUBLIC_MERCHANT_ID || "",
  },
});

apiInstance.interceptors.request.use((config) => {
  const method = config.method?.toUpperCase() || "GET";
  const url = config.url || "/";

  // PENTING: Untuk GET/DELETE gunakan params, untuk POST/PUT/PATCH gunakan data
  let body: Record<string, any> | undefined;

  if (method === "GET" || method === "DELETE") {
    // Untuk GET/DELETE, query params ada di config.params
    body = config.params || undefined;
  } else {
    // Untuk POST/PUT/PATCH, body ada di config.data
    body = config.data || undefined;
  }

  const signatureHeaders = createSignature({
    method: method as "GET" | "DELETE" | "POST" | "PUT",
    url,
    body,
    clientSecret: process.env.NEXT_PUBLIC_API_CLIENT_KEY || "",
    isMultipart: config.headers?.get("Content-Type") === "multipart/form-data",
  });

  Object.entries(signatureHeaders).forEach(([key, value]) => {
    config.headers?.set?.(key, value);
  });

  return config;
});

apiInstance.interceptors.response.use(
  (res: AxiosResponse<APIResponse<any>>) => {
    // sukses → return res full (status, headers, data, dll)
    res.data.rc = res.status;
    return res;
  },
  async (err: AxiosError) => {
    let formattedError: APIErrorResponse = {
      rc: err.response?.status || 500,
      success: false,
      message: "Unknown error",
    };

    if (err.response) {
      formattedError = {
        ...(err.response.data as APIErrorResponse),
        rc: err.response.status,
      };
    } else if (err.request) {
      formattedError.message = "No response from server";
    } else {
      formattedError.message = err.message;
    }

    return Promise.reject(formattedError);
  }
);

export const apiClient = {
  get: <T = any>(
    path: string,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> =>
    apiInstance.get<APIResponse<T>>(path, config),

  post: <T = any>(
    path: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> =>
    apiInstance.post<APIResponse<T>>(path, data, config),

  put: <T = any>(
    path: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> =>
    apiInstance.put<APIResponse<T>>(path, data, config),

  patch: <T = any>(
    path: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> =>
    apiInstance.patch<APIResponse<T>>(path, data, config),

  delete: <T = any>(
    path: string,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<APIResponse<T>>> =>
    apiInstance.delete<APIResponse<T>>(path, config),
};
