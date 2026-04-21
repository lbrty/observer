import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { FilterBar } from "@/components/forms/filter-bar";
import { DownloadSimpleIcon } from "@/components/ui/icons";
import { SelectedTagChips, TagFilter } from "@/components/tags/tag-filter";

interface PetFilterBarProps {
  projectId: string;
  dateFrom: string;
  dateTo: string;
  tagIds: string[];
  canExport: boolean;
  exporting: boolean;
  onDateFromChange: (value: string) => void;
  onDateToChange: (value: string) => void;
  onTagsChange: (ids: string[]) => void;
  onExport: () => void;
}

export function PetFilterBar({
  projectId,
  dateFrom,
  dateTo,
  tagIds,
  canExport,
  exporting,
  onDateFromChange,
  onDateToChange,
  onTagsChange,
  onExport,
}: PetFilterBarProps) {
  const { t } = useTranslation();

  const filters = [
    {
      type: "date-range" as const,
      fromValue: dateFrom,
      toValue: dateTo,
      onFromChange: onDateFromChange,
      onToChange: onDateToChange,
      fromPlaceholder: t("common.dateFrom"),
      toPlaceholder: t("common.dateTo"),
    },
  ];

  return (
    <>
      <FilterBar
        filters={filters}
        trailing={
          <div className="flex items-center gap-2">
            <TagFilter projectId={projectId} selectedIds={tagIds} onChange={onTagsChange} />
            {canExport && (
              <Button
                variant="secondary"
                icon={<DownloadSimpleIcon size={16} />}
                onClick={onExport}
                disabled={exporting}
              >
                {t("common.export")}
              </Button>
            )}
          </div>
        }
      />
      <SelectedTagChips projectId={projectId} selectedIds={tagIds} onChange={onTagsChange} />
    </>
  );
}
