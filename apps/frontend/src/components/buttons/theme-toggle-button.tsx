"use client";

import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTheme } from "next-themes";
import { useMounted } from "@/hooks/use-mounted";
import { Skeleton } from "../ui/skeleton";

export function ThemeToggleButton({ className }: { className?: string }) {
  const { theme, setTheme } = useTheme();
  const mounted = useMounted();

  if (!mounted) return <Skeleton className="w-8 h-8" />;

  return (
    <Button
      variant="ghost"
      size="icon"
      aria-label="Ganti tema"
      title="Ganti tema"
      onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
      className={className}
    >
      {theme === "dark" ? (
        <Sun className="h-5 w-5" />
      ) : (
        <Moon className="h-5 w-5" />
      )}
    </Button>
  );
}
