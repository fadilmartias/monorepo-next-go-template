import Script from "next/script";

export default function BrevoChatbot() {
  return (
    <>
       <Script id="brevo-conversations" strategy="afterInteractive">
          {`
          (function(d, w, c) {
              w.BrevoConversationsID = '68c3793e76303d0cc601cbd3';
              w[c] = w[c] || function() {
                  (w[c].q = w[c].q || []).push(arguments);
              };
              var s = d.createElement('script');
              s.async = true;
              s.src = 'https://conversations-widget.brevo.com/brevo-conversations.js';
              if (d.head) d.head.appendChild(s);
          })(document, window, 'BrevoConversations');
        `}
        </Script>
    </>
  );
}
