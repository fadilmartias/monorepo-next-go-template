"use client";

import { toast } from "sonner";

export function copyToClipboard(text: string): void {
    navigator.clipboard.writeText(text.trim());
    toast.success("Copied to clipboard");
}

export async function pasteFromClipboard() : Promise<string> {
    try {
      const text = await navigator.clipboard.readText();
      toast.success("Pasted from clipboard");
      return text?.trim();
    } catch (err) {
      toast.error("Gagal mengambil teks dari clipboard");
      console.error("Gagal mengambil teks dari clipboard:", err);
      return "";
    }
}