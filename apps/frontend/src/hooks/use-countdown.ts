"use client";
import { useEffect, useState } from "react";

type TimeLeft = {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  total: number;
};

function getDiff(expiredAt: string): TimeLeft {
  const total = Math.max(0, new Date(expiredAt).getTime() - Date.now());
  const days = Math.floor(total / (1000 * 60 * 60 * 24));
  const hours = Math.floor((total / (1000 * 60 * 60)) % 24);
  const minutes = Math.floor((total / 1000 / 60) % 60);
  const seconds = Math.floor((total / 1000) % 60);
  return { days, hours, minutes, seconds, total };
}

export function useCountdown(expiredAt: string): TimeLeft {
  const [timeLeft, setTimeLeft] = useState(() => getDiff(expiredAt));

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeLeft(getDiff(expiredAt));
    }, 1000);
    return () => clearInterval(interval);
  }, [expiredAt]);

  return timeLeft;
}
