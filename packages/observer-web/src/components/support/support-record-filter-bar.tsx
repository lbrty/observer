import type { ReactNode } from "react";

import { useTranslation } from "react-i18next";

import type { FilterDef } from "@/components/forms/filter-bar";
import { FilterBar } from "@/components/forms/filter-bar";
import { sphereKeys } from "@/constants/i18n";

const sphereValues = [
  "housing_assistance",
  "document_recovery",
  "social_benefits",
  "property_rights",
  "employment_rights",
  "family_law",
  "healthcare_access",
  "education_access",
  "financial_aid",
  "psychological_support",
  "other",
] as const;

interface SupportRecordFilterBarProps {
  sphere: string;
  onSphereChange: (value: string) => void;
  dateFrom: string;
  dateTo: string;
  onDateFromChange: (value: string) => void;
  onDateToChange: (value: string) => void;
  trailing?: ReactNode;
}

export function SupportRecordFilterBar({
  sphere,
  onSphereChange,
  dateFrom,
  dateTo,
  onDateFromChange,
  onDateToChange,
  trailing,
}: SupportRecordFilterBarProps) {
  const { t } = useTranslation();

  const sphereOptions = [
    { label: t("project.supportRecords.allSpheres"), value: "" },
    ...sphereValues.map((s) => ({
      label: sphereKeys[s] ? t(sphereKeys[s]) : s,
      value: s,
    })),
  ];

  const filters: FilterDef[] = [
    {
      type: "select",
      value: sphere,
      onValueChange: onSphereChange,
      options: sphereOptions,
      placeholder: t("project.supportRecords.allSpheres"),
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

  return <FilterBar filters={filters} trailing={trailing} />;
}
