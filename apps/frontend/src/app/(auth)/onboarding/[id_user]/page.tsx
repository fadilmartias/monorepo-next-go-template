import { OnboardingForm } from "./onboarding-form"

type Params = Promise<{ id_user: string }>;

export const metadata = {
  title: "Onboarding",
}
export default async function Onboarding({ params }: { params: Params }) {
  const { id_user } = await params;
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-md">
        <OnboardingForm id_user={id_user} />
      </div>
    </div>
  )
}
