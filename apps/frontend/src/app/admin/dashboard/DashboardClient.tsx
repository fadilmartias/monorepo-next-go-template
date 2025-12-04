"use client";

import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiClient } from "@/lib/api-client"; // axios instance
import { useState, useEffect } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { formatRupiah } from "@/utils/format";
import { useBreadcrumbStore } from "@/stores/breadcrumb";
import { Loader2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

export default function DashboardClient() {
  const breadcrumbs = [{
    title: "Dashboard",
    url: "/admin/dashboard"
  }]
  useEffect(() => {
    useBreadcrumbStore.setState({ breadcrumbs })
  }, [])
  
  const [saldoDigi, setSaldoDigi] = useState(0);

  // Fetch saldo Digiflazz
  const { data: digiflazzBalance, isLoading: loadingDigi, refetch: refetchDigi } = useQuery({
    queryKey: ["digiflazz-balance"],
    queryFn: async () => {
      const res = await apiClient.get("/v1/digiflazz/balance");
      setSaldoDigi(res.data.data.balance);
      return res.data;
    },
    enabled: false,
  });

  // Fetch Dashboard
  const { data: dataDashboard, isLoading: loadingDashboard } = useQuery({
    queryKey: ["dashboard"],
    queryFn: async () => {
      const res = await apiClient.get("/v1/dashboard");
      setSaldoDigi(res.data.data.balance_digiflazz);
      return res.data;
    },
  });

  return (
    <div className="space-y-6">
      {/* Cards Saldo */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Saldo Digiflazz</CardTitle>
            <Button variant="outline" size="sm" onClick={() => refetchDigi()}>
              Refresh
            </Button>
          </CardHeader>
          <CardContent>
            {loadingDashboard || loadingDigi ? (
                <Skeleton className="w-40 h-8" />
            ) : (
              <p className="text-2xl font-semibold">
                {formatRupiah(saldoDigi)}
              </p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Total Penjualan</CardTitle>
          </CardHeader>
          <CardContent>
            {loadingDashboard ? (
             <Skeleton className="w-40 h-8" />
            ) : (
              <p className="text-2xl font-semibold">{formatRupiah(dataDashboard?.data?.total_penjualan)}</p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Total Transaksi Berhasil</CardTitle>
          </CardHeader>
          <CardContent>
            {loadingDashboard ? (
             <Skeleton className="w-40 h-8" />
            ) : (
              <p className="text-2xl font-semibold">{dataDashboard?.data?.total_transaksi || 0}</p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Total Keuntungan</CardTitle>
          </CardHeader>
          <CardContent>
            {loadingDashboard ? (
             <Skeleton className="w-40 h-8" />
            ) : (
              <p className="text-2xl font-semibold">{formatRupiah(dataDashboard?.data?.total_profit)}</p>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {/* Tabel Transaksi Terakhir */}
        <Card className="gap-0">
          <CardHeader>
            <CardTitle>5 Transaksi Terakhir</CardTitle>
          </CardHeader>
          <CardContent>
            <DataTable
            isPagination={false}
            columns={[
                {
                  header: "Order ID",
                  accessorKey: "id",
                },
                {
                  header: "Tanggal Transaksi",
                  accessorKey: "created_at",
                  cell: (info: any) => (
                    <span className="font-semibold">
                      {new Date(info.row.original.created_at).toLocaleString("id-ID")}
                    </span>
                  ),
                },
                {
                  header: "Product",
                  accessorKey: "product.name",
                },
                {
                  header: "Variant",
                  accessorKey: "variant",
                  cell: (info: any) => (
                    <span className="font-semibold">
                      {
                        info.row.original.product_order_details[0]
                          .product_variant.name
                      }
                    </span>
                  ),
                },
                {
                  header: "Status",
                  accessorKey: "status",
                  cell: (info: any) => (
                    <span
                      className={`font-semibold capitalize inline-flex items-center justify-center px-2 py-1 rounded-md text-xs ${
                        info.row.original.status === "waiting payment" ||
                        info.row.original.status === "waiting delivery"
                          ? "bg-yellow-100 text-yellow-800"
                          : info.row.original.status === "succeeded"
                          ? "bg-success/10 text-success"
                          : "bg-destructive/10 text-destructive"
                      }`}
                    >
                      {info.row.original.status}
                    </span>
                  ),
                },
                {
                  header: "Total",
                  accessorKey: "final_price",
                  cell: (info: any) => (
                    <span className="font-semibold">
                      {formatRupiah(info.row.original.final_price)}
                    </span>
                  ),
                },
              ]}
              data={loadingDashboard ? null : dataDashboard?.data?.last_transactions || []}
            />
          </CardContent>
        </Card>

        {/* Top Produk Paling Laris */}
        <Card className="gap-0">
          <CardHeader>
            <CardTitle>Top 5 Produk Paling Laris</CardTitle>
          </CardHeader>
          <CardContent>
            <DataTable
            isPagination={false}
              columns={[
                {
                  header: "Product",
                  accessorKey: "product_name",
                },
                {
                  header: "Variant",
                  accessorKey: "variant_name",
                },
                {
                  header: "Total Terjual",
                  accessorKey: "total_sales",
                },
              ]}
              data={loadingDashboard ? null : dataDashboard?.data?.top_products || []}
            />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
