import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import EditProfileForm from "@/components/forms/edit-profile";
import EditPasswordForm from "@/components/forms/edit-password";

export default async function EditProfilePage() {
  return (
    <section className="container-custom mx-auto py-6 flex flex-col gap-6">
      <Button variant="ghost" className="w-fit" asChild>
        <Link href="/profile">
          <ArrowLeft className="w-5 h-5" />
          Kembali
        </Link>
      </Button>
      <Card className="bg-card border border-border shadow-xl rounded-2xl">
        <CardHeader>
          <h2 className="text-2xl font-bold text-card-foreground">
            Edit Profil
          </h2>
        </CardHeader>

        <CardContent>
          <Tabs defaultValue="profile" className="w-full">
            <TabsList className="grid grid-cols-2 bg-muted">
              <TabsTrigger value="profile">Data Pribadi</TabsTrigger>
              <TabsTrigger value="password">Ganti Password</TabsTrigger>
            </TabsList>

            {/* Data Pribadi */}
            <TabsContent value="profile" className="space-y-5 mt-6">
              <EditProfileForm/>
            </TabsContent>

            {/* Ganti Password */}
            <TabsContent value="password" className="space-y-5 mt-6">
              <EditPasswordForm/>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </section>
  );
}
