"use client";

import * as React from "react";
import {
  BookOpen,
  Command,
  DollarSign,
  Home,
  AlertCircle,
  Bell,
  Settings,
  FileText,
  Image,
  HelpCircle,
  Tag,
  Package,
  Percent,
  Ticket,
  Users,
  History,
} from "lucide-react";

import { NavMain } from "@/components/partials/admin/sidebar/nav-main";
import { NavSecondary } from "@/components/partials/admin/sidebar/nav-secondary";
import { NavUser } from "@/components/partials/admin/sidebar/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import Link from "next/link";
import Logo from "@/components/logo";

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  menu: [
    {
      title: "Dasbor",
      url: "/admin/dashboard",
      icon: Home,
      isActive: true,
    },
    {
      title: "Transaksi",
      url: "/admin/transactions",
      icon: BookOpen,
      items: [
        {
          title: "Riwayat Transaksi",
          url: "/admin/transactions/history",
          icon: History,
        },
        {
          title: "Riwayat Deposit",
          url: "/admin/transactions/deposit",
          icon: DollarSign,
        },
        {
          title: "Permintaan Refund",
          url: "/admin/transactions/refund",
          icon: AlertCircle,
        },
        {
          title: "Tambah Transaksi DigiFlazz",
          url: "/admin/transactions/create-digiflazz-transaction",
          icon: AlertCircle,
        },
      ],
    },
  ],
  secondary: [
    {
      title: "Log Aktifitas",
      url: "/admin/activity-logs",
      icon: History,
    },
    {
      title: "Pemberitahuan",
      url: "/admin/notifications",
      icon: Bell,
    },
    {
      title: "Pengaturan",
      url: "/admin/setting",
      icon: Settings,
    },
  ],
  cms: [
    {
      title: "Artikel",
      url: "/admin/articles",
      icon: FileText,
    },
    {
      title: "Banner",
      url: "/admin/banners",
      icon: Image,
    },
    {
      title: "FAQ",
      url: "/admin/faq",
      icon: HelpCircle,
    },
  ],
  masters: [
    {
      title: "Kategori",
      url: "/admin/categories",
      icon: Tag,
    },
    {
      title: "Produk",
      url: "/admin/products",
      icon: Package,
      items: [
        {
          title: "Prabayar",
          url: "/admin/products/prepaid",
          icon: Package,
        },
        {
          title: "Pascabayar",
          url: "/admin/products/postpaid",
          icon: Package,
        },
      ],
    },
    {
      title: "Diskon",
      url: "/admin/discounts",
      icon: Percent,
    },
    {
      title: "Voucher",
      url: "/admin/vouchers",
      icon: Ticket,
    },
    {
      title: "Pengguna",
      url: "/admin/users",
      icon: Users,
    },
  ],
};

export function AppSidebar({user, ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="inset" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <div className="flex items-center justify-center">
                <Logo href="/admin/dashboard" />
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.menu} title="Menu" />
        <NavMain items={data.masters} title="Master" />
        <NavMain items={data.cms} title="CMS" />
        <NavSecondary items={data.secondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
    </Sidebar>
  );
}
