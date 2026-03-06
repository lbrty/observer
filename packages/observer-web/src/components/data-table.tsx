import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";

import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";

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
  emptyState?: ReactNode;
}

function toColumnDefs<T>(columns: Column<T>[]): ColumnDef<T, unknown>[] {
  return columns.map((col) => ({
    id: col.key,
    header: () => col.header,
    cell: ({ row }) => col.render(row.original),
    meta: { className: col.className },
  }));
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
  emptyState,
}: DataTableProps<T>) {
  const { t } = useTranslation();

  const table = useReactTable({
    data,
    columns: toColumnDefs(columns),
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row) => keyExtractor(row),
  });

  return (
    <div className="overflow-hidden rounded-xl border border-border-secondary bg-bg-secondary shadow-card">
      <table className="w-full text-left text-sm">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr
              key={headerGroup.id}
              className="border-b border-border-secondary bg-bg-tertiary/60"
            >
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  className={`px-4 py-2.5 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary ${(header.column.columnDef.meta as { className?: string })?.className ?? ""}`}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-border-secondary/60">
          {isLoading ? (
            <SkeletonRows cols={columns.length} />
          ) : table.getRowModel().rows.length === 0 ? (
            <tr>
              <td colSpan={columns.length}>
                {emptyState ?? (
                  <p className="px-4 py-12 text-center text-fg-tertiary">
                    {t("admin.common.noData")}
                  </p>
                )}
              </td>
            </tr>
          ) : (
            table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                onClick={
                  onRowClick ? () => onRowClick(row.original) : undefined
                }
                className={`transition-colors ${onRowClick ? "cursor-pointer hover:bg-bg-tertiary/40" : ""}`}
              >
                {row.getVisibleCells().map((cell) => (
                  <td
                    key={cell.id}
                    className={`px-4 py-3 ${(cell.column.columnDef.meta as { className?: string })?.className ?? ""}`}
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
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
