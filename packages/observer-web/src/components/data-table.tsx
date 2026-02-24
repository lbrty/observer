import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";

export interface Column<T> {
  key: string;
  header: string;
  render: (item: T) => ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyExtractor: (item: T) => string;
  onRowClick?: (item: T) => void;
  isLoading?: boolean;
}

function SkeletonRows({ cols, rows = 5 }: { cols: number; rows?: number }) {
  return Array.from({ length: rows }, (_, i) => (
    <tr key={i}>
      {Array.from({ length: cols }, (_, j) => (
        <td key={j} className="px-4 py-3">
          <div className="h-4 w-3/4 animate-pulse rounded bg-bg-tertiary" />
        </td>
      ))}
    </tr>
  ));
}

export function DataTable<T>({
  columns,
  data,
  keyExtractor,
  onRowClick,
  isLoading,
}: DataTableProps<T>) {
  const { t } = useTranslation();

  return (
    <div className="overflow-x-auto rounded-lg border border-border-secondary">
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b border-border-secondary bg-bg-tertiary">
            {columns.map((col) => (
              <th
                key={col.key}
                className={`px-4 py-2.5 font-medium text-fg-secondary ${col.className ?? ""}`}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-border-secondary">
          {isLoading ? (
            <SkeletonRows cols={columns.length} />
          ) : data.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="px-4 py-8 text-center text-fg-tertiary"
              >
                {t("admin.common.noData")}
              </td>
            </tr>
          ) : (
            data.map((item) => (
              <tr
                key={keyExtractor(item)}
                onClick={onRowClick ? () => onRowClick(item) : undefined}
                className={`bg-bg-secondary ${onRowClick ? "cursor-pointer hover:bg-bg-tertiary" : ""}`}
              >
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className={`px-4 py-3 ${col.className ?? ""}`}
                  >
                    {col.render(item)}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
