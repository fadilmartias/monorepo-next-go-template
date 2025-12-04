import { Button } from "@/components/ui/button";
import { FaGoogle, FaDiscord, FaSteam, FaTwitch, FaFacebook } from "react-icons/fa";

type OAuthConfig = {
  google?: boolean;
  discord?: boolean;
  steam?: boolean;
  twitch?: boolean;
  facebook?: boolean;
};

interface OAuthButtonsProps {
  config: OAuthConfig;
  isPending?: boolean;
  handleOAuth: (provider: keyof OAuthConfig) => void;
}

export function OAuthButtons({ config, isPending, handleOAuth }: OAuthButtonsProps) {
  const providers = [
    {
      key: "google" as const,
      label: "Lanjutkan dengan Google",
      icon: <FaGoogle className="w-4 h-4" />,
    },
    {
      key: "discord" as const,
      label: "Lanjutkan dengan Discord",
      icon: <FaDiscord className="w-4 h-4" />,
    },
    {
      key: "steam" as const,
      label: "Lanjutkan dengan Steam",
      icon: <FaSteam className="w-4 h-4" />,
    },
    {
      key: "twitch" as const,
      label: "Lanjutkan dengan Twitch",
      icon: <FaTwitch className="w-4 h-4" />,
    },
    {
      key: "facebook" as const,
      label: "Lanjutkan dengan Facebook",
      icon: <FaFacebook className="w-4 h-4" />,
    },
  ];

  return (
    <div className="flex flex-col gap-4">
      {providers.map(
        (p) =>
          config[p.key] && (
            <Button
              key={p.key}
              type="button"
              variant="outline"
              className="w-full"
              disabled={isPending}
              onClick={() => handleOAuth(p.key)}
            >
              {p.icon}
              {p.label}
            </Button>
          )
      )}
    </div>
  );
}
