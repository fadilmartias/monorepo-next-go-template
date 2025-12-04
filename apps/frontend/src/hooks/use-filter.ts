"use client";
import { useState, useMemo } from "react";
import { debounce } from "@/utils/debounce";

export default function useFilter(delay: number = 500) {
  const [filter, setFilterState] = useState<{ [key: string]: any }>({});

  // ✅ Memoized debounced setter
  const setFilter = useMemo(
    () =>
      debounce((newFilter: { [key: string]: any }) => {
        setFilterState(newFilter);
      }, delay),
    [delay]
  );

  return {
    filter,
    setFilter,
  };
}
