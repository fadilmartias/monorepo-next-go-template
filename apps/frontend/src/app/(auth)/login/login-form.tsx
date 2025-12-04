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
import { toast } from "sonner";
import { useRouter } from "nextjs-toploader/app";
import { useState } from "react";
import { apiClient } from "@/lib/api-client";
import Typography from "@/components/typography";
import Logo from "@/components/logo";
import { useAuthStore } from "@/stores/auth-store";
import { OAuthButtons } from "@/components/buttons/oauth-buttons";
import { motion } from "framer-motion";
import useReCaptcha from "@/hooks/use-recaptcha";
const loginSchema = z.object({
  credential: z.string().email("Email tidak valid"),
  password: z.string().min(8, "Kata Sandi minimal 8 karakter"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

const ACTION: "LOGIN" | "REGISTER" = "LOGIN";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [isPending, setIsPending] = useState(false);
  const loginForm = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      credential: "",
      password: "",
    },
  });
  const router = useRouter();
  const { getToken, verifyRecaptcha } = useReCaptcha();
  const onSubmit = async (formData: LoginFormValues) => {
    loginForm.clearErrors();
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
      // const user = await loginAction(formData); // panggil server action
      const res = await apiClient.post("/v1/auth/login", formData);
      const data = res.data.data;
      if (data.is_2fa_enabled) {
        router.push(`/login/2fa/${data.temp_token}`);
      }
      useAuthStore.setState({ user: data, isLoggedIn: true });
      toast.success(`Selamat datang kembali, ${data.name}`, { duration: 5000 });

      if (data.role === "admin") {
        router.push("/admin/dashboard");
      } else {
        router.push("/");
      }
      router.refresh();
    } catch (err: any) {
      console.error("ERR", err);
      loginForm.setError("credential", {
        type: "manual",
        message: err.message,
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

  // Variants terpisah
  const containerVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        staggerChildren: 0.1, // animasi children muncul bergantian
      },
    },
  };

  const childVariants = {
    hidden: { opacity: 0, y: 10 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <motion.div
      className={cn("flex justify-center flex-col gap-6", className)}
      initial="hidden"
      animate="visible"
      variants={containerVariants}
    >
      {/* Logo */}
      <motion.div className="flex justify-center" variants={childVariants}>
        <Logo size="2xl" />
      </motion.div>

      {/* Card */}
      <motion.div
        variants={{
          hidden: { opacity: 0, y: 15 },
          visible: { opacity: 1, y: 0 },
        }}
      >
        <Card>
          <CardHeader>
            <CardTitle>
              <Typography variant="h6" as="h1" className="text-center">
                Selamat Datang
              </Typography>
            </CardTitle>
            <CardDescription>
              <Typography variant="b3" as="h2" className="text-center">
                Masuk untuk melanjutkan
              </Typography>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className="space-y-4"
              onSubmit={loginForm.handleSubmit(onSubmit)}
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

              <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                <span className="bg-card text-muted-foreground relative z-10 px-2">
                  Atau
                </span>
              </div>

              <motion.div
                className="flex flex-col gap-6"
                variants={{
                  hidden: { opacity: 0, y: 15 },
                  visible: { opacity: 1, y: 0 },
                }}
              >
                <div className="grid gap-3">
                  <Label htmlFor="credential">Email</Label>
                  <Controller
                    name="credential"
                    control={loginForm.control}
                    render={({ field, formState: { errors } }) => (
                      <Input
                        id="credential"
                        type="email"
                        placeholder="user@nexttemplate.com"
                        required
                        error={errors.credential?.message}
                        className={cn(
                          errors.credential && "border-destructive"
                        )}
                        {...field}
                      />
                    )}
                  />
                </div>
                <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Kata Sandi</Label>
                    <Link
                      href="/forgot-password"
                      tabIndex={-1}
                      className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                    >
                      Lupa Kata Sandi?
                    </Link>
                  </div>
                  <Controller
                    name="password"
                    control={loginForm.control}
                    render={({ field, formState: { errors } }) => (
                      <Input
                        id="password"
                        type="password"
                        placeholder="********"
                        required
                        error={errors.password?.message}
                        className={cn(errors.password && "border-destructive")}
                        {...field}
                      />
                    )}
                  />
                </div>
                <div className="flex flex-col gap-3">
                  <Button
                    type="submit"
                    className="w-full"
                    isLoading={isPending}
                  >
                    Masuk
                  </Button>
                </div>
              </motion.div>

              <motion.div
                className="mt-4 text-center text-sm"
                variants={{
                  hidden: { opacity: 0, y: 15 },
                  visible: { opacity: 1, y: 0 },
                }}
              >
                Belum punya akun?{" "}
                <Link href="/register" className="underline underline-offset-4">
                  Daftar
                </Link>
              </motion.div>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </motion.div>
  );
}
