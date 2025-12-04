import Image from "next/image";
import Link from "next/link";
import Typography from "@/components/typography";
import { getApiImage } from "@/lib/api-asset";
import { ChevronRight, Clock } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Article {
  slug: string;
  title: string;
  img: string;
  excerpt?: string;
  published_at: string;
  reading_time?: number;
}

export default function RelatedArticles({
  articles,
}: {
  articles: Article[];
}) {
  if (!articles || articles.length === 0) {
    return null;
  }

  return (
    <section className="mt-12">
      <Typography variant="h4" as="h2" className="mb-6">
        📰 Artikel Terkait
      </Typography>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {articles.slice(0, 3).map((article) => (
          <Link
            key={article.slug}
            href={`/articles/${article.slug}`}
            className="group bg-card border rounded-xl overflow-hidden hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
          >
            <div className="relative aspect-video overflow-hidden">
              <Image
                src={getApiImage("articles", article.img)}
                alt={article.title}
                fill
                sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                className="object-cover group-hover:scale-105 transition-transform duration-300"
              />
            </div>
            <div className="p-4">
              <Typography
                variant="b1"
                className="font-semibold line-clamp-2 mb-2 group-hover:text-primary transition-colors"
              >
                {article.title}
              </Typography>
              {article.excerpt && (
                <Typography
                  variant="b3"
                  className="text-muted-foreground line-clamp-2 mb-3"
                >
                  {article.excerpt}
                </Typography>
              )}
              <div className="flex items-center gap-3 text-xs text-muted-foreground">
                <time dateTime={article.published_at}>{new Date(article.published_at).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" })}</time>
                {article.reading_time && (
                  <>
                    <span>•</span>
                    <div className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      <span>{article.reading_time} menit</span>
                    </div>
                  </>
                )}
              </div>
            </div>
          </Link>
        ))}
      </div>
      
      {/* View All Articles Button */}
      <div className="mt-8 text-center">
       <Button size="xl" asChild>
         <Link
          href="/articles"
          className=""
        >
          Lihat Semua Artikel
          <ChevronRight className="w-4 h-4" />
        </Link>
       </Button>
      </div>
    </section>
  );
}