"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { toast } from "sonner";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import z from "zod";
import dynamic from "next/dynamic";
import { Loader2, Plus, Trash2 } from "lucide-react";
import ImageCropper from "../image-cropper";
import { Label } from "../ui/label";
import { Switch } from "../ui/switch";
import Link from "next/link";
import { useBreadcrumbStore } from "@/stores/breadcrumb";
import { roundNumber } from "@/utils/calculate";

const TinyEditor = dynamic(() => import("../tiny-editor"), {
  ssr: false,
  loading: () => <Loader2 className="animate-spin" />,
});

const ReactSelect = dynamic(() => import("../react-select"), {
  ssr: false,
  loading: () => <Loader2 className="animate-spin" />,
});

const productSchema = z.object({
  categories: z.array(z.string()),
  iak_id: z.string().optional(),
  check_account_sku: z.string().optional(),
  name: z.string().min(1, "Nama produk tidak boleh kosong"),
  additional_fields: z.array(
    z.object({
      name: z.string().min(1, "Nama field tidak boleh kosong"),
      type: z.string().min(1, "Type field tidak boleh kosong"),
      label: z.string().min(1, "Label field tidak boleh kosong"),
      check_name: z.string().optional(),
      copy: z.boolean().optional(),
      required: z.boolean().optional(),
      guide_title: z.string().optional(),
      guide_content: z.string().optional(),
      min_length: z
        .string()
        .regex(/^\d+$/, "Panjang minimal harus angka")
        .refine(
          (v) => parseInt(v, 10) >= 0,
          "Panjang minimal tidak boleh kurang dari 0"
        )
        .optional(),
      order: z
        .string()
        .regex(/^\d+$/, "Urutan harus angka")
        .refine((v) => parseInt(v, 10) >= 0, "Urutan tidak boleh kurang dari 0")
        .optional(),
      view_order: z
        .string()
        .regex(/^\d+$/, "Urutan untuk tampilan harus angka")
        .refine(
          (v) => parseInt(v, 10) >= 0,
          "Urutan untuk tampilan tidak boleh kurang dari 0"
        )
        .optional(),
    })
  ),
  additional_fields_separator: z.string().optional(),
  desc: z.string().optional(),
  slug: z.string().min(1, "Slug tidak boleh kosong"),
  img: z.string().min(1, "Gambar tidak boleh kosong"),
  banner_mobile: z.string().min(1, "Banner mobile tidak boleh kosong"),
  banner_desktop: z.string().min(1, "Banner desktop tidak boleh kosong"),
  developer: z.string().min(1, "Developer tidak boleh kosong"),
  publisher: z.string().min(1, "Publisher tidak boleh kosong"),
  product_type_id: z.string().min(1, "Type tidak boleh kosong"),
  variants: z.array(
    z.object({
      id: z.string().optional(),
      categories: z.array(z.string()),
      name: z.string().min(1, "Nama varian tidak boleh kosong"),
      provider_sku: z.string().min(1, "Harga varian tidak boleh kosong"),

      cost_price: z
        .string()
        .regex(/^\d+(\.\d+)?$/, "Harga varian harus angka")
        .refine(
          (v) => parseFloat(v) >= 0,
          "Harga modal varian tidak boleh kurang dari 0"
        ),

      margin_fixed: z
        .string()
        .regex(/^\d+(\.\d+)?$/, "Margin harus angka")
        .refine((v) => parseFloat(v) >= 0, "Margin tidak boleh kurang dari 0"),

      margin_percent: z
        .string()
        .regex(/^\d+(\.\d+)?$/, "Margin persen harus angka")
        .refine(
          (v) => parseFloat(v) >= 0,
          "Margin persen tidak boleh kurang dari 0"
        )
        .refine(
          (v) => parseFloat(v) <= 100,
          "Margin persen tidak boleh lebih dari 100"
        ),

      price: z
        .string()
        .regex(/^\d+(\.\d+)?$/, "Harga harus angka")
        .refine(
          (v) => parseFloat(v) >= 0,
          "Harga varian tidak boleh kurang dari 0"
        ),

      fee_admin: z
        .string()
        .regex(/^\d+(\.\d+)?$/, "Fee admin harus angka")
        .refine(
          (v) => parseFloat(v) >= 0,
          "Fee admin tidak boleh kurang dari 0"
        ),

      img: z.string().min(1, "Gambar varian tidak boleh kosong"),

      stock: z
        .string()
        .regex(/^\d+$/, "Stok harus angka")
        .refine(
          (v) => parseInt(v, 10) >= 0,
          "Stok varian tidak boleh kurang dari 0"
        ),
      is_unlimited: z.boolean(),
      is_multi: z.boolean(),
      start_cut_off: z.string().optional(),
      end_cut_off: z.string().optional(),
      min_qty: z
        .string()
        .regex(/^\d+$/, "Minimal qty harus angka")
        .refine(
          (v) => parseInt(v, 10) >= 0,
          "Minimal qty varian tidak boleh kurang dari 0"
        ),

      max_qty: z
        .string()
        .regex(/^\d+$/, "Maksimal qty harus angka")
        .refine(
          (v) => parseInt(v, 10) >= 0,
          "Maksimal qty varian tidak boleh kurang dari 0"
        ),

      is_active: z.boolean(),
      is_ready: z.boolean(),
    })
  ),
});

const calculatePrice = (
  cost_price: number,
  margin_fixed: number,
  margin_percent: number
) => {
  const price = cost_price + margin_fixed + (cost_price * margin_percent) / 100;
  return String(Math.ceil(price));
};

const productTypeData = [
  { value: "1", label: "Prepaid" },
  { value: "2", label: "Postpaid" },
];

const additionalFieldDataType = [
  { value: "text", label: "Text" },
  { value: "number", label: "Number" },
];

export type ProductFormValues = z.infer<typeof productSchema>;

export default function FormActionProduct({ id }: { id?: string }) {
  const breadcrumbs = [
    {
      title: "Produk",
      url: "#",
    },
  ];

  if (!id) {
    breadcrumbs.push({
      title: "Tambah Produk",
      url: "#",
    });
  } else {
    breadcrumbs.push({
      title: "Ubah Produk",
      url: "#",
    });
  }

  useEffect(() => {
    useBreadcrumbStore.setState({ breadcrumbs });
  }, [breadcrumbs]);
  const router = useRouter();

  const getProductTypeData = async () => {
    try {
      const res = await apiClient.get("/v0/categories/", {
        params: {
          filters: {
            type: "product-type",
            is_active: 1,
          },
        limit: 0
      },
    });
    return res.data.data.map((item: any) => {
      return {
        value: item.id,
        label: item.name,
      }
    })
  } catch (error) {
    console.error(error);
    toast.error("Terjadi kesalahan, coba lagi");
  }
}

  const getProductCategories = async () => {
    try {
      const res = await apiClient.get("/v0/categories/", {
        params: {
          filters: {
            type: "product",
            is_active: 1,
          },
          limit: 0
        },
      });
    return res.data.data.map((item: any) => {
      return {
        value: item.id,
        label: item.name,
      }
    })
  } catch (error) {
    console.error(error);
    toast.error("Terjadi kesalahan, coba lagi");
  }
}

  const getProductSubCategories = async () => {
    try {
      const res = await apiClient.get("/v0/categories/", {
        params: {
          filters: {
            type: "product-subcategory",
            is_active: 1,
          },
          limit: 0
        },
      });
    return res.data.data.map((item: any) => {
      return {
        value: item.id,
        label: item.name,
      }
    })
  } catch (error) {
    console.error(error);
    toast.error("Terjadi kesalahan, coba lagi");
  }
}

  const getProduct = async (id?: string) => {
    try {
      const res = await apiClient.get("/v0/products/" + id, {
        params: {
          joins: "Categories,ProductVariants.Categories"
        }
      });
      const formattedData: ProductFormValues = {
        variants: res.data.data.product_variants?.map((variant: any) => {
          return {
            id: variant.id,
            categories: variant.categories?.flatMap((field:any) => [field.id]) || [],
            name: variant.name || "",
            price: String(variant.price) || "0",
            img: variant.img || "",
            provider_sku: variant.provider_sku || "",
            cost_price: String(variant.cost_price) || "0",
            margin_fixed: String(variant.margin_fixed) || "0",
            margin_percent: String(roundNumber(variant.margin_percent*100,2)) || "0",
            fee_admin: String(variant.fee_admin) || "0",
            stock: String(variant.stock) || "0",
            is_unlimited: variant.is_unlimited || false,
            is_multi: variant.is_multi || false,
            start_cut_off: variant.start_cut_off || "",
            end_cut_off: variant.end_cut_off || "",
            min_qty: String(variant.min_qty) || "0",
            max_qty: String(variant.max_qty) || "5",
            is_active: variant.is_active || false,
            is_ready: variant.is_ready || false,
          }
        }) || [],
        additional_fields: res.data.data.additional_fields?.map((field: any) => {
          return {
            name: field.name || "",
            type: field.type || "",
            label: field.label || "",
            check_name: field.check_name || "",
            copy: field.copy || false,
            required: field.required || false,
            min_length: String(field.min_length) || "0",
            guide_title: field.guide_title || "",
            guide_content: field.guide_content || "",
            order: String(field.order) || "0",
            view_order: String(field.view_order) || "0",
          }
        }) || [],
        categories: res.data.data.categories?.flatMap((field:any) => [field.id]) || [],
        name: res.data.data.name || "",
        additional_fields_separator: res.data.data.additional_fields_separator || "",
        desc: res.data.data.desc || "",
        slug: res.data.data.slug || "",
        img: res.data.data.img || "",
        banner_mobile: res.data.data.banner_mobile || "",
        banner_desktop: res.data.data.banner_desktop || "",
        developer: res.data.data.developer || "",
        publisher: res.data.data.publisher || "",
        product_type_id: res.data.data.categories?.find((item:any) => item.type === "product-type")?.id || "",
        iak_id: res.data.data.iak_id || "",
        check_account_sku: res.data.data.check_account_sku || ""
      }
      return formattedData;
    } catch (error: any) {
      console.error(error);
      toast.error("Terjadi kesalahan, coba lagi");
      throw error;
    }
  };

  const { data, isSuccess, isLoading } = useQuery({
    queryKey: ["product", id],
    queryFn: () => getProduct(id),
    enabled: !!id,
  });

  const { data: productTypeData, isSuccess: isSuccessProductType } = useQuery({
    queryKey: ["product-type"],
    queryFn: () => getProductTypeData(),
  });

  const { data: productCategoriesData, isSuccess: isSuccessProductCategories } =
    useQuery({
      queryKey: ["product-categories"],
      queryFn: () => getProductCategories(),
    });

  const {
    data: productSubCategoriesData,
    isSuccess: isSuccessProductSubCategories,
  } = useQuery({
    queryKey: ["product-subcategories"],
    queryFn: () => getProductSubCategories(),
  });

  useEffect(() => {
    if (isSuccess && data) {
      reset(data);
    }
  }, [isSuccess, data]);

  const variantDefaultValue: ProductFormValues["variants"][number] = {
    id: "",
    categories: [],
    name: "",
    price: "",
    img: "",
    provider_sku: "",
    cost_price: "",
    margin_fixed: "",
    margin_percent: "",
    fee_admin: "",
    stock: "",
    is_unlimited: false,
    is_multi: false,
    start_cut_off: "",
    end_cut_off: "",
    min_qty: "",
    max_qty: "",
    is_active: false,
    is_ready: false,
  };

  const additionalFieldDefaultValue: ProductFormValues["additional_fields"][number] =
    {
      name: "",
      type: "text",
      label: "",
      check_name: "",
      copy: false,
      required: false,
      min_length: "",
      guide_title: "",
      guide_content: "",
      order: "",
      view_order: "",
    };

  const methods = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      variants: [variantDefaultValue],
      additional_fields: [additionalFieldDefaultValue],
    },
  });

  const { handleSubmit, setValue, reset, control, getValues } = methods;

  // const { fields, append, remove } = useFieldArray({
  //   control,
  //   name: "additional_fields",
  // });

  const {
    fields: variants,
    append: appendVariant,
    remove: removeVariant,
  } = useFieldArray({
    control,
    name: "variants",
  });

  const {
    fields: additionalFields,
    append: appendAdditionalField,
    remove: removeAdditionalField,
  } = useFieldArray({
    control,
    name: "additional_fields",
  });

  const addVariant = () => {
    appendVariant(variantDefaultValue);
  };

  const removeVariantFn = (index: number) => {
    if (variants.length > 1) {
      removeVariant(index);
    }
  };

  const addAdditionalField = () => {
    appendAdditionalField(additionalFieldDefaultValue);
  };

  const removeAdditionalFieldFn = (index: number) => {
    if (additionalFields.length > 1) {
      removeAdditionalField(index);
    }
  };

  const onSubmit = async (data: ProductFormValues) => {
    try {
      let res;
      if (id) {
        res = await apiClient.put("/v1/products/" + id, data);
      } else {
        res = await apiClient.post("/v1/products", data);
      }
      toast.success("Produk berhasil disimpan");
      router.push("/admin/products/prepaid");
    } catch (error) {
      console.error(error);
      toast.error("Terjadi kesalahan, coba lagi");
    }
  };

  return (
    <section className="">
      <Button className="font-semibold mb-6" asChild>
        <Link href="#variants" className="block">
          To variant
        </Link>
      </Button>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Informasi Produk</CardTitle>
                <CardDescription>Tambahkan informasi produk</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            <Controller
              name="name"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input
                  {...field}
                  required
                  label="Nama Produk"
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
              name="slug"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input
                  {...field}
                  required
                  label="Slug"
                  error={error?.message}
                />
              )}
            />
            <Controller
              name={`categories`}
              control={control}
              render={({ field, fieldState: { error } }) => (
                <ReactSelect
                  {...field}
                  label="Kategori"
                  id="categories"
                  options={productCategoriesData}
                  error={error?.message}
                  required
                  isMulti
                  // ubah selected ke bentuk object agar tampil di select
                  value={productCategoriesData?.filter((option: any) =>
                    field.value?.includes(option.value)
                  )}
                  // ubah hasil pilih agar ke form masuk array string
                  onChange={(selected: any) => {
                    const values = selected
                      ? selected.map((opt: any) => opt.value)
                      : [];
                    field.onChange(values);
                  }}
                />
              )}
            />
            <Controller
              name="product_type_id"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <ReactSelect
                  {...field}
                  label="Type"
                  id="product_type_id"
                  required
                  options={productTypeData}
                  onChange={(value: any) => field.onChange(value.value)}
                  value={
                    productTypeData?.find(
                      (cat: any) => cat.value === field.value
                    ) || null
                  }
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="desc"
              control={control}
              render={({ field, fieldState: { error } }) => (
                <TinyEditor
                  {...field}
                  label="Deskripsi"
                  id="desc"
                  onEditorChange={(content: string) => field.onChange(content)}
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="developer"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input
                  {...field}
                  required
                  label="Developer"
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="publisher"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input
                  {...field}
                  required
                  label="Publisher"
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="img"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <ImageCropper
                  {...field}
                  pathName="products"
                  label="Image"
                  id="img"
                  required
                  onCropped={(name: string) => field.onChange(name)}
                  initialUrl={data?.img || ""}
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="banner_mobile"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <ImageCropper
                  {...field}
                  pathName="products"
                  label="Banner Mobile"
                  id="banner_mobile"
                  required
                  onCropped={(name: string) => field.onChange(name)}
                  initialUrl={data?.banner_mobile || ""}
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="banner_desktop"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <ImageCropper
                  {...field}
                  pathName="products"
                  label="Banner Desktop"
                  id="banner_desktop"
                  required
                  onCropped={(name: string) => field.onChange(name)}
                  initialUrl={data?.banner_desktop}
                  error={error?.message}
                />
              )}
            />
            <Controller
              name="iak_id"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input {...field} label="Iak ID" error={error?.message} />
              )}
            />
            <Controller
              name="check_account_sku"
              defaultValue=""
              control={control}
              render={({ field, fieldState: { error } }) => (
                <Input
                  {...field}
                  label="Check Account SKU"
                  error={error?.message}
                />
              )}
            />
          </CardContent>
        </Card>

        {/* Additional Fields */}
        <Card className="mt-6" id="additional_fields">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Field Tambahan</CardTitle>
                <CardDescription>
                  Tambahkan berbagai field tambahan
                </CardDescription>
              </div>
              <Button
                type="button"
                onClick={addAdditionalField}
                variant="outline"
                size="sm"
                className="gap-2 bg-transparent"
              >
                <Plus className="h-4 w-4" />
                Tambah Field
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <Controller
                name={`additional_fields_separator`}
                defaultValue=""
                control={control}
                render={({ field, fieldState: { error } }) => (
                  <Input {...field} label="Pemisah" error={error?.message} />
                )}
              />
              {additionalFields?.map((variant, index) => (
                <div key={index} className="p-4 border rounded-lg space-y-4">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">Field {index + 1}</h4>
                    {additionalFields.length > 1 && (
                      <Button
                        type="button"
                        onClick={() => removeAdditionalFieldFn(index)}
                        variant="ghost"
                        size="sm"
                        className="text-destructive hover:text-destructive"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <Controller
                      name={`additional_fields.${index}.required`}
                      defaultValue={false}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="required">Required</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.copy`}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="copy">Copy</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.name`}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Nama Field"
                          error={error?.message}
                          required
                        />
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.check_name`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Check Name"
                          error={error?.message}
                        />
                      )}
                    />

                    <Controller
                      name={`additional_fields.${index}.label`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Label"
                          error={error?.message}
                          required
                        />
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.type`}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <ReactSelect
                          {...field}
                          label="Type"
                          id="type"
                          options={additionalFieldDataType}
                          error={error?.message}
                          required
                          onChange={(value: any) => field.onChange(value.value)}
                          value={
                            additionalFieldDataType?.find(
                              (cat: any) => cat.value === field.value
                            ) || null
                          }
                        />
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.min_length`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Panjang Minimal"
                          error={error?.message}
                          type="number"
                          required
                        />
                      )}
                    />

                    <Controller
                      name={`additional_fields.${index}.order`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Urutan"
                          error={error?.message}
                          type="number"
                          required
                        />
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.view_order`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Urutan Tampilan"
                          error={error?.message}
                          type="number"
                          required
                        />
                      )}
                    />
                    <Controller
                      name={`additional_fields.${index}.guide_title`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Judul Guide"
                          error={error?.message}
                        />
                      )}
                    />
                    <div className="col-span-2">
                      <Controller
                        name={`additional_fields.${index}.guide_content`}
                        defaultValue=""
                        control={control}
                        render={({ field, fieldState: { error } }) => (
                          <TinyEditor
                            {...field}
                            label="Konten Guide"
                            id="guide_content"
                            onEditorChange={(content: string) =>
                              field.onChange(content)
                            }
                            error={error?.message}
                          />
                        )}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Product Variants */}
        <Card className="mt-6" id="variants">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Varian Produk</CardTitle>
                <CardDescription>
                  Tambahkan berbagai denominasi dan harga
                </CardDescription>
              </div>
              <Button
                type="button"
                onClick={addVariant}
                variant="outline"
                size="sm"
                className="gap-2 bg-transparent"
              >
                <Plus className="h-4 w-4" />
                Tambah Varian
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {variants?.map((variant, index) => (
                <div key={index} className="p-4 border rounded-lg space-y-4">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">Varian {index + 1}</h4>
                    {variants.length > 1 && (
                      <Button
                        type="button"
                        onClick={() => removeVariantFn(index)}
                        variant="ghost"
                        size="sm"
                        className="text-destructive hover:text-destructive"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <Controller
                      name={`variants.${index}.categories`}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <ReactSelect
                          {...field}
                          label="Kategori"
                          id="categories"
                          options={productSubCategoriesData}
                          error={error?.message}
                          required
                          isMulti
                          // ubah selected ke bentuk object agar tampil di select
                          value={productSubCategoriesData?.filter(
                            (option: any) => field.value?.includes(option.value)
                          )}
                          // ubah hasil pilih agar ke form masuk array string
                          onChange={(selected: any) => {
                            const values = selected
                              ? selected.map((opt: any) => opt.value)
                              : [];
                            field.onChange(values);
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.name`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Nama Varian"
                          required
                          error={error?.message}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.provider_sku`}
                      defaultValue=""
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="SKU Provider"
                          error={error?.message}
                          required
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.cost_price`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Harga Modal (Rp)"
                          required
                          error={error?.message}
                          type="decimal"
                          onChange={(e) => {
                            const newPrice = calculatePrice(
                              Number(e.target.value),
                              Number(
                                getValues(`variants.${index}.margin_fixed`)
                              ),
                              Number(
                                getValues(`variants.${index}.margin_percent`)
                              )
                            );
                            setValue(`variants.${index}.price`, newPrice);
                            field.onChange(e.target.value);
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.margin_fixed`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Harga Margin Tetap (Rp)"
                          required
                          error={error?.message}
                          type="decimal"
                          onChange={(e) => {
                            const newPrice = calculatePrice(
                              Number(getValues(`variants.${index}.cost_price`)),
                              Number(e.target.value),
                              Number(
                                getValues(`variants.${index}.margin_percent`)
                              )
                            );
                            setValue(`variants.${index}.price`, newPrice);
                            field.onChange(e.target.value);
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.margin_percent`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Harga Margin Persen (%)"
                          required
                          error={error?.message}
                          type="decimal"
                          onChange={(e) => {
                            const newPrice = calculatePrice(
                              Number(getValues(`variants.${index}.cost_price`)),
                              Number(
                                getValues(`variants.${index}.margin_fixed`)
                              ),
                              Number(e.target.value)
                            );
                            setValue(`variants.${index}.price`, newPrice);
                            field.onChange(e.target.value);
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.price`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Harga Final (Rp)"
                          required
                          error={error?.message}
                          type="decimal"
                          readOnly
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.fee_admin`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Fee Admin (Rp)"
                          required
                          error={error?.message}
                          type="decimal"
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.img`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <ImageCropper
                          pathName="product-variants"
                          onCropped={(name: string) => field.onChange(name)}
                          initialUrl={variant.img || ""}
                          label="Gambar"
                          id="img"
                          required
                          error={error?.message}
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.stock`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Stok"
                          required
                          error={error?.message}
                          type="number"
                        />
                      )}
                    />

                    <Controller
                      name={`variants.${index}.start_cut_off`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Start Cut Off"
                          error={error?.message}
                          type="datetime-local"
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.end_cut_off`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="End Cut Off"
                          error={error?.message}
                          type="datetime-local"
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.min_qty`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Min Qty"
                          required
                          error={error?.message}
                          type="number"
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.max_qty`}
                      defaultValue={""}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <Input
                          {...field}
                          label="Max Qty"
                          required
                          error={error?.message}
                          type="number"
                        />
                      )}
                    />
                    <Controller
                      name={`variants.${index}.is_unlimited`}
                      defaultValue={false}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="isUnlimited">Unlimited</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                    <Controller
                      name={`variants.${index}.is_multi`}
                      defaultValue={false}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="isMulti">Multi</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                    <Controller
                      name={`variants.${index}.is_active`}
                      defaultValue={false}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="isActive">Active</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                    <Controller
                      name={`variants.${index}.is_ready`}
                      defaultValue={false}
                      control={control}
                      render={({ field, fieldState: { error } }) => (
                        <div className="flex items-center justify-between gap-2">
                          <Label htmlFor="isReady">Ready</Label>
                          <Switch
                            checked={field.value} // boolean state dari form
                            onCheckedChange={field.onChange} // update form ketika toggle
                          />
                        </div>
                      )}
                    />
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
        <div className="flex gap-3 justify-end mt-6">
          <Button onClick={() => router.back()} type="button" variant="outline">
            Batal
          </Button>
          <Button type="submit">Simpan</Button>
        </div>
      </form>
    </section>
  );
}
