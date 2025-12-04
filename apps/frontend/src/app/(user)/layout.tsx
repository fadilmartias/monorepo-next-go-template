import Navbar from "@/components/partials/navbar";
import Footer from "@/components/partials/footer";
import { headers } from "next/headers";
import { getUserFromToken, isLoggedIn } from "@/lib/auth";
import GoogleOneTapLogin from "@/components/google-one-tap-login";
import YellowAIChatbot from "@/components/widget/chatbot/yellow-ai";

export default async function UserLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = await headers();
  const path = pathname.get("X-Pathname") || "";
  const isAuth = await isLoggedIn();
  const user = await getUserFromToken()
  const isLocal = process.env.NEXT_PUBLIC_APP_ENV === "local";
  

  const hideLayout = path.startsWith("/payment");
  return (
    <>
      {!hideLayout && <Navbar isLoggedIn={isAuth} user={user} />}
      {!hideLayout && !isAuth && !isLocal && <GoogleOneTapLogin />}
      <main
        className={`relative min-h-screen max-w-7xl mx-auto ${
          !hideLayout ? "mb-12" : "my-8"
        }`}
      >
        {children}
      </main>
      <YellowAIChatbot/>
      {!hideLayout && <Footer />}
    </>
  );
}
