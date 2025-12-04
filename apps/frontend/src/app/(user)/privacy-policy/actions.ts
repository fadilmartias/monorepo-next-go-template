import { createApiServer } from "@/lib/api-server";

export async function getSetting(key: string) {
    try {
        const apiServer = await createApiServer();
        const response = await apiServer.get(`/v1/settings/${key}`);
        return response.data;
    } catch (error) {
        console.error('Error fetching settings:', error);
        throw error;
    }
}