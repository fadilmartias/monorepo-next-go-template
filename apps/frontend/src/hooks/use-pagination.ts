"use client"
import { APIPaginationResponse } from "@/types/api";
import { useState } from "react";

export default function usePagination(initialPage = 1, perPage = 10) {
  const [pagination, setPagination] = useState({
    page: initialPage,
    page_size: perPage,
    total_items: 0,
    total_pages: 0,
    has_more: false,
    from: 0,
    to: 0,
  } as APIPaginationResponse);

  const setPageSize = (pageSize: number) =>
    setPagination(prev => ({ ...prev, page_size: pageSize }));

  const nextPage = () =>
    setPagination(prev => ({ ...prev, page: prev.page + 1 }));

  const prevPage = () =>
    setPagination(prev => ({ ...prev, page: Math.max(1, prev.page - 1) }));

  const goToPage = (page: number) =>
    setPagination(prev => ({ ...prev, page: Math.max(1, page) }));

  const goToLastPage = () =>
    setPagination(prev => ({ ...prev, page: prev.total_pages }));

  const goToFirstPage = () =>
    setPagination(prev => ({ ...prev, page: 1 }));

  const resetPagination = () =>
    setPagination({
      page: initialPage,
      page_size: perPage,
      total_items: 0,
      total_pages: 0,
      has_more: false,
      from: 0,
      to: 0,
    });

  return {
    pagination,
    setPagination,
    nextPage,
    prevPage,
    goToPage,
    goToLastPage,
    goToFirstPage,
    resetPagination,
    setPageSize
  };
}
