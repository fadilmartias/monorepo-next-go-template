import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getSetting } from "../privacy-policy/actions";
import Typography from "@/components/typography";

export const revalidate = 86400;

export const metadata = {
  title: "Syarat dan Ketentuan",
  description: "Syarat dan Ketentuan Next Template",
  alternates: {
    canonical: process.env.NEXT_PUBLIC_APP_BASE_URL + "/terms-and-conditions",
  },
};

export default async function TermsAndConditionsPage() {
  const tnc = await getSetting("terms-and-conditions");

  return (
    <div className="container mx-auto max-w-4xl py-10">
      <Card>
        <CardHeader>
          <CardTitle className="text-center">
            <Typography variant="h3" as="h1" className="font-bold">
              Syarat dan Ketentuan
            </Typography>
            <Typography variant="b1" className="text-center">
              Terakhir diperbarui pada {new Date(tnc.updated_at).toLocaleDateString("id-ID")}
            </Typography>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div
            className="prose dark:prose-invert max-w-none"
            dangerouslySetInnerHTML={{ __html: tnc.value }}
          />
        </CardContent>
      </Card>
    </div>
  );
}
