import Script from "next/script";

export default function YellowAIChatbot() {
  if (process.env.NEXT_PUBLIC_APP_ENV !== "production") {
    return (
      <Script id="yellow-ai-config" strategy="beforeInteractive">
        {`
          window.ymConfig = {
            bot: "x1758600598807",
            host: "https://r2.cloud.yellow.ai"
          };
          (function() {
            var w = window,
              ic = w.YellowMessenger;
            if (typeof ic === "function") {
              ic("reattach_activator");
              ic("update", ymConfig);
            } else {
              var d = document,
                i = function() { i.c(arguments); };
              function l() {
                var e = d.createElement("script");
                e.type = "text/javascript";
                e.async = true;
                e.src = "https://cdn.yellowmessenger.com/plugin/widget-v2/latest/dist/main.min.js";
                var t = d.getElementsByTagName("script")[0];
                t.parentNode.insertBefore(e, t);
              }
              i.q = [];
              i.c = function(e) { i.q.push(e); };
              w.YellowMessenger = i;
              if (w.attachEvent) w.attachEvent("onload", l);
              else w.addEventListener("load", l, false);
            }
          })();
        `}
      </Script>
    );
  } else {
    return (
      <Script id="yellow-ai-widget" strategy="beforeInteractive">
        {`
          window.ymConfig = {
            bot: "x1758600601144",
            host: "https://r2.cloud.yellow.ai"
          };
          (function() {
            var w = window,
                ic = w.YellowMessenger;
            if (typeof ic === "function") {
              ic("reattach_activator");
              ic("update", ymConfig);
            } else {
              var d = document,
                  i = function() { i.c(arguments); };
              function l() {
                var e = d.createElement("script");
                e.type = "text/javascript";
                e.async = true;
                e.src = "https://cdn.yellowmessenger.com/plugin/widget-v2/latest/dist/main.min.js";
                var t = d.getElementsByTagName("script")[0];
                t.parentNode.insertBefore(e, t);
              }
              i.q = [];
              i.c = function(e) { i.q.push(e); };
              w.YellowMessenger = i;
              if (w.attachEvent) w.attachEvent("onload", l);
              else w.addEventListener("load", l, false);
            }
          })();
        `}
      </Script>
    );
  }
}
