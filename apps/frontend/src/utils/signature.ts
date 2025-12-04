import CryptoJS from "crypto-js";

type CreateSignatureOptions = {
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  url: string; // full or relative URL path (e.g. /api/topup)

  body?: any; // JSON object (will be stringified)
  timestamp?: string; // Optional ISO string, default = now
  clientSecret: string;
  isMultipart?: boolean;
};

/**
 * Generates signature headers for secured API calls.
 */
export function createSignature({
  method,
  url,
  body,
  timestamp = new Date().toISOString(),
  clientSecret,
  isMultipart = false,
}: CreateSignatureOptions): Record<string, string> {
  let bodyHash = "";
  let raw = "";

  if (body && Object.keys(body).length > 0) {
    if (method.toLowerCase() === "get") {
      const flat = flattenQuery(body);

      const keys = Object.keys(flat).sort();

      const parts = keys.map((k) => `${k}=${flat[k]}`);

      raw = parts.join("&"); // Sama seperti Go
    } else if (isMultipart && body) {
      // Body harus object key → value biasa (file diabaikan)
      const fieldsOnly: Record<string, any> = {};

      Object.keys(body).forEach((key) => {
        const val = body[key];

        // Abaikan File / Blob
        if (val instanceof File || val instanceof Blob) return;

        fieldsOnly[key] = val;
      });

      const keys = Object.keys(fieldsOnly).sort();
      const parts = keys.map((k) => `${k}=${fieldsOnly[k]}`);

      raw = parts.join("&"); // SAMA PERSIS seperti Go
    } else if (body && Object.keys(body).length > 0) {
      // Untuk POST/PUT/PATCH: JSON stringify
      raw = JSON.stringify(body);
    }

    if (raw !== "") {
      const hash = CryptoJS.SHA256(raw).toString(CryptoJS.enc.Hex);
      bodyHash = hash.toLowerCase();
    }
  }

  // Construct stringToSign (gunakan url yang sama dengan endpoint di backend)
  const stringToSign = `${method.toUpperCase()}:${url}:${bodyHash}:${timestamp}`;
  console.log("stringToSign", stringToSign);

  // Create HMAC-SHA512 signature
  const hmac = CryptoJS.HmacSHA512(stringToSign, clientSecret);
  const signature = CryptoJS.enc.Base64.stringify(hmac);

  return {
    "X-Tenant-Id": process.env.NEXT_PUBLIC_MERCHANT_ID || "",
    "X-Signature": signature,
    "X-Timestamp": timestamp,
  };
}

function flattenQuery(obj: any, prefix = ""): Record<string, string> {
  let res: Record<string, string> = {};

  Object.keys(obj).forEach((key) => {
    const value = obj[key];
    const pref = prefix ? `${prefix}[${key}]` : key;

    if (typeof value === "object" && value !== null) {
      Object.assign(res, flattenQuery(value, pref));
    } else {
      res[pref] = String(value);
    }
  });

  return res;
}
