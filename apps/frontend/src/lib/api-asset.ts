export function getApiImage(type: string, name: string): string {
    const url = `${process.env.NEXT_PUBLIC_CDN_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL}/uploads/images/${type}/${name}`;
    return url;
}
