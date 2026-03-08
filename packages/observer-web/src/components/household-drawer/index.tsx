import { type SyntheticEvent, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import {
  useAddHouseholdMember,
  useCreateHousehold,
  useHousehold,
  useRemoveHouseholdMember,
  useUpdateHousehold,
} from "@/hooks/use-households";
import { usePerson } from "@/hooks/use-people";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";

import type { Relationship } from "@/types/household";

import { HeadSection } from "./head-section";
import { MembersSection } from "./members-section";

interface HouseholdDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  householdId: string | null;
}

const emptyForm = {
  reference_number: "",
  head_person_id: "",
};

const emptyMemberForm = {
  person_id: "",
  relationship: "head" as Relationship,
};

export function HouseholdDrawer({
  open,
  onOpenChange,
  projectId,
  householdId,
}: HouseholdDrawerProps) {
  const { t } = useTranslation();
  const isEdit = householdId !== null;

  const { data: household } = useHousehold(projectId, householdId ?? "");
  const qc = useQueryClient();
  const createHousehold = useCreateHousehold(projectId);
  const updateHousehold = useUpdateHousehold(projectId);
  const addMember = useAddHouseholdMember(projectId);
  const removeMember = useRemoveHouseholdMember(projectId);

  const toast = useToast();
  const { form, set, error, setError, editingId, setEditingId } =
    useDrawerForm({
      initial: emptyForm,
      open,
      isEdit,
      data: household,
      mapData: (d) => ({
        reference_number: (d.reference_number as string) ?? "",
        head_person_id: (d.head_person_id as string) ?? "",
      }),
    });

  const [memberForm, setMemberForm] = useState(emptyMemberForm);
  const [headPersonName, setHeadPersonName] = useState("");
  const [memberPersonName, setMemberPersonName] = useState("");

  const { data: headPerson } = usePerson(
    projectId,
    isEdit && form.head_person_id && !headPersonName ? form.head_person_id : "",
  );

  const resolvedHeadName = headPersonName
    || (headPerson ? `${headPerson.first_name} ${headPerson.last_name ?? ""}`.trim() : "");
  const headPersonLabel = resolvedHeadName || form.head_person_id;

  const isPending = createHousehold.isPending || updateHousehold.isPending;

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    try {
      if (isEdit && editingId) {
        await updateHousehold.mutateAsync({
          id: editingId,
          data: {
            reference_number: form.reference_number || undefined,
            head_person_id: form.head_person_id || undefined,
          },
        });
        await qc.invalidateQueries({ queryKey: ["households", projectId] });
        toast.success(t("project.households.saved"));
      } else {
        const created = await createHousehold.mutateAsync({
          reference_number: form.reference_number || undefined,
          head_person_id: form.head_person_id || undefined,
        });
        await qc.invalidateQueries({ queryKey: ["households", projectId] });
        setEditingId(created.id);
        toast.success(t("project.households.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleAddMember() {
    if (!editingId || !memberForm.person_id) return;
    setError("");

    try {
      await addMember.mutateAsync({
        householdId: editingId,
        data: {
          person_id: memberForm.person_id,
          relationship: memberForm.relationship,
        },
      });
      await qc.invalidateQueries({
        queryKey: ["households", projectId, editingId],
      });
      setMemberForm(emptyMemberForm);
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleRemoveMember(personId: string) {
    if (!editingId) return;
    setError("");

    try {
      await removeMember.mutateAsync({ householdId: editingId, personId });
      await qc.invalidateQueries({
        queryKey: ["households", projectId, editingId],
      });
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

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
        onHeadPersonSelect={(id, name) => {
          set("head_person_id", id);
          setHeadPersonName(name);
        }}
        onHeadPersonClear={() => {
          set("head_person_id", "");
          setHeadPersonName("");
        }}
      />

      {isEdit && editingId && (
        <MembersSection
          editingId={editingId}
          household={household}
          projectId={projectId}
          memberForm={memberForm}
          memberPersonName={memberPersonName}
          addMemberPending={addMember.isPending}
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
