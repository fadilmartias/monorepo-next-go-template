export function disableConsoleInProduction() {
    if (process.env.NODE_ENV === "production" || process.env.NEXT_PUBLIC_APP_ENV === "production") {
      console.log = () => {};
      console.info = () => {};
      console.debug = () => {};
    }
  }
  