import { type SyntheticEvent, useEffect, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { SectionHeading } from "@/components/section-heading";
import { TagPicker } from "@/components/tag-picker";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useCreatePet, usePet, useUpdatePet } from "@/hooks/use-pets";
import { usePetTags, useReplacePetTags } from "@/hooks/use-tags";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";

import type { CreatePetInput, UpdatePetInput } from "@/types/pet";

import { InfoSection } from "./info-section";

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
  const { data: petTagsData } = usePetTags(projectId, petId ?? "");
  const qc = useQueryClient();
  const createPet = useCreatePet(projectId);
  const updatePet = useUpdatePet(projectId);
  const replacePetTags = useReplacePetTags(projectId);

  const [ownerName, setOwnerName] = useState("");
  const [tagIds, setTagIds] = useState<string[]>([]);

  useEffect(() => {
    if (petTagsData) setTagIds(petTagsData.tag_ids ?? []);
  }, [petTagsData]);

  useEffect(() => {
    if (!open) setTagIds([]);
  }, [open]);

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

  const isPending = createPet.isPending || updatePet.isPending || replacePetTags.isPending;

  async function handleSubmit(e: SyntheticEvent) {
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
        await replacePetTags.mutateAsync({ petId: editingId, ids: tagIds });
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
        if (tagIds.length > 0) {
          await replacePetTags.mutateAsync({ petId: created.id, ids: tagIds });
        }
        await qc.invalidateQueries({ queryKey: ["pets", projectId] });
        setEditingId(created.id);
        toast.success(t("project.pets.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

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

      <InfoSection
        name={form.name}
        status={form.status}
        ownerId={form.owner_id}
        ownerName={ownerName}
        registrationId={form.registration_id}
        notes={form.notes}
        projectId={projectId}
        onNameChange={(v) => set("name", v)}
        onStatusChange={(v) => set("status", v)}
        onOwnerSelect={(id, name) => {
          set("owner_id", id);
          setOwnerName(name);
        }}
        onOwnerClear={() => {
          set("owner_id", "");
          setOwnerName("");
        }}
        onRegistrationIdChange={(v) => set("registration_id", v)}
        onNotesChange={(v) => set("notes", v)}
      />

      <SectionHeading>{t("project.tags.title")}</SectionHeading>
      <TagPicker projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
    </DrawerShell>
  );
}
