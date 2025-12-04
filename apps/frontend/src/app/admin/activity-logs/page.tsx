import ActivityLogsClient from "./activity-logs-client";

export const metadata = {
  title: "Activity Logs",
}

export default async function ActivityLogsPage() {
  return (
    <ActivityLogsClient />
  );
}
