"use client";

import { toast } from "sonner";
import { apiClient } from "./api-client";
import { useRouter } from "nextjs-toploader/app";
import { useAuthStore } from "@/stores/auth-store";

export function useLogoutActionClient() {
  const router = useRouter();

  async function logoutActionClient(): Promise<void> {
    try {
      const res = await apiClient.post("/v1/auth/logout");

      if (res.data.success) {
        toast.success("Logout berhasil");
        useAuthStore.setState({ user: null, isLoggedIn: false });
        router.push("/");
        router.refresh();
      }
    } catch (error) {
      console.error("Error logging out:", error);
      toast.error("Gagal logout");
      throw error;
    }
  }

  return { logoutActionClient };
}
