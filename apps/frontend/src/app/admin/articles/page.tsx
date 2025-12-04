import ArticlesClient from "./articles-client";
import { useBreadcrumbStore } from "@/stores/breadcrumb";

export const metadata = {
  title: "Articles",
}

export default async function ArticlesPage() {
  return (
    <ArticlesClient />
  );
}
