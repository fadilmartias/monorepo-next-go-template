import TwoFaForm from "./two-fa-form"

export const metadata = {
  title: "Autentikasi 2 Langkah",
}
type Params = Promise<{ token: string }>;

export default async function TwoFA(props: { params: Params }) {
  const { token } = await props.params;
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-md">
        <TwoFaForm token={token} />
      </div>
    </div>
  )
}
