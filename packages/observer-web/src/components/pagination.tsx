import { CaretLeft, CaretRight } from "@phosphor-icons/react";
import { useTranslation } from "react-i18next";

interface PaginationProps {
  page: number;
  perPage: number;
  total: number;
  onChange: (page: number) => void;
}

export function Pagination({
  page,
  perPage,
  total,
  onChange,
}: PaginationProps) {
  const { t } = useTranslation();
  const totalPages = Math.max(1, Math.ceil(total / perPage));
  const hasPrev = page > 1;
  const hasNext = page < totalPages;

  return (
    <div className="flex items-center justify-between pt-4">
      <span className="text-sm text-fg-tertiary">
        {t("admin.common.pagination", { page, totalPages })}
      </span>
      <div className="flex gap-1">
        <button
          type="button"
          disabled={!hasPrev}
          onClick={() => onChange(page - 1)}
          className="inline-flex cursor-pointer items-center rounded-md border border-border-secondary px-2 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary disabled:cursor-not-allowed disabled:opacity-40"
        >
          <CaretLeft size={16} />
        </button>
        <button
          type="button"
          disabled={!hasNext}
          onClick={() => onChange(page + 1)}
          className="inline-flex cursor-pointer items-center rounded-md border border-border-secondary px-2 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary disabled:cursor-not-allowed disabled:opacity-40"
        >
          <CaretRight size={16} />
        </button>
      </div>
    </div>
  );
}
