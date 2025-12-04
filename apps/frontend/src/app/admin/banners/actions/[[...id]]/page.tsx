import FormActionBanner from "@/components/forms/action-banner";

type Params = Promise<{ id: string }>;
export const generateMetadata = async ({ params }: { params: Params }) => {
  const { id } = await params;
  return {
    title: id ? "Ubah Banner" : "Tambah Banner",
  };
};

export default async function BannerActionPage(props: { params: Params }) {
  const { id } = await props.params;

  return <FormActionBanner id={id} />;
}
