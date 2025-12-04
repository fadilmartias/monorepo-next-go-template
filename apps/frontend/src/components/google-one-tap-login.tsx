"use client";

import { apiClient } from "@/lib/api-client";
import { useEffect } from "react";
import { toast } from "sonner";
import { useRouter } from "nextjs-toploader/app";
import { useAuthStore } from "@/stores/auth-store";

export default function GoogleOneTapLogin() {
  const router = useRouter();
  useEffect(() => {
    if (window.google) {
      window.google.accounts.id.initialize({
        client_id: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID!,
        callback: handleCredentialResponse,
        auto_select: false, // kalau user udah pernah pilih akun, langsung login
        use_fedcm_for_prompt: true,
      });

      // Ini yang munculin popup di pojok kanan
      window.google.accounts.id.prompt((notification: any) => {
        if (notification.isNotDisplayed()) {
          console.warn(
            "One Tap not displayed:",
            notification.getNotDisplayedReason()
          );
        }
        if (notification.isSkippedMoment()) {
          console.warn("One Tap skipped:", notification.getSkippedReason());
        }
        if (notification.isDismissedMoment()) {
          console.warn("One Tap dismissed:", notification.getDismissedReason());
        }
      });
    }
  }, []);

  const handleCredentialResponse = async (response: any) => {
    try {
      const res = await apiClient.post("/v1/auth/google/one-tap", {
        credential: response.credential,
      });
      const data = res.data.data;
      if (data.is_2fa_enabled) {
        router.push(`/login/2fa/${data.temp_token}`);
        return;
      }
      useAuthStore.setState({ user: data });
      toast(`Selamat datang kembali, ${data.name}`, { duration: 5000 });

      if (data.role === "admin") {
        router.push("/admin/dashboard");
      } else {
        router.push("/");
      }
      router.refresh();
    } catch (error) {
      toast.error("Terjadi kesalahan, coba lagi");
      console.error("Login error:", error);
    }
  };

  // ga perlu render button lagi
  return null;
}
