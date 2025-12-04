import { create } from "zustand";

interface BreadcrumbState {
    breadcrumbs: { title: string; url: string; subtitle?: string; backButton?: boolean }[];
    setBreadcrumbs: (breadcrumbs: { title: string; url: string; subtitle?: string; backButton?: boolean }[]) => void;
}

export const useBreadcrumbStore = create<BreadcrumbState>((set) => ({
    breadcrumbs: [],
    setBreadcrumbs: (breadcrumbs) => set({ breadcrumbs }),
}));
