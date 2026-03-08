import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormField } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { caseStatusKeys } from "@/constants/person";

interface CaseSectionProps {
  form: {
    case_status: string;
    external_id: string;
    office_id: string;
    consent_given: boolean;
    consent_date: string;
  };
  // biome-ignore lint: section receives parent's polymorphic set
  set: (key: string, value: any) => void;
  officeOptions: { label: string; value: string }[];
}

export function CaseSection({ form, set, officeOptions }: CaseSectionProps) {
  const { t } = useTranslation();

  const caseStatusOptions = Object.entries(caseStatusKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.people.case")}>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.caseStatusLabel")}
        </Field.Label>
        <UISelect
          value={form.case_status}
          onValueChange={(v) => set("case_status", v)}
          options={caseStatusOptions}
          fullWidth
        />
      </Field.Root>

      <FormField
        label={t("project.people.externalId")}
        value={form.external_id}
        onChange={(v) => set("external_id", v)}
      />

      {officeOptions.length > 0 && (
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.office")}
          </Field.Label>
          <UISelect
            value={form.office_id}
            onValueChange={(v) => set("office_id", v)}
            options={officeOptions}
            fullWidth
          />
        </Field.Root>
      )}

      <div className="col-span-full space-y-4">
        <UISwitch
          checked={form.consent_given}
          onCheckedChange={(v) => set("consent_given", v)}
          label={t("project.people.consentGiven")}
        />

        {form.consent_given && (
          <div>
            <span className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.people.consentDate")}
            </span>
            <DatePicker value={form.consent_date} onChange={(v) => set("consent_date", v)} />
          </div>
        )}
      </div>
    </FormSection>
  );
}
