import { Column } from "@tanstack/react-table"
import { ArrowDown, ArrowUp, ChevronsUpDown, EyeOff } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

interface DataTableColumnHeaderProps<TData, TValue>
  extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>
  title: string
}

interface DataTableColumnHeaderProps<TData, TValue>
  extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>
  title: string
  serverSide?: boolean
  sort?: string
  setSort?: (sort: string) => void
}

export function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  className,
  serverSide = false,
  sort,
  setSort,
}: DataTableColumnHeaderProps<TData, TValue>) {
  const currentSort = column.getIsSorted() as false | "asc" | "desc";

  const handleSort = (direction: "asc" | "desc") => {
    if (serverSide && setSort) {
      const sortKey = `${column.id}:${direction}`;
      setSort(sortKey);
    } else {
      column.toggleSorting(direction === "desc");
    }
  };

  const handleHide = () => {
    column.toggleVisibility(false);
  };

  const sortIcon =
    currentSort === "desc" ? <ArrowDown /> :
    currentSort === "asc"  ? <ArrowUp /> :
    <ChevronsUpDown />;

  return (
    <div className={cn("flex items-center gap-2", className)}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            className="data-[state=open]:bg-accent -ml-3 h-8"
          >
            <span>{title}</span>
            {column.getCanSort() && sortIcon}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start">
          {/* Tampilkan menu sorting hanya jika column bisa sort */}
          {column.getCanSort() && (
            <>
              <DropdownMenuItem onClick={() => handleSort("asc")}>
                <ArrowUp className="mr-2 h-4 w-4" />
                Asc
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleSort("desc")}>
                <ArrowDown className="mr-2 h-4 w-4" />
                Desc
              </DropdownMenuItem>
              <DropdownMenuSeparator />
            </>
          )}

          {/* Tampilkan menu hide hanya jika column bisa hide */}
          {column.getCanHide() && (
            <DropdownMenuItem onClick={handleHide}>
              <EyeOff className="mr-2 h-4 w-4" />
              Hide
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}


