"use client"

import { ElementType, JSX, ComponentPropsWithoutRef, ReactNode } from "react";

const fontSizeConfig = {
  h1: "font-bold leading-tight text-[32px] md:text-[36px] lg:text-[48px] md:leading-[44px] lg:leading-[56px]",

  h2: "font-bold leading-tight text-[28px] md:text-[32px] lg:text-[40px] md:leading-[40px] lg:leading-[50px]",

  h3: "font-bold leading-snug text-[24px] md:text-[28px] lg:text-[36px] md:leading-[36px] lg:leading-[46px]",

  h4: "font-bold leading-normal text-[22px] md:text-[24px] lg:text-[32px] md:leading-[32px] lg:leading-[40px]",

  h5: "font-bold leading-normal text-[20px] md:text-[22px] lg:text-[28px] md:leading-[28px] lg:leading-[36px]",

  h6: "font-bold leading-normal text-[18px] md:text-[20px] lg:text-[24px] md:leading-[26px] lg:leading-[32px]",

  b1: "leading-relaxed text-[16px] md:text-[17px] lg:text-[20px] md:leading-[24px] lg:leading-[28px]",

  b2: "leading-relaxed text-[14px] md:text-[15px] lg:text-[18px] md:leading-[22px] lg:leading-[26px]",

  b3: "leading-normal text-[12px] md:text-[13px] lg:text-[15px] md:leading-[20px] lg:leading-[22px]",

  c1: "leading-normal text-[11px] md:text-[12px] lg:text-[13px] md:leading-[16px] lg:leading-[18px]",

  c2: "leading-snug text-[9px] md:text-[10px] lg:text-[11px] md:leading-[14px] lg:leading-[16px]",

  c3: "leading-tight text-[8px] md:text-[9px] lg:text-[10px] md:leading-[12px] lg:leading-[14px]",
};


const tags = {
  h1: "h1",
  h2: "h2",
  h3: "h3",
  h4: "h4",
  h5: "h5",
  h6: "h6",
  b1: "p",
  b2: "p",
  b3: "p",
  c1: "span",
  c2: "span",
  c3: "span",
};

type TypographyProps<T extends ElementType = 'p'> = {
    variant: keyof typeof fontSizeConfig;
    as?: T;
    className?: string;
    children?: ReactNode;
    dangerouslySetInnerHTML?: { __html: string };
  } & ComponentPropsWithoutRef<T>; // mewarisi semua props dari elemen HTML

  const Typography = <T extends ElementType = 'p'>({
    variant,
    as,
    children,
    className,
    dangerouslySetInnerHTML,
    ...props
  }: TypographyProps<T>): JSX.Element => {
    const Tag = as || tags[variant] || 'p';
    const fontSizeClasses = fontSizeConfig[variant] || 'body-3';
  
    return (
      <Tag className={`${fontSizeClasses} ${className || ''}`.trim()} {...props}>
        {dangerouslySetInnerHTML ? (
          <span dangerouslySetInnerHTML={dangerouslySetInnerHTML} />
        ) : (
          children
        )}
      </Tag>
    );
  };

export default Typography;
