import { useTranslation } from "react-i18next";

import { FormField } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { PersonCombobox } from "@/components/person-combobox";

interface HeadSectionProps {
  referenceNumber: string;
  headPersonId: string;
  headPersonLabel: string;
  projectId: string;
  onReferenceNumberChange: (v: string) => void;
  onHeadPersonSelect: (id: string, name: string) => void;
  onHeadPersonClear: () => void;
}

export function HeadSection({
  referenceNumber,
  headPersonId,
  headPersonLabel,
  projectId,
  onReferenceNumberChange,
  onHeadPersonSelect,
  onHeadPersonClear,
}: HeadSectionProps) {
  const { t } = useTranslation();

  return (
    <FormSection title={t("project.households.info")}>
      <FormField
        label={t("project.households.referenceNumber")}
        value={referenceNumber}
        onChange={onReferenceNumberChange}
      />
      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.households.headPerson")}
        </span>
        {headPersonId ? (
          <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
            <span className="flex-1 truncate text-sm text-fg">
              {headPersonLabel}
            </span>
            <button
              type="button"
              onClick={onHeadPersonClear}
              className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
            >
              ×
            </button>
          </div>
        ) : (
          <PersonCombobox
            projectId={projectId}
            onSelect={(p) =>
              onHeadPersonSelect(p.id, `${p.first_name} ${p.last_name ?? ""}`.trim())
            }
          />
        )}
      </div>
    </FormSection>
  );
}
