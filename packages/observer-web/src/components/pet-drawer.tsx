import { type FormEvent, useState } from "react";
import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { FormField, FormTextarea } from "@/components/form-field";
import { PersonCombobox } from "@/components/person-combobox";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useCreatePet, usePet, useUpdatePet } from "@/hooks/use-pets";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";
import type { CreatePetInput, UpdatePetInput } from "@/types/pet";

interface PetDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  petId: string | null;
}

const emptyForm = {
  name: "",
  status: "registered",
  owner_id: "",
  registration_id: "",
  notes: "",
};

export function PetDrawer({ open, onOpenChange, projectId, petId }: PetDrawerProps) {
  const { t } = useTranslation();
  const isEdit = petId !== null;

  const { data: pet } = usePet(projectId, petId ?? "");
  const qc = useQueryClient();
  const createPet = useCreatePet(projectId);
  const updatePet = useUpdatePet(projectId);

  const [ownerName, setOwnerName] = useState("");

  const toast = useToast();
  const { form, set, error, setError, editingId, setEditingId } = useDrawerForm({
    initial: emptyForm,
    open,
    isEdit,
    data: pet,
    mapData: (d) => ({
      name: (d.name as string) ?? "",
      status: (d.status as string) ?? "registered",
      owner_id: (d.owner_id as string) ?? "",
      registration_id: (d.registration_id as string) ?? "",
      notes: (d.notes as string) ?? "",
    }),
  });

  const isPending = createPet.isPending || updatePet.isPending;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    try {
      if (isEdit && editingId) {
        const data: UpdatePetInput = {
          name: form.name || undefined,
          status: form.status as UpdatePetInput["status"],
          owner_id: form.owner_id || undefined,
          registration_id: form.registration_id || undefined,
          notes: form.notes || undefined,
        };
        await updatePet.mutateAsync({ id: editingId, data });
        await qc.invalidateQueries({ queryKey: ["pets", projectId] });
        toast.success(t("project.pets.saved"));
      } else {
        const input: CreatePetInput = {
          name: form.name,
          ...(form.status && {
            status: form.status as CreatePetInput["status"],
          }),
          ...(form.owner_id && { owner_id: form.owner_id }),
          ...(form.registration_id && {
            registration_id: form.registration_id,
          }),
          ...(form.notes && { notes: form.notes }),
        };
        const created = await createPet.mutateAsync(input);
        await qc.invalidateQueries({ queryKey: ["pets", projectId] });
        setEditingId(created.id);
        toast.success(t("project.pets.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const statusOptions = [
    { label: t("project.pets.statusRegistered"), value: "registered" },
    { label: t("project.pets.statusAdopted"), value: "adopted" },
    { label: t("project.pets.statusOwnerFound"), value: "owner_found" },
    { label: t("project.pets.statusNeedsShelter"), value: "needs_shelter" },
    { label: t("project.pets.statusUnknown"), value: "unknown" },
  ];

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? t("project.pets.editTitle") : t("project.pets.formTitle")}
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.pets.save")}
      savingLabel={t("project.pets.saving")}
    >
      <ErrorBanner message={error} />

      <fieldset className="space-y-3">
        <SectionHeading>{t("project.pets.title")}</SectionHeading>
        <div className="rounded-xl border border-border-secondary bg-bg p-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={t("project.pets.name")}
              value={form.name}
              onChange={(v) => set("name", v)}
              required
            />

            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.pets.status")}
              </Field.Label>
              <UISelect
                value={form.status}
                onValueChange={(v) => set("status", v)}
                options={statusOptions}
                fullWidth
              />
            </Field.Root>

            <div>
              <span className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.pets.ownerId")}
              </span>
              {form.owner_id ? (
                <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
                  <span className="flex-1 truncate text-sm text-fg">
                    {ownerName || form.owner_id}
                  </span>
                  <button
                    type="button"
                    onClick={() => {
                      set("owner_id", "");
                      setOwnerName("");
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
                    set("owner_id", p.id);
                    setOwnerName(`${p.first_name} ${p.last_name ?? ""}`.trim());
                  }}
                />
              )}
            </div>

            <FormField
              label={t("project.pets.registrationId")}
              value={form.registration_id}
              onChange={(v) => set("registration_id", v)}
            />
          </div>

          <div className="mt-4">
            <FormTextarea
              label={t("project.pets.notes")}
              value={form.notes}
              onChange={(v) => set("notes", v)}
            />
          </div>
        </div>
      </fieldset>
    </DrawerShell>
  );
}
