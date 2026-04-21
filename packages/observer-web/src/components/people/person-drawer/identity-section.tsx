import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormField } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { ageGroupKeys, sexKeys } from "@/constants/person";

interface IdentitySectionProps {
  form: {
    first_name: string;
    last_name: string;
    patronymic: string;
    sex: string;
    birth_date: string;
    age_group: string;
    primary_phone: string;
    email: string;
  };
  set: (key: string, value: string) => void;
}

export function IdentitySection({ form, set }: IdentitySectionProps) {
  const { t } = useTranslation();

  const sexOptions = Object.entries(sexKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  const ageGroupOptions = Object.entries(ageGroupKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.people.identity")}>
      <FormField
        label={t("project.people.firstName")}
        value={form.first_name}
        onChange={(v) => set("first_name", v)}
        required
      />
      <FormField
        label={t("project.people.lastName")}
        value={form.last_name}
        onChange={(v) => set("last_name", v)}
      />
      <FormField
        label={t("project.people.patronymic")}
        value={form.patronymic}
        onChange={(v) => set("patronymic", v)}
      />

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.sexLabel")}
        </Field.Label>
        <UISelect
          value={form.sex}
          onValueChange={(v) => set("sex", v)}
          options={sexOptions}
          fullWidth
        />
      </Field.Root>

      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.birthDate")}
        </span>
        <DatePicker
          value={form.birth_date}
          onChange={(v) => set("birth_date", v)}
          captionLayout="dropdown"
          clearable
        />
      </div>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.ageGroup")}
        </Field.Label>
        <UISelect
          value={form.age_group}
          onValueChange={(v) => set("age_group", v)}
          options={ageGroupOptions}
          fullWidth
          clearable
        />
      </Field.Root>

      <FormField
        label={t("project.people.phone")}
        value={form.primary_phone}
        onChange={(v) => set("primary_phone", v)}
      />
      <FormField
        label={t("project.people.email")}
        value={form.email}
        onChange={(v) => set("email", v)}
        type="email"
      />
    </FormSection>
  );
}
