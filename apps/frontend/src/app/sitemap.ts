import type { MetadataRoute } from 'next'
import { getLatestArticles, getTrendingProducts } from './(user)/actions'

export const dynamic = "force-dynamic";
 
export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
    const articles = await getLatestArticles(999)
    const products = await getTrendingProducts()
    const baseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL || 'https://nexttemplate.com'
    
    const articleUrls = articles.map((article: any) => ({
      url: `${baseUrl}/articles/${article.slug}`,
      lastModified: new Date(article.updated_at),
      changeFrequency: 'weekly',
      priority: 0.8,
    }))

    const productUrls = products.map((product: any) => ({
      url: `${baseUrl}/topup/${product.slug}`,
      lastModified: new Date(product.updated_at),
      changeFrequency: 'weekly',
      priority: 0.9,
    }))
  return [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 1,
    },
    {
      url: `${baseUrl}/terms-and-conditions`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.8,
    },
    {
      url: `${baseUrl}/check-transaction`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.7,
    },
    {
      url: `${baseUrl}/privacy-policy`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.8,
    },
    {
      url: `${baseUrl}/articles`,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 0.8,
    },
    ...articleUrls,
    ...productUrls,
  ]
}