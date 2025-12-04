import FormActionArticle from "@/components/forms/action-article";
import { useBreadcrumbStore } from "@/stores/breadcrumb";
type Params = Promise<{ id: string }>;
export const generateMetadata = async ({ params }: { params: Params }) => {
  const { id } = await params;
  return {
    title: id ? "Ubah Artikel" : "Tambah Artikel",
  };
};

export default async function ArticleActionPage(props: { params: Params }) {
  const { id } = await props.params;

  return <FormActionArticle id={id} />;
}
