"use client";

import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  useReactTable,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { DataTablePagination } from "./data-table-pagination";
import { useEffect, useState } from "react";
import { DataTableSearch } from "./data-table-search";
import { DataTableViewOptions } from "./data-table-view-options";
import { Loader2 } from "lucide-react";
import { APIPaginationResponse } from "@/types/api";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  isLoading?: boolean;
  isPagination?: boolean;
  columnsSearch?: string[];
  serverSide?: boolean;
  pagination?: APIPaginationResponse;
  sort?: string;
  setSort?: (sort: string) => void;
  filter?: { [key: string]: any };
  setFilter?: (filter: { [key: string]: any }) => void;
  nextPage?: () => void;
  prevPage?: () => void;
  goToPage?: (page: number) => void;
  goToLastPage?: () => void;
  goToFirstPage?: () => void;
  setPageSize?: (pageSize: number) => void;
}

export function DataTable<TData, TValue>({
  columns,
  isLoading,
  data,
  isPagination = true,
  columnsSearch,
  serverSide = false,
  pagination,
  sort,
  setSort,
  filter,
  setFilter,
  nextPage,
  prevPage,
  goToPage,
  goToLastPage,
  goToFirstPage,
  setPageSize,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});
  const [lastData, setLastData] = useState<TData[]>([]);
  const table = useReactTable({
    data: lastData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    manualPagination: serverSide,
    manualSorting: serverSide,
    manualFiltering: serverSide,
    pageCount: serverSide ? pagination?.total_pages : undefined,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  });

  useEffect(() => {
    data && setLastData(data);
  }, [data]);

  return (
    <>
      <div className="flex items-center py-4">
        {columnsSearch && columnsSearch.length > 0 && (
          <DataTableSearch
            columns={columnsSearch}
            disabled={!data}
            table={table}
            serverSide={serverSide}
            filter={filter}
            setFilter={setFilter}
          />
        )}
        <DataTableViewOptions disabled={!data} table={table} />
      </div>
      <div className="relative">
        {!data && (
          <div className="absolute inset-0 flex items-center justify-center rounded-md opacity-100 bg-white/10 z-10">
            {isLoading && (
              <Loader2 className="animate-spin w-8 h-8 text-white" />
            )}
          </div>
        )}
        <div
          className={`overflow-hidden rounded-md border mb-6 ${
            !data ? "opacity-50" : "opacity-100"
          }`}
        >
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHead key={header.id}>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                      </TableHead>
                    );
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    No results.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
        {isPagination && lastData.length > 0 && (
          <DataTablePagination
            serverSide={serverSide}
            table={table}
            pagination={pagination}
            nextPage={nextPage}
            prevPage={prevPage}
            goToPage={goToPage}
            goToLastPage={goToLastPage}
            goToFirstPage={goToFirstPage}
            setPageSize={setPageSize}
          />
        )}
      </div>
    </>
  );
}
