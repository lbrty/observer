import { DrawerPreview as Drawer } from "@base-ui/react/drawer";
import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { useAuth } from "@/stores/auth";

import { DatePicker } from "@/components/date-picker";
import { CheckIcon, WarningIcon, XIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { HTTPError } from "@/lib/api";
import { useOffices } from "@/hooks/use-offices";
import {
  useCreateSupportRecord,
  useSupportRecord,
  useUpdateSupportRecord,
} from "@/hooks/use-support-records";
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

  const [form, setForm] = useState(emptyForm);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setForm({
        ...emptyForm,
        person_id: personId ?? "",
        consultant_id: user?.id ?? "",
        office_id: user?.office_id ?? "",
      });
      setSaved(false);
      setError("");
      setEditingId(null);
      return;
    }
    if (isEdit && record) {
      setForm({
        type: record.type,
        sphere: record.sphere ?? "",
        provided_at: record.provided_at ?? "",
        person_id: record.person_id,
        referral_status: record.referral_status ?? "",
        referred_to_office: record.referred_to_office ?? "",
        consultant_id: record.consultant_id ?? "",
        office_id: record.office_id ?? "",
        notes: record.notes ?? "",
      });
      setEditingId(record.id);
    } else if (personId) {
      setForm((f) => ({ ...f, person_id: personId }));
    }
  }, [open, isEdit, record]);

  const { data: offices } = useOffices();

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }));
    setSaved(false);
    setError("");
  }

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

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";
  const textareaClass =
    "block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent";

  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs transition-opacity duration-200 data-ending-style:opacity-0 data-starting-style:opacity-0" />
        <Drawer.Viewport className="fixed inset-0 z-50">
          <Drawer.Popup className="fixed top-0 right-0 flex h-dvh w-full max-w-[840px] flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full">
            <div className="flex shrink-0 items-center justify-between border-b border-border-secondary px-6 py-4">
              <Drawer.Title className="font-serif text-lg font-semibold text-fg">
                {isEdit
                  ? t("project.supportRecords.editTitle")
                  : t("project.supportRecords.formTitle")}
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
                    {t("project.supportRecords.saved")}
                  </div>
                )}
                {error && (
                  <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
                    <WarningIcon size={16} weight="bold" className="shrink-0" />
                    {error}
                  </div>
                )}

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("project.supportRecords.recordInfo")}
                </h3>
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
                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.supportRecords.personId")} *
                      </Field.Label>
                      <Field.Control
                        required={!isEdit}
                        value={form.person_id}
                        onChange={(e) => set("person_id", e.target.value)}
                        disabled={isEdit}
                        className={inputClass}
                      />
                    </Field.Root>
                  )}
                </div>

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("project.supportRecords.referral")}
                </h3>
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

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.supportRecords.consultantId")}
                    </Field.Label>
                    <div className={`${inputClass} flex items-center text-fg-secondary`}>
                      {user ? `${user.first_name} ${user.last_name}`.trim() || user.email : ""}
                    </div>
                    <input type="hidden" name="consultant_id" value={form.consultant_id} />
                  </Field.Root>

                  {officeOptions.length > 0 && (
                    <Field.Root>
                      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                        {t("project.supportRecords.office")}
                      </Field.Label>
                      <UISelect
                        value={form.office_id}
                        onValueChange={(v) => set("office_id", v)}
                        options={officeOptions}
                        placeholder={t("project.supportRecords.selectOffice")}
                        fullWidth
                      />
                    </Field.Root>
                  )}
                </div>

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("project.supportRecords.notesSection")}
                </h3>
                <textarea
                  value={form.notes}
                  onChange={(e) => set("notes", e.target.value)}
                  rows={4}
                  className={textareaClass}
                />
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
                    ? t("project.supportRecords.saving")
                    : t("project.supportRecords.save")}
                </button>
              </div>
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
