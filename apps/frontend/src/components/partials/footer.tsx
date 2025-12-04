import Link from "next/link";
import { appConfig } from "@/config/app";
import Typography from "../typography";
import {
  Mail,
  Phone,
  MapPin,
  Clock,
  Shield,
  Award,
  Facebook,
  Instagram,
  Twitter,
  Youtube,
  CheckCircleIcon,
} from "lucide-react";
import Logo from "../logo";
import { getApiImage } from "@/lib/api-asset";
import Image from "next/image";

export default async function Footer() {
  const socialLinks = [
    // { name: "Facebook", icon: Facebook, href: "#" },
    {
      name: "Instagram",
      icon: Instagram,
      href: "https://instagram.com/next_template",
    },
    // { name: "Twitter", icon: Twitter, href: "#" },
    // { name: "Youtube", icon: Youtube, href: "#" },
  ];

  return (
    <footer className="bg-linear-to-b from-background to-muted/30 border-t border-border">
      {/* Trust Bar */}
      <div className="bg-primary/5 border-b border-primary/10">
        <div className="container-custom py-6">
          <div className="grid grid-cols-2 md:grid-cols-4 justify-center gap-6 text-center">
            {/* Aman */}
            <div className="flex flex-col items-center gap-2">
              <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                <Shield className="w-6 h-6 text-primary" />
              </div>
              <Typography variant="b2" className="font-semibold">
                100% Aman
              </Typography>
              <Typography variant="c1" className="text-muted-foreground">
                Transaksi terenkripsi & bebas penipuan.
              </Typography>
            </div>

            {/* Cepat */}
            <div className="flex flex-col items-center gap-2">
              <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                <Clock className="w-6 h-6 text-primary" />
              </div>
              <Typography variant="b2" className="font-semibold">
                Proses Super Cepat
              </Typography>
              <Typography variant="c1" className="text-muted-foreground">
                Topup langsung masuk tanpa ribet.
              </Typography>
            </div>

            {/* Terpercaya */}
            <div className="flex flex-col items-center gap-2">
              <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                <Award className="w-6 h-6 text-primary" />
              </div>
              <Typography variant="b2" className="font-semibold">
                Terpercaya
              </Typography>
              <Typography variant="c1" className="text-muted-foreground">
                Dipilih banyak gamer karena cepat & aman.
              </Typography>
            </div>

            {/* Support */}
            <div className="flex flex-col items-center gap-2">
              <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                <Phone className="w-6 h-6 text-primary" />
              </div>
              <Typography variant="b2" className="font-semibold">
                Support 24/7
              </Typography>
              <Typography variant="c1" className="text-muted-foreground">
                Tim kami siap bantu kapan pun kamu butuh.
              </Typography>
            </div>
          </div>
        </div>
      </div>

      {/* Main Footer Content */}
      <div className="container-custom py-12">
        <div className="grid grid-cols-1 min-[460px]:grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-12">
          {/* Kolom 1 - Brand & Description */}
          <div className="lg:col-span-1">
            <Logo size="3xl" className="mb-6" />
            <Typography
              variant="b3"
              className="text-muted-foreground mb-6 leading-relaxed"
            >
              Platform top-up game terpercaya di Indonesia. Kami menyediakan
              layanan cepat, aman, dan harga terjangkau untuk semua kebutuhan
              gaming kamu.
            </Typography>

            {/* Social Media */}
            <div className="flex gap-3">
              {socialLinks.map((social) => {
                const Icon = social.icon;
                return (
                  <Link
                    key={social.name}
                    href={social.href}
                    className="w-10 h-10 rounded-full bg-muted hover:bg-primary hover:text-primary-foreground flex items-center justify-center transition-all duration-300 hover:scale-110"
                    aria-label={social.name}
                  >
                    <Icon className="w-5 h-5" />
                  </Link>
                );
              })}
            </div>
          </div>

          {/* Kolom 2 - Navigasi */}
          <div>
            <Typography variant="h6" className="font-semibold mb-4">
              Navigasi
            </Typography>
            <ul className="space-y-3">
              {[
                { label: "Beranda", href: "/" },
                { label: "Cek Transaksi", href: "/check-transaction" },
                { label: "Artikel & Tips", href: "/articles" },
                // { label: "Promo", href: "/promo" },
                // { label: "Tentang Kami", href: "/about" },
              ].map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-muted-foreground hover:text-primary transition-colors inline-flex items-center gap-2 group"
                  >
                    <Typography variant="c1" className="w-1.5 h-1.5 rounded-full bg-primary/40 group-hover:bg-primary transition-colors" />
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Kolom 3 - Bantuan & Legal */}
          <div>
            <Typography variant="h6" className="font-semibold mb-4">
              Bantuan & Legal
            </Typography>
            <ul className="space-y-3">
              {[
                // { label: "FAQ", href: "/faq" },
                // { label: "Cara Pemesanan", href: "/how-to-order" },
                { label: "Kebijakan Privasi", href: "/privacy-policy" },
                { label: "Syarat & Ketentuan", href: "/terms-and-conditions" },
                { label: "Kebijakan Refund", href: "/refund-policy" },
              ].map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-muted-foreground hover:text-primary transition-colors inline-flex items-center gap-2 group"
                  >
                    <Typography variant="c1" className="w-1.5 h-1.5 rounded-full bg-primary/40 group-hover:bg-primary transition-colors" />
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Kolom 4 - Kontak */}
          <div>
            <Typography variant="h6" className="font-semibold mb-4">
              Hubungi Kami
            </Typography>
            <ul className="space-y-4">
              <li className="flex items-start gap-3 text-muted-foreground">
                <Phone className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <div>
                  <Typography variant="b3" className="block mb-1">
                    WhatsApp
                  </Typography>
                  <Link
                    href={`https://wa.me/${appConfig.whatsappNumber}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="hover:text-primary transition-colors font-medium"
                  >
                    +{appConfig.whatsappNumber}
                  </Link>
                </div>
              </li>

              <li className="flex items-start gap-3 text-muted-foreground">
                <Mail className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <div>
                  <Typography variant="b3" className="block mb-1">
                    Email
                  </Typography>
                  <a
                    href="mailto:support@nexttemplate.com"
                    className="hover:text-primary transition-colors font-medium"
                  >
                    support@nexttemplate.com
                  </a>
                </div>
              </li>

              <li className="flex items-start gap-3 text-muted-foreground">
                <MapPin className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <div>
                  <Typography variant="b3" className="block mb-1">
                    Alamat
                  </Typography>
                  <Typography variant="b3">Riau, Indonesia</Typography>
                </div>
              </li>

              <li className="flex items-start gap-3 text-muted-foreground">
                <Clock className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <div>
                  <Typography variant="b3" className="block mb-1">
                    Jam Operasional
                  </Typography>
                  <Typography variant="b3">24/7 (Otomatis)</Typography>
                </div>
              </li>
            </ul>
          </div>
        </div>

        {/* Payment Methods */}
        <div className="mt-12 pt-8 border-t border-border">
          <Typography variant="b2" className="font-semibold mb-4 text-center">
            Metode Pembayaran
          </Typography>
          <div className="flex flex-wrap justify-center items-center gap-3 py-2">
            
          </div>

          {/* Trust Badges */}
          <div className="flex flex-wrap justify-center items-center gap-6 mt-6">
            <div className="flex items-center gap-2 text-muted-foreground">
              <Shield className="w-5 h-5 text-green-500" />
              <Typography variant="c1">SSL Encrypted</Typography>
            </div>
            <div className="flex items-center gap-2 text-muted-foreground">
              <CheckCircleIcon className="w-5 h-5 text-blue-500" />
              <Typography variant="c1">Payment Verified</Typography>
            </div>
            <div className="flex items-center gap-2 text-muted-foreground">
              <Award className="w-5 h-5 text-yellow-500" />
              <Typography variant="c1">Trusted Partner</Typography>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom Bar */}
      <div className="border-t border-border bg-muted/20">
        <div className="container-custom py-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-center md:text-left">
            <Typography variant="b3" className="text-muted-foreground">
              © {new Date().getFullYear()} PT. DilZ Space Digital. Hak Cipta
              Dilindungi Undang-Undang.
            </Typography>
            <Typography variant="b3" className="text-muted-foreground">
              Made with ❤️ in Indonesia
            </Typography>
          </div>
        </div>
      </div>
    </footer>
  );
}
