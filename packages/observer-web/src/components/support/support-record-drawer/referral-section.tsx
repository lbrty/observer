import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { referralKeys } from "@/constants/support";

interface ReferralSectionProps {
  form: {
    referral_status: string;
    referred_to_office: string;
  };
  set: (key: string, value: string) => void;
  officeOptions: { label: string; value: string }[];
}

export function ReferralSection({ form, set, officeOptions }: ReferralSectionProps) {
  const { t } = useTranslation();

  const referralStatusOptions = Object.entries(referralKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.supportRecords.referral")}>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.supportRecords.referralStatus")}
        </Field.Label>
        <UISelect
          value={form.referral_status}
          onValueChange={(v) => set("referral_status", v)}
          options={referralStatusOptions}
          placeholder={t("project.supportRecords.selectReferralStatus")}
          fullWidth
        />
      </Field.Root>

      {officeOptions.length > 0 && (
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.referredToOffice")}
          </Field.Label>
          <UISelect
            value={form.referred_to_office}
            onValueChange={(v) => set("referred_to_office", v)}
            options={officeOptions}
            placeholder={t("project.supportRecords.selectOffice")}
            fullWidth
          />
        </Field.Root>
      )}
    </FormSection>
  );
}
