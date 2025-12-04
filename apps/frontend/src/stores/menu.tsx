import { create } from "zustand";
import { Home, Newspaper, Ticket } from "lucide-react";

interface MenuItem {
  label: string;
  href: string;
  icon: React.ReactNode;
}

interface MenuState {
  menus: MenuItem[];
  setMenus: (menus: MenuItem[]) => void;
}

export const useMenuStore = create<MenuState>((set) => ({
  menus: [
    {
      label: "Beranda",
      href: "/",
      icon: <Home className="w-[18px] h-[18px]" />,
    },
    {
      label: "Cek Transaksi",
      href: "/check-transaction",
      icon: <Ticket className="w-[18px] h-[18px]" />,
    },
    {
      label: "Artikel",
      href: "/articles",
      icon: <Newspaper className="w-[18px] h-[18px]" />,
    },
    {
      label: "Leaderboard",
      href: "/leaderboard",
      icon: <Newspaper className="w-[18px] h-[18px]" />,
    },
  ],
  setMenus: (menus) => set({ menus }),
}));
