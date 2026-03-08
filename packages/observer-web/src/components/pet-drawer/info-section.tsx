import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { FormField, FormTextarea } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { PersonCombobox } from "@/components/person-combobox";
import { PersonName } from "@/components/person-name";
import { UISelect } from "@/components/ui-select";
import { petStatusKeys } from "@/constants/pet";

interface InfoSectionProps {
  name: string;
  status: string;
  ownerId: string;
  ownerName: string;
  registrationId: string;
  notes: string;
  projectId: string;
  onNameChange: (v: string) => void;
  onStatusChange: (v: string) => void;
  onOwnerSelect: (id: string, name: string) => void;
  onOwnerClear: () => void;
  onRegistrationIdChange: (v: string) => void;
  onNotesChange: (v: string) => void;
}

export function InfoSection({
  name,
  status,
  ownerId,
  ownerName,
  registrationId,
  notes,
  projectId,
  onNameChange,
  onStatusChange,
  onOwnerSelect,
  onOwnerClear,
  onRegistrationIdChange,
  onNotesChange,
}: InfoSectionProps) {
  const { t } = useTranslation();

  const statusOptions = Object.entries(petStatusKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <fieldset className="space-y-3">
      <FormSection title={t("project.pets.title")}>
        <div className="col-span-full rounded-xl border border-border-secondary bg-bg p-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={t("project.pets.name")}
              value={name}
              onChange={onNameChange}
              required
            />

            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.pets.status")}
              </Field.Label>
              <UISelect
                value={status}
                onValueChange={onStatusChange}
                options={statusOptions}
                fullWidth
              />
            </Field.Root>

            <div>
              <span className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.pets.ownerId")}
              </span>
              {ownerId ? (
                <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
                  <span className="flex-1 truncate text-sm text-fg">
                    {ownerName || <PersonName projectId={projectId} personId={ownerId} />}
                  </span>
                  <button
                    type="button"
                    onClick={onOwnerClear}
                    className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
                  >
                    ×
                  </button>
                </div>
              ) : (
                <PersonCombobox
                  projectId={projectId}
                  onSelect={(p) =>
                    onOwnerSelect(p.id, `${p.first_name} ${p.last_name ?? ""}`.trim())
                  }
                />
              )}
            </div>

            <FormField
              label={t("project.pets.registrationId")}
              value={registrationId}
              onChange={onRegistrationIdChange}
            />
          </div>

          <div className="mt-4">
            <FormTextarea
              label={t("project.pets.notes")}
              value={notes}
              onChange={onNotesChange}
            />
          </div>
        </div>
      </FormSection>
    </fieldset>
  );
}
