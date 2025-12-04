"use client";

import { useState } from "react";
import { Star } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { ReviewType } from "@/app/(user)/payment/[order_id]/payment-page-content";
import { Controller, UseFormReturn, useWatch } from "react-hook-form";

interface ReviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  reviewForm: UseFormReturn<ReviewType>;
  onSubmit: (review: ReviewType) => void;
}

export default function ReviewModal({
  isOpen,
  onClose,
  reviewForm,
  onSubmit,
}: ReviewModalProps) {

  const handleSubmit = (data: ReviewType) => {
    onSubmit(data);
  };

  const rating = useWatch({
    control: reviewForm.control,
    name: "rating",
  });

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[500px] bg-card border-2 border-primary/20">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold bg-linear-to-r from-primary to-accent bg-clip-text text-transparent">
            Berikan Ulasan
          </DialogTitle>
          <DialogDescription className="text-muted-foreground">
            Bagikan pengalaman kamu dengan memberikan rating dan komentar
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={reviewForm.handleSubmit(handleSubmit)} className="space-y-6 mt-4">
          {/* Rating Stars */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">Rating <span className="text-destructive">*</span></Label>
            <div className="flex items-center gap-2 p-4 bg-linear-to-r from-primary/5 to-accent/5 border border-primary/20 rounded-lg">
            {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  onClick={() => reviewForm.setValue("rating", star,{shouldValidate: true, shouldDirty: true, shouldTouch: true})}
                  className="star-button transition-transform duration-200 hover:scale-125 focus:outline-none"
                  data-star={star}
                  data-rated={rating}
                >
                  <Star
                    className={`w-10 h-10 transition-colors duration-200 ${
                      star <= rating
                        ? "fill-primary text-primary drop-shadow-[0_0_8px_rgba(var(--primary),0.5)]"
                        : "text-muted-foreground"
                    }`}
                  />
                </button>
              ))}
              {rating > 0 && (
                <span className="ml-3 text-lg font-bold text-primary">
                  {rating}/5
                </span>
              )}
            </div>
            {reviewForm.formState.errors.rating && (
              <p className="text-sm text-destructive">{reviewForm.formState.errors.rating.message}</p>
            )}
          </div>

          {/* Name Input */}
          <div className="space-y-2">
          <Controller
            name={"name"}
            control={reviewForm.control}
            render={({ field, fieldState: { error } }) => (
              <Input
                {...field}
                label="Nama"
                error={error?.message}
                required
                placeholder="Nama"
              />
            )}
          />
          </div>

          {/* Comment Textarea */}
          <div className="space-y-2">
            <Label htmlFor="comment" className="text-sm font-medium">
              Komentar
            </Label>
            <Controller
              name="comment"
              control={reviewForm.control}
              render={({ field, fieldState: { error } }) => (
                <Textarea
                  {...field}
                  error={error?.message}
                  required
                  placeholder="Tulis komentar Anda..."
                  rows={4}
                />
              )}
            />
            {reviewForm.formState.errors.comment && (
              <p className="text-sm text-destructive">{reviewForm.formState.errors.comment.message}</p>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              className="flex-1 border-border hover:bg-muted"
            >
              Batal
            </Button>
            <Button
              type="submit"
              className="flex-1"
            >
              Kirim Ulasan
            </Button>
          </div>
        </form>
      </DialogContent>
      <style jsx>{`
        /* CSS-only hover effect untuk menghindari re-render */
        .star-button:hover ~ .star-button svg {
          fill: transparent !important;
          color: hsl(var(--muted-foreground)) !important;
          filter: none !important;
        }
        
        .star-button:hover svg {
          fill: hsl(var(--primary)) !important;
          color: hsl(var(--primary)) !important;
          filter: drop-shadow(0 0 8px hsl(var(--primary) / 0.5)) !important;
        }
      `}</style>
    </Dialog>
  );
}