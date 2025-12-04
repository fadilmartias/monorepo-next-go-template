import Select from "react-select";
import { useEffect, useState } from "react";
import Typography from "./typography";
import { Label } from "./ui/label";

export default function ReactSelect({
  options,
  onChange,
  value,
  label,
  id,
  required,
  requiredSign = true,
  error,
  isMulti = false,
}: {
  options: any;
  onChange: any;
  value: any;
  label?: string;
  id?: string;
  required?: boolean;
  requiredSign?: boolean;
  error?: string;
  isMulti?: boolean;
}) {
  const [themeVars, setThemeVars] = useState({
    background: "#fff",
    foreground: "#000",
    input: "#ccc",
    ring: "#3b82f6",
    accent: "#e0e7ff",
    accentForeground: "#1e40af",
    mutedForeground: "#6b7280",
    error: "#dc2626",
  });

  useEffect(() => {
    const root = getComputedStyle(document.documentElement);
    setThemeVars({
      background: root.getPropertyValue("--background")?.trim() || "#fff",
      foreground: root.getPropertyValue("--foreground")?.trim() || "#000",
      input: root.getPropertyValue("--input")?.trim() || "#ccc",
      ring: root.getPropertyValue("--ring")?.trim() || "#3b82f6",
      accent: root.getPropertyValue("--accent")?.trim() || "#e0e7ff",
      accentForeground:
        root.getPropertyValue("--accent-foreground")?.trim() || "#1e40af",
      mutedForeground:
        root.getPropertyValue("--muted-foreground")?.trim() || "#6b7280",
      error: root.getPropertyValue("--destructive")?.trim() || "#dc2626",
    });
  }, []);

  const customStyles = {
    control: (base: any, state: any) => ({
      ...base,
      backgroundColor: themeVars.input, // ← pakai warna input, bukan accent
      borderColor: state.isFocused ? themeVars.ring : error ? themeVars.error : themeVars.input,
      boxShadow: state.isFocused ? `0 0 0 2px ${themeVars.ring}` : "none",
      borderRadius: "0.5rem",
      padding: "2px",
      minHeight: "40px",
      fontSize: "0.875rem",
      "&:hover": {
        borderColor: state.isFocused ? themeVars.ring : themeVars.input,
      },
    }),
    menu: (base: any) => ({
      ...base,
      backgroundColor: themeVars.background,
      zIndex: 9999,
    }),
    menuPortal: (base: any) => ({
      ...base,
      backgroundColor: themeVars.background,
      zIndex: 9999,
    }),
    option: (base: any, state: any) => ({
      ...base,
      backgroundColor: state.isFocused
        ? themeVars.accent
        : themeVars.background,
      color: state.isSelected
        ? themeVars.accentForeground
        : themeVars.foreground,
      padding: "0.5rem 0.75rem",
      fontSize: "0.875rem",
    }),
    input: (base: any) => ({
      ...base,
      color: themeVars.foreground,
    }),
    singleValue: (base: any) => ({
      ...base,
      color: themeVars.foreground,
    }),
    placeholder: (base: any) => ({
      ...base,
      color: themeVars.mutedForeground,
    }),
  };

  return (
    <div className="space-y-2">
      {label && <Label htmlFor={id || label || ""}>{label}{required && requiredSign && <span className="text-destructive">*</span>}</Label>}
      <Select
        id={id || label || ""}
        options={options}
        value={value}
        onChange={onChange}
        isSearchable
        menuPortalTarget={document.body}
        styles={customStyles}
        className={error ? "border-destructive" : ""}
        isMulti={isMulti}
      />
      {error && (
        <Typography variant="c1" className="text-destructive">
          {error}
        </Typography>
      )}
    </div>
  );
}
