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

export type ActivityLog = {
  id: string;
  log_name: string;
  event: string;
  subject_type: string;
  desc: string;
  subject_id: string;
  causer_type: string;
  causer_id: string;
  properties: string;
  trace_id: string;
  ip: string;
  user_agent: string;
  created_at: string;
};

interface ActivityLogsProps {
  serverSide?: boolean;
  sort?: string; // contoh: "title:asc"
  setSort?: (sort: string) => void;
  pagination?: APIPaginationResponse;
}

export function getActivityLogsColumns({
  serverSide = false,
  sort,
  setSort,
  pagination,
}: ActivityLogsProps): ColumnDef<ActivityLog>[] {
  return [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
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
        <DataTableColumnHeader column={column} title="No" />
      ),
      cell: ({ row }) => {
        console.log(pagination);
        const pageSize = pagination?.page_size || 10;
        const pageIndex = (pagination?.page || 1) - 1;
        return <span>{pageIndex * pageSize + row.index + 1}</span>;
      },
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "log_name",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Nama Log"
        />
      ),
    },
    {
      accessorKey: "event",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Event"
        />
      ),
    },
    {
      accessorKey: "desc",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Deskripsi"
        />
      ),
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          serverSide={serverSide}
          sort={sort}
          setSort={setSort}
          title="Tanggal"
        />
      ),
      cell: ({ row }) => {
        const activity_log = row.original;
        return <span>{formatTanggal(activity_log.created_at)}</span>;
      },
    },
    {
      id: "actions",
      header: "Aksi",
      cell: ({ row }) => {
        const article = row.original;

        return (
          <>
            <Button size="sm" variant="outline" asChild>
              <Link
                className="w-fit"
                href={`/admin/activity-logs/details/${article.id}`}
              >
                Details
              </Link>
            </Button>
          </>
        );
      },
    },
  ];
}
