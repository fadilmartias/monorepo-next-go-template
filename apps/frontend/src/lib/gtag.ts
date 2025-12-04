// lib/gtag.ts
export const GA_TRACKING_ID = "G-CBN60MXFVF"

// Kirim pageview ke GA
export const pageview = (url: string) => {
  window.gtag("config", GA_TRACKING_ID, {
    page_path: url,
  })
}

// Kirim event custom (misal klik tombol, topup sukses, dll)
export const event = ({
  action,
  category,
  label,
  value,
}: {
  action: string
  category?: string
  label?: string
  value?: number
}) => {
  window.gtag("event", action, {
    event_category: category,
    event_label: label,
    value: value,
  })
}
