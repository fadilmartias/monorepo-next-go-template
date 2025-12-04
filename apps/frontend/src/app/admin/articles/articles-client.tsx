"use client";
import React, { useEffect } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { getArticleColumns } from "./columns";
import { apiClient } from "@/lib/api-client";
import { useQuery } from "@tanstack/react-query";
import useServerSideDataTable from "@/hooks/use-server-side-data-table";
import { useBreadcrumbStore } from "@/stores/breadcrumb";

const ArticlesClient = () => {
  const breadcrumbs = [
    {
      title: "Artikel",
      url: "/admin/articles",
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

  const getArticles = async () => {
    try {
      const res = await apiClient.get("/v0/articles", {
        params: {
          page: pagination.page,
          limit: pagination.page_size,
          orders: sort,
          joins: "category",
          filters: filter,
        },
      });
      res.data.data.map((item: any) => {
        item.category_name = item.category.name;
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
    queryKey: ["articles", pagination.page, pagination.page_size, sort, filter],
    queryFn: () => getArticles(),
    
  });

  const columns = getArticleColumns({
    serverSide: true,
    sort,
    setSort,
    pagination,
  });

  return (
    <>
      <Button className="">
        <Link href="/admin/articles/actions">Tambah Artikel</Link>
      </Button>
        <DataTable
          columnsSearch={["title"]}
          isLoading={isLoading}
          columns={columns}
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

export default ArticlesClient;
