"use client";
import { useEffect } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { apiClient } from "@/lib/api-client";
import { useQuery } from "@tanstack/react-query";
import useServerSideDataTable from "@/hooks/use-server-side-data-table";
import { useBreadcrumbStore } from "@/stores/breadcrumb";
import { getActivityLogsColumns } from "./columns";

const ActivityLogsClient = () => {
  const breadcrumbs = [
    {
      title: "Activity Logs",
      url: "/admin/activity-logs",
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

  const getActivityLogs = async () => {
    try {
      const res = await apiClient.get("/v0/activity-logs", {
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
      console.log("err",error);
      return error;
    }
  };
  const { data, isSuccess, isLoading } = useQuery({
    queryKey: ["activity-logs", pagination.page, pagination.page_size, sort, filter],
    queryFn: () => getActivityLogs(),
    
  });

  const columns = getActivityLogsColumns({
    serverSide: true,
    sort,
    setSort,
    pagination,
  });

  return (
    <>
        <DataTable
          columnsSearch={["log_name","desc"]}
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

export default ActivityLogsClient;
