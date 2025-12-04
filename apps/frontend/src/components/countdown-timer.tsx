"use client";
import { useCountdown } from "@/hooks/use-countdown";
import { useMounted } from "@/hooks/use-mounted";

export default function CountdownTimer({ expiredAt }: { expiredAt: string }) {
  const timeLeft = useCountdown(expiredAt);
  const isMounted = useMounted();

  if (!isMounted) {
    // ⛔ avoid hydration mismatch
    return <span className="font-mono text-2xl">00:00:00</span>;
  }

  if (timeLeft.total <= 0) {
    return (
      <span className="text-destructive text-center text-lg font-semibold">
        Waktu Habis
      </span>
    );
  }

  return (
    <span className="font-mono text-2xl">
      {String(timeLeft.hours).padStart(2, "0")}:
      {String(timeLeft.minutes).padStart(2, "0")}:
      {String(timeLeft.seconds).padStart(2, "0")}
    </span>
  );
}
