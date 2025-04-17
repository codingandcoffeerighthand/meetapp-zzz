import React from "react";
import Image from "next/image";
import { cn } from "@/lib/utils";

const getMaskUrl = variant => {
  const maskMap = {
    shape1: "/mask-shape-1.svg",
    shape2: "/mask-shape-2.svg",
    shape3: "/mask-shape-3.svg",
    shape4: "/mask-shape-4.svg",
    shape5: "/mask-shape-5.svg",
    shape6: "/mask-shape-6.svg",
  };

  return maskMap[variant];
};

export const MaskedImage = ({
  src,
  alt,
  width,
  height,
  variant = "shape1",
  className,
}) => {
  const maskUrl = getMaskUrl(variant);

  return (
    <Image
      className={cn("mask-no-repeat object-cover mask-size-[100%_100%] mask-center", className)}
      style={{ maskImage: `url(${maskUrl})` }}
      src={src}
      alt={alt}
      width={width}
      height={height} />
  );
};
