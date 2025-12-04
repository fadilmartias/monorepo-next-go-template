"use client";

import { FaWhatsapp } from "react-icons/fa";
import { useEffect, useState } from "react";
import { appConfig } from "@/config/app";
import Link from "next/link";
import { Button } from "../ui/button";

export default function FloatingContactButton() {
  const [rightOffset, setRightOffset] = useState(16); // jarak dari kanan (px)
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    function updateOffset() {
      const container = document.querySelector(".container-custom");
      if (!container) return;

      const windowWidth = window.innerWidth;
      const containerWidth = container.clientWidth;
      const sideSpace = (windowWidth - containerWidth) / 2;

      // 16px = 1rem = tailwind bottom-4/right-4
      setRightOffset(sideSpace + 16);
    }

    updateOffset();
    setMounted(true);
    window.addEventListener("resize", updateOffset);
    return () => window.removeEventListener("resize", updateOffset);
  }, []);

  if (!mounted) return null;

  return (
      <Link
        href={`https://wa.me/${appConfig.whatsappNumber}`}
        role="button"
        target="_blank"
        aria-label="Hubungi kami di WhatsApp"
        title="Hubungi kami di WhatsApp"
        rel="noopener noreferrer"
        style={{ right: rightOffset }}
        className="fixed bottom-4 z-50"
      >
        <div className="w-10 h-10 md:w-12 md:h-12 lg:w-14 lg:h-14 bg-green-500 text-white rounded-full flex items-center justify-center shadow-lg hover:scale-110 transition-transform">
          <FaWhatsapp size={28} />
        </div>
      </Link>
  );
}
