import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";

import { HeadSection } from "./head-section";
import { MembersSection } from "./members-section";
import { useHouseholdDrawerForm } from "./use-household-drawer-form";

interface HouseholdDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  householdId: string | null;
}

export function HouseholdDrawer({
  open,
  onOpenChange,
  projectId,
  householdId,
}: HouseholdDrawerProps) {
  const { t } = useTranslation();
  const {
    isEdit,
    form,
    set,
    error,
    isPending,
    handleSubmit,
    handleAddMember,
    handleRemoveMember,
    handleHeadPersonSelect,
    handleHeadPersonClear,
    memberForm,
    setMemberForm,
    headPersonLabel,
    memberPersonName,
    setMemberPersonName,
    addMemberPending,
    editingId,
    household,
  } = useHouseholdDrawerForm({ open, projectId, householdId });

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? t("project.households.editTitle") : t("project.households.formTitle")}
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.households.save")}
      savingLabel={t("project.households.saving")}
    >
      {error && <ErrorBanner message={error} />}

      <HeadSection
        referenceNumber={form.reference_number}
        headPersonId={form.head_person_id}
        headPersonLabel={headPersonLabel}
        projectId={projectId}
        onReferenceNumberChange={(v) => set("reference_number", v)}
        onHeadPersonSelect={handleHeadPersonSelect}
        onHeadPersonClear={handleHeadPersonClear}
      />

      {isEdit && editingId && (
        <MembersSection
          editingId={editingId}
          household={household}
          projectId={projectId}
          memberForm={memberForm}
          memberPersonName={memberPersonName}
          addMemberPending={addMemberPending}
          onMemberFormChange={setMemberForm}
          onMemberPersonNameChange={setMemberPersonName}
          onAddMember={handleAddMember}
          onRemoveMember={handleRemoveMember}
          onCloseDrawer={() => onOpenChange(false)}
        />
      )}
    </DrawerShell>
  );
}
