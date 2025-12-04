"use client";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Save } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuthStore } from "@/stores/auth-store";
import { useEffect } from "react";
import { apiClient } from "@/lib/api-client";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { setFormErrors } from "@/utils/errors";

const schema = z.object({
  name: z.string().min(1, "Nama tidak boleh kosong"),
  phone: z.string().min(1, "Nomor HP tidak boleh kosong"),
  email: z.string().email("Email tidak valid"),
});
export type ProfileType = z.infer<typeof schema>;

export default function EditProfileForm() {
  const { user } = useAuthStore();

  const profileForm = useForm<ProfileType>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      email: "",
      phone: "",
    },
  });

  useEffect(() => {
    if (user) {
      profileForm.reset({
        name: user.name,
        email: user.email,
        phone: user.phone,
      });
    }
  }, [user]);

  const mutateProfile = useMutation({
    mutationFn: async (data: ProfileType) => {
      const res = await apiClient.patch("/v1/users/profile", data);
      return res;
    },
    onSuccess: (res) => {
      toast.success("Profile berhasil diperbarui");
      useAuthStore.setState({ user: res.data.data });
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

  const onSubmit = async (data: ProfileType) => {
    mutateProfile.mutate(data);
  };

  return (
    <form onSubmit={profileForm.handleSubmit(onSubmit)} className="space-y-6">
      <Controller
        name="name"
        control={profileForm.control}
        render={({ field, fieldState: { error } }) => (
          <Input
            {...field}
            label="Nama Lengkap"
            id="name"
            required
            error={error?.message}
          />
        )}
      />
      <Controller
        name="email"
        control={profileForm.control}
        render={({ field, fieldState: { error } }) => (
          <Input
            {...field}
            label="Email"
            id="email"
            required
            error={error?.message}
          />
        )}
      />
      <div>
        <Controller
          name="phone"
          control={profileForm.control}
          render={({ field, fieldState: { error } }) => (
            <Input
              {...field}
              label="Nomor HP"
              id="phone"
              required
              error={error?.message}
            />
          )}
        />
      </div>
      <div className="flex justify-end">
        <Button isLoading={mutateProfile.isPending} className="mt-4 bg-primary text-primary-foreground hover:bg-accent hover:text-accent-foreground">
          <Save className="w-4 h-4" /> Simpan Perubahan
        </Button>
      </div>
    </form>
  );
}
