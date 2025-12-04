"use client";

import { Controller, useForm } from "react-hook-form";
import { useMutation, useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Input } from "@/components/ui/input";
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

const bannerSchema = z.object({
  id: z.string().optional(),
  title: z.string().min(3, "Judul minimal 3 karakter"),
  link: z.string().min(3, "Link minimal 3 karakter"),
  img: z.string(),
  content: z.string().nullable(),
  valid_from: z.string().nonempty("Valid from is required"),
  valid_until: z.string().nonempty("Valid until is required"),
  is_active: z.boolean(),
  order: z.number(),
});

export type BannerFormValues = z.infer<typeof bannerSchema>;

export default function FormActionBanner({ id }: { id?: string }) {
  const breadcrumbs = [
    {
      title: "Banner",
      url: "/admin/banners",
    },
  ];

  if (!id) {
    breadcrumbs.push({
      title: "Tambah Banner",
      url: "#",
    });
  } else {
    breadcrumbs.push({
      title: "Ubah Banner",
      url: "#",
    });
  }
  useEffect(() => {
    useBreadcrumbStore.setState({ breadcrumbs });
  }, [breadcrumbs]);

  const router = useRouter();
  const methods = useForm<BannerFormValues>({
    resolver: zodResolver(bannerSchema),
    defaultValues: {
      is_active: true,
      order: 0,
      content: "",
    },
  });

  const { register, handleSubmit, setValue, reset, formState, control, watch } =
    methods;

  const { data, isSuccess } = useQuery({
    queryKey: ["banner", id],
    queryFn: () => apiClient.get(`/v0/banners/${id}`).then((res) => res.data),
    enabled: !!id,
  });

  useEffect(() => {
    if (isSuccess && data) {
      reset(data?.data);
      setValue("valid_from", toDatetimeLocal(data?.data.valid_from));
      setValue("valid_until", toDatetimeLocal(data?.data.valid_until));
    }
  }, [isSuccess, data]);

  const mutation = useMutation({
    mutationFn: (formData: BannerFormValues) =>
      apiClient.put("/v1/banners", formData),
    onSuccess: () => {
      toast.success("Banner saved successfully");
      router.push("/admin/banners");
    },
    onError: (error: any) => {
      console.error(error);
      shadSwal.fire({
        text: error.message,
      });
    },
  });

  const onSubmit = handleSubmit((values) => {
    values.valid_from = new Date(values.valid_from).toISOString();
    values.valid_until = new Date(values.valid_until).toISOString();
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
              />
            )}
          />
          <Controller
            name="link"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Link"
                id="link"
                required
                error={error?.message}
              />
            )}
          />
          <Controller
            name="order"
            defaultValue={0}
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Urutan"
                id="order"
                type="number"
                required
                value={field.value || ""}
                error={error?.message}
                onChange={(e) => field.onChange(e.target.value === "" ? "" : Number(e.target.value))}
              />
            )}
          />
          <Controller
            name="valid_from"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                type="datetime-local"
                label="Tanggal Mulai"
                id="valid_from"
                required
                error={error?.message}
              />
            )}
          />
          <Controller
            name="valid_until"
            defaultValue=""
            control={control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                type="datetime-local"
                label="Tanggal Selesai"
                id="valid_until"
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
                value={field.value || ""}
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
        </div>
      </div>
      <div className="flex justify-end w-full space-y-4">
        <Button type="submit">Submit</Button>
      </div>
    </form>
  );
}
