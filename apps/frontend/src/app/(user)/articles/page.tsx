import { getLatestArticles } from "@/app/(user)/actions";
import { Metadata } from "next";
import LatestArticles from "@/app/(user)/latest-articles";
import Typography from "@/components/typography";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: {
    default: "Artikel Game, Tips Bermain & Top Up Murah",
    template: `%s | ${process.env.NEXT_PUBLIC_APP_NAME}`,
  },
  description:
    "Baca artikel seputar game online, tips & trik bermain, panduan top up murah, review game terbaru, dan informasi promo eksklusif hanya di " +
    process.env.NEXT_PUBLIC_APP_NAME,
  keywords: [
    "artikel game",
    "berita game",
    "tips game online",
    "top up game murah",
    "panduan top up",
    "promo game online",
    "voucher game",
    "tutorial game",
    "review game",
    "diskon top up game",
  ],
  alternates: {
    canonical: process.env.NEXT_PUBLIC_APP_BASE_URL + "/articles",
  },
};

export default async function Home() {
  const latestArticles = await getLatestArticles(999);

  return (
    <div className="container-custom flex-col flex gap-8">
      <h1 className="sr-only">Artikel Seputar Game & Top Up Game</h1>
      <div className="mt-8">
        {latestArticles.length > 0 ? (
          <LatestArticles data={latestArticles} isShowFeatured />
        ) : (
          <div className="flex flex-col items-center justify-center">
            <Typography variant="h5" className="text-center">
              Belum ada artikel
            </Typography>
          </div>
        )}
      </div>
    </div>
  );
}
