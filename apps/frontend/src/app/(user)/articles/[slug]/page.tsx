import dynamic from "next/dynamic";
import { Suspense } from "react";
const ArticleShare = dynamic(() => import("./article-share"));
const RelatedArticles = dynamic(() => import("./related-articles"));
const TableOfContents = dynamic(() => import("./table-of-contents"));
const ReadingProgress = dynamic(() => import("@/components/reading-progress"));

import { getArticleDetail, getPopularArticles } from "./actions";
import { getApiImage } from "@/lib/api-asset";
import Image from "next/image";
import Typography from "@/components/typography";
import { Clock, User, Calendar } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { calculateReadingTime } from "@/utils/calculate";
import JsonLd from "@/components/seo/json-ld";

type Params = Promise<{ slug: string }>;

export async function generateMetadata(props: { params: Params }) {
  const { slug } = await props.params;
  const dataArticle = await getArticleDetail(slug);

  if (!dataArticle) {
    return {
      title: "Artikel Tidak Ditemukan | Next Template",
      description: "Artikel yang kamu cari tidak ditemukan.",
    };
  }

  const baseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL;
  const appName = process.env.NEXT_PUBLIC_APP_NAME || "Next Template";

  const title = `${dataArticle?.title} | ${appName} - Artikel Game & Top Up`;
  const description =
    dataArticle?.meta_description?.length > 0
      ? dataArticle.meta_description
      : `Baca artikel menarik seputar ${dataArticle?.title} hanya di ${appName}. Dapatkan tips, berita, dan panduan top up game murah terpercaya!`;

  const url = `${baseUrl}/articles/${slug}`;
  const image = getApiImage("articles", dataArticle?.img);

  const dataKeywords = dataArticle?.meta_keywords || "";

  const arrKeywords = dataKeywords
    .split(",")
    .map((k: string) => k.trim())
    .filter((k: string) => k.length > 0);

  const keywords = [
    dataArticle?.title,
    ...arrKeywords,
    "artikel game",
    "tips game online",
    "top up game murah",
    "voucher game",
    "review game",
    "berita game terbaru",
    "panduan top up",
    "top up game instan",
    "top up game terpercaya",
    "harga top up game 2025",
    "cara top up game cepat",
    "top up diamond game Indonesia",
    "beli voucher game resmi",
    "game populer Indonesia",
    "top up game termurah",
    "top up game legal Indonesia",
    "promo top up game",
    "diskon top up game",
    "top up game 24 jam",
    "metode pembayaran top up game Indonesia",
  ];

  return {
    title,
    description,
    keywords,
    alternates: { canonical: url },
    openGraph: {
      title,
      description,
      url,
      siteName: appName,
      images: [
        {
          url: image,
          width: 1200,
          height: 630,
          alt: dataArticle?.title,
        },
      ],
      locale: "id_ID",
      type: "article",
      publishedTime: dataArticle?.published_at,
      authors: [dataArticle?.author || "Tim Next Template"],
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: [image],
    },
    other: {
      "article:published_time": dataArticle?.published_at,
      "article:author": dataArticle?.author,
    },
  };
}

export default async function ArticleDetail(props: { params: Params }) {
  const { slug } = await props.params;
  const dataArticle = await getArticleDetail(slug);
  const dataPopularArticles = await getPopularArticles(slug);

  const title = dataArticle?.title;
  const content = dataArticle?.content;
  const shareUrl = `${process.env.NEXT_PUBLIC_APP_BASE_URL}/articles/${slug}`;
  const readingTime = calculateReadingTime(content);

  // Struktur data JSON-LD untuk SEO
  const jsonLdSchema = {
    "@context": "https://schema.org",
    "@type": "Article",
    headline: title,
    image: getApiImage("articles", dataArticle?.img),
    author: {
      "@type": "Person",
      name: dataArticle?.author,
    },
    publisher: {
      "@type": "Organization",
      name: process.env.NEXT_PUBLIC_APP_NAME,
      logo: {
        "@type": "ImageObject",
        url: `${process.env.NEXT_PUBLIC_APP_BASE_URL}/assets/images/logos/logo-dark.png`,
      },
    },
    datePublished: dataArticle?.published_at,
    dateModified: dataArticle?.updated_at,
    description: dataArticle?.meta_description,
  };

  return (
    <>
      {/* JSON-LD Schema */}
      <JsonLd data={jsonLdSchema} />

      {/* Reading Progress Bar */}
      <ReadingProgress />

      <article className="relative">
        <div className="container-custom py-8 lg:py-12">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
            {/* Sidebar - Table of Contents (Desktop) */}
            <aside className="hidden lg:block lg:col-span-3 xl:col-span-2">
              <div className="sticky top-20">
                <Suspense
                  fallback={
                    <div className="h-64 animate-pulse bg-muted rounded-lg" />
                  }
                >
                  <TableOfContents content={content} />
                </Suspense>
              </div>
            </aside>

            {/* Main Content */}
            <main className="lg:col-span-9 xl:col-span-7">
              {/* Breadcrumb */}
              <nav className="mb-6 text-sm" aria-label="Breadcrumb">
                <ol className="flex items-center gap-2 flex-wrap">
                  <li>
                    <a
                      href="/"
                      className="text-muted-foreground hover:text-foreground transition-colors"
                    >
                      Home
                    </a>
                  </li>
                  <li className="text-muted-foreground">/</li>
                  <li>
                    <a
                      href="/articles"
                      className="text-muted-foreground hover:text-foreground transition-colors"
                    >
                      Artikel
                    </a>
                  </li>
                  <li className="text-muted-foreground">/</li>
                  <li
                    className="text-foreground font-medium truncate max-w-[200px] sm:max-w-[400px]"
                    title={title}
                  >
                    {title}
                  </li>
                </ol>
              </nav>

              {/* Header */}
              <header className="mb-8">
                {/* Category Badge */}
                {dataArticle?.category && (
                  <Badge variant="secondary" className="mb-4">
                    {dataArticle.category.name}
                  </Badge>
                )}

                <Typography variant="h1" as="h1" className="mb-4 leading-tight">
                  {title}
                </Typography>

                {/* Meta Info */}
                <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground mb-6">
                  <div className="flex items-center gap-2">
                    <User className="w-4 h-4" />
                    <span>{dataArticle?.author}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4" />
                    <time dateTime={dataArticle?.published_at}>
                      {dataArticle?.published_at}
                    </time>
                  </div>
                  <div className="flex items-center gap-2">
                    <Clock className="w-4 h-4" />
                    <span>{readingTime} menit baca</span>
                  </div>
                </div>

                {/* Featured Image */}
                <div className="relative w-full overflow-hidden aspect-16/7 rounded-xl shadow-lg mb-8">
                  <Image
                    src={getApiImage("articles", dataArticle?.img)}
                    alt={dataArticle?.title}
                    fill
                    priority
                    sizes="(max-width: 768px) 100vw, (max-width: 1200px) 70vw, 900px"
                    className="object-cover"
                  />
                </div>
              </header>

              {/* Article Content */}
              <section
                className="prose prose-lg dark:prose-invert max-w-none 
                           prose-headings:scroll-mt-24 
                           prose-a:text-primary prose-a:no-underline hover:prose-a:underline
                           prose-img:rounded-xl prose-img:shadow-md
                           prose-code:bg-muted prose-code:px-1 prose-code:py-0.5 prose-code:rounded
                           prose-pre:bg-muted prose-pre:border
                           mb-12"
              >
                <div dangerouslySetInnerHTML={{ __html: content }} />
              </section>

              {/* Tags */}
              {dataArticle?.tags && dataArticle.tags.length > 0 && (
                <div className="mb-8">
                  <Typography variant="h6" className="mb-3">
                    Tags:
                  </Typography>
                  <div className="flex flex-wrap gap-2">
                    {dataArticle.tags.map((tag: string, index: number) => (
                      <a
                        key={index}
                        href={`/articles?tag=${encodeURIComponent(tag)}`}
                        className="px-3 py-1 text-sm bg-muted hover:bg-muted/80 rounded-full transition-colors"
                      >
                        #{tag}
                      </a>
                    ))}
                  </div>
                </div>
              )}

              {/* Share Buttons */}
              <ArticleShare title={title} shareUrl={shareUrl} />

              {/* Author Bio */}
              {dataArticle?.author_bio && (
                <div className="bg-muted/50 rounded-xl p-6 mb-12">
                  <Typography variant="h6" className="mb-3">
                    Tentang Penulis
                  </Typography>
                  <div className="flex gap-4">
                    <div className="shrink-0">
                      <div className="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center text-2xl font-bold">
                        {dataArticle.author.charAt(0).toUpperCase()}
                      </div>
                    </div>
                    <div>
                      <Typography variant="b1" className="font-semibold mb-1">
                        {dataArticle.author}
                      </Typography>
                      <Typography
                        variant="b3"
                        className="text-muted-foreground"
                      >
                        {dataArticle.author_bio}
                      </Typography>
                    </div>
                  </div>
                </div>
              )}
            </main>

            {/* Right Sidebar - CTA / Sticky Elements (Desktop) */}
            <aside className="hidden xl:block xl:col-span-3">
              <div className="sticky top-20 space-y-6">
                {/* Newsletter Signup */}
                {/* <div className="bg-linear-to-br from-primary/10 to-primary/5 rounded-xl p-6 border border-primary/20">
                  <Typography variant="h6" className="mb-2">
                    📧 Newsletter
                  </Typography>
                  <Typography
                    variant="b3"
                    className="text-muted-foreground mb-4"
                  >
                    Dapatkan artikel terbaru langsung ke email Anda
                  </Typography>
                  <form className="space-y-3">
                    <input
                      type="email"
                      placeholder="Email Anda"
                      className="w-full px-4 py-2 rounded-lg border bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                    <button
                      type="submit"
                      className="w-full bg-primary text-primary-foreground px-4 py-2 rounded-lg font-medium hover:bg-primary/90 transition-colors"
                    >
                      Subscribe
                    </button>
                  </form>
                </div> */}

                {/* Popular Articles Widget */}
                <div className="">
                  <Typography variant="h6" className="mb-4">
                    🔥 Artikel Populer
                  </Typography>
                  <div className="space-y-3">
                    {/* Placeholder - ganti dengan data real */}
                    {dataPopularArticles.map((article: any) => (
                      <a
                        key={article.slug}
                        href={`/articles/${article.slug}`}
                        className="block group rounded-lg p-4 transition-colors hover:bg-primary/10"
                      >
                        <div className="relative aspect-video overflow-hidden rounded-md">
                          <Image
                            src={getApiImage("articles", article.img)}
                            alt="Thumbnail"
                            fill
                            className="object-cover"
                          />
                        </div>
                        <div className="mt-2 flex justify-between">
                          <Typography
                            variant="b2"
                            className="font-semibold line-clamp-2 group-hover:text-primary transition-colors"
                          >
                            {article.title}
                          </Typography>
                        </div>
                        <div className="flex mt-2 items-center gap-2 text-accent/80">
                          <Typography
                            variant="c1"
                            as="time"
                            dateTime="2021-01-01"
                            className="text-accent/80 group-hover:text-accent transition-colors"
                          >
                            {new Date(article.published_at).toLocaleDateString(
                              "id-ID",
                              {
                                day: "numeric",
                                month: "short",
                                year: "numeric",
                              }
                            )}
                          </Typography>
                          {"•"}
                          <Typography
                            variant="c1"
                            className="text-accent/80 group-hover:text-accent transition-colors"
                          >
                            {calculateReadingTime(article.content)} menit baca
                          </Typography>
                        </div>
                      </a>
                    ))}
                  </div>
                </div>
              </div>
            </aside>
          </div>
          {/* Related Articles */}
          <Suspense
            fallback={
              <div className="space-y-4">
                <div className="h-8 w-48 bg-muted animate-pulse rounded" />
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {[1, 2].map((i) => (
                    <div
                      key={i}
                      className="h-48 bg-muted animate-pulse rounded-lg"
                    />
                  ))}
                </div>
              </div>
            }
          >
            <RelatedArticles articles={dataArticle.related_articles} />
          </Suspense>
        </div>
      </article>
    </>
  );
}
