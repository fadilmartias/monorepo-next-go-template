import { create } from "zustand";
import { apiClient } from "@/lib/api-client";
import { persist } from "zustand/middleware";

export interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: string;
  tenant_id: string;
  avatar?: string;
  email_verified_at?: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  is_2fa_enabled: boolean;
  level: number;
  exp: number;
  total_spent: number;
  total_exp: number;
  exp_for_level: number;
  referral_code: string;
  active_title?: {
    id: string;
    name: string;
    rarity: string;
  };
}

interface AuthState {
  user: User | null;
  fetchUser: () => Promise<void>;
  isLoggedIn: boolean;
  session: any;
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isLoggedIn: false,
      session: null,
      fetchUser: async () => {
        try {
          const res = await apiClient.get("/v1/auth/me");
          const data = res.data.data;
          set({
            user: data,
            isLoggedIn: true,
          });
        } catch (err) {
          console.error("Failed to fetch user", err);
          set({ user: null, isLoggedIn: false, session: null });
        }
      },

      logout: () =>
        set({
          user: null,
          isLoggedIn: false,
          session: null,
        }),
    }),
    {
      name: "nexttemplate-user-storage", // key dalam localStorage
    }
  )
);
