"use client";

import * as React from "react";
import { GoogleReCaptchaProvider } from "react-google-recaptcha-v3";

export function RecaptchaProvider({
  children,
  reCaptchaKey,
}: React.ComponentProps<typeof GoogleReCaptchaProvider>) {
  return (
    <GoogleReCaptchaProvider
      reCaptchaKey={reCaptchaKey}
    >
      {children}
    </GoogleReCaptchaProvider>
  );
}
