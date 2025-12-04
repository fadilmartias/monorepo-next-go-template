"use client";

import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Lock } from "lucide-react";
import { Controller } from "react-hook-form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import { useMutation } from "@tanstack/react-query";
import { setFormErrors } from "@/utils/errors";

const schema = z
  .object({
    current_password: z.string().min(8, "Password lama minimal 8 karakter"),
    new_password: z.string().min(8, "Password baru minimal 8 karakter"),
    new_password_confirmation: z
      .string()
      .min(8, "Konfirmasi password minimal 8 karakter"),
  })
  .refine((data) => data.new_password === data.new_password_confirmation, {
    path: ["new_password_confirmation"], // error akan diarahkan ke field confirm_password
    message: "Password tidak cocok",
  });

export type ProfilePasswordType = z.infer<typeof schema>;

export default function EditPasswordForm() {
  const profileForm = useForm<ProfilePasswordType>({
    resolver: zodResolver(schema),
    defaultValues: {
      current_password: "",
      new_password: "",
      new_password_confirmation: "",
    },
  });

  const mutateProfile = useMutation({
    mutationFn: async (data: ProfilePasswordType) => {
      const res = await apiClient.patch("/v1/users/password", data);
      return res;
    },
    onSuccess: () => {
      toast.success("Password berhasil diperbarui");
      profileForm.reset();
    },
    onError: (error: any) => {
      const errors = error.errors;
      if (errors) {
        setFormErrors(profileForm.setError, errors);
      } else {
        toast.error(error.message);
      }
    },
  });

  const onSubmit = async (data: ProfilePasswordType) => {
    mutateProfile.mutate(data);
  };

  return (
    <form onSubmit={profileForm.handleSubmit(onSubmit)} className="space-y-6">
      <Controller
        name="current_password"
        control={profileForm.control}
        render={({ field, fieldState: { error } }) => (
          <Input
            {...field}
            label="Password Lama"
            type="password"
            id="current_password"
            required
            error={error?.message}
          />
        )}
      />

      <Controller
        name="new_password"
        control={profileForm.control}
        render={({ field, fieldState: { error } }) => (
          <Input
            {...field}
            label="Password Baru"
            type="password"
            id="new_password"
            required
            error={error?.message}
          />
        )}
      />

      <Controller
        name="new_password_confirmation"
        control={profileForm.control}
        render={({ field, fieldState: { error } }) => (
          <Input
            {...field}
            label="Konfirmasi Password"
            type="password"
            id="new_password_confirmation"
            required
            error={error?.message}
          />
        )}
      />

      <div className="flex justify-end">
        <Button isLoading={mutateProfile.isPending} className="mt-4 bg-destructive text-destructive-foreground hover:bg-accent hover:text-accent-foreground">
          <Lock className="w-4 h-4" /> Update Password
        </Button>
      </div>
    </form>
  );
}
