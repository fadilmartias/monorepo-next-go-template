"use client";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Controller } from "react-hook-form";
import Typography from "@/components/typography";
import Logo from "@/components/logo";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { useState } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import router from "next/router";

const resetPasswordSchema = z
  .object({
    token: z.string(),
    password: z.string().min(6, "Kata Sandi minimal 6 karakter"),
    password_confirmation: z
      .string()
      .min(6, "Konfirmasi Kata Sandi minimal 6 karakter"),
  })
  .superRefine((data, ctx) => {
    if (data.password !== data.password_confirmation) {
      ctx.addIssue({
        path: ["password_confirmation"],
        code: z.ZodIssueCode.custom,
        message: "Kata Sandi tidak cocok",
      });
    }
  });

type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>;

export function ResetPasswordForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [isPending, setIsPending] = useState(false);
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  const resetPasswordForm = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      token: token || "",
      password: "",
      password_confirmation: "",
    },
  });
  const onSubmit = async (data: ResetPasswordFormValues) => {
    resetPasswordForm.clearErrors();
    try {
      setIsPending(true);
      const res = await apiClient.post("/v1/auth/reset-password", data);
      toast.success(`Berhasil mengatur ulang kata sandi`, {
        duration: 5000,
      });
      router.push("/login");
    } catch (err: any) {
      console.error(err);
      toast.error(
        err.message || "Gagal mengatur ulang kata sandi",
        { duration: 5000 }
      );
    } finally {
      setIsPending(false);
    }
  };
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <div className="flex justify-center">
        <Logo size="2xl" />
      </div>
      <Card>
        <CardHeader>
          <CardTitle>
            <Typography variant="h6" as="h1" className="text-center">
              Atur Ulang Kata Sandi
            </Typography>
          </CardTitle>
          <CardDescription>
            <Typography variant="b3" as="h2" className="text-center">
              Masukkan Kata Sandi baru
            </Typography>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={resetPasswordForm.handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="password">Kata Sandi</Label>
                <Controller
                  name="password"
                  control={resetPasswordForm.control}
                  render={({ field, formState: { errors } }) => (
                    <Input
                      id="password"
                      type="password"
                      required
                      error={errors.password?.message}
                      className={cn(errors.password && "border-destructive")}
                      {...field}
                    />
                  )}
                />
              </div>
              <div className="grid gap-3">
                <Label htmlFor="password_confirmation">
                  Konfirmasi Kata Sandi
                </Label>
                <Controller
                  name="password_confirmation"
                  control={resetPasswordForm.control}
                  render={({ field, formState: { errors } }) => (
                    <Input
                      id="password_confirmation"
                      type="password"
                      required
                      error={errors.password_confirmation?.message}
                      className={cn(
                        errors.password_confirmation && "border-destructive"
                      )}
                      {...field}
                    />
                  )}
                />
              </div>
              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full" isLoading={isPending}>
                  Atur Ulang Kata Sandi
                </Button>
              </div>
            </div>
            <div className="mt-4 text-center text-sm">
              Sudah ingat kata sandi?{" "}
              <Link href="/login" className="underline underline-offset-4">
                Masuk
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
