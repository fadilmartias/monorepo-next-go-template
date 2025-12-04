"use client";

import { Controller, useForm } from "react-hook-form";
import { useMutation, useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";
import ImageCropper from "../image-cropper";
import dynamic from "next/dynamic";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useRouter } from "nextjs-toploader/app";
import { shadSwal } from "../shad-swal";
import { toDatetimeLocal } from "@/utils/format";
import { useBreadcrumbStore } from "@/stores/breadcrumb";

const TinyEditor = dynamic(() => import("../tiny-editor"), {
  ssr: false,
  loading: () => <Loader2 className="animate-spin" />,
});

const ReactSelect = dynamic(() => import("../react-select"), {
  ssr: false,
  loading: () => <Loader2 className="animate-spin" />,
});

const articleSchema = z.object({
  id: z.string().optional(),
  title: z.string().min(3, "Judul minimal 3 karakter"),
  slug: z.string().min(3, "Slug minimal 3 karakter"),
  img: z.string(),
  category_id: z.string(),
  content: z.string().min(10, "Konten minimal 10 karakter"),
  meta_description: z.string().min(10, "Meta description minimal 10 karakter"),
  meta_img: z.string(),
  author: z.string(),
  meta_keywords: z.string(),
  published_at: z.string().nonempty("Published at is required"),
});

export type ArticleFormValues = z.infer<typeof articleSchema>;

export default function FormActionArticle({ id }: { id?: string }) {
  const breadcrumbs = [
    {
      title: "Artikel",
      url: "/admin/articles",
    },
  ];

  if (!id) {
    breadcrumbs.push({
      title: "Tambah Artikel",
      url: "#",
    });
  } else {
    breadcrumbs.push({
      title: "Ubah Artikel",
      url: "#",
    });
  }
  useEffect(() => {
    useBreadcrumbStore.setState({ breadcrumbs });
  }, [breadcrumbs]);

  const router = useRouter();
  const methods = useForm<ArticleFormValues>({
    resolver: zodResolver(articleSchema),
  });

  const { register, handleSubmit, setValue, reset, formState, control, watch } =
    methods;

  const { data, isSuccess } = useQuery({
    queryKey: ["article", id],
    queryFn: () => apiClient.get(`/v0/articles/${id}`).then((res) => res.data),
    enabled: !!id,
  });

  const { data: categories, error: errorCategories } = useQuery({
    queryKey: ["categories"],
    queryFn: () =>
      apiClient
        .get("/v0/categories", {
          params: {
            filters: {
              is_active: 1,
              type: "article",
            },
          },
        })
        .then((res) => res.data),
  });

  const categoriesData = categories?.data.map((item: any) => ({
    value: item.id,
    label: item.name,
  }));

  useEffect(() => {
    if (isSuccess && data) {
      reset(data?.data);
      setValue("published_at", toDatetimeLocal(data?.data.published_at));
    }
  }, [isSuccess, data]);

  const mutation = useMutation({
    mutationFn: (formData: ArticleFormValues) =>
      apiClient.put("/v1/articles", formData),
    onSuccess: () => {
      toast.success("Article saved successfully");
      router.push("/admin/articles");
    },
    onError: (error: any) => {
      console.error(error);
      shadSwal.fire({
        text: error.message,
      });
    },
  });

  const onSubmit = handleSubmit((values) => {
    values.published_at = new Date(values.published_at).toISOString();
    mutation.mutate(values);
  });

  return (
    <form onSubmit={onSubmit} className="space-y-8">
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-10">
        <div className="lg:col-span-8 space-y-4">
          <Controller
            name="title"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Judul"
                id="title"
                required
                error={error?.message}
                onChange={(e) => {
                  setValue(
                    "slug",
                    e.target.value.toLowerCase().replace(/[^a-z0-9]/g, "-"),
                    {
                      shouldValidate: true,
                      shouldTouch: true,
                      shouldDirty: true,
                    }
                  );
                  field.onChange(e.target.value);
                }}
              />
            )}
          />
          <Controller
            name="author"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Penulis"
                id="author"
                required
                error={error?.message}
              />
            )}
          />
          <Controller
            name="slug"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Slug"
                id="slug"
                required
                error={error?.message}
              />
            )}
          />
          <Controller
            name="category_id"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <ReactSelect
                {...field}
                label="Kategori"
                id="category"
                required
                options={categoriesData}
                onChange={(value: any) => field.onChange(value.value)}
                value={
                  categoriesData?.find(
                    (cat: any) => cat.value === field.value
                  ) || null
                }
                error={error?.message}
              />
            )}
          />
          <Controller
            name="published_at"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                type="datetime-local"
                label="Tanggal Terbit"
                id="published_at"
                required
                error={error?.message}
              />
            )}
          />
          <Controller
            name="content"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <TinyEditor
                label="Konten"
                id="content"
                required
                value={field.value}
                onEditorChange={(content: string) => field.onChange(content)}
                error={error?.message}
              />
            )}
          />
        </div>

        <div className="lg:col-span-4 space-y-4">
          <Controller
            name="img"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <ImageCropper
                {...field}
                pathName="articles"
                label="Banner"
                id="img"
                required
                onCropped={(name: string) => field.onChange(name)}
                initialUrl={data?.data.img}
                error={error?.message}
              />
            )}
          />
          <Controller
            name="meta_img"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <ImageCropper
                {...field}
                pathName="articles"
                label="Meta Image"
                id="meta_img"
                onCropped={(name: string) => field.onChange(name)}
                initialUrl={data?.data.meta_img}
                error={error?.message}
              />
            )}
          />
          <Controller
            name="meta_description"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Textarea
                {...field}
                error={error?.message}
                label="Meta Description"
                id="metaDescription"
              />
            )}
          />

          <Controller
            name="meta_keywords"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                error={error?.message}
                label="Meta Keywords"
                id="metaKeywords"
              />
            )}
          />
        </div>
      </div>
      <div className="flex justify-end w-full space-y-4">
        <Button type="submit">Submit</Button>
      </div>
    </form>
  );
}
