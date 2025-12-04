import BannersClient from "./banners-client";

export const metadata = {
  title: "Banners",
}

export default async function BannersPage() {
  return (
    <BannersClient />
  );
}
