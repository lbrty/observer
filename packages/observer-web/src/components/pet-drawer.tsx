import { DrawerPreview as Drawer } from "@base-ui/react/drawer";
import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { CheckIcon, WarningIcon, XIcon } from "@/components/icons";
import { HTTPError } from "@/lib/api";
import { UISelect } from "@/components/ui-select";
import { useCreatePet, usePet, useUpdatePet } from "@/hooks/use-pets";
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

export function PetDrawer({
  open,
  onOpenChange,
  projectId,
  petId,
}: PetDrawerProps) {
  const { t } = useTranslation();
  const isEdit = petId !== null;

  const { data: pet } = usePet(projectId, petId ?? "");
  const qc = useQueryClient();
  const createPet = useCreatePet(projectId);
  const updatePet = useUpdatePet(projectId);

  const [form, setForm] = useState(emptyForm);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setForm(emptyForm);
      setSaved(false);
      setError("");
      setEditingId(null);
      return;
    }
    if (isEdit && pet) {
      setForm({
        name: pet.name,
        status: pet.status,
        owner_id: pet.owner_id ?? "",
        registration_id: pet.registration_id ?? "",
        notes: pet.notes ?? "",
      });
      setEditingId(pet.id);
    }
  }, [open, isEdit, pet]);

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }));
    setSaved(false);
    setError("");
  }

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
        setSaved(true);
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

  const statusOptions = [
    { label: t("project.pets.statusRegistered"), value: "registered" },
    { label: t("project.pets.statusAdopted"), value: "adopted" },
    { label: t("project.pets.statusOwnerFound"), value: "owner_found" },
    { label: t("project.pets.statusNeedsShelter"), value: "needs_shelter" },
    { label: t("project.pets.statusUnknown"), value: "unknown" },
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
                {isEdit
                  ? t("project.pets.editTitle")
                  : t("project.pets.formTitle")}
              </Drawer.Title>
              <Drawer.Close className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg">
                <XIcon size={18} />
              </Drawer.Close>
            </div>

            <form
              onSubmit={handleSubmit}
              className="flex min-h-0 flex-1 flex-col"
            >
              <div className="flex-1 space-y-6 overflow-y-auto px-6 py-5">
                {saved && (
                  <div className="flex items-center gap-2 rounded-lg border border-foam/20 bg-foam/8 px-3 py-2.5 text-sm font-medium text-foam">
                    <CheckIcon size={16} weight="bold" className="shrink-0" />
                    {t("project.pets.saved")}
                  </div>
                )}
                {error && (
                  <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
                    <WarningIcon size={16} weight="bold" className="shrink-0" />
                    {error}
                  </div>
                )}

                <Section title={t("project.pets.title")}>
                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.pets.name")} *
                      </Field.Label>
                      <Field.Control
                        required
                        value={form.name}
                        onChange={(e) => set("name", e.target.value)}
                        className={inputClass}
                      />
                    </Field.Root>

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

                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.pets.ownerId")}
                      </Field.Label>
                      <Field.Control
                        value={form.owner_id}
                        onChange={(e) => set("owner_id", e.target.value)}
                        className={inputClass}
                      />
                    </Field.Root>

                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.pets.registrationId")}
                      </Field.Label>
                      <Field.Control
                        value={form.registration_id}
                        onChange={(e) =>
                          set("registration_id", e.target.value)
                        }
                        className={inputClass}
                      />
                    </Field.Root>
                  </div>

                  <div className="mt-4">
                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.pets.notes")}
                      </Field.Label>
                      <textarea
                        value={form.notes}
                        onChange={(e) => set("notes", e.target.value)}
                        rows={4}
                        className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
                      />
                    </Field.Root>
                  </div>
                </Section>
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
                  {isPending
                    ? t("project.pets.saving")
                    : t("project.pets.save")}
                </button>
              </div>
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}

function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <fieldset className="space-y-3">
      <legend className="text-sm font-semibold text-fg">{title}</legend>
      <div className="rounded-xl border border-border-secondary bg-bg p-4">
        {children}
      </div>
    </fieldset>
  );
}
