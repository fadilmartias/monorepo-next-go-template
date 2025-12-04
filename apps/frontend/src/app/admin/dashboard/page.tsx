import { AppSidebar } from "@/components/partials/admin/sidebar/app-sidebar"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import DashboardClient from "./DashboardClient"
import { useBreadcrumbStore } from "@/stores/breadcrumb"

export const metadata = {
  title: "Dashboard",
}

export default function Page() {
  

  return (
    <DashboardClient />
  )
}
