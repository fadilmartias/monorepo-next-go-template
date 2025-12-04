export function setFormErrors<T extends Record<string, any>>(
    setError: (name: keyof T, error: { type: string; message: string }) => void,
    errors: Record<string, string>
  ) {
    Object.entries(errors).forEach(([field, message]) => {
      setError(field as keyof T, {
        type: "manual",
        message,
      })
    })
  }
  