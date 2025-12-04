"use client"

import { useTransition } from "react"
import { useRouter } from "nextjs-toploader/app"
import { cn } from "@/lib/utils"
import { useLogoutActionClient } from "@/lib/auth-client"
import { Button } from "../ui/button"

export function LogoutButton({
  className,
  ...props
}: React.ComponentProps<"button">) {
  const router = useRouter()
  const { logoutActionClient } = useLogoutActionClient()
  const [isPending, startTransition] = useTransition()

  return (
    <Button
      onClick={logoutActionClient}
      disabled={isPending}
      className={cn("w-full", className)}
      {...props}
    >
      {isPending ? "Logging out..." : "Logout"}
    </Button>
  )
}
