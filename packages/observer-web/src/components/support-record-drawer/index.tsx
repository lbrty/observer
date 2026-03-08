import { type SyntheticEvent, useEffect, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/use-offices";
import { usePerson } from "@/hooks/use-people";
import {
  useCreateSupportRecord,
  useSupportRecord,
  useUpdateSupportRecord,
} from "@/hooks/use-support-records";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";
import { useToast } from "@/stores/toast";

import type { CreateSupportRecordInput, UpdateSupportRecordInput } from "@/types/support-record";

import { InfoSection } from "./info-section";
import { ReferralSection } from "./referral-section";

interface SupportRecordDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  recordId: string | null;
  personId?: string;
}

const emptyForm = {
  type: "general",
  sphere: "",
  provided_at: "",
  person_id: "",
  referral_status: "",
  referred_to_office: "",
  consultant_id: "",
  office_id: "",
  notes: "",
};

export function SupportRecordDrawer({
  open,
  onOpenChange,
  projectId,
  recordId,
  personId,
}: SupportRecordDrawerProps) {
  const { t } = useTranslation();
  const { user } = useAuth();
  const isEdit = recordId !== null;

  const { data: record } = useSupportRecord(projectId, recordId ?? "");
  const qc = useQueryClient();
  const createRecord = useCreateSupportRecord(projectId);
  const updateRecord = useUpdateSupportRecord(projectId);

  const initial = {
    ...emptyForm,
    person_id: personId ?? "",
    consultant_id: user?.id ?? "",
    office_id: user?.office_id ?? "",
  };

  const toast = useToast();
  const { form, set, error, setError, editingId, setEditingId, setForm } =
    useDrawerForm({
      initial,
      open,
      isEdit,
      data: record,
      mapData: (data) => ({
        type: (data.type as string) ?? "general",
        sphere: (data.sphere as string) ?? "",
        provided_at: (data.provided_at as string) ?? "",
        person_id: (data.person_id as string) ?? "",
        referral_status: (data.referral_status as string) ?? "",
        referred_to_office: (data.referred_to_office as string) ?? "",
        consultant_id: (data.consultant_id as string) ?? "",
        office_id: (data.office_id as string) ?? "",
        notes: (data.notes as string) ?? "",
      }),
    });

  const [personName, setPersonName] = useState("");
  const { data: personData } = usePerson(projectId, form.person_id);

  useEffect(() => {
    if (personData) {
      setPersonName(`${personData.first_name} ${personData.last_name ?? ""}`.trim());
    }
  }, [personData]);

  useEffect(() => {
    if (open && !isEdit && personId) {
      setForm((f) => ({ ...f, person_id: personId }));
    }
    if (!open) {
      setPersonName("");
    }
  }, [open, isEdit, personId]);

  const { data: offices } = useOffices();

  const isPending = createRecord.isPending || updateRecord.isPending;

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    try {
      if (isEdit && editingId) {
        const data: UpdateSupportRecordInput = {
          type: form.type as UpdateSupportRecordInput["type"],
          sphere: (form.sphere || undefined) as UpdateSupportRecordInput["sphere"],
          provided_at: form.provided_at || undefined,
          consultant_id: form.consultant_id || undefined,
          office_id: form.office_id || undefined,
          referred_to_office: form.referred_to_office || undefined,
          referral_status: (form.referral_status ||
            undefined) as UpdateSupportRecordInput["referral_status"],
          notes: form.notes || undefined,
        };
        await updateRecord.mutateAsync({ id: editingId, data });
        await qc.invalidateQueries({
          queryKey: ["support-records", projectId],
        });
        toast.success(t("project.supportRecords.saved"));
      } else {
        const input: CreateSupportRecordInput = {
          person_id: form.person_id,
          type: form.type as CreateSupportRecordInput["type"],
          ...(form.sphere && {
            sphere: form.sphere as CreateSupportRecordInput["sphere"],
          }),
          ...(form.provided_at && { provided_at: form.provided_at }),
          ...(form.consultant_id && { consultant_id: form.consultant_id }),
          ...(form.office_id && { office_id: form.office_id }),
          ...(form.referred_to_office && {
            referred_to_office: form.referred_to_office,
          }),
          ...(form.referral_status && {
            referral_status: form.referral_status as CreateSupportRecordInput["referral_status"],
          }),
          ...(form.notes && { notes: form.notes }),
        };
        const created = await createRecord.mutateAsync(input);
        await qc.invalidateQueries({
          queryKey: ["support-records", projectId],
        });
        setEditingId(created.id);
        toast.success(t("project.supportRecords.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const officeOptions = (offices ?? []).map((o) => ({
    label: o.name,
    value: o.id,
  }));

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

      <ReferralSection form={form} set={(k, v) => set(k as keyof typeof form, v)} officeOptions={officeOptions} />

      <SectionHeading>{t("project.supportRecords.notesSection")}</SectionHeading>
      <FormTextarea
        label=""
        value={form.notes}
        onChange={(v) => set("notes", v)}
        rows={4}
      />
    </DrawerShell>
  );
}
