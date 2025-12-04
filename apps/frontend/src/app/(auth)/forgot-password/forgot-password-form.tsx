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
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Controller } from "react-hook-form";
import Typography from "@/components/typography";
import Logo from "@/components/logo";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { useState } from "react";

const forgotPasswordSchema = z.object({
  email: z.string().email("Email tidak valid"),
});

type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>;

export function ForgotPasswordForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [isPending, setIsPending] = useState(false);
  const forgotPasswordForm = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  });
  const onSubmit = async (formData: ForgotPasswordFormValues) => {
    forgotPasswordForm.clearErrors();
    try {
      setIsPending(true);
      const res = await apiClient.post("/v1/auth/forgot-password", formData);
      toast.success(`Permintaan berhasil, silakan periksa emailmu`, {
        duration: 5000,
      });
    } catch (err: any) {
      console.error("ERR", err);
      toast.error(err.message || "Gagal mengirim permintaan", {
        duration: 5000,
      });
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
              Lupa Kata Sandi
            </Typography>
          </CardTitle>
          <CardDescription>
            <Typography variant="b3" as="h2" className="text-center">
              Masukkan email kamu untuk mengirimkan link atur ulang kata sandi
            </Typography>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={forgotPasswordForm.handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Controller
                  name="email"
                  control={forgotPasswordForm.control}
                  render={({ field, formState: { errors } }) => (
                    <Input
                      id="email"
                      type="email"
                      placeholder="m@example.com"
                      required
                      error={errors.email?.message?.toString()}
                      className={cn(errors.email && "border-destructive")}
                      {...field}
                    />
                  )}
                />
              </div>
              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full" isLoading={isPending}>
                  Kirim
                </Button>
              </div>
            </div>
            <div className="mt-4 text-center text-sm">
              Sudah punya akun?{" "}
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
