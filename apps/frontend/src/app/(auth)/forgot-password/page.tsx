import { ForgotPasswordForm } from "./forgot-password-form"
import { Suspense } from "react"
import { Loader2 } from "lucide-react"
import { connection } from 'next/server'
export const metadata = {
  title: "Lupa Password",
}
export default async function ForgotPassword() {
  await connection()
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
      <Suspense fallback={<Loader2 className="h-4 w-4 animate-spin" />}>
        <ForgotPasswordForm />
      </Suspense>
      </div>
    </div>
  )
}
