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
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "nextjs-toploader/app";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import Typography from "@/components/typography";
import Logo from "@/components/logo";
import { OAuthButtons } from "@/components/buttons/oauth-buttons";
import { motion, Variants } from "framer-motion";
import useReCaptcha from "@/hooks/use-recaptcha";
import { User } from "@/stores/auth-store";
import { useQuery } from "@tanstack/react-query";

const containerVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.08 },
  },
};

const childVariants: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

const onboardingSchema = z
  .object({
    name: z.string().min(3, "Nama minimal 3 karakter"),
    email: z.string().email("Email tidak valid"),
    phone: z.string().min(10, "No WA minimal 10 karakter"),
    referral_code: z
      .string()
      .optional()
      .refine((code) => !code || code?.length === 8, {
        message: "Kode Referral harus 8 karakter",
      }),
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

type OnboardingFormValues = z.infer<typeof onboardingSchema>;

const ACTION = "ONBOARDING";

export function OnboardingForm({
  className,
  id_user,
  ...props
}: React.ComponentProps<"div"> & { id_user: string }) {


  const { data: user, isSuccess } = useQuery({
    queryKey: ["user", id_user],
    queryFn: () => apiClient.get(`/v0/users/${id_user}`).then((res) => res.data),
    enabled: !!id_user,
  });

  useEffect(() => {
    if (isSuccess && user) {
      onboardingForm.reset(user?.data);
    }
  }, [isSuccess, user]);

  const [isPending, setIsPending] = useState(false);
  const onboardingForm = useForm<OnboardingFormValues>({
    resolver: zodResolver(onboardingSchema),
    defaultValues: {
      name: user?.data.name || "",
      email: user?.data.email || "",
      phone: user?.data.phone || "",
      referral_code: "",
      password: "",
      password_confirmation: "",
    },
  });
  const router = useRouter();
  const { getToken, verifyRecaptcha } = useReCaptcha();

  const onSubmit = async (formData: OnboardingFormValues) => {
    onboardingForm.clearErrors();
    try {
      setIsPending(true);
      const captchaToken = await getToken(ACTION);
      if (!captchaToken) {
        toast.error("Error validasi captcha");
        return;
      }
      const verifyRecaptchaResult = await verifyRecaptcha(captchaToken, ACTION);
      if (!verifyRecaptchaResult.success) {
        toast.error(verifyRecaptchaResult.message);
        return;
      }
      const res = await apiClient.post("/v1/auth/onboarding", formData);
      toast.success(`Pendaftaran berhasil, silakan login`, { duration: 5000 });
      router.push("/login");
    } catch (err: any) {
      toast.error(err.message || "Gagal mendaftar", {
        duration: 5000,
      });
    } finally {
      setIsPending(false);
    }
  };

  return (
    <motion.div
      className={cn("flex flex-col gap-6", className)}
      initial="hidden"
      animate="visible"
      variants={containerVariants}
    >
      <motion.div variants={childVariants} className="flex justify-center">
        <Logo size="2xl" />
      </motion.div>

      <motion.div variants={childVariants}>
        <Card>
          <CardHeader>
            <CardTitle>
              <Typography variant="h6" as="h1" className="text-center">
                Onboarding
              </Typography>
            </CardTitle>
            <CardDescription>
              <Typography variant="b3" as="h2" className="text-center">
                Masukkan data kamu untuk onboarding
              </Typography>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className="space-y-4"
              onSubmit={onboardingForm.handleSubmit(onSubmit)}
            >
              <motion.div
                variants={childVariants}
                className="flex flex-col gap-6"
              >
                {[
                  "name",
                  "phone",
                  "email",
                  "referral_code",
                  "password",
                  "password_confirmation",
                ].map((field) => (
                  <motion.div
                    variants={childVariants}
                    key={field}
                    className="grid gap-3"
                  >
                    <Controller
                      name={field as keyof OnboardingFormValues}
                      control={onboardingForm.control}
                      render={({ field: f, formState: { errors } }) => (
                        <Input
                          id={field}
                          label={
                            field === "name"
                              ? "Nama"
                              : field === "phone"
                              ? "No WA"
                              : field === "email"
                              ? "Email"
                              : field === "password"
                              ? "Kata Sandi"
                              : field === "referral_code"
                              ? "Kode Referral"
                              : "Konfirmasi Kata Sandi"
                          }
                          type={
                            field.includes("password") ? "password" : "text"
                          }
                          placeholder={
                            field === "name"
                              ? "Nama"
                              : field === "phone"
                              ? "No WA"
                              : field === "email"
                              ? "m@example.com"
                              : field === "referral_code"
                              ? "REF-1234"
                              : field.includes("password")
                              ? "********"
                              : ""
                          }
                          {...(field !== "referral_code" && { required: true })}
                          error={
                            errors[field as keyof OnboardingFormValues]?.message
                          }
                          className={cn(
                            errors[field as keyof OnboardingFormValues] &&
                              "border-destructive"
                          )}
                          {...f}
                        />
                      )}
                    />
                  </motion.div>
                ))}
              </motion.div>

              <motion.div
                variants={childVariants}
                className="flex flex-col gap-3"
              >
                <Button type="submit" className="w-full" isLoading={isPending}>
                  Lanjutkan
                </Button>
              </motion.div>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </motion.div>
  );
}
