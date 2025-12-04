import ActivityLogsDetailsClient from "./activy-logs-details-client";

export const metadata = {
  title: "Activity Logs Details",
};

type Params = Promise<{ id: string }>;

export default async function ActivityLogsDetailsPage({
  params,
}: {
  params: Params;
}) {
  const { id } = await params;
  return <ActivityLogsDetailsClient id={id} />;
}
