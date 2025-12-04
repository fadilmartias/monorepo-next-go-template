export function formatRupiah(
  value: number | string,
  withPrefix = true,
  short = false
): string {
  const number = typeof value === "string" ? parseFloat(value) : value;
  if (isNaN(number)) return withPrefix ? "Rp 0" : "0";

  if (short) {
    const abs = Math.abs(number);
    let formatted: string;

    if (abs >= 1_000_000_000_000) {
      formatted = (number / 1_000_000_000_000).toFixed(1).replace(/\.0$/, "") + "T";
    } else if (abs >= 1_000_000_000) {
      formatted = (number / 1_000_000_000).toFixed(1).replace(/\.0$/, "") + "M";
    } else if (abs >= 1_000_000) {
      formatted = (number / 1_000_000).toFixed(1).replace(/\.0$/, "") + "jt";
    } else if (abs >= 1_000) {
      formatted = (number / 1_000).toFixed(1).replace(/\.0$/, "") + "rb";
    } else {
      formatted = number.toString();
    }

    return withPrefix ? `Rp ${formatted}` : formatted;
  }

  // default: full currency format
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  })
    .format(number)
    .replace(withPrefix ? "" : /^Rp\s?/, "");
}


  export function toDatetimeLocal(dateString: string) {
    const date = new Date(dateString);
    const pad = (n: number) => n.toString().padStart(2, "0");
  
    return (
      date.getFullYear() +
      "-" +
      pad(date.getMonth() + 1) +
      "-" +
      pad(date.getDate()) +
      "T" +
      pad(date.getHours()) +
      ":" +
      pad(date.getMinutes())
    );
  }

  export function formatTanggal(dateString: string) {
    const date = new Date(dateString);
    const pad = (n: number) => n.toString().padStart(2, "0");
  
    const formattedDate = pad(date.getDate()) +
      "/" +
      pad(date.getMonth() + 1) +
      "/" +
      date.getFullYear();

    const formattedTime = pad(date.getHours()) +
      ":" +
      pad(date.getMinutes());

    return `${formattedDate}, ${formattedTime}`;
  }
  
  