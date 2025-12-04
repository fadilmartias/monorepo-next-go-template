"use client";

import { useEffect, useState, useRef } from "react";
import Cropper from "react-cropper";
import "cropperjs/dist/cropper.css";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog";
import { apiClient } from "@/lib/api-client";
import { getApiImage } from "@/lib/api-asset";
import { Label } from "./ui/label";
import Typography from "./typography";

export default function ImageCropper({
  pathName,
  onCropped,
  initialUrl,
  label,
  id,
  required,
  requiredSign = true,
  error,
  cropShape = "rect",
  aspect = 4 / 3,
}: {
  pathName: string;
  onCropped: (url: string) => void;
  initialUrl?: string;
  label?: string;
  id?: string;
  required?: boolean;
  requiredSign?: boolean;
  error?: string;
  cropShape?: "round" | "rect";
  aspect?: number;
}) {
  const [image, setImage] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState("");
  const [open, setOpen] = useState(false);
  const [currentAspect, setCurrentAspect] = useState<number | undefined>(aspect);
  const cropperRef = useRef<Cropper>(null);

  useEffect(() => {
    if (initialUrl) {
      setPreviewUrl(getApiImage(pathName, initialUrl));
    }
  }, [initialUrl, pathName]);

  const handleUpload = async () => {
    if (!image || !cropperRef.current) return;

    try {
      const croppedCanvas = cropperRef.current.cropper.getCroppedCanvas();
      const croppedBlob = await new Promise<Blob>((resolve) => {
        croppedCanvas.toBlob((blob: Blob) => {
          if (blob) resolve(blob);
        }, "image/webp");
      });

      const formData = new FormData();
      formData.append("image", croppedBlob);

      const res = await apiClient.put("/v0/uploads", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      setPreviewUrl(res.data.data.file_url);
      onCropped(res.data.data.file_name);
      resetTransformations();
      setOpen(false);
    } catch (error) {
      console.error("Upload failed:", error);
    }
  };

  const resetTransformations = () => {
    setImage(null);
    setCurrentAspect(aspect);
    if (cropperRef.current) {
      cropperRef.current.cropper.reset();
    }
  };

  const setAspectRatio = (aspectRatio: number | undefined) => {
    setCurrentAspect(aspectRatio);
    if (cropperRef.current) {
      cropperRef.current.cropper.setAspectRatio(aspectRatio || 0);
    }
  };

  const rotateLeft = () => {
    if (cropperRef.current) {
      cropperRef.current.cropper.rotate(-90);
    }
  };

  const rotateRight = () => {
    if (cropperRef.current) {
      cropperRef.current.cropper.rotate(90);
    }
  };

  const flipHorizontal = () => {
    if (cropperRef.current) {
      const data = cropperRef.current.cropper.getData();
      cropperRef.current.cropper.scaleX(-data.scaleX || -1);
    }
  };

  const flipVertical = () => {
    if (cropperRef.current) {
      const data = cropperRef.current.cropper.getData();
      cropperRef.current.cropper.scaleY(-data.scaleY || -1);
    }
  };

  return (
    <div className="space-y-2">
      {label && (
        <Label htmlFor={id || label || ""}>
          {label}
          {required && requiredSign && (
            <span className="text-destructive">*</span>
          )}
        </Label>
      )}
      {previewUrl && (
        <img src={previewUrl} alt="preview" className="rounded-md max-h-48" />
      )}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button variant="outline">Upload & Crop</Button>
        </DialogTrigger>
        <DialogContent showCloseButton={false} className="w-[90vw] max-w-3xl">
          <DialogHeader>
            <DialogTitle>Upload & Crop</DialogTitle>
            <DialogDescription>
              Upload and crop your image. All transformations (crop, rotate, flip) are applied in real-time.
            </DialogDescription>
          </DialogHeader>
          <Input
            type="file"
            accept="image/*"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) {
                setImage(file);
                setCurrentAspect(aspect);
              }
            }}
          />
          {image && (
            <>
              <div className="relative w-full h-80 bg-muted">
                <Cropper
                  src={URL.createObjectURL(image)}
                  style={{ height: "100%", width: "100%" }}
                  initialAspectRatio={aspect}
                  aspectRatio={currentAspect}
                  guides={true}
                  viewMode={1}
                  minCropBoxHeight={50}
                  minCropBoxWidth={50}
                  background={true}
                  responsive={true}
                  autoCropArea={0.8}
                  checkOrientation={false}
                  cropBoxMovable={true}
                  cropBoxResizable={true}
                  dragMode="move"
                  scalable={true}
                  zoomable={true}
                  rotatable={true}
                  ref={cropperRef as any}
                />
              </div>
              <div className="flex flex-wrap justify-center gap-2 mt-4">
                <Button
                  variant={currentAspect === 1 ? "default" : "secondary"}
                  onClick={() => setAspectRatio(1)}
                >
                  1:1
                </Button>
                <Button
                  variant={currentAspect === 4 / 3 ? "default" : "secondary"}
                  onClick={() => setAspectRatio(4 / 3)}
                >
                  4:3
                </Button>
                <Button
                  variant={currentAspect === 16 / 9 ? "default" : "secondary"}
                  onClick={() => setAspectRatio(16 / 9)}
                >
                  16:9
                </Button>
                <Button
                  variant={currentAspect === undefined ? "default" : "secondary"}
                  onClick={() => setAspectRatio(undefined)}
                >
                  Free
                </Button>
                <Button variant="secondary" onClick={rotateLeft}>
                  Rotate Left
                </Button>
                <Button variant="secondary" onClick={rotateRight}>
                  Rotate Right
                </Button>
                <Button variant="secondary" onClick={flipHorizontal}>
                  Mirror Horizontal
                </Button>
                <Button variant="secondary" onClick={flipVertical}>
                  Mirror Vertical
                </Button>
              </div>
            </>
          )}
          <DialogFooter>
            <DialogClose asChild onClick={resetTransformations}>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button
              onClick={handleUpload}
              disabled={!image || !cropperRef.current}
            >
              Crop & Upload
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      {error && (
        <Typography variant="c1" className="text-destructive block">
          {error}
        </Typography>
      )}
    </div>
  );
}