import { debounce } from "@/utils/debounce";
import { Input } from "../ui/input";
import { useCallback, useState } from "react";
import { SearchIcon } from "lucide-react";

interface MultiColumnFilterProps {
  columns: string[];
  table: any;
  disabled?: boolean;
  serverSide?: boolean;
  filter?: { [key: string]: any };
  setFilter?: (filter: { [key: string]: any }) => void;
}

export const DataTableSearch: React.FC<MultiColumnFilterProps> = ({
  columns,
  table,
  disabled,
  serverSide,
  filter,
  setFilter,
}) => {
  // Local state for input value
  const [inputValue, setInputValue] = useState(
    serverSide
      ? filter?.[columns[0]] || ""
      : columns
          .map((col) => table.getColumn(col)?.getFilterValue() as string)
          .find((v) => !!v) ?? ""
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setInputValue(value);

    if (serverSide && setFilter) {
      if (value) {
        const orFilter = columns.reduce(
          (acc, col) => ({
            ...acc,
            [col]: { like: `%${value}%` },
          }),
          {}
        );

        setFilter({
          or: orFilter,
        });
      } else {
        setFilter({});
      }
    } else {
      // client-side
      columns.forEach((col) => {
        table.getColumn(col)?.setFilterValue({ like: `%${value}%` });
      });
    }
  };

  return (
    <Input
      leftIcon={<SearchIcon className="w-4 h-4 opacity-50" />}
      placeholder={`Search...`}
      value={inputValue}
      onChange={handleChange}
      className="max-w-sm"
      disabled={disabled}
    />
  );
};
