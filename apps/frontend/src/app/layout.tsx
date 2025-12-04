import type { Metadata } from "next";
import { Geist, Geist_Mono, Open_Sans} from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import NextTopLoader from "nextjs-toploader";
import { disableConsoleInProduction } from "@/lib/disable-console";
disableConsoleInProduction();
import { ReactQueryProvider } from "@/components/providers/ReactQueryProviders";
import NetworkStatus from "@/components/network-status";
import Script from "next/script";
import GoogleAnalyticsProvider from "./providers/GoogleAnalyticsProvider";
import { Suspense } from "react";
import { Loader2 } from "lucide-react";
import { connection } from "next/server";

const openSans = Open_Sans({
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: {
    default: `${
      process.env.NEXT_PUBLIC_APP_NAME || "Next Template"
    } - Top Up Game Murah, Instan, dan 100% Aman 24 Jam`,
    template: `%s | ${
      process.env.NEXT_PUBLIC_APP_NAME || "Next Template"
    } - Top Up Game Murah, Instan, dan 100% Aman 24 Jam`,
  },
  description:
    "Next Template adalah platform top up game online, voucher digital termurah di Indonesia. Nikmati proses instan, harga bersahabat, dan layanan 24 jam untuk Mobile Legends, Free Fire, Genshin Impact, Honkai Star Rail, dan banyak lagi.",
  keywords: [
    "next template",
    "dilz top up",
    "nexttemplate",
    "top up game murah",
    "top up game instan",
    "top up game terpercaya",
    "top up game aman",
    "top up game 24 jam",
    "voucher game resmi",
    "top up mobile legends murah",
    "top up ml instan",
    "isi diamond mobile legends cepat",
    "top up free fire murah",
    "top up ff instan",
    "voucher free fire resmi",
    "top up genshin impact murah",
    "top up genshin impact instan",
    "voucher genshin impact resmi",
    "top up honkai star rail murah",
    "top up honkai star rail instan",
    "voucher honkai star rail resmi",
    "top up valorant murah",
    "top up valorant cepat",
    "voucher valorant resmi",
    "top up pubg mobile murah",
    "top up pubg mobile instan",
    "voucher pubg mobile resmi",
    "top up call of duty mobile murah",
    "top up cod mobile cepat",
    "voucher cod mobile resmi",
    "beli voucher game online murah",
    "top up game legal indonesia",
    "top up game resmi indonesia",
    "platform top up game terbaik 2025",
    "promo top up game Indonesia",
    "diskon top up game",
    "metode pembayaran top up game lengkap",
    "top up game tanpa login akun",
    "top up game global server indonesia",
  ],
  robots: {
    index: process.env.NEXT_PUBLIC_APP_ENV === "production",
    follow: process.env.NEXT_PUBLIC_APP_ENV === "production",
  },
  authors: [{ name: "M. Fadil Martias", url: "https://fadilmartias.my.id" }],
  creator: "M. Fadil Martias",
  publisher: "DilZ Space Digital",
  openGraph: {
    title: `${
      process.env.NEXT_PUBLIC_APP_NAME || "Next Template"
    } - Top Up Game Murah, Instan, dan 100% Aman 24 Jam`,
    description:
      "Top up game termurah dan tercepat di Indonesia. Proses otomatis 24 jam, harga bersahabat, dan dijamin aman di Next Template. Tersedia MLBB, FF, Genshin Impact, HOK, dan lainnya.",
    url: process.env.NEXT_PUBLIC_APP_URL || "https://nexttemplate.com",
    siteName: process.env.NEXT_PUBLIC_APP_NAME || "Next Template",
    images: [
      {
        url:
          process.env.NEXT_PUBLIC_APP_OG_IMAGE ||
          `${process.env.NEXT_PUBLIC_APP_BASE_URL}/assets/images/logos/logo-nexttemplate.jpg`,
        width: 800,
        height: 800,
        alt: "Next Template - Top Up Game Murah dan Aman",
      },
    ],
    locale: "id_ID",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: `${
      process.env.NEXT_PUBLIC_APP_NAME || "Next Template"
    } - Top Up Game Murah & Aman 24 Jam`,
    description:
      "Top up game favoritmu dengan cepat, murah, dan aman di Next Template. Proses instan dan tersedia 24 jam setiap hari!",
    images: [
      process.env.NEXT_PUBLIC_APP_OG_IMAGE ||
        `${process.env.NEXT_PUBLIC_APP_BASE_URL}/assets/images/logos/logo-nexttemplate.jpg`,
    ],
    creator: "@FadilMartias",
  },
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_APP_URL || "https://nexttemplate.com"
  ),
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  await connection();
  return (
    <html
      lang="id"
      suppressHydrationWarning
      className={`scroll-smooth ${openSans.className} antialiased`}
      data-scroll-behavior="smooth"
    >
      <head>
        {/* Apple Touch Icons */}
        <link
          rel="apple-touch-icon"
          sizes="57x57"
          href="/assets/images/icons/apple-icon-57x57.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="60x60"
          href="/assets/images/icons/apple-icon-60x60.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="72x72"
          href="/assets/images/icons/apple-icon-72x72.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="76x76"
          href="/assets/images/icons/apple-icon-76x76.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="114x114"
          href="/assets/images/icons/apple-icon-114x114.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="120x120"
          href="/assets/images/icons/apple-icon-120x120.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="144x144"
          href="/assets/images/icons/apple-icon-144x144.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="152x152"
          href="/assets/images/icons/apple-icon-152x152.png"
        />
        <link
          rel="apple-touch-icon"
          sizes="180x180"
          href="/assets/images/icons/apple-icon-180x180.png"
        />

        {/* Favicon for Browsers */}
        <link
          rel="icon"
          type="image/png"
          sizes="192x192"
          href="/assets/images/icons/android-icon-192x192.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="32x32"
          href="/assets/images/icons/favicon-32x32.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="96x96"
          href="/assets/images/icons/favicon-96x96.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="16x16"
          href="/assets/images/icons/favicon-16x16.png"
        />

        {/* Web App Manifest */}
        <link rel="manifest" href="/assets/images/icons/manifest.json" />

        {/* Microsoft Windows Tiles */}
        <meta
          name="msapplication-TileColor"
          content="#ffffff"
          media="(prefers-color-scheme: light)"
        />
        <meta
          name="msapplication-TileColor"
          content="#000000"
          media="(prefers-color-scheme: dark)"
        />
        <meta
          name="msapplication-TileImage"
          content="/assets/images/icons/ms-icon-144x144.png"
        />

        <meta
          name="theme-color"
          content="#ffffff"
          media="(prefers-color-scheme: light)"
        />
        <meta
          name="theme-color"
          content="#000000"
          media="(prefers-color-scheme: dark)"
        />
        {process.env.NEXT_PUBLIC_APP_ENV !== "local" && (
          <Script
            src="https://accounts.google.com/gsi/client"
            strategy="afterInteractive"
            async
            defer
          />
        )}

        {process.env.NEXT_PUBLIC_APP_ENV === "production" && (
          <>
            <Script
              src="https://www.googletagmanager.com/gtag/js?id=G-CBN60MXFVF"
              strategy="afterInteractive"
            />
            <Script id="google-analytics" strategy="afterInteractive">
              {`
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', 'G-CBN60MXFVF', {
              page_path: window.location.pathname,
            });
          `}
            </Script>
          </>
        )}

        {/* <BrevoChatbot/> */}
      </head>
      <body
      >
        <NextTopLoader color="oklch(0.72 0.23 55)" showSpinner={false} />
        <NetworkStatus />
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          // enableSystem
          disableTransitionOnChange
        >
          {process.env.NEXT_PUBLIC_APP_ENV === "production" && (
            <Suspense fallback={""}>
              <GoogleAnalyticsProvider />
            </Suspense>
          )}
          <ReactQueryProvider>{children}</ReactQueryProvider>
          <Toaster position="top-right" duration={5000} />
        </ThemeProvider>
        {/* <TawkToChatbot/> */}
      </body>
    </html>
  );
}
