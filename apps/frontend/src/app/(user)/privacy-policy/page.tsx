import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getSetting } from "./actions";
import Typography from "@/components/typography";

export const revalidate = 86400;

export const metadata = {
  title: "Kebijakan Privasi",
  description: "Kebijakan Privasi Next Template",
  alternates: {
    canonical: process.env.NEXT_PUBLIC_APP_BASE_URL + "/privacy-policy",
  },
};

export default async function PrivacyPolicyPage() {
  const policy = await getSetting("privacy-policy");

  return (
    <div className="container mx-auto max-w-4xl py-10">
      <Card>
        <CardHeader>
          <CardTitle className="text-center">
            <Typography variant="h3" as="h1" className="font-bold">
              Kebijakan Privasi
            </Typography>
            <Typography variant="b1" className="text-center">
              Terakhir diperbarui pada {new Date(policy.updated_at).toLocaleDateString("id-ID")}
            </Typography>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div
            className="prose dark:prose-invert max-w-none"
            dangerouslySetInnerHTML={{ __html: policy.value }}
          />
        </CardContent>
      </Card>
    </div>
  );
}
