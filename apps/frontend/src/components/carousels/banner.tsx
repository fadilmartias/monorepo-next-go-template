"use client";
import { useRef } from "react";
import { Card, CardContent } from "@/components/ui/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import Autoplay from "embla-carousel-autoplay";
import { Foreach } from "../utils/foreach";
import { getApiImage } from "@/lib/api-asset";
import Image from "next/image";
import Link from "next/link";
import { motion } from "framer-motion";

export default function Banner({ data }: { data: Array<any> }) {
  const plugin = useRef(Autoplay({ delay: 5000 }));

  // Variants ringan
  const variants = {
    hidden: { opacity: 0, scale: 0.98 },
    visible: { opacity: 1, scale: 1 },
  };

  return (
    <Carousel className="w-full rounded-md" plugins={[plugin.current]}>
      <CarouselContent className="rounded-md">
        <Foreach
          of={data}
          render={(item) => (
            <CarouselItem key={item.id} className="rounded-md">
              <Link href={item.link || "#"}>
                <motion.div
                  initial="hidden"
                  animate="visible"
                  exit="hidden"
                  variants={variants}
                  transition={{ duration: 0.5, ease: "easeInOut" }}
                  className="w-full rounded-md overflow-hidden"
                >
                  <Image
                    src={getApiImage("banners", item.img)}
                    alt={item.title || "Banner Next Template"}
                    width={1024}
                    height={576}
                    priority
                    className="w-full aspect-[16/7] object-cover rounded-md"
                  />
                </motion.div>
              </Link>
            </CarouselItem>
          )}
        />
      </CarouselContent>
      <CarouselPrevious className="hidden lg:flex" />
      <CarouselNext className="hidden lg:flex" />
    </Carousel>
  );
}
