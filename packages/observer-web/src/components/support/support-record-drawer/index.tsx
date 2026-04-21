import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/ui/alert-banner";
import { DrawerShell } from "@/components/drawer/drawer-shell";
import { FormTextarea } from "@/components/forms/form-field";
import { SectionHeading } from "@/components/layout/section-heading";

import { InfoSection } from "./info-section";
import { ReferralSection } from "./referral-section";
import { useSupportRecordForm } from "./use-support-record-form";

interface SupportRecordDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  recordId: string | null;
  personId?: string;
}

export function SupportRecordDrawer({
  open,
  onOpenChange,
  projectId,
  recordId,
  personId,
}: SupportRecordDrawerProps) {
  const { t } = useTranslation();
  const {
    isEdit,
    form,
    set,
    error,
    isPending,
    handleSubmit,
    personName,
    setPersonName,
    officeOptions,
  } = useSupportRecordForm({ open, projectId, recordId, personId });

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={
        isEdit
          ? `${t("project.supportRecords.editTitle")}${personName ? ` — ${personName}` : ""}`
          : t("project.supportRecords.formTitle")
      }
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.supportRecords.save")}
      savingLabel={t("project.supportRecords.saving")}
    >
      <ErrorBanner message={error} />

      <InfoSection
        form={form}
        set={(k, v) => set(k as keyof typeof form, v)}
        projectId={projectId}
        personId={personId}
        isEdit={isEdit}
        personName={personName}
        onSelectPerson={(p) => {
          set("person_id", p.id);
          setPersonName(`${p.first_name} ${p.last_name ?? ""}`.trim());
        }}
        onClearPerson={() => {
          set("person_id", "");
          setPersonName("");
        }}
      />

      <ReferralSection
        form={form}
        set={(k, v) => set(k as keyof typeof form, v)}
        officeOptions={officeOptions}
      />

      <SectionHeading>{t("project.supportRecords.notesSection")}</SectionHeading>
      <FormTextarea label="" value={form.notes} onChange={(v) => set("notes", v)} rows={4} />
    </DrawerShell>
  );
}
