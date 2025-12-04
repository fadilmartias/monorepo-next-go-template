"use client";
import { useState } from "react";

export default function useSort() {
  const [sort, setSort] = useState("created_at:desc");

  return {
    sort,
    setSort,
  };
}