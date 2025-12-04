import { createApiServer } from "@/lib/api-server";
import { cache } from "react";

export const getTrendingProducts = async () => {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.get('/v1/products/trending');
    return response.data.data;
  } catch (error) {
    console.error('Error fetching products:', error);
    throw error;
  }
};

export const getFlashSales = async () => {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.get('/v1/promotions/flash-sales');
    return response.data.data;
  } catch (error) {
    console.error('Error fetching flash sales:', error);
    throw error;
  }
};

export const getLatestArticles = async (limit: number = 6) => {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.get(`/v1/articles`, {
      params: {
        limit,
      },
    });
    if(response.data && response.data.data.length > 0) {
      response.data.data.forEach((article: any) => {
        article.published_at = new Date(article.published_at)
        .toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" })
      });
    }
    return response.data.data;
  } catch (error) {
    console.error('Error fetching articles:', error);
    throw error;
  }
};

export const getBanners = async () => {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.get('/v1/banners');
    return response.data.data;
  } catch (error) {
    console.error('Error fetching banners:', error);
    throw error;
  }
};

export const getPaymentMethods = cache(async ({
  cache = false,
  cache_ttl = 60 * 60 * 24,
  cache_key = 'active-payment-methods',
}: {
  cache?: boolean;
  cache_ttl?: number;
  cache_key?: string;
}) => {
  try {
    const apiServer = await createApiServer();
    const response = await apiServer.get('/v1/payment-methods', {
      params: {
        cache,
        cache_ttl,
        cache_key,
      }
    });
    return response.data.data;
  } catch (error) {
    console.error('Error fetching payment methods:', error);
    throw error;
  }
});