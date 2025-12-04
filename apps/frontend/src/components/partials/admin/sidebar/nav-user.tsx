"use client";

import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  Gamepad2,
  LogOut,
  Sparkles,
} from "lucide-react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { truncateText } from "@/utils/str";
import { navUserAdminDropdownContent } from "@/data/dropdown-content/nav-user";
import { useAuthStore } from "@/stores/auth-store";
import { useEffect } from "react";
import Link from "next/link";

export function NavUser({ user }: { user: any }) {
  const { isMobile } = useSidebar();
  const { user: authUser, fetchUser } = useAuthStore();

  useEffect(() => {
    if (user && !authUser) {
      fetchUser();
    }
  }, []);
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="h-8 w-8 rounded-lg flex items-center justify-center">
                <Gamepad2 className="w-6 h-6 text-primary" />
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">
                  {truncateText(authUser?.name, 16)}
                </span>
                <span className="truncate text-xs">{authUser?.email}</span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            {navUserAdminDropdownContent.map((item, index) => {
              if (item.component && !item.separator) {
                return (
                  <DropdownMenuItem key={index + "-nav-user-component"}>
                    <item.component />
                  </DropdownMenuItem>
                );
              }
              if (item.separator) {
                return (
                  <DropdownMenuSeparator key={index + "-nav-user-separator"} />
                );
              }
              return (
                <DropdownMenuItem key={index + "-nav-user-link"} asChild>
                  <Link
                    href={item.url || "#"}
                    className="flex items-center gap-2 px-2 py-1 rounded-lg hover:bg-muted/50"
                  >
                    {item.icon && <item.icon className="mr-2 h-4 w-4" />}
                    {item.title}
                  </Link>
                </DropdownMenuItem>
              );
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
