"use client";

import { useBreadcrumbStore } from "@/stores/breadcrumb";
import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import EmptyState from "@/components/states/empty-state";
import LoadingState from "@/components/states/loading-state";
import { ActivityLog } from "../../columns";
import { format } from "date-fns";
import { id } from "date-fns/locale";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Activity,
  User,
  Calendar,
  Globe,
  Fingerprint,
  FileJson,
  Clock,
  Hash,
  Monitor,
  MapPin,
  ChevronDown,
  ChevronUp,
  Copy,
  Check,
  ArrowLeft,
  Info,
  AlertCircle,
} from "lucide-react";
import { formatTanggal } from "@/utils/format";
import Typography from "@/components/typography";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

export default function ActivityLogsDetailsClient({ id }: { id: string }) {
  const { setBreadcrumbs } = useBreadcrumbStore();
  const router = useRouter();
  const [copiedField, setCopiedField] = useState<string | null>(null);
  const [isPropertiesExpanded, setIsPropertiesExpanded] = useState(true);

  const getDetails = async () => {
    const res = await apiClient.get(`/v0/activity-logs/${id}`);
    return res.data.data as ActivityLog;
  };

  const {
    data: log,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["activity-logs", id],
    queryFn: getDetails,
  });

  useEffect(() => {
    setBreadcrumbs([
      { title: "Activity Logs", url: "/admin/activity-logs" },
      {
        title: "Detail Log",
        url: `/admin/activity-logs/details/${id}`,
        subtitle: "Informasi lengkap aktivitas sistem",
        backButton: true,
      },
    ]);
  }, [id, setBreadcrumbs]);

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text);
    setCopiedField(field);
    setTimeout(() => setCopiedField(null), 2000);
  };

  if (isLoading) return <LoadingState />;
  if (error || !log)
    return (
      <EmptyState
        title="Error"
        description="Gagal memuat detail activity log"
      />
    );

  // Parse properties jika berupa JSON string
  let propertiesObj: any = {};
  let hasOldValues = false;
  let hasNewValues = false;
  
  try {
    propertiesObj = JSON.parse(log.properties || "{}");
    hasOldValues = propertiesObj.old && Object.keys(propertiesObj.old).length > 0;
    hasNewValues = propertiesObj.attributes && Object.keys(propertiesObj.attributes).length > 0;
  } catch {
    propertiesObj = { raw: log.properties };
    hasOldValues = propertiesObj.raw.old && Object.keys(propertiesObj.raw.old).length > 0;
    hasNewValues = propertiesObj.raw.new && Object.keys(propertiesObj.raw.new).length > 0;
  }
  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const getEventConfig = (event: string) => {
    const configs: Record<string, { variant: any; icon: any; color: string; label: string }> = {
      created: {
        variant: "default",
        icon: "✨",
        color: "text-green-600 dark:text-green-400",
        label: "Dibuat"
      },
      updated: {
        variant: "secondary",
        icon: "✏️",
        color: "text-blue-600 dark:text-blue-400",
        label: "Diperbarui"
      },
      deleted: {
        variant: "destructive",
        icon: "🗑️",
        color: "text-red-600 dark:text-red-400",
        label: "Dihapus"
      },
      restored: {
        variant: "outline",
        icon: "♻️",
        color: "text-purple-600 dark:text-purple-400",
        label: "Dipulihkan"
      },
    };
    return configs[event.toLowerCase()] || {
      variant: "secondary",
      icon: "📝",
      color: "text-gray-600 dark:text-gray-400",
      label: event.charAt(0).toUpperCase() + event.slice(1)
    };
  };

  const eventConfig = getEventConfig(log.event);

  const InfoCard = ({ icon: Icon, label, value, copyable = false, mono = false }: any) => (
    <div className="group relative rounded-lg border bg-card p-4 hover: transition-colors">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 space-y-1">
          <div className="flex items-center gap-2">
            <Icon className="w-4 h-4 text-muted-foreground" />
            <Typography variant="c1" className="text-muted-foreground font-medium">
              {label}
            </Typography>
          </div>
          <Typography 
            variant="b3" 
            className={cn(
              "break-all",
              mono && "font-mono text-sm"
            )}
          >
            {value || "-"}
          </Typography>
        </div>
        {copyable && value && (
          <Button
            variant="ghost"
            size="sm"
            className="opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8 p-0"
            onClick={() => copyToClipboard(value, label)}
          >
            {copiedField === label ? (
              <Check className="w-4 h-4 text-green-600" />
            ) : (
              <Copy className="w-4 h-4" />
            )}
          </Button>
        )}
      </div>
    </div>
  );

  const renderPropertyChanges = () => {
    if (!hasOldValues && !hasNewValues) {
      return (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <Info className="w-12 h-12 text-muted-foreground/50 mb-3" />
          <Typography variant="b3" className="text-muted-foreground">
            Tidak ada perubahan data yang tercatat
          </Typography>
        </div>
      );
    }

    return (
      <div className="space-y-4">
        {hasOldValues && hasNewValues && (
          <div className="grid md:grid-cols-2 gap-4">
            {/* Old Values */}
            <div className="space-y-2">
              <div className="flex items-center gap-2 px-3 py-2 bg-red-50 dark:bg-red-950/20 rounded-md">
                <AlertCircle className="w-4 h-4 text-red-600 dark:text-red-400" />
                <Typography variant="b3" className="font-semibold text-red-900 dark:text-red-100">
                  Nilai Sebelumnya
                </Typography>
              </div>
              <div className="bg-muted/50 rounded-lg p-4 max-h-64 overflow-auto">
                <pre className="text-xs font-mono">
                  {JSON.stringify(propertiesObj.old || propertiesObj.raw.old, null, 2)}
                </pre>
              </div>
            </div>

            {/* New Values */}
            <div className="space-y-2">
              <div className="flex items-center gap-2 px-3 py-2 bg-green-50 dark:bg-green-950/20 rounded-md">
                <Check className="w-4 h-4 text-green-600 dark:text-green-400" />
                <Typography variant="b3" className="font-semibold text-green-900 dark:text-green-100">
                  Nilai Baru
                </Typography>
              </div>
              <div className="bg-muted/50 rounded-lg p-4 max-h-64 overflow-auto">
                <pre className="text-xs font-mono">
                  {JSON.stringify(propertiesObj.new || propertiesObj.raw.new, null, 2)}
                </pre>
              </div>
            </div>
          </div>
        )}

        {/* Full Properties View */}
        {Object.keys(propertiesObj).length > 0 && (
          <div className="mt-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsPropertiesExpanded(!isPropertiesExpanded)}
              className="w-full"
            >
              {isPropertiesExpanded ? (
                <>
                  <ChevronUp className="w-4 h-4 mr-2" />
                  Sembunyikan Detail Lengkap
                </>
              ) : (
                <>
                  <ChevronDown className="w-4 h-4 mr-2" />
                  Tampilkan Detail Lengkap
                </>
              )}
            </Button>
            
            {isPropertiesExpanded && (
              <div className="mt-3 bg-muted/50 rounded-lg p-4 max-h-96 overflow-auto">
                <pre className="text-xs font-mono">
                  {JSON.stringify(propertiesObj, null, 2)}
                </pre>
              </div>
            )}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-accent/10 border">
                <span className="text-xl">{eventConfig.icon}</span>
                  <span className="text-primary font-semibold">{eventConfig.label}</span>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>Tipe aktivitas: {log.event}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Overview Card */}
          <Card className="border-2">
            <CardHeader className="">
              <CardTitle className="flex items-center gap-2">
                <Activity className="w-5 h-5" />
                Overview Aktivitas
              </CardTitle>
              <CardDescription>
                Informasi ringkas tentang aktivitas yang terjadi
              </CardDescription>
            </CardHeader>
            <CardContent className="pt-6">
              <div className="grid gap-3">
                <InfoCard
                  icon={Hash}
                  label="Log ID"
                  value={log.id}
                  copyable
                  mono
                />
                <InfoCard
                  icon={FileJson}
                  label="Log Name"
                  value={log.log_name}
                />
                <InfoCard
                  icon={Info}
                  label="Deskripsi"
                  value={log.desc}
                />
              </div>
            </CardContent>
          </Card>

          {/* Subject Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Fingerprint className="w-5 h-5" />
                Subject Information
              </CardTitle>
              <CardDescription>
                Entitas yang terpengaruh oleh aktivitas ini
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-3">
                <InfoCard
                  icon={FileJson}
                  label="Subject Type"
                  value={log.subject_type}
                  mono
                />
                <InfoCard
                  icon={Hash}
                  label="Subject ID"
                  value={log.subject_id}
                  copyable
                  mono
                />
              </div>
            </CardContent>
          </Card>

          {/* Properties / Changes */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileJson className="w-5 h-5" />
                Perubahan Data
              </CardTitle>
              <CardDescription>
                Detail perubahan atribut dalam aktivitas ini
              </CardDescription>
            </CardHeader>
            <CardContent>
              {renderPropertyChanges()}
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Timeline Card */}
          <Card className="border-2 gap-0 border-primary/20">
            <CardHeader className="">
              <CardTitle className="text-lg flex items-center gap-2">
                <Clock className="w-5 h-5" />
                Timeline
              </CardTitle>
            </CardHeader>
            <CardContent className="">
              <div className="flex items-center gap-3 p-3 rounded-lg ">
                <Calendar className="w-5 h-5 text-primary" />
                <div className="flex-1">
                  <Typography variant="c1" className="text-muted-foreground text-xs">
                    Tanggal
                  </Typography>
                  <Typography variant="b3" className="font-semibold">
                    {new Date(log.created_at).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" })}
                  </Typography>
                </div>
              </div>
              
              <div className="flex items-center gap-3 p-3 rounded-lg ">
                <Clock className="w-5 h-5 text-primary" />
                <div className="flex-1">
                  <Typography variant="c1" className="text-muted-foreground text-xs">
                    Waktu
                  </Typography>
                  <Typography variant="b3" className="font-semibold font-mono">
                    {new Date(log.created_at).toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit", second: "2-digit" })} WIB
                  </Typography>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Causer Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <User className="w-5 h-5" />
                Pelaku
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-4 p-4 rounded-lg  border">
                <Avatar className="h-12 w-12 ring-2 ring-primary/20">
                  <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                    {log.causer_type ? getInitials(log.causer_type) : <User className="w-6 h-6" />}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <Typography variant="b3" className="font-semibold">
                    {log.causer_type || "System"}
                  </Typography>
                  <Typography variant="c1" className="text-muted-foreground font-mono">
                    ID: {log.causer_id || "Auto"}
                  </Typography>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Device & Network Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Monitor className="w-5 h-5" />
                Informasi Perangkat
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <MapPin className="w-4 h-4 text-muted-foreground" />
                  <Typography variant="b3" className="font-medium">
                    IP Address
                  </Typography>
                </div>
                <div className="flex items-center justify-between gap-2 p-3  rounded-lg border group">
                  <Typography as="code" variant="c1" className="font-mono text-sm">
                    {log.ip || "-"}
                  </Typography>
                  {log.ip && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                      onClick={() => copyToClipboard(log.ip, "IP")}
                    >
                      {copiedField === "IP" ? (
                        <Check className="w-3 h-3 text-green-600" />
                      ) : (
                        <Copy className="w-3 h-3" />
                      )}
                    </Button>
                  )}
                </div>
              </div>

              <Separator />

              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Globe className="w-4 h-4 text-muted-foreground" />
                  <Typography variant="b3" className="font-medium">
                    User Agent
                  </Typography>
                </div>
                <div className="p-3  rounded-lg border">
                  <Typography as="code" variant="c1" className="text-xs font-mono break-all">
                    {log.user_agent || "-"}
                  </Typography>
                </div>
              </div>

              {log.trace_id && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <Fingerprint className="w-4 h-4 text-muted-foreground" />
                      <Typography variant="b3" className="font-medium">
                        Trace ID
                      </Typography>
                    </div>
                    <div className="flex items-center justify-between gap-2 p-3  rounded-lg border group">
                      <Typography as="code" variant="c1" className="text-xs font-mono break-all">
                        {log.trace_id}
                      </Typography>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={() => copyToClipboard(log.trace_id!, "Trace ID")}
                      >
                        {copiedField === "Trace ID" ? (
                          <Check className="w-3 h-3 text-green-600" />
                        ) : (
                          <Copy className="w-3 h-3" />
                        )}
                      </Button>
                    </div>
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}