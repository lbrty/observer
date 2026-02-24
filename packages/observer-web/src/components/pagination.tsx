import { CaretLeftIcon, CaretRightIcon } from "@/components/icons";
import { useTranslation } from "react-i18next";

interface PaginationProps {
  page: number;
  perPage: number;
  total: number;
  onChange: (page: number) => void;
}

function getPageNumbers(current: number, total: number): (number | "...")[] {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const pages: (number | "...")[] = [1];

  if (current > 3) pages.push("...");

  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  if (current < total - 2) pages.push("...");

  pages.push(total);
  return pages;
}

export function Pagination({
  page,
  perPage,
  total,
  onChange,
}: PaginationProps) {
  const { t } = useTranslation();
  const totalPages = Math.max(1, Math.ceil(total / perPage));

  if (totalPages <= 1) return null;

  const from = (page - 1) * perPage + 1;
  const to = Math.min(page * perPage, total);
  const pages = getPageNumbers(page, totalPages);

  return (
    <div className="flex items-center justify-between pt-4">
      <span className="text-sm text-fg-tertiary">
        {t("admin.common.paginationRange", { from, to, total })}
      </span>
      <div className="flex items-center gap-1">
        <button
          type="button"
          disabled={page <= 1}
          onClick={() => onChange(page - 1)}
          className="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg text-sm text-fg-secondary hover:bg-bg-tertiary disabled:cursor-not-allowed disabled:opacity-30"
        >
          <CaretLeftIcon size={14} />
        </button>
        {pages.map((p, i) =>
          p === "..." ? (
            <span
              key={`ellipsis-${i}`}
              className="inline-flex size-8 items-center justify-center text-sm text-fg-tertiary"
            >
              ...
            </span>
          ) : (
            <button
              key={p}
              type="button"
              onClick={() => onChange(p)}
              className={`inline-flex size-8 cursor-pointer items-center justify-center rounded-lg text-sm font-medium transition-colors ${
                p === page
                  ? "bg-accent text-accent-fg shadow-card"
                  : "text-fg-secondary hover:bg-bg-tertiary"
              }`}
            >
              {p}
            </button>
          ),
        )}
        <button
          type="button"
          disabled={page >= totalPages}
          onClick={() => onChange(page + 1)}
          className="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg text-sm text-fg-secondary hover:bg-bg-tertiary disabled:cursor-not-allowed disabled:opacity-30"
        >
          <CaretRightIcon size={14} />
        </button>
      </div>
    </div>
  );
}
