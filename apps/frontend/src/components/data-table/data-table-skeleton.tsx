import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table"

export default function DataTableSkeleton({ columns = 5, rows = 5 }) {
  return (
    <div>
      <div className="overflow-hidden rounded-md border mb-6">
        <Table>
          {/* 🔹 Header skeleton */}
          <TableHeader>
            <TableRow>
              {Array.from({ length: columns }).map((_, i) => (
                <TableHead key={i}>
                  <div className="h-4 w-20 bg-muted animate-pulse rounded" />
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>

          {/* 🔹 Body skeleton */}
          <TableBody>
            {Array.from({ length: rows }).map((_, rowIndex) => (
              <TableRow key={rowIndex}>
                {Array.from({ length: columns }).map((_, colIndex) => (
                  <TableCell key={colIndex} className="py-4">
                    <div className="h-4 w-full bg-muted animate-pulse rounded" />
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* 🔹 Pagination skeleton */}
      <div className="flex justify-between items-center">
        <div className="h-5 w-32 bg-muted animate-pulse rounded" /> 
        <div className="flex gap-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-9 w-9 bg-muted animate-pulse rounded-md" />
          ))}
        </div>
      </div>
    </div>
  )
}
