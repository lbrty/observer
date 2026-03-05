import { DrawerPreview as Drawer } from "@base-ui/react/drawer";
import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { CheckIcon, TrashIcon, WarningIcon, XIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { HTTPError } from "@/lib/api";
import {
  useAddHouseholdMember,
  useCreateHousehold,
  useHousehold,
  useRemoveHouseholdMember,
  useUpdateHousehold,
} from "@/hooks/use-households";
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

  const [form, setForm] = useState(emptyForm);
  const [memberForm, setMemberForm] = useState(emptyMemberForm);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setForm(emptyForm);
      setMemberForm(emptyMemberForm);
      setSaved(false);
      setError("");
      setEditingId(null);
      return;
    }
    if (isEdit && household) {
      setForm({
        reference_number: household.reference_number ?? "",
        head_person_id: household.head_person_id ?? "",
      });
      setEditingId(household.id);
    }
  }, [open, isEdit, household]);

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }));
    setSaved(false);
    setError("");
  }

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
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
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
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
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
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
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

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";

  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs transition-opacity duration-200 data-ending-style:opacity-0 data-starting-style:opacity-0" />
        <Drawer.Viewport className="fixed inset-0 z-50">
          <Drawer.Popup className="fixed top-0 right-0 flex h-dvh w-full max-w-[840px] flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full">
            <div className="flex shrink-0 items-center justify-between border-b border-border-secondary px-6 py-4">
              <Drawer.Title className="font-serif text-lg font-semibold text-fg">
                {isEdit ? t("project.households.editTitle") : t("project.households.formTitle")}
              </Drawer.Title>
              <Drawer.Close className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg">
                <XIcon size={18} />
              </Drawer.Close>
            </div>

            <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col">
              <div className="flex-1 space-y-5 overflow-y-auto px-6 py-5">
                {saved && (
                  <div className="flex items-center gap-2 rounded-lg border border-foam/20 bg-foam/8 px-3 py-2.5 text-sm font-medium text-foam">
                    <CheckIcon size={16} weight="bold" className="shrink-0" />
                    {t("project.households.saved")}
                  </div>
                )}
                {error && (
                  <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
                    <WarningIcon size={16} weight="bold" className="shrink-0" />
                    {error}
                  </div>
                )}

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("project.households.info")}
                </h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.households.referenceNumber")}
                    </Field.Label>
                    <Field.Control
                      value={form.reference_number}
                      onChange={(e) => set("reference_number", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.households.headPerson")}
                    </Field.Label>
                    <Field.Control
                      value={form.head_person_id}
                      onChange={(e) => set("head_person_id", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>
                </div>

                {isEdit && editingId && (
                  <>
                    <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                      {t("project.households.members")}
                    </h3>
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
                                  <button
                                    type="button"
                                    onClick={() => handleRemoveMember(m.person_id)}
                                    className="cursor-pointer rounded-lg p-1 text-fg-tertiary hover:bg-bg-tertiary hover:text-rose"
                                  >
                                    <TrashIcon size={14} />
                                  </button>
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
                      <Field.Root className="flex-1">
                        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                          {t("project.households.personId")}
                        </Field.Label>
                        <Field.Control
                          value={memberForm.person_id}
                          onChange={(e) =>
                            setMemberForm((f) => ({
                              ...f,
                              person_id: e.target.value,
                            }))
                          }
                          className={inputClass}
                        />
                      </Field.Root>

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

                      <button
                        type="button"
                        onClick={handleAddMember}
                        disabled={addMember.isPending || !memberForm.person_id}
                        className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
                      >
                        {t("project.households.addMember")}
                      </button>
                    </div>
                  </>
                )}
              </div>

              <div className="flex shrink-0 items-center justify-end gap-2 border-t border-border-secondary px-6 py-4">
                <Drawer.Close className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary shadow-card hover:bg-bg-tertiary">
                  {t("admin.common.cancel")}
                </Drawer.Close>
                <button
                  type="submit"
                  disabled={isPending}
                  className="cursor-pointer rounded-lg bg-accent px-5 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
                >
                  {isPending ? t("project.households.saving") : t("project.households.save")}
                </button>
              </div>
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
