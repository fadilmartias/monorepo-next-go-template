"use client";
import { ArrowLeft } from "lucide-react";
import { Button } from "./ui/button";
import { CardTitle } from "./ui/card";
import { useBreadcrumbStore } from "@/stores/breadcrumb";

export default function CardTitleAdminLayout() {
  const initialBreadcrumbs = useBreadcrumbStore((state) => state.breadcrumbs);
  return (
    <CardTitle className="text-2xl flex items-center gap-3 font-semibold">
      {initialBreadcrumbs[initialBreadcrumbs.length - 1]?.backButton && (
        <Button
          variant="ghost"
          size="lg"
          className=""
          onClick={() => window.history.back()}
        >
          <ArrowLeft className="w-8 h-8" />
        </Button>
      )}
      <div className="flex flex-col gap-1">
      {initialBreadcrumbs[initialBreadcrumbs.length - 1]?.title}
      {initialBreadcrumbs[initialBreadcrumbs.length - 1]?.subtitle && (
        <span className="text-muted-foreground text-sm">
          {initialBreadcrumbs[initialBreadcrumbs.length - 1]?.subtitle}
        </span>
      )}
      </div>
    </CardTitle>
  );
}
