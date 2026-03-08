import type { ReactNode } from "react";

import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { type FilterDef, FilterBar } from "@/components/filter-bar";
import type { Icon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";

interface PaginationConfig {
  page: number;
  perPage: number;
  total: number;
  onChange: (page: number) => void;
}

interface DataTablePageProps<T> {
  title: string;
  columns: Column<T>[];
  data: T[];
  keyExtractor: (item: T) => string;
  isLoading?: boolean;
  onRowClick?: (item: T) => void;
  pagination?: PaginationConfig;
  filters?: FilterDef[];
  filterTrailing?: ReactNode;
  emptyIcon?: Icon;
  emptyTitle?: string;
  emptyDescription?: string;
  emptyAction?: ReactNode;
  createAction?: ReactNode;
  children?: ReactNode;
}

export function DataTablePage<T>({
  title, columns, data, keyExtractor, isLoading, onRowClick,
  pagination, filters, filterTrailing,
  emptyIcon, emptyTitle, emptyDescription, emptyAction,
  createAction, children,
}: DataTablePageProps<T>) {
  return (
    <div>
      <PageHeader title={title} action={createAction} />

      {filters && filters.length > 0 && (
        <FilterBar filters={filters} trailing={filterTrailing} />
      )}

      <DataTable
        columns={columns}
        data={data}
        keyExtractor={keyExtractor}
        onRowClick={onRowClick}
        isLoading={isLoading}
        emptyState={
          emptyIcon && emptyTitle ? (
            <EmptyState
              icon={emptyIcon}
              title={emptyTitle}
              description={emptyDescription}
              action={emptyAction}
            />
          ) : undefined
        }
      />

      {pagination && (
        <Pagination
          page={pagination.page}
          perPage={pagination.perPage}
          total={pagination.total}
          onChange={pagination.onChange}
        />
      )}

      {children}
    </div>
  );
}
