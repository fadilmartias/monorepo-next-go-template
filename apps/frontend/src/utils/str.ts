export function snakeToTitle(str: string): string {
    return str
      .split('_')
      .map(w => w.charAt(0).toUpperCase() + w.slice(1))
      .join(' ');
  }

export function truncateText(text: string | undefined, max: number): string {
  if (!text) return "";
  return text.length > max ? text.substring(0, max) + "..." : text;
}

  
  