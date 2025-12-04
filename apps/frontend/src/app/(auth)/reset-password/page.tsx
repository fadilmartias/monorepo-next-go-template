import { Suspense } from "react"
import { ResetPasswordForm } from "./reset-password-form"
import { Loader2 } from "lucide-react"
import { connection } from 'next/server'
export const metadata = {
  title: "Atur Ulang Kata Sandi",
}
export default async function ResetPassword() {
  await connection()
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-md">
        <Suspense fallback={<Loader2 className="h-4 w-4 animate-spin" />}>
          <ResetPasswordForm />
        </Suspense>
      </div>
    </div>
  )
}
