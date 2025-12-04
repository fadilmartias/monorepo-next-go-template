import { Table } from "@tanstack/react-table"
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { APIPaginationResponse } from "@/types/api"

interface DataTablePaginationProps<TData> {
  table: Table<TData>
  serverSide?: boolean
  pagination?:APIPaginationResponse
  nextPage?: () => void
  prevPage?: () => void
  goToPage?: (page: number) => void
  goToLastPage?: () => void
  goToFirstPage?: () => void
  setPageSize?: (pageSize: number) => void
}

export function DataTablePagination<TData>({
  table,
  serverSide,
  pagination,
  nextPage,
  prevPage,
  goToPage,
  goToLastPage,
  goToFirstPage,
  setPageSize,
}: DataTablePaginationProps<TData>) {
  const toFirstPage = () => {
    goToFirstPage ? goToFirstPage() : table.setPageIndex(0)
  }
  const toLastPage = () => {
    goToLastPage ? goToLastPage() : table.setPageIndex(table.getPageCount() - 1)
  }
  const toPreviousPage = () => {
    prevPage ? prevPage() : table.previousPage()
  }
  const toNextPage = () => {
    nextPage ? nextPage() : table.nextPage()
  }
  return (
    <div className="flex items-center justify-between px-2">
      <div className="text-muted-foreground flex-1 text-sm">
        {table.getFilteredSelectedRowModel().rows.length} of{" "}
        {table.getFilteredRowModel().rows.length} row(s) selected.
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select
            value={`${setPageSize ? pagination?.page_size : table.getState().pagination.pageSize}`}
            onValueChange={(value) => {
              setPageSize ? setPageSize(Number(value)) : table.setPageSize(Number(value))
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={table.getState().pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 25, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {serverSide ? pagination?.page : table.getState().pagination.pageIndex + 1} of{" "}
          {table.getPageCount()}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="icon"
            className="hidden size-8 lg:flex"
            onClick={toFirstPage}
            disabled={serverSide ? pagination?.page === 1 : !table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to first page</span>
            <ChevronsLeft />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="size-8"
            onClick={toPreviousPage}
            disabled={serverSide ? pagination?.page === 1 : !table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeft />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="size-8"
            onClick={toNextPage}
            disabled={serverSide ? pagination?.page === pagination?.total_pages : !table.getCanNextPage()}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRight />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="hidden size-8 lg:flex"
            onClick={toLastPage}
            disabled={serverSide ? pagination?.page === pagination?.total_pages : !table.getCanNextPage()}
          >
            <span className="sr-only">Go to last page</span>
            <ChevronsRight />
          </Button>
        </div>
      </div>
    </div>
  )
}
