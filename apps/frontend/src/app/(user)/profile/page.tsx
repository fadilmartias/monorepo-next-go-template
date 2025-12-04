import ProfileClient from "./profile-client";
import { Card } from "@/components/ui/card";

export default async function ProfilePage() {
  return (
    <section className="container-custom mx-auto max-w-3xl py-10">
      <Card className="bg-card border border-border shadow-xl rounded-2xl">
        <ProfileClient/>
      </Card>
    </section>
  );
}
