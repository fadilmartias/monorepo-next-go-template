import * as React from "react"

import { cn } from "@/lib/utils"
import { Label } from "./label"
import Typography from "../typography"

function Textarea({ className, label, id, error, required, requiredSign = true, ...props }: React.ComponentProps<"textarea"> & { error?: string, label?: string, id?: string, required?: boolean, requiredSign?: boolean }) {
  return (
    <div className="flex flex-col gap-2 w-full">
    {label && <Label htmlFor={id || label || ""}>{label}{required && requiredSign && <span className="text-destructive">*</span>}</Label>}
    <textarea
      id={id || label || ""}
      data-slot="textarea"
      className={cn(
        "border-input placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:bg-input/30 flex field-sizing-content min-h-16 w-full rounded-md border bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
        className,
        error && "border-destructive"
      )}
      {...props}
    />
     {error && <Typography variant="c1" className="text-destructive">{error}</Typography>}
  </div>
  )
}

export { Textarea }
