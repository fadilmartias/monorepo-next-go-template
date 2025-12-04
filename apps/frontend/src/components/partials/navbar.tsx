"use client";

import * as React from "react";

import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu";
import { usePathname } from "next/navigation";
import { Button } from "../ui/button";
import Link from "next/link";
import { useMenuStore } from "@/stores/menu";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Bell,
  ChevronsUpDown,
  Gamepad2,
  MenuIcon,
  User,
  XIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { LogoutButton } from "../buttons/logout-button";
import { ThemeToggleButton } from "../buttons/theme-toggle-button";
import { AdminLogoutButton } from "../buttons/admin-logout-button";
import { truncateText } from "@/utils/str";
import Logo from "../logo";
import { navUserDropdownContent } from "@/data/dropdown-content/nav-user";
import { useAuthStore } from "@/stores/auth-store";
import { useEffect } from "react";
import { Skeleton } from "../ui/skeleton";
import { Badge } from "../ui/badge";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";
import Typography from "../typography";

export default function Navbar({
  isLoggedIn,
  user,
}: {
  isLoggedIn: boolean;
  user: any;
}) {
  const { user: authUser, fetchUser } = useAuthStore();
  const pathname = usePathname();
  const { menus } = useMenuStore();

  useEffect(() => {
    if (isLoggedIn && !authUser) {
      fetchUser();
    }
  }, [isLoggedIn]);

  return (
    <header className="container-custom navbar py-3 flex items-center justify-between sticky top-0 z-20 bg-background/80 backdrop-blur">
      <Logo size="lg" />
      <NavigationMenu className="w-full hidden sm:flex" viewport={false}>
        <NavigationMenuList>
          {menus.map((menu) => (
            <NavigationMenuItem key={menu.label}>
              <NavigationMenuLink asChild>
                <Link
                  href={menu.href}
                  className={cn(
                    "block text-sm focus:bg-muted focus:text-primary font-medium px-3 py-2 rounded-md border border-transparent",
                    "hover:bg-muted hover:text-primary transition-all duration-200",
                    pathname === menu.href && "text-primary bg-muted font-bold"
                  )}
                >
                  {menu.label}
                </Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
          ))}
        </NavigationMenuList>
      </NavigationMenu>
      <div className="flex items-center gap-3">
        <ThemeToggleButton className="hidden sm:flex" />
        {isLoggedIn ? (
          <div className="hidden sm:flex">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  className="flex items-center gap-2 px-2 py-1 rounded-lg hover:bg-muted/50"
                >
                  <Avatar className="h-8 w-8 rounded-lg flex items-center justify-center">
                    <Gamepad2 className="w-8 h-8 text-primary" />
                  </Avatar>
                  <div className="hidden sm:grid text-left text-sm leading-tight">
                    <span className="truncate font-medium">
                      {!authUser ? (
                        <Skeleton className="w-24 h-6" />
                      ) : (
                        <div className="flex flex-col">
                          <Typography
                            variant="b3"
                            as="p"
                            className="truncate font-medium text-sm"
                          >
                            {truncateText(authUser?.name, 12)}
                          </Typography>
                          {/* {authUser?.active_title && (
                            <p className="text-xs text-muted-foreground">
                              {authUser.active_title.name}
                            </p>
                          )} */}
                          <Typography
                            variant="c1"
                            as="p"
                            className="text-muted-foreground"
                          >
                            Level {authUser?.level}
                          </Typography>
                        </div>
                      )}
                    </span>
                  </div>
                </Button>
              </DropdownMenuTrigger>

              <DropdownMenuContent
                className="min-w-56 rounded-lg"
                side="bottom"
                align="end"
                sideOffset={4}
              >
                {navUserDropdownContent
                  .filter((item) => {
                    // Kalau item gak punya role, berarti boleh untuk semua
                    if (!item.role) return true;
                    // Kalau user belum login, skip item yang butuh role
                    if (!authUser?.role) return false;
                    // Hanya tampil kalau role user termasuk role item
                    return item.role.includes(authUser.role);
                  })
                  .map((item, index) => {
                    if (item.component && !item.separator) {
                      return (
                        <DropdownMenuItem key={index + "-nav-user-component"}>
                          <item.component />
                        </DropdownMenuItem>
                      );
                    }

                    if (item.separator) {
                      return (
                        <DropdownMenuSeparator
                          key={index + "-nav-user-separator"}
                        />
                      );
                    }

                    return (
                      <DropdownMenuItem
                        key={index + "-nav-user-link"}
                        disabled={item.disabled}
                        asChild
                      >
                        <Link
                          href={item.url || "#"}
                          className="flex items-center gap-2 px-2 py-1 rounded-lg hover:bg-muted/50"
                        >
                          {item.icon && <item.icon className="mr-2 h-4 w-4" />}
                          {item.title}
                          {item.badge && (
                            <Badge variant={item.badge} className="ml-2">
                              {item.badgeContent}
                            </Badge>
                          )}
                        </Link>
                      </DropdownMenuItem>
                    );
                  })}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        ) : (
          <div className="space-x-3 hidden sm:flex">
            <Button asChild variant="outline">
              <Link href="/register">Daftar</Link>
            </Button>
            <Button asChild variant="default">
              <Link href="/login">Masuk</Link>
            </Button>
          </div>
        )}
      </div>

      <div className="items-center gap-2 flex sm:hidden">
        <div className="flex sm:hidden">
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="hover:bg-muted">
                <MenuIcon className="w-5 h-5" />
              </Button>
            </SheetTrigger>
            <VisuallyHidden>
              <SheetTitle>Menu</SheetTitle>
            </VisuallyHidden>

            <SheetContent
              side="right"
              className="w-64 p-6 flex flex-col bg-background"
              aria-describedby="menu"
              aria-description="Menu"
            >
              {/* Main Menu */}
              <div className="flex-1 py-6 space-y-1">
                {menus.map((menu, index) => (
                  <SheetClose asChild key={index + "-menu"}>
                    <Link
                      href={menu.href}
                      className={cn(
                        "block text-sm font-medium px-3 py-2 rounded-md",
                        "hover:bg-muted hover:text-primary transition-colors duration-200",
                        pathname === menu.href && "bg-muted text-primary"
                      )}
                    >
                      {menu.label}
                    </Link>
                  </SheetClose>
                ))}
                {isLoggedIn &&
                  navUserDropdownContent
                    .filter((item) => {
                      // Kalau item gak punya role, berarti boleh untuk semua
                      if (!item.role) return true;
                      // Kalau user belum login, skip item yang butuh role
                      if (!authUser?.role) return false;
                      // Hanya tampil kalau role user termasuk role item
                      return item.role.includes(authUser.role);
                    })
                    .map((item, index) => (
                      <SheetClose asChild key={index + "-nav-user"}>
                        <Link
                          href={item.url || "#"}
                          className={cn(
                            "block text-sm font-medium px-3 py-2 rounded-md",
                            "hover:bg-muted hover:text-primary transition-colors duration-200",
                            pathname === item.url && "bg-muted text-primary"
                          )}
                        >
                          {item.title}
                        </Link>
                      </SheetClose>
                    ))}
              </div>

              {/* User Section */}
              <ThemeToggleButton className="" />
              <div className="border-t pt-4">
                {isLoggedIn ? (
                  <div className="space-y-4">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8 rounded-lg flex items-center justify-center">
                        <Gamepad2 className="w-8 h-8 text-primary" />
                      </Avatar>
                      <div className="flex flex-col">
                        <span className="text-sm font-medium truncate max-w-[140px] line-clamp-1">
                          {authUser?.name}
                        </span>
                        <span className="text-xs text-muted-foreground line-clamp-1">
                          {authUser?.email}
                        </span>
                      </div>
                    </div>
                    <LogoutButton className="w-full" />
                  </div>
                ) : (
                  <div className="flex flex-col space-y-3">
                    <Button asChild variant="default" size="sm">
                      <Link href="/login">Masuk</Link>
                    </Button>
                    <Button asChild variant="outline" size="sm">
                      <Link href="/register">Daftar</Link>
                    </Button>
                  </div>
                )}
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  );
}
