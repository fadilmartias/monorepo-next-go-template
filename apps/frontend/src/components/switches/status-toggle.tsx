"use client";
import { Switch } from "@/components/ui/switch";
import { apiClient } from "@/lib/api-client";
import { useState } from "react";
import { toast } from "sonner";

export default function StatusToggle({ endpoint, id, initialActive }: { endpoint: string; id: string; initialActive: boolean }) {
    const [active, setActive] = useState(initialActive);
    const [loading, setLoading] = useState(false);
    const handleToggle = async (checked: boolean) => {
      setActive(checked);
      setLoading(true);
      try {
        const res = await apiClient.patch(`${endpoint}/${id}`, {
            is_active: checked,
        });
  
        toast.success("Status berhasil diupdate");
      } catch (error) {
        console.error(error);
        setActive(!checked); // revert kalau gagal
        toast.error("Gagal update status");
      } finally {
        setLoading(false);
      }
    };
  
    return <Switch checked={active} onCheckedChange={handleToggle} disabled={loading} />;
  }
  