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
import { RadioGroup, RadioGroupItem } from "../ui/radio-group";
import { Input } from "../ui/input";
import { Card } from "../ui/card";
import { Label } from "../ui/label";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "../ui/input-otp";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "nextjs-toploader/app";
import { FaEnvelope, FaRegEnvelope, FaWhatsapp } from "react-icons/fa";

interface OTPModalProps {
  isOpen: boolean;
  onClose: () => void;
  openModal2FA: () => void;
  is2FAEnabled: boolean | undefined;
}

interface OTPModalResponse {
  totp_secret: string;
  qr_string: string;
}

type DeliveryMethodOption = {
  value: string;
  label: string;
  icon: React.ReactNode;
  disabled: boolean;
};

const DELIVERY_METHODS: DeliveryMethodOption[] = [
  {
    value: "email",
    label: "Email",
    icon: <FaRegEnvelope size={18} />,
    disabled: true,
  },
  {
    value: "wa",
    label: "WhatsApp",
    icon: <FaWhatsapp size={20} />,
    disabled: false,
  },
];

export default function OTPModal({
  isOpen,
  onClose,
  openModal2FA,
  is2FAEnabled,
}: OTPModalProps) {
  const [step, setStep] = useState("select"); // select, otp
  const [deliveryMethod, setDeliveryMethod] = useState("");

  useEffect(() => {
    if (step === "loading") {
      const timer = setTimeout(() => setStep("otp"), 2000);
      return () => clearTimeout(timer);
    }
  }, [step]);

  const sendOTP = async () => {
    const res = await apiClient.post("/v1/auth/send-otp", {
      target_type: deliveryMethod,
      subject:
        "Verifikasi identitas Anda untuk " +
        (is2FAEnabled ? "menonaktifkan" : "mengaktifkan") +
        " autentikasi dua faktor (2FA)",
    });
    return res.data;
  };

  const verifyOTP = async (otp: string) => {
    const res = await apiClient.post("/v1/auth/verify-otp", {
      target_type: deliveryMethod,
      subject:
        "Verifikasi identitas Anda untuk " +
        (is2FAEnabled ? "menonaktifkan" : "mengaktifkan") +
        " autentikasi dua faktor (2FA)",
      otp,
    });
    return res.data;
  };

  const disable2FA = async () => {
    const res = await apiClient.post("/v1/auth/disable-2fa");
    return res.data;
  };

  const { mutateAsync: disable2FAAsync, isPending: disable2FAPending } =
    useMutation({
      mutationFn: disable2FA,
      mutationKey: ["disable-2fa"],
      onSuccess: () => {
        toast.success("Autentikasi 2 Langkah berhasil dinonaktifkan");
        onClose();
        window.location.reload();
      },
      onError: (e: any) => {
        toast.error(e.response.data.message || "Gagal dinonaktifkan");
      },
    });

  const { mutateAsync: verifyOTPAsync, isPending: verifyOTPPending } =
    useMutation({
      mutationFn: verifyOTP,
      mutationKey: ["verify-otp"],
      onSuccess: () => {
        if (!is2FAEnabled) {
          toast.success("OTP berhasil diverifikasi");
          openModal2FA();
        } else {
          disable2FAAsync();
        }
        onClose();
      },
      onError: (e: any) => {
        toast.error(e.response.data.message || "Gagal verifikasi OTP");
      },
    });

  const { mutateAsync, isPending } = useMutation({
    mutationFn: sendOTP,
    mutationKey: ["send-otp"],
    onSuccess: () => {
      toast.success("OTP berhasil dikirim");
      setStep("otp");
    },
    onError: () => {
      toast.error("Gagal mengirim OTP");
    },
  });

  const handleSubmit = (e?: React.FormEvent) => {
    e?.preventDefault();
    mutateAsync();
  };

  const handleOtpSubmit = (val: string) => {
    verifyOTPAsync(val);
  };

  if (!isOpen) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent
        className="sm:max-w-[425px]"
        onInteractOutside={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <div
          className={`overlay absolute top-0 left-0 w-full h-full bg-white z-[2] rounded-2xl transition-opacity duration-300 ${
            isPending || verifyOTPPending || disable2FAPending
              ? "opacity-70 pointer-events-auto "
              : "opacity-0 pointer-events-none"
          }`}
        ></div>
        <DialogHeader>
          {step === "select" ? (
            <>
              <DialogTitle>Verifikasi Identitas Anda</DialogTitle>
              <DialogDescription>
                Kami perlu memverifikasi identitas kamu untuk{" "}
                {is2FAEnabled ? "menonaktifkan" : "mengaktifkan"} autentikasi
                dua faktor (2FA). Silakan pilih metode pengiriman kode OTP yang
                kamu inginkan.
              </DialogDescription>
            </>
          ) : step === "otp" ? (
            <>
              <DialogTitle>Masukkan Kode OTP</DialogTitle>
              <DialogDescription>
                Kami telah mengirimkan kode OTP ke{" "}
                {deliveryMethod === "email" ? "email" : "WhatsApp"} kamu.
              </DialogDescription>
            </>
          ) : null}
        </DialogHeader>
        {step === "select" && (
          <form onSubmit={handleSubmit}>
            <RadioGroup
              value={deliveryMethod}
              onValueChange={setDeliveryMethod}
            >
              {DELIVERY_METHODS.map((method) => (
                <Label
                  key={method.value}
                  htmlFor={method.value}
                  className={`cursor-pointer w-full ${
                    method.disabled ? "opacity-50 cursor-not-allowed" : ""
                  }`}
                >
                  <Card
                    className={`p-4  rounded-xl border-2 transition w-full ${
                      deliveryMethod === method.value
                        ? "border-primary shadow-md"
                        : "border-muted"
                    }`}
                  >
                    <div className="flex items-center space-x-3">
                      <RadioGroupItem
                        id={method.value}
                        value={method.value}
                        className="hidden"
                        disabled={method.disabled}
                      />
                      <span className="font-medium flex gap-2 items-center">
                        {method.icon} {method.label}
                      </span>
                    </div>
                  </Card>
                </Label>
              ))}
            </RadioGroup>
            <DialogFooter className="mt-6">
              <Button type="submit" disabled={!deliveryMethod}>Kirim Kode OTP</Button>
            </DialogFooter>
          </form>
        )}

        {step === "otp" && (
          <div className="flex flex-col gap-6 justify-center items-center">
            <InputOTP
              maxLength={6}
              onChange={(val) => {
                if (val.length === 6) {
                  handleOtpSubmit(val);
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
