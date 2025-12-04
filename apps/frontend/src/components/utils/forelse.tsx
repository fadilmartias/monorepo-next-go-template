type ForelseProps<T> = {
    of: T[];
    render: (item: T, index: number) => React.ReactNode;
    empty: React.ReactNode;
  };
  
  export function Forelse<T>({ of, render, empty }: ForelseProps<T>) {
    return <>{of.length > 0 ? of.map(render) : empty}</>;
  }
  