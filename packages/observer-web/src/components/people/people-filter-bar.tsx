import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import type { FilterDef } from "@/components/filter-bar";
import { FilterBar } from "@/components/filter-bar";
import { DownloadSimpleIcon } from "@/components/icons";
import { SelectedTagChips, TagFilter } from "@/components/tag-filter";

interface PeopleFilterBarProps {
  projectId: string;
  search: string;
  onSearchChange: (value: string) => void;
  sex: string;
  onSexChange: (value: string) => void;
  ageGroup: string;
  onAgeGroupChange: (value: string) => void;
  dateFrom: string;
  onDateFromChange: (value: string) => void;
  dateTo: string;
  onDateToChange: (value: string) => void;
  tagIds: string[];
  onTagIdsChange: (ids: string[]) => void;
  canExport: boolean;
  exporting: boolean;
  onExport: () => void;
}

export function PeopleFilterBar({
  projectId,
  search,
  onSearchChange,
  sex,
  onSexChange,
  ageGroup,
  onAgeGroupChange,
  dateFrom,
  onDateFromChange,
  dateTo,
  onDateToChange,
  tagIds,
  onTagIdsChange,
  canExport,
  exporting,
  onExport,
}: PeopleFilterBarProps) {
  const { t } = useTranslation();

  const sexOptions = [
    { label: t("project.people.allSex"), value: "" },
    { label: t("project.people.sexMale"), value: "male" },
    { label: t("project.people.sexFemale"), value: "female" },
    { label: t("project.people.sexOther"), value: "other" },
    { label: t("project.people.sexUnknown"), value: "unknown" },
  ];

  const ageGroupOptions = [
    { label: t("project.people.allAgeGroups"), value: "" },
    { label: t("project.people.ageInfant"), value: "infant" },
    { label: t("project.people.ageToddler"), value: "toddler" },
    { label: t("project.people.agePreSchool"), value: "pre_school" },
    { label: t("project.people.ageMiddleChildhood"), value: "middle_childhood" },
    { label: t("project.people.ageYoungTeen"), value: "young_teen" },
    { label: t("project.people.ageTeenager"), value: "teenager" },
    { label: t("project.people.ageYoungAdult"), value: "young_adult" },
    { label: t("project.people.ageEarlyAdult"), value: "early_adult" },
    { label: t("project.people.ageMiddleAgedAdult"), value: "middle_aged_adult" },
    { label: t("project.people.ageOldAdult"), value: "old_adult" },
  ];

  const filters: FilterDef[] = [
    {
      type: "search",
      placeholder: t("project.people.search"),
      value: search,
      onChange: onSearchChange,
    },
    {
      type: "select",
      value: sex,
      onValueChange: onSexChange,
      options: sexOptions,
      placeholder: t("project.people.allSex"),
    },
    {
      type: "select",
      value: ageGroup,
      onValueChange: onAgeGroupChange,
      options: ageGroupOptions,
      placeholder: t("project.people.allAgeGroups"),
    },
    {
      type: "date-range",
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
            <TagFilter projectId={projectId} selectedIds={tagIds} onChange={onTagIdsChange} />
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
      <SelectedTagChips projectId={projectId} selectedIds={tagIds} onChange={onTagIdsChange} />
    </>
  );
}
