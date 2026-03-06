import { type FormEvent, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner, SuccessBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { DrawerShell } from "@/components/drawer-shell";
import { FormField } from "@/components/form-field";
import { TrashIcon } from "@/components/icons";
import { PersonCombobox } from "@/components/person-combobox";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import {
  useAddHouseholdMember,
  useCreateHousehold,
  useHousehold,
  useRemoveHouseholdMember,
  useUpdateHousehold,
} from "@/hooks/use-households";
import { handleApiError } from "@/lib/form-error";

import type { Relationship } from "@/types/household";

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

  const { form, set, saved, setSaved, error, setError, editingId, setEditingId } =
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

  const headPersonLabel = headPersonName || form.head_person_id;
  const memberPersonLabel = memberPersonName || memberForm.person_id;

  const isPending = createHousehold.isPending || updateHousehold.isPending;

  async function handleSubmit(e: FormEvent) {
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
        setSaved(true);
      } else {
        const created = await createHousehold.mutateAsync({
          reference_number: form.reference_number || undefined,
          head_person_id: form.head_person_id || undefined,
        });
        await qc.invalidateQueries({ queryKey: ["households", projectId] });
        setEditingId(created.id);
        setSaved(true);
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

  const relationshipOptions = [
    { label: t("project.households.relationshipHead"), value: "head" },
    { label: t("project.households.relationshipSpouse"), value: "spouse" },
    { label: t("project.households.relationshipChild"), value: "child" },
    { label: t("project.households.relationshipParent"), value: "parent" },
    { label: t("project.households.relationshipSibling"), value: "sibling" },
    {
      label: t("project.households.relationshipGrandchild"),
      value: "grandchild",
    },
    {
      label: t("project.households.relationshipGrandparent"),
      value: "grandparent",
    },
    {
      label: t("project.households.relationshipOtherRelative"),
      value: "other_relative",
    },
    {
      label: t("project.households.relationshipNonRelative"),
      value: "non_relative",
    },
  ];

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
      {saved && <SuccessBanner message={t("project.households.saved")} />}
      {error && <ErrorBanner message={error} />}

      <SectionHeading>{t("project.households.info")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField
          label={t("project.households.referenceNumber")}
          value={form.reference_number}
          onChange={(v) => set("reference_number", v)}
        />
        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.households.headPerson")}
          </span>
          {form.head_person_id ? (
            <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
              <span className="flex-1 truncate font-mono text-xs text-fg">
                {headPersonLabel}
              </span>
              <button
                type="button"
                onClick={() => {
                  set("head_person_id", "");
                  setHeadPersonName("");
                }}
                className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
              >
                ×
              </button>
            </div>
          ) : (
            <PersonCombobox
              projectId={projectId}
              onSelect={(p) => {
                set("head_person_id", p.id);
                setHeadPersonName(`${p.first_name} ${p.last_name ?? ""}`.trim());
              }}
            />
          )}
        </div>
      </div>

      {isEdit && editingId && (
        <>
          <SectionHeading>{t("project.households.members")}</SectionHeading>
          {household?.members && household.members.length > 0 ? (
            <div className="overflow-hidden rounded-lg border border-border-secondary">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border-secondary bg-bg-secondary">
                    <th className="px-3 py-2 text-left font-medium text-fg-secondary">
                      {t("project.households.personId")}
                    </th>
                    <th className="px-3 py-2 text-left font-medium text-fg-secondary">
                      {t("project.households.relationship")}
                    </th>
                    <th className="w-10 px-3 py-2" />
                  </tr>
                </thead>
                <tbody>
                  {household.members.map((m) => (
                    <tr
                      key={m.person_id}
                      className="border-b border-border-secondary last:border-b-0"
                    >
                      <td className="px-3 py-2 font-mono text-xs text-fg">
                        {m.person_id}
                      </td>
                      <td className="px-3 py-2 text-fg-secondary">
                        {t(
                          `project.households.relationship${m.relationship.charAt(0).toUpperCase()}${m.relationship.slice(1).replace(/_(\w)/g, (_, c: string) => c.toUpperCase())}`,
                        )}
                      </td>
                      <td className="px-3 py-2">
                        <Button
                          variant="ghost"
                          className="p-1 hover:text-rose"
                          onClick={() => handleRemoveMember(m.person_id)}
                        >
                          <TrashIcon size={14} />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-sm text-fg-tertiary">
              {t("project.households.noMembers")}
            </p>
          )}

          <p className="text-xs font-medium text-fg-tertiary">
            {t("project.households.addMember")}
          </p>
          <div className="flex items-end gap-3">
            <div className="flex-1">
              <label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.households.personId")}
              </label>
              {memberForm.person_id ? (
                <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
                  <span className="flex-1 truncate font-mono text-xs text-fg">
                    {memberPersonLabel}
                  </span>
                  <button
                    type="button"
                    onClick={() => {
                      setMemberForm(emptyMemberForm);
                      setMemberPersonName("");
                    }}
                    className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
                  >
                    ×
                  </button>
                </div>
              ) : (
                <PersonCombobox
                  projectId={projectId}
                  excludeIds={household?.members?.map((m) => m.person_id) ?? []}
                  onSelect={(p) => {
                    setMemberForm((f) => ({ ...f, person_id: p.id }));
                    setMemberPersonName(`${p.first_name} ${p.last_name ?? ""}`.trim());
                  }}
                />
              )}
            </div>

            <div className="flex-1">
              <label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.households.relationship")}
              </label>
              <UISelect
                value={memberForm.relationship}
                onValueChange={(v) =>
                  setMemberForm((f) => ({
                    ...f,
                    relationship: v as Relationship,
                  }))
                }
                options={relationshipOptions}
                fullWidth
              />
            </div>

            <Button
              onClick={() => {
                handleAddMember();
                setMemberPersonName("");
              }}
              disabled={addMember.isPending || !memberForm.person_id}
            >
              {t("project.households.addMember")}
            </Button>
          </div>
        </>
      )}
    </DrawerShell>
  );
}
