import { createApiServer } from "@/lib/api-server";
import { cache } from 'react'

export const getArticleDetail = cache(async (slug: string) => {
    try {
        const apiServer = await createApiServer();
        const response = await apiServer.get(`/v1/articles/${slug}`, {
            params: {
                related: true
            }
        });
        response.data.data.published_at = new Date(response.data.data.published_at)
        .toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" })
        return response.data.data;
    } catch (error) {
        console.error('Error fetching article detail:', error);
        throw error;
    }
})

export const getPopularArticles = cache(async (slug: string) => {
    try {
        const apiServer = await createApiServer();
        const response = await apiServer.get(`/v1/articles/popular`, {
            params: {
                except: slug,
                limit: 3
            }
        });
        return response.data.data;
    } catch (error) {
        console.error('Error fetching popular articles:', error);
        throw error;
    }
})