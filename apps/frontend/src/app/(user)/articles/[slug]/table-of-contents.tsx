"use client";

import { useEffect, useState } from "react";
import Typography from "@/components/typography";

interface TocItem {
  id: string;
  text: string;
  level: number;
}

export default function TableOfContents({ content }: { content: string }) {
  const [headings, setHeadings] = useState<TocItem[]>([]);
  const [activeId, setActiveId] = useState<string>("");

  useEffect(() => {
    // Tunggu sebentar agar DOM sudah ter-render
    const timer = setTimeout(() => {
      // Ambil heading langsung dari DOM artikel yang sudah di-render
      const articleElement = document.querySelector("article .prose");
      if (!articleElement) return;

      const headingElements = articleElement.querySelectorAll("h2, h3");
      
      const items: TocItem[] = [];
      
      Array.from(headingElements).filter((heading) => heading.textContent).forEach((heading, index) => {
        const headingEl = heading as HTMLElement;
        
        // Generate ID dari text atau gunakan index
        const text = headingEl.textContent || "";
        const id = `heading-${index}-${text
          .toLowerCase()
          .replace(/[^a-z0-9]+/g, "-")
          .substring(0, 50)}`;
        
        // Set ID ke element di DOM (ini yang penting!)
        headingEl.id = id;
        
        items.push({
          id,
          text,
          level: parseInt(headingEl.tagName.substring(1)),
        });
      });

      setHeadings(items);

      // Fungsi untuk menentukan active heading
      const determineActiveHeading = () => {
        // Offset dari top browser (adjust sesuai kebutuhan)
        const scrollOffset = 150;
        const scrollTop = window.scrollY + scrollOffset;

        // Array untuk store heading positions
        const headingPositions: { id: string; top: number }[] = [];

        headingElements.forEach((heading) => {
          const element = heading as HTMLElement;
          if (element.id) {
            headingPositions.push({
              id: element.id,
              top: element.offsetTop,
            });
          }
        });

        // Sort by position (ascending)
        headingPositions.sort((a, b) => a.top - b.top);

        // Cari heading yang tepat
        let currentActiveId = headingPositions[0]?.id || "";

        for (let i = 0; i < headingPositions.length; i++) {
          const heading = headingPositions[i];
          
          // Jika scroll position sudah melewati heading ini
          if (scrollTop >= heading.top) {
            currentActiveId = heading.id;
          } else {
            // Sudah ketemu heading yang belum di-pass, stop loop
            break;
          }
        }

        setActiveId(currentActiveId);
      };

      // Initial check
      determineActiveHeading();

      // Scroll listener dengan throttle
      let isThrottled = false;
      const handleScroll = () => {
        if (isThrottled) return;
        
        isThrottled = true;
        requestAnimationFrame(() => {
          determineActiveHeading();
          isThrottled = false;
        });
      };

      window.addEventListener("scroll", handleScroll, { passive: true });

      return () => {
        window.removeEventListener("scroll", handleScroll);
      };
    }, 300);

    return () => clearTimeout(timer);
  }, [content]);

  if (headings.length === 0) return null;

  const handleClick = (id: string) => {
    const element = document.getElementById(id);
    
    if (element) {
      const offset = 100;
      const elementPosition = element.getBoundingClientRect().top;
      const offsetPosition = elementPosition + window.pageYOffset - offset;

      window.scrollTo({
        top: offsetPosition,
        behavior: "smooth",
      });
      
      // Update active state immediately
      setActiveId(id);
    }
  };

  return (
    <nav className="bg-muted/50 rounded-xl p-4">
      <Typography variant="b2" className="font-semibold mb-3 block">
        📑 Daftar Isi
      </Typography>
      <ul className="space-y-2 text-sm">
        {headings.map((heading) => (
          <li
            key={heading.id}
            style={{ paddingLeft: `${(heading.level - 2) * 12}px` }}
          >
            <button
              onClick={() => handleClick(heading.id)}
              className={`text-left w-full hover:text-primary transition-colors line-clamp-2 cursor-pointer ${
                activeId === heading.id
                  ? "text-primary font-medium"
                  : "text-muted-foreground"
              }`}
            >
              {heading.text}
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
}