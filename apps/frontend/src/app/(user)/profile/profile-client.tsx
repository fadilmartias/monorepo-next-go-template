"use client";

import { CardHeader, CardContent } from "@/components/ui/card";
import { Gamepad2, Trophy, Shield, Zap, Share, Copy } from "lucide-react";
import { Mail, Phone } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useAuthStore } from "@/stores/auth-store";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar } from "@/components/ui/avatar";
import { useState } from "react";
import TwoFAModal from "@/components/modals/two-fa-modal";
import OTPModal from "@/components/modals/otp-modal";
import Typography from "@/components/typography";
import { toast } from "sonner";

export default function ProfileClient() {
  const { user } = useAuthStore();
  const [isOpenModalOTP, setIsOpenModalOTP] = useState(false);
  const [isOpenModal2FA, setIsOpenModal2FA] = useState(false);

  const expPercentage = user ? (user.exp / user.exp_for_level) * 100 : 0;

  return (
    <>
      <CardHeader className="flex flex-col items-center gap-6 pb-8 relative overflow-hidden">
        {/* Decorative background elements */}
        <div className="absolute top-0 left-0 w-full h-full opacity-5">
          <div className="absolute top-4 left-4 w-32 h-32 border-2 border-primary rotate-45" />
          <div className="absolute bottom-4 right-4 w-24 h-24 border-2 border-primary rotate-12" />
        </div>

        {/* Avatar with hexagonal shape effect */}
        <div className="relative">
          <div className="absolute inset-0 bg-linear-to-br from-primary via-accent to-primary opacity-20 blur-xl rounded-full animate-pulse" />
          <Avatar className="h-32 w-32 relative z-10 flex items-center justify-center border-4 border-primary shadow-lg shadow-primary/50 bg-linear-to-br from-card to-card/80">
            <Gamepad2 className="w-16 h-16 text-primary" />
          </Avatar>
          {/* Level badge */}
          <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 z-20 bg-primary text-primary-foreground px-4 py-1 rounded-full text-xs font-bold shadow-lg border-2 border-background">
            LV {user?.level || 0}
          </div>
        </div>

        {/* User info */}
        <div className="flex flex-col gap-2 items-center relative z-10 w-full">
          <h2 className="text-3xl font-extrabold text-card-foreground tracking-wider uppercase bg-linear-to-r from-primary to-accent bg-clip-text text-transparent">
            {!user ? <Skeleton className="w-48 h-9" /> : user.name}
          </h2>

          {/* Title badge */}
          {user?.active_title?.name ? (
            <div className="flex items-center gap-2 px-4 py-1.5 bg-accent/10 border border-accent/30 rounded-lg">
              <Shield className="w-4 h-4 text-accent" />
              <p className="text-sm font-medium text-accent">
                {user.active_title.name}
              </p>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground italic">
              Title belum dipilih
            </p>
          )}

          {/* XP Progress Bar */}
          <div className="w-full max-w-md mt-4 space-y-2">
            <div className="flex justify-between items-center text-xs font-medium">
              <span className="text-muted-foreground flex items-center gap-1">
                <Zap className="w-3 h-3 text-primary" />
                XP Progress
              </span>
              <span className="text-primary">
                {user?.exp || 0} / {user?.exp_for_level || 0}
              </span>
            </div>
            <div className="relative h-3 w-full bg-muted rounded-full overflow-hidden border border-border">
              <div
                className="absolute top-0 left-0 h-full bg-linear-to-r from-primary to-accent transition-all duration-500 ease-out rounded-full"
                style={{ width: `${expPercentage}%` }}
              ></div>
              <div className="absolute w-full inset-0 bg-linear-to-r from-transparent via-white/20 to-transparent animate-shimmer" />
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-6 text-foreground">
        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-4">
          <div className="flex flex-col items-center gap-2 p-4 bg-linear-to-br from-primary/5 to-accent/5 border border-primary/20 rounded-lg hover:border-primary/40 transition-colors">
            <Trophy className="w-6 h-6 text-primary" />
            <Typography variant="b3" className="text-muted-foreground">
              Level
            </Typography>
            <span className="text-2xl font-bold text-primary">
              {user?.level || 0}
            </span>
          </div>
          <div className="flex flex-col items-center gap-2 p-4 bg-linear-to-br from-accent/5 to-primary/5 border border-accent/20 rounded-lg hover:border-accent/40 transition-colors">
            <Share className="w-6 h-6 text-accent" />
            <Typography variant="b3" className="text-muted-foreground">
              Kode Referral
            </Typography>
            <div className="flex gap-3 items-center">
              <span className="text-2xl font-bold text-accent">
                {user?.referral_code}
              </span>
              <Button
                variant="outline"
                className=""
                size="sm"
                onClick={() => {
                  navigator.clipboard.writeText(user?.referral_code || "");
                  toast.success("Kode referral berhasil disalin", {
                    duration: 5000,
                  });
                }}
              >
                <Copy />
              </Button>
            </div>
          </div>
        </div>

        <Separator className="my-6 bg-border" />

        {/* Contact Info */}
        <div className="space-y-3">
          <div className="flex items-center gap-3 p-3 bg-card/50 border border-border rounded-lg hover:border-primary/50 transition-colors">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Mail className="w-5 h-5 text-primary" />
            </div>
            <div className="flex-1">
              <span className="text-xs text-muted-foreground block">Email</span>
              <span className="font-medium text-sm">
                {!user ? <Skeleton className="w-36 h-5" /> : user.email}
              </span>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 bg-card/50 border border-border rounded-lg hover:border-primary/50 transition-colors">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Phone className="w-5 h-5 text-primary" />
            </div>
            <div className="flex-1">
              <span className="text-xs text-muted-foreground block">
                Nomor HP
              </span>
              <span className="font-medium text-sm">
                {!user ? <Skeleton className="w-36 h-5" /> : user.phone}
              </span>
            </div>
          </div>
        </div>

        <Separator className="my-6 bg-border" />

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row justify-center items-center gap-3">
          <Button
            className="w-full sm:w-auto bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/30 transition-all hover:shadow-primary/50"
            asChild
          >
            <Link href="/profile/edit">Edit Profil</Link>
          </Button>
          {user?.is_2fa_enabled ? (
            <Button
              className="w-full sm:w-auto bg-accent text-accent-foreground hover:bg-accent/90 shadow-lg shadow-accent/30 transition-all hover:shadow-accent/50"
              onClick={() => setIsOpenModalOTP(true)}
            >
              Nonaktifkan 2FA
            </Button>
          ) : (
            <Button
              onClick={() => setIsOpenModalOTP(true)}
              disabled={!user?.phone || !user?.email}
              className="w-full sm:w-auto bg-accent text-accent-foreground hover:bg-accent/90 shadow-lg shadow-accent/30 transition-all hover:shadow-accent/50 disabled:opacity-50"
            >
              Aktifkan 2FA
            </Button>
          )}
        </div>
      </CardContent>

      <TwoFAModal
        isOpen={isOpenModal2FA}
        onClose={() => setIsOpenModal2FA(false)}
      />
      <OTPModal
        isOpen={isOpenModalOTP}
        onClose={() => setIsOpenModalOTP(false)}
        openModal2FA={() => setIsOpenModal2FA(true)}
        is2FAEnabled={user?.is_2fa_enabled}
      />

      <style jsx>{`
        @keyframes shimmer {
          0% {
            transform: translateX(-100%);
          }
          100% {
            transform: translateX(100%);
          }
        }
        .animate-shimmer {
          animation: shimmer 2s infinite;
        }
      `}</style>
    </>
  );
}
