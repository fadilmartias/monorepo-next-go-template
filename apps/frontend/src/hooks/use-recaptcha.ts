"use client";

import { apiClient } from "@/lib/api-client";
import { useGoogleReCaptcha } from "react-google-recaptcha-v3";

export default function useReCaptcha() {
  const { executeRecaptcha } = useGoogleReCaptcha();

  const getToken = async (action: string) => {
    if (!executeRecaptcha) {
      console.warn("reCAPTCHA not yet initialized");
      return null;
    }

    try {
      const token = await executeRecaptcha(action);
      return token;
    } catch (err) {
      console.error("Failed to execute reCAPTCHA:", err);
      return null;
    }
  };

  interface BackendRecaptchaResponse {
    success: boolean;
    message?: string;
  }

  const verifyRecaptcha = async (
    token: string,
    action: string
  ): Promise<BackendRecaptchaResponse> => {
    try {
      if (!token) {
        return {
          success: false,
          message: `Missing reCAPTCHA token for ${action}`,
        };
      }

      const res = await apiClient.post(`/v1/auth/verify-recaptcha`, {
        token,
        action,
      });

      if (!res.data.success) {
        return {
          success: false,
          message: res.data.message,
        };
      }

      return {
        success: res.data.success,
        message: res.data.message,
      };
    } catch (error) {
      console.error("verifyRecaptcha error:", error);
      return { success: false, message: `Internal server error for ${action}` };
    }
  };

  return { getToken, verifyRecaptcha };
}
