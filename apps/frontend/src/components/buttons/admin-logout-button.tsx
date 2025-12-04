"use client"

import { useTransition } from "react"
import { cn } from "@/lib/utils"
import { LogOut } from "lucide-react"
import { useLogoutActionClient } from "@/lib/auth-client"

export function AdminLogoutButton({
  className,
  ...props
}: React.ComponentProps<"button">) {
  const {logoutActionClient} = useLogoutActionClient()
  const [isPending, startTransition] = useTransition()

  return (
    <button
      onClick={logoutActionClient}
      disabled={isPending}
      className={cn("w-full flex items-center gap-4 cursor-pointer", className)}
      {...props}
    >
        <LogOut />
      {isPending ? "Mengakhiri Sesi..." : "Keluar"}
    </button>
  )
}
