'use client';

import Link from 'next/link';
import Image from 'next/image';
import { useTheme } from 'next-themes';
import { useMounted } from '@/hooks/use-mounted';
import { cn } from '@/lib/utils';

export default function Logo({
  size = 'xl',
  href = '/',
  className = '',
}: {
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl';
  href?: string;
  className?: string;
}) {
  const appName = process.env.NEXT_PUBLIC_APP_NAME || 'Next Template';
  const { theme } = useTheme();
  const mounted = useMounted();

  const sizeStyles = {
    sm: { width: 48, height: 48 },
    md: { width: 72, height: 72 },
    lg: { width: 84, height: 84 },
    xl: { width: 96, height: 96 },
    '2xl': { width: 120, height: 120 },
    '3xl': { width: 144, height: 144 },
  };
  // pilih logo sesuai theme

  if (!mounted) return (
    <Image
        src="/assets/images/logos/logo-light.png"
        alt={`${appName}'s Logo`}
        width={sizeStyles[size].width}
        height={sizeStyles[size].height}
        className={cn('h-fit', className)}
        priority
      />
  );

  const logoSrc =
    theme === 'light'
      ? '/assets/images/logos/logo-dark.png'
      : '/assets/images/logos/logo-light.png';

  return (
    <Link href={href} className="inline-block">
      <Image
        src={logoSrc}
        alt={`${appName}'s Logo`}
        width={sizeStyles[size].width}
        height={sizeStyles[size].height}
        className={cn('h-fit', className)}
        priority
      />
    </Link>
  );
}
