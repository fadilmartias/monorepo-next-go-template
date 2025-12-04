import * as React from "react";
import Typography from "../typography";
import { cn } from "@/lib/utils";
import { Label } from "./label";
import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

function Input({
  className,
  label,
  id,
  type,
  error,
  required,
  requiredSign = true,
  onChange,
  onScroll,
  leftIcon,
  rightIcon,
  labelClassName,
  ...props
}: React.ComponentProps<"input"> & {
  error?: string;
  label?: string;
  id?: string;
  required?: boolean;
  requiredSign?: boolean;
  type?: string;
  labelClassName?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onScroll?: (e: React.WheelEvent<HTMLInputElement>) => void;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
}) {
  const [showPassword, setShowPassword] = useState(false);

  const isPassword = type === "password";
  const inputType = isPassword && showPassword ? "text" : type;
  const handleScroll = (e: React.WheelEvent<HTMLInputElement>) => {
    if (type === "number") {
      e.preventDefault();
      e.stopPropagation();
    }
    if(onScroll) {
      onScroll(e);
    }
  };
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let value = e.target.value;

    if (type === "number") {
      value = value.replace(/[^0-9]/g, ""); // hanya angka

      // hapus leading zero (kecuali kalau memang "0")
      if (value.length > 1 && value.startsWith("0")) {
        value = value.replace(/^0+/, "");
      }

      e.target.value = value;
    }

    if (type === "decimal") {
      value = value.replace(/[^0-9.,]/g, "");

      // cuma boleh ada satu koma/titik
      value = value.replace(/([.,].*)[.,]/g, "$1");

      // hapus leading zero kecuali kalau ada koma/titik
      if (value.length > 1 && value.startsWith("0") && !/[.,]/.test(value)) {
        value = value.replace(/^0+/, "");
      }

      // normalize: ubah koma jadi titik
      const normalized = value.replace(",", ".");

      e.target.value = normalized;
    }

    if (onChange) {
      onChange(e);
    }
  };
  return (
    <div className="flex flex-col gap-2 w-full">
      {label && (
        <Label htmlFor={id || label || ""} className={labelClassName}>
          {label}
          {required && requiredSign && (
            <span className="text-destructive">*</span>
          )}
        </Label>
      )}
      <div className="relative">
        {leftIcon && (
          <span className={cn("absolute left-3 top-1/2 -translate-y-1/2", props.disabled && "opacity-50")}>
            {leftIcon}
          </span>
        )}
        <input
          id={id || label || ""}
          type={inputType}
          autoComplete={isPassword ? "off" : props.autoComplete}
          data-slot="input"
          className={cn(
            "file:text-foreground placeholder:text-muted-foreground selection:bg-primary placeholder:text-sm selection:text-primary-foreground dark:bg-input/30 border-input flex h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
            "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
            "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
            leftIcon && "pl-10",
            rightIcon && "pr-10",
            className,
            error && "border-destructive"
          )}
          onChange={handleChange}
          {...props}
        />

        {rightIcon && (
          <span className={cn("absolute right-3 top-1/2 -translate-y-1/2", props.disabled && "opacity-50")}>
            {rightIcon}
          </span>
        )}

        {isPassword && (
          <button
            type="button"
            className={cn("absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground", props.disabled && "opacity-50")}
            onClick={() => setShowPassword(!showPassword)}
            tabIndex={-1}
          >
            {showPassword ? (
              <EyeOff className="w-4 h-4" />
            ) : (
              <Eye className="w-4 h-4" />
            )}
          </button>
        )}
      </div>

      {error && (
        <Typography variant="c1" className="text-destructive">
          {error}
        </Typography>
      )}
    </div>
  );
}

export { Input };
