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
import { useState } from "react";
import { useRouter } from "nextjs-toploader/app";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import Typography from "@/components/typography";
import Logo from "@/components/logo";
import { OAuthButtons } from "@/components/buttons/oauth-buttons";
import { motion, Variants } from "framer-motion";
import useReCaptcha from "@/hooks/use-recaptcha";

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

const registerSchema = z
  .object({
    name: z.string().min(3, "Nama minimal 3 karakter"),
    email: z.string().email("Email tidak valid"),
    phone: z.string().min(10, "No WA minimal 10 karakter"),
    referral_code: z
      .string()
      .optional()
      .refine((code) => code?.length === 8, {
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

type RegisterFormValues = z.infer<typeof registerSchema>;

const ACTION: "LOGIN" | "REGISTER" = "REGISTER";

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [isPending, setIsPending] = useState(false);
  const registerForm = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      email: "",
      phone: "",
      referral_code: "",
      password: "",
      password_confirmation: "",
    },
  });
  const router = useRouter();
  const { getToken, verifyRecaptcha } = useReCaptcha();

  const onSubmit = async (formData: RegisterFormValues) => {
    registerForm.clearErrors();
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
      const res = await apiClient.post("/v1/auth/register", formData);
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

  const handleOAuth = (
    provider: "google" | "discord" | "facebook" | "steam" | "twitch"
  ) => {
    window.location.href =
      process.env.NEXT_PUBLIC_API_BASE_URL + `/v1/auth/${provider}/redirect`;
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
                Daftar
              </Typography>
            </CardTitle>
            <CardDescription>
              <Typography variant="b3" as="h2" className="text-center">
                Masukkan data kamu untuk mendaftar
              </Typography>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className="space-y-4"
              onSubmit={registerForm.handleSubmit(onSubmit)}
            >
              <motion.div
                variants={childVariants}
                className="flex flex-col gap-4"
              >
                <OAuthButtons
                  config={{
                    google: true,
                    discord: true,
                    steam: false,
                    twitch: false,
                    facebook: false,
                  }}
                  handleOAuth={handleOAuth}
                />
              </motion.div>

              <motion.div
                variants={childVariants}
                className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t"
              >
                <span className="bg-card text-muted-foreground relative z-10 px-2">
                  Atau
                </span>
              </motion.div>

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
                      name={field as keyof RegisterFormValues}
                      control={registerForm.control}
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
                          required
                          error={
                            errors[field as keyof RegisterFormValues]?.message
                          }
                          className={cn(
                            errors[field as keyof RegisterFormValues] &&
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
                  Daftar
                </Button>
              </motion.div>

              <motion.div
                variants={childVariants}
                className="mt-4 text-center text-sm"
              >
                Sudah punya akun?{" "}
                <Link href="/login" className="underline underline-offset-4">
                  Masuk
                </Link>
              </motion.div>
            </form>
          </CardContent>
        </Card>
      </motion.div>

      <motion.div variants={childVariants}>
        <Typography variant="b3" as="p" className="text-center text-sm">
          Dengan mengklik Daftar, Anda menyetujui{" "}
          <Link
            href="/terms-and-conditions"
            target="_blank"
            rel="noopener noreferrer"
            className="underline underline-offset-4"
          >
            Syarat dan Ketentuan
          </Link>{" "}
          dan{" "}
          <Link
            href="/privacy-policy"
            target="_blank"
            rel="noopener noreferrer"
            className="underline underline-offset-4"
          >
            Kebijakan Privasi
          </Link>
          .
        </Typography>
      </motion.div>
    </motion.div>
  );
}
