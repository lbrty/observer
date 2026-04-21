import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { reasonKeys, housingKeys } from "@/constants/migration";

interface DetailsSectionProps {
  migrationDate: string;
  movementReason: string;
  housingAtDestination: string;
  onDateChange: (v: string) => void;
  onReasonChange: (v: string) => void;
  onHousingChange: (v: string) => void;
}

export function DetailsSection({
  migrationDate,
  movementReason,
  housingAtDestination,
  onDateChange,
  onReasonChange,
  onHousingChange,
}: DetailsSectionProps) {
  const { t } = useTranslation();

  const reasonOptions = Object.entries(reasonKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  const housingOptions = Object.entries(housingKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.migrationRecords.details")}>
      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.migrationRecords.date")}
        </span>
        <DatePicker value={migrationDate} onChange={onDateChange} />
      </div>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.migrationRecords.reason")}
        </Field.Label>
        <UISelect
          value={movementReason}
          onValueChange={onReasonChange}
          options={reasonOptions}
          fullWidth
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.migrationRecords.housing")}
        </Field.Label>
        <UISelect
          value={housingAtDestination}
          onValueChange={onHousingChange}
          options={housingOptions}
          fullWidth
        />
      </Field.Root>
    </FormSection>
  );
}
