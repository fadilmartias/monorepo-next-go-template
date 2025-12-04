"use client";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import QrisViewer from "../../app/(user)/payment/[order_id]/qris-viewer";
import { apiClient } from "@/lib/api-client";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "../ui/input-otp";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "nextjs-toploader/app";

interface TwoFAModalProps {
  isOpen: boolean;
  onClose: () => void;
}

interface TwoFAResponse {
  totp_secret: string;
  qr_string: string;
}

export default function TwoFAModal({ isOpen, onClose }: TwoFAModalProps) {
  const [twoFAData, setTwoFAData] = useState<TwoFAResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState("scan");
  const router = useRouter();

  const verifyCode = async (code: string) => {
    const res = await apiClient.post("/v1/auth/verify-2fa", {
      code,
    });
    return res.data;
  };

  const { mutateAsync: verifyCodeAsync, isPending: verifyCodePending } =
    useMutation({
      mutationFn: verifyCode,
      mutationKey: ["verify-2fa"],
      onSuccess: () => {
        toast.success("Autentikasi 2 Langkah Berhasil Diaktifkan");
        onClose();
        window.location.reload();
      },
      onError: (e: any) => {
        toast.error(e.response.data.message || "Gagal verifikasi kode");
      },
    });

  useEffect(() => {
    if (isOpen) {
      const fetchTwoFAData = async () => {
        setLoading(true);
        try {
          const response = await apiClient.post("/v1/auth/register-2fa");
          const data: TwoFAResponse = response.data.data;
          (data);
          setTwoFAData(data);
        } catch (error) {
          toast.error("Gagal mengambil data 2FA");
        } finally {
          setLoading(false);
        }
      };
      fetchTwoFAData();
    }
  }, [isOpen]);

  const handleCodeSubmit = (val: string) => {
    verifyCodeAsync(val);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent
        className="sm:max-w-[425px]"
        onInteractOutside={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle>Aktifkan Autentikasi 2 Langkah</DialogTitle>
          {step === "scan" ? (
            <DialogDescription>
              Pindai kode QR dibawah dengan aplikasi authenticator (e.g., Google
              Authenticator) untuk mengaktifkan Autentikasi 2 Langkah
            </DialogDescription>
          ) : step === "verify" ? (
            <DialogDescription>
              Masukkan 6 digit kode di aplikasi authenticator kamu untuk
              verifikasi pertama kali
            </DialogDescription>
          ) : null}
        </DialogHeader>

        {step === "scan" && (
          <>
            <div className="flex justify-center py-4">
              {loading ? (
                <p>Memuat QR code...</p>
              ) : twoFAData ? (
                <QrisViewer qrString={twoFAData.qr_string} />
              ) : (
                <p>No QR code available</p>
              )}
            </div>
            {twoFAData && (
              <div className="text-center">
                <p className="text-sm text-muted-foreground">
                  Atau masukkan kode ini secara manual:{" "}
                  <strong>
                    {loading ? "Memuat kode..." : twoFAData.totp_secret}
                  </strong>
                </p>
              </div>
            )}
            <DialogFooter>
              <Button onClick={() => setStep("verify")}>Verifikasi</Button>
            </DialogFooter>
          </>
        )}
        {step === "verify" && (
          <div className="flex flex-col gap-6 justify-center items-center">
            <InputOTP
              maxLength={6}
              onChange={(val) => {
                if (val.length === 6) {
                  handleCodeSubmit(val);
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
        )}
      </DialogContent>
    </Dialog>
  );
}
