import type { TFunction } from "i18next";

import { Button } from "@/components/ui/button";
import { DataTable, type Column } from "@/components/table/data-table";
import { EmptyState } from "@/components/ui/empty-state";
import { PlusIcon, UserCircleIcon } from "@/components/ui/icons";
import { Pagination } from "@/components/table/pagination";
import type { Person } from "@/types/person";

interface PeopleTableProps {
  columns: Column<Person>[];
  data: Person[];
  isLoading: boolean;
  onRowClick: (person: Person) => void;
  canWrite: boolean;
  onCreateClick: () => void;
  t: TFunction;
  pagination?: { page: number; perPage: number; total: number; onChange: (page: number) => void };
}

export function PeopleTable({
  columns,
  data,
  isLoading,
  onRowClick,
  canWrite,
  onCreateClick,
  t,
  pagination,
}: PeopleTableProps) {
  return (
    <>
      <DataTable
        columns={columns}
        data={data}
        keyExtractor={(p) => p.id}
        onRowClick={onRowClick}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={UserCircleIcon}
            title={t("project.people.emptyTitle")}
            description={t("project.people.emptyDescription")}
            action={
              canWrite ? (
                <Button onClick={onCreateClick} icon={<PlusIcon size={16} />}>
                  {t("project.people.register")}
                </Button>
              ) : undefined
            }
          />
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
    </>
  );
}
