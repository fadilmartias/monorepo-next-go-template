import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { AppSidebar } from "@/components/partials/admin/sidebar/app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import React from "react";
import { getUserFromToken } from "@/lib/auth";
import { useBreadcrumbStore } from "@/stores/breadcrumb";
import BreadcrumbNav from "@/components/breadcrumb-nav";
import CardTitleAdminLayout from "@/components/card-title-admin-layout";

export const metadata = {
  robots: {
    index: false,
    follow: false,
  },
};

export default async function AdminLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const user = await getUserFromToken();

  return (
    <>
      <SidebarProvider>
        <AppSidebar user={user} />
        <SidebarInset className="overflow-hidden">
          <header className="flex h-16 shrink-0 items-center gap-2">
            <div className="flex items-center gap-2 px-4">
              <SidebarTrigger className="-ml-1" />
              <Separator
                orientation="vertical"
                className="mr-2 data-[orientation=vertical]:h-4"
              />
              <BreadcrumbNav />
            </div>
          </header>
          <main className="p-4">
            <Card className="overflow-hidden!">
              <CardHeader>
                <CardTitleAdminLayout />
              </CardHeader>
              <CardContent>{children}</CardContent>
            </Card>
          </main>
        </SidebarInset>
      </SidebarProvider>
    </>
  );
}
