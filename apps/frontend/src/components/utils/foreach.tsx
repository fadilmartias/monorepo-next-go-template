type ForeachProps<T> = {
    of: T[];
    render: (item: T, index: number) => React.ReactNode;
  };
  
  export function Foreach<T>({ of, render }: ForeachProps<T>) {
    return <>{of.map(render)}</>;
  }
  