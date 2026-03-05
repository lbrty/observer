import { type FormEvent, useEffect } from "react";

import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner, SuccessBanner } from "@/components/alert-banner";
import { DatePicker } from "@/components/date-picker";
import { DrawerShell } from "@/components/drawer-shell";
import { FormField, FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/use-offices";
import {
  useCreateSupportRecord,
  useSupportRecord,
  useUpdateSupportRecord,
} from "@/hooks/use-support-records";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";

import type { CreateSupportRecordInput, UpdateSupportRecordInput } from "@/types/support-record";

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

  const { form, set, saved, setSaved, error, setError, editingId, setEditingId, setForm } =
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

  useEffect(() => {
    if (open && !isEdit && personId) {
      setForm((f) => ({ ...f, person_id: personId }));
    }
  }, [open, isEdit, personId]);

  const { data: offices } = useOffices();

  const isPending = createRecord.isPending || updateRecord.isPending;

  async function handleSubmit(e: FormEvent) {
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
        setSaved(true);
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
        setSaved(true);
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const officeOptions = (offices ?? []).map((o) => ({
    label: o.name,
    value: o.id,
  }));

  const typeOptions = [
    {
      label: t("project.supportRecords.typeHumanitarian"),
      value: "humanitarian",
    },
    { label: t("project.supportRecords.typeLegal"), value: "legal" },
    { label: t("project.supportRecords.typeSocial"), value: "social" },
    {
      label: t("project.supportRecords.typePsychological"),
      value: "psychological",
    },
    { label: t("project.supportRecords.typeMedical"), value: "medical" },
    { label: t("project.supportRecords.typeGeneral"), value: "general" },
  ];

  const sphereOptions = [
    {
      label: t("project.supportRecords.sphereHousing"),
      value: "housing_assistance",
    },
    {
      label: t("project.supportRecords.sphereDocumentRecovery"),
      value: "document_recovery",
    },
    {
      label: t("project.supportRecords.sphereSocialBenefits"),
      value: "social_benefits",
    },
    {
      label: t("project.supportRecords.spherePropertyRights"),
      value: "property_rights",
    },
    {
      label: t("project.supportRecords.sphereEmploymentRights"),
      value: "employment_rights",
    },
    {
      label: t("project.supportRecords.sphereFamilyLaw"),
      value: "family_law",
    },
    {
      label: t("project.supportRecords.sphereHealthcareAccess"),
      value: "healthcare_access",
    },
    {
      label: t("project.supportRecords.sphereEducationAccess"),
      value: "education_access",
    },
    {
      label: t("project.supportRecords.sphereFinancialAid"),
      value: "financial_aid",
    },
    {
      label: t("project.supportRecords.spherePsychologicalSupport"),
      value: "psychological_support",
    },
    { label: t("project.supportRecords.sphereOther"), value: "other" },
  ];

  const referralStatusOptions = [
    {
      label: t("project.supportRecords.referralPending"),
      value: "pending",
    },
    {
      label: t("project.supportRecords.referralAccepted"),
      value: "accepted",
    },
    {
      label: t("project.supportRecords.referralCompleted"),
      value: "completed",
    },
    {
      label: t("project.supportRecords.referralDeclined"),
      value: "declined",
    },
    {
      label: t("project.supportRecords.referralNoResponse"),
      value: "no_response",
    },
  ];

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={
        isEdit
          ? t("project.supportRecords.editTitle")
          : t("project.supportRecords.formTitle")
      }
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.supportRecords.save")}
      savingLabel={t("project.supportRecords.saving")}
    >
      {saved && <SuccessBanner message={t("project.supportRecords.saved")} />}
      <ErrorBanner message={error} />

      <SectionHeading>{t("project.supportRecords.recordInfo")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.type")} *
          </Field.Label>
          <UISelect
            value={form.type}
            onValueChange={(v) => set("type", v)}
            options={typeOptions}
            fullWidth
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.sphere")}
          </Field.Label>
          <UISelect
            value={form.sphere}
            onValueChange={(v) => set("sphere", v)}
            options={sphereOptions}
            placeholder={t("project.supportRecords.selectSphere")}
            fullWidth
          />
        </Field.Root>

        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.providedAt")}
          </span>
          <DatePicker value={form.provided_at} onChange={(v) => set("provided_at", v)} />
        </div>

        {!personId && (
          <FormField
            label={t("project.supportRecords.personId")}
            required={!isEdit}
            value={form.person_id}
            onChange={(v) => set("person_id", v)}
            disabled={isEdit}
          />
        )}
      </div>

      <SectionHeading>{t("project.supportRecords.referral")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.supportRecords.referralStatus")}
          </Field.Label>
          <UISelect
            value={form.referral_status}
            onValueChange={(v) => set("referral_status", v)}
            options={referralStatusOptions}
            placeholder={t("project.supportRecords.selectReferralStatus")}
            fullWidth
          />
        </Field.Root>

        {officeOptions.length > 0 && (
          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.supportRecords.referredToOffice")}
            </Field.Label>
            <UISelect
              value={form.referred_to_office}
              onValueChange={(v) => set("referred_to_office", v)}
              options={officeOptions}
              placeholder={t("project.supportRecords.selectOffice")}
              fullWidth
            />
          </Field.Root>
        )}

      </div>

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
