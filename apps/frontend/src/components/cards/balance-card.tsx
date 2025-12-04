// components/BalanceCard.tsx
"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Wallet, Star, ArrowUpRight, ArrowDownLeft } from "lucide-react";
import { motion } from "framer-motion";

export default function BalanceCard() {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
      className="w-full"
    >
      <Card className="w-full max-w-sm shadow-md rounded-2xl bg-linear-to-br from-primary/90 to-primary/70 text-primary-foreground">
        <CardContent className="p-3 sm:p-4">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <p className="text-[10px] sm:text-xs opacity-90">Saldo</p>
              <div className="flex items-center gap-1">
                <motion.span
                  key="balance"
                  initial={{ scale: 0.9, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="text-lg sm:text-xl font-bold"
                >
                  Rp 1.250.000
                </motion.span>
                <Badge
                  variant="secondary"
                  className="text-[10px] sm:text-xs px-1 py-0 bg-white/20 text-primary-foreground"
                >
                  Aktif
                </Badge>
              </div>
            </div>
            <div className="bg-white/20 p-2 rounded-lg">
              <Wallet className="h-4 w-4 sm:h-5 sm:w-5" />
            </div>
          </div>

          {/* Points Section */}
          <div className="flex items-center justify-between mt-3 border-t border-white/20 pt-2">
            <div className="flex items-center gap-2">
              <div className="bg-yellow-400/20 p-1.5 rounded-md">
                <Star className="h-3.5 w-3.5 text-yellow-300" />
              </div>
              <div>
                <p className="text-[10px] opacity-80">Poin</p>
                <p className="text-sm font-semibold">1.200</p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="h-6 px-2 text-[11px] hover:bg-white/10"
            >
              Tukar
            </Button>
          </div>

          {/* Actions */}
          <div className="grid grid-cols-2 gap-2 mt-3">
            <motion.div whileTap={{ scale: 0.95 }}>
              <Button
                variant="secondary"
                className="h-8 text-xs bg-white/20 hover:bg-white/30 border-white/30 w-full"
              >
                <ArrowDownLeft className="mr-1 h-3 w-3" />
                Top Up
              </Button>
            </motion.div>
            <motion.div whileTap={{ scale: 0.95 }}>
              <Button
                variant="outline"
                className="h-8 text-xs bg-white hover:bg-white/90 text-primary border-white w-full"
              >
                <ArrowUpRight className="mr-1 h-3 w-3" />
                Transfer
              </Button>
            </motion.div>
          </div>

          {/* Footer */}
          <div className="mt-2 flex justify-between text-[11px] opacity-75">
            <span>Update: 5m lalu</span>
            <span className="underline cursor-pointer">Riwayat</span>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}
