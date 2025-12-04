// app/articles/[slug]/_components/article-share.tsx
"use client"

import { Copy, Share2 } from "lucide-react"
import { FaFacebook, FaTwitter } from "react-icons/fa"
import Link from "next/link"
import { toast } from "sonner"
import { useEffect } from "react"
import Typography from "@/components/typography"

type Props = {
  title: string
  shareUrl: string
}

export default function ArticleShare({ title, shareUrl }: Props) {
  // Auto-append saat user manual copy
  useEffect(() => {
    const handleCopy = (e: ClipboardEvent) => {
      const selection = document.getSelection()
      if (!selection) return

      const copiedText = selection.toString()
      const footer = `\n\n🎮 Baca selengkapnya di: ${shareUrl}`
      e.clipboardData?.setData("text/plain", copiedText + footer)
      e.preventDefault()
    }

    document.addEventListener("copy", handleCopy)
    return () => document.removeEventListener("copy", handleCopy)
  }, [shareUrl])

  return (
    <footer className="mt-10 border-t border-white/10 pt-6">
      <Typography variant="h5" as="h2" className="flex items-center gap-2 mb-4">
        <Share2 className="w-5 h-5" />
        Bagikan Artikel
      </Typography>
      <div className="flex flex-wrap gap-4">
        <Link
          href={`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(shareUrl)}`}
          target="_blank"
          className="bg-[#1877f2] hover:bg-[#166fe5] px-4 py-2 rounded-lg text-white flex items-center gap-2"
        >
          <FaFacebook className="w-4 h-4" />
          Facebook
        </Link>
        <Link
          href={`https://twitter.com/intent/tweet?url=${encodeURIComponent(shareUrl)}&text=${encodeURIComponent(title)}`}
          target="_blank"
          className="bg-[#1da1f2] hover:bg-[#0d95e8] px-4 py-2 rounded-lg text-white flex items-center gap-2"
        >
          <FaTwitter className="w-4 h-4" />
          Twitter / X
        </Link>
        <button
          className="bg-gray-800 hover:bg-gray-700 px-4 py-2 rounded-lg text-white flex items-center gap-2"
          onClick={() => {
            navigator.clipboard.writeText(shareUrl)
            toast.success("Link artikel disalin ke clipboard!", {
              description: "Link siap dibagikan 🎯",
            })
          }}
        >
          <Copy className="w-4 h-4" />
          Salin Link
        </button>
      </div>
    </footer>
  )
}
