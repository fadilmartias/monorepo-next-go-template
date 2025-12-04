export const getCroppedImg = async (
  imageSrc: string,
  pixelCrop: any,
  rotation: number = 0,
  flipX: boolean = false,
  flipY: boolean = false
): Promise<Blob> => {
  const image = new Image();
  image.src = imageSrc;
  await new Promise((resolve) => (image.onload = resolve));

  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");

  if (!ctx) {
    throw new Error("Failed to get canvas context");
  }

  const rotRad = (rotation * Math.PI) / 180;

  // Set canvas size to match cropped area
  canvas.width = pixelCrop.width;
  canvas.height = pixelCrop.height;

  // Apply transformations
  ctx.translate(pixelCrop.width / 2, pixelCrop.height / 2);
  ctx.rotate(rotRad);
  ctx.scale(flipX ? -1 : 1, flipY ? -1 : 1);

  // Draw the cropped image
  ctx.drawImage(
    image,
    pixelCrop.x,
    pixelCrop.y,
    pixelCrop.width,
    pixelCrop.height,
    -pixelCrop.width / 2,
    -pixelCrop.height / 2,
    pixelCrop.width,
    pixelCrop.height
  );

  return new Promise((resolve) => {
    canvas.toBlob((blob) => {
      if (blob) resolve(blob);
    }, "image/webp");
  });
};