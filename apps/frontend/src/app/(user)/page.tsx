import Banner from "@/components/carousels/banner";
import TrendingGame from "@/app/(user)/trending-game";
import {
  getBanners,
  getFlashSales,
  getLatestArticles,
  getTrendingProducts,
} from "@/app/(user)/actions";
import { Metadata } from "next";
import LatestArticles from "./latest-articles";
import JsonLd from "@/components/seo/json-ld";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  alternates: {
    canonical: process.env.NEXT_PUBLIC_APP_BASE_URL,
  },
};

export default async function Home() {
  const jsonLdSchema = {
    "@context": "https://schema.org",
    "@type": "Organization",
    name: process.env.NEXT_PUBLIC_APP_NAME,
    url: process.env.NEXT_PUBLIC_APP_BASE_URL,
    logo: `${process.env.NEXT_PUBLIC_APP_BASE_URL}/assets/images/logos/logo-dark.png`,
    sameAs: [
      "https://www.instagram.com/next_template",
    ],
    contactPoint: {
      "@type": "ContactPoint",
      telephone: "+" + process.env.NEXT_PUBLIC_WHATSAPP_NUMBER,
      contactType: "customer service",
      availableLanguage: ["Indonesian", "English"],
    },
  };

  return (
    <>
      <JsonLd data={jsonLdSchema} />
      <div className="container-custom flex-col flex gap-8">
        <h1 className="sr-only">
          {process.env.NEXT_PUBLIC_APP_NAME} - Top Up Game Termurah, Tercepat, Terpercaya dan Aman
        </h1>
        tes
      </div>
    </>
  );
}
