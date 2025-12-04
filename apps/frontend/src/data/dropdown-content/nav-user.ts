import { User, ShoppingCart, Home } from "lucide-react";
import { AdminLogoutButton } from "../../components/buttons/admin-logout-button";
import { DropdownMenuSeparator } from "@/components/ui/dropdown-menu";

type NavUserDropdownContent = {
    title?: string;
    url?: string;
    icon?: any;
    disabled?: boolean;
    badge?: "default" | "destructive" | "outline" | "secondary" | "success" | "warning";
    badgeContent?: string;
    separator?: boolean;
    component?: any;
    role?: string[];
}

export const navUserDropdownContent: NavUserDropdownContent[] = [
    {
        title: "Profil",
        url: "/profile",
        icon: User,
        role: ['user', 'admin']
    },
    {
        title: "Pesananku",
        url: "/orders",
        icon: ShoppingCart,
        role: ['user', 'admin']
    },
    {
        title: "Admin Dashboard",
        url: "/admin/dashboard",
        icon: Home,
        role: ['admin']
    },
    {
        separator: true,
        component: DropdownMenuSeparator,
    },
    {
        component: AdminLogoutButton,
    },
];

export const navUserAdminDropdownContent: NavUserDropdownContent[] = [
    {
        title: "Profil",
        url: "/admin/profile",
        icon: User,
        role: ['admin']
    },
    {
        title: "Beranda",
        url: "/",
        icon: Home,
        role: ['admin']
    },
    {
        separator: true,
        component: DropdownMenuSeparator,
    },
    {
        component: AdminLogoutButton,
    },
];