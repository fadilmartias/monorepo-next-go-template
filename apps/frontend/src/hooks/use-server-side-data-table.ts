"use client";
import { useState } from "react";
import usePagination from "./use-pagination";
import useSort from "./use-sort";
import useFilter from "./use-filter";

export default function useServerSideDataTable() {
  const {
    pagination,
    setPagination,
    nextPage,
    prevPage,
    goToPage,
    goToLastPage,
    goToFirstPage,
    resetPagination,
    setPageSize
  } = usePagination(1, 10);
  const { sort, setSort } = useSort();
  const { filter, setFilter } = useFilter();
  return {
    pagination,
    setPagination,
    nextPage,
    prevPage,
    goToPage,
    goToLastPage,
    goToFirstPage,
    resetPagination,
    sort,
    setSort,
    filter,
    setFilter,
    setPageSize
  };
}
