"use client";
import React, { useEffect } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { getBannerColumns } from "./columns";
import { apiClient } from "@/lib/api-client";
import { useQuery } from "@tanstack/react-query";
import useServerSideDataTable from "@/hooks/use-server-side-data-table";
import { useBreadcrumbStore } from "@/stores/breadcrumb";

const BannersClient = () => {
  const breadcrumbs = [
    {
      title: "Banner",
      url: "/admin/banners",
    },
  ];
  useEffect(() => {
    useBreadcrumbStore.setState({ breadcrumbs })
  }, [])
  const {
    pagination,
    setPagination,
    sort,
    setSort,
    filter,
    setFilter,
    setPageSize,
    nextPage,
    prevPage,
    goToPage,
    goToLastPage,
    goToFirstPage,
    resetPagination,
  } = useServerSideDataTable();

  const getBanners = async () => {
    try {
      const res = await apiClient.get("/v0/banners", {
        params: {
          page: pagination.page,
          limit: pagination.page_size,
          orders: sort,
          filters: filter,
        },
      });
      setPagination(prev => ({
        ...prev,
        ...res.data.pagination,
      }));
      return res.data;
    } catch (error: any) {
      if (error.response) {
        return error.response.data;
      }
      return error;
    }
  };
  const { data, isSuccess, isLoading } = useQuery({
    queryKey: ["banners", pagination.page, pagination.page_size, sort, filter],
    queryFn: () => getBanners(),
    
  });

  const columns = getBannerColumns({
    serverSide: true,
    sort,
    setSort,
    pagination,
  });

  return (
    <>
      <Button>
        <Link href="/admin/banners/actions">Tambah Banner</Link>
      </Button>
        <DataTable
          columnsSearch={["title"]}
          columns={columns}
          isLoading={isLoading}
          data={data?.data}
          pagination={pagination}
          serverSide={true}
          sort={sort}
          setSort={setSort}
          filter={filter}
          setFilter={setFilter}
          nextPage={nextPage}
          prevPage={prevPage}
          goToPage={goToPage}
          goToLastPage={goToLastPage}
          goToFirstPage={goToFirstPage}
          setPageSize={setPageSize}

        />
    </>
  );
};

export default BannersClient;
