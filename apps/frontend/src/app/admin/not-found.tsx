"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Home, Gamepad2 } from "lucide-react";
import Typography from "@/components/typography";

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground px-4 text-center">
      {/* Icon */}
      <div className="flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 text-primary mb-6">
        <Gamepad2 className="w-10 h-10" />
      </div>

      {/* Text */}
      <Typography variant="h1" as="h1" className="text-4xl font-bold text-primary mb-2">404</Typography>
      <Typography variant="b1" as="p" className="text-lg text-muted-foreground mb-6">
        Halaman yang kamu cari tidak ditemukan.
        <br />
        Mungkin sudah dihapus atau alamatnya salah.
      </Typography>

      {/* CTA */}
      <div className="flex flex-col items-center sm:flex-row gap-3 w-full sm:w-auto">
        <Button asChild variant="default" className="w-full sm:w-auto bg-primary text-primary-foreground hover:bg-primary/90">
          <Link href="/" className="w-full sm:w-auto flex items-center gap-2">
            <Home className="w-4 h-4" />
            Kembali ke Beranda
          </Link>
        </Button>

        <Button asChild variant="outline" className="w-full sm:w-auto border-primary text-primary hover:bg-primary/10 hover:text-primary" >
          <Link href="/articles" className="w-full sm:w-auto flex items-center gap-2">
            <Gamepad2 className="w-4 h-4" />
            Lihat Artikel & Promo
          </Link>
        </Button>
      </div>
    </div>
  );
}
