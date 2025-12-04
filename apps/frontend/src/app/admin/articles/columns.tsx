"use client";

import { ColumnDef } from "@tanstack/react-table";
import clsx from "clsx";
import { MoreHorizontal } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Checkbox } from "@/components/ui/checkbox";
import { DataTableColumnHeader } from "@/components/data-table/data-table-column-header";
import Link from "next/link";
import { copyToClipboard } from "@/utils/copy";
import StatusToggle from "@/components/switches/status-toggle";
import { APIPaginationResponse } from "@/types/api";
import { formatTanggal } from "@/utils/format";

export type Article = {
  id: string;
  tenant_id: string;
  category_name: string;
  title: string;
  content: string;
  author: string;
  slug: string;
  img: string;
  published_at: string;
  is_active: string;
};

interface ArticleColumnsProps {
  serverSide?: boolean;
  sort?: string; // contoh: "title:asc"
  setSort?: (sort: string) => void;
  pagination?: APIPaginationResponse;
  setPagination?: (pagination: APIPaginationResponse) => void;
}

export function getArticleColumns({
  serverSide = false,
  sort,
  setSort,
  pagination,
  setPagination,
}: ArticleColumnsProps): ColumnDef<Article>[] {
  return [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) =>
            table.toggleAllPageRowsSelected(!!value)
          }
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "no",
      id: "no",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title="No"
        />
      ),
      cell: ({ row }) => {
        console.log(pagination)
        const pageSize = pagination?.page_size || 10
        const pageIndex = (pagination?.page || 1) - 1
        return (
          <span>
            {pageIndex * pageSize + row.index + 1}
          </span>
        )
      },
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "category_name",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Kategori"
        />
      ),
    },
    {
      accessorKey: "title",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Judul"
        />
      ),
    },
    {
      accessorKey: "author",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Penulis"
        />
      ),
    },
    {
      accessorKey: "published_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Tanggal Publikasi"
        />
      ),
      cell: ({ row }) => {
        return (
          <span>
            {formatTanggal(row.original.published_at)}
          </span>
        )
      },
    },
    {
      accessorKey: "is_active",
      id: "Status",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Status"
        />
      ),
      cell: ({ row }) => {
        return (
          <StatusToggle
            endpoint="/v0/articles"
            id={row.original.id}
            initialActive={row.original.is_active == "1"}
          />
        );
      },
    },
    {
      id: "actions",
      header: "Aksi",
      cell: ({ row }) => {
        const article = row.original;

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem
                onClick={() =>
                  copyToClipboard(`${process.env.NEXT_PUBLIC_APP_BASE_URL}/articles/${article.slug}`)
                }
              >
                Copy URL
              </DropdownMenuItem>
              {article.is_active ? (
                <DropdownMenuItem>Unpublish</DropdownMenuItem>
              ) : (
                <DropdownMenuItem>Publish</DropdownMenuItem>
              )}
              <DropdownMenuItem>Preview</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <Link
                  className="w-full"
                  href={`/admin/articles/actions/${article.id}`}
                >
                  Edit
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem>Delete</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];
}
