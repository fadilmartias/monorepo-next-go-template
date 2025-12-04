"use client";

import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Typography from "@/components/typography";
import Logo from "@/components/logo";

import { motion } from "framer-motion";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { apiClient } from "@/lib/api-client";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";
import { useRouter } from "nextjs-toploader/app";
import Link from "next/link";
import { appConfig } from "@/config/app";

export default function TwoFaForm({ token }: { token: string }) {
  const router = useRouter();
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

  const submit2FA = async (formData: string) => {
    const res = await apiClient.post("/v1/auth/login-2fa", {
      temp_token: token,
      code: formData,
    });
    return res.data.data;
  };

  const { mutateAsync: submit2FAAsync, isPending: submit2FAPending } =
    useMutation({
      mutationFn: submit2FA,
      mutationKey: ["submit-2fa"],
      onSuccess: (data) => {
        useAuthStore.setState({ user: data });
        toast(`Selamat datang kembali, ${data.name}`, { duration: 5000 });

        if (data.role === "admin") {
          router.push("/admin/dashboard");
        } else {
          router.push("/");
        }
        router.refresh();
      },
      onError: (e: any) => {
        console.error(e);
        toast.error(
          e.response.data.message ||
            "Gagal verifikasi kode autentikasi 2 langkah"
        );
        if (
          e.response.data.error_code === "ERR_TOKEN_EXPIRED" ||
          e.response.data.error_code === "ERR_TOKEN_INVALID"
        ) {
          router.push("/login");
        }
      },
    });

  const handleSubmit = (val: string) => {
    submit2FAAsync(val);
  };
  return (
    <motion.div
      className={cn("flex justify-center flex-col gap-6")}
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
                Autentikasi 2 Langkah
              </Typography>
            </CardTitle>
            <CardDescription>
              <Typography variant="b3" as="h2" className="text-center">
                Masukkan 6 digit kode yang ada di aplikasi autentikasi kamu
              </Typography>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-6 justify-center items-center">
              <InputOTP
                maxLength={6}
                onChange={(val) => {
                  if (val.length === 6) {
                    handleSubmit(val);
                  }
                }}
              >
                <InputOTPGroup className="flex gap-2 justify-center">
                  {Array.from({ length: 6 }).map((_, i) => (
                    <InputOTPSlot key={i} index={i} />
                  ))}
                </InputOTPGroup>
              </InputOTP>
            </div>
          </CardContent>
        </Card>
        <motion.div
          variants={{
            hidden: { opacity: 0, y: 15 },
            visible: { opacity: 1, y: 0 },
          }}
          className="mt-3 space-y-1"
        >
          <Typography variant="b3" as="h2" className="text-center">
            <Link href={`https://wa.me/${appConfig.whatsappNumber}?text=Saya kehilangan akses kode autentikasi 2 langkah`} target="_blank" rel="noopener noreferrer" className="underline underline-offset-4">
              Kehilangan akses kode?
            </Link>
          </Typography>
          <Typography variant="b3" as="h2" className="text-center">
            <Link href="/login" className="underline underline-offset-4">
              Logout
            </Link>
          </Typography>
        </motion.div>
      </motion.div>
    </motion.div>
  );
}
