import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormSection } from "@/components/form-section";
import { CopySimpleIcon } from "@/components/icons";
import { PersonCombobox } from "@/components/person-combobox";
import { Tooltip } from "@/components/tooltip";
import { UISelect } from "@/components/ui-select";
import { sphereKeys, typeKeys } from "@/constants/support";
import { useToast } from "@/stores/toast";

interface InfoSectionProps {
  form: {
    type: string;
    sphere: string;
    provided_at: string;
    person_id: string;
  };
  set: (key: string, value: string) => void;
  projectId: string;
  personId?: string;
  isEdit: boolean;
  personName: string;
  onSelectPerson: (p: { id: string; first_name: string; last_name?: string | null }) => void;
  onClearPerson: () => void;
}

export function InfoSection({
  form,
  set,
  projectId,
  personId,
  isEdit,
  personName,
  onSelectPerson,
  onClearPerson,
}: InfoSectionProps) {
  const { t } = useTranslation();
  const toast = useToast();

  const typeOptions = Object.entries(typeKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  const sphereOptions = Object.entries(sphereKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.supportRecords.recordInfo")}>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.supportRecords.type")} *
        </Field.Label>
        <UISelect
          value={form.type}
          onValueChange={(v) => set("type", v)}
          options={typeOptions}
          fullWidth
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.supportRecords.sphere")}
        </Field.Label>
        <UISelect
          value={form.sphere}
          onValueChange={(v) => set("sphere", v)}
          options={sphereOptions}
          placeholder={t("project.supportRecords.selectSphere")}
          fullWidth
        />
      </Field.Root>

      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.supportRecords.providedAt")}
        </span>
        <DatePicker value={form.provided_at} onChange={(v) => set("provided_at", v)} />
      </div>

      {!personId && (
        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.person")} {!isEdit && " *"}
          </span>
          {form.person_id ? (
            <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
              <span className="flex-1 truncate text-sm text-fg">
                {personName || form.person_id}
              </span>
              <Tooltip label={t("admin.common.copyId")}>
                <button
                  type="button"
                  onClick={() => {
                    navigator.clipboard.writeText(form.person_id);
                    toast.success(t("admin.common.copied"));
                  }}
                  className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
                >
                  <CopySimpleIcon size={14} />
                </button>
              </Tooltip>
              {!isEdit && (
                <button
                  type="button"
                  onClick={onClearPerson}
                  className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
                >
                  ×
                </button>
              )}
            </div>
          ) : (
            <PersonCombobox projectId={projectId} onSelect={onSelectPerson} />
          )}
        </div>
      )}
    </FormSection>
  );
}
