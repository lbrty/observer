import { type SyntheticEvent, useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DatePicker } from "@/components/date-picker";
import { PlusIcon, WarningIcon } from "@/components/icons";
import { StatusBadge } from "@/components/status-badge";
import { UISelect } from "@/components/ui-select";
import { useOffices } from "@/hooks/use-offices";
import { usePerson } from "@/hooks/use-people";
import { useCreateSupportRecord } from "@/hooks/use-support-records";
import { HTTPError } from "@/lib/api";
import { useToast } from "@/stores/toast";

import type { SupportSphere, SupportType } from "@/types/support-record";

export const Route = createFileRoute("/_app/projects/$projectId/people/$personId/")({
  component: PersonOverview,
});

function Detail({ label, value }: { label: string; value?: string | null }) {
  if (!value) return null;
  return (
    <div>
      <dt className="text-xs font-medium text-fg-tertiary">{label}</dt>
      <dd className="mt-0.5 text-sm text-fg">{value}</dd>
    </div>
  );
}

const inputClass =
  "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg";

const typeOptions: SupportType[] = [
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
];

const sphereOptions: SupportSphere[] = [
  "housing_assistance",
  "document_recovery",
  "social_benefits",
  "property_rights",
  "employment_rights",
  "family_law",
  "healthcare_access",
  "education_access",
  "financial_aid",
  "psychological_support",
  "other",
];

const typeKeyMap: Record<SupportType, string> = {
  humanitarian: "project.supportRecords.typeHumanitarian",
  legal: "project.supportRecords.typeLegal",
  social: "project.supportRecords.typeSocial",
  psychological: "project.supportRecords.typePsychological",
  medical: "project.supportRecords.typeMedical",
  general: "project.supportRecords.typeGeneral",
};

const sphereKeyMap: Record<SupportSphere, string> = {
  housing_assistance: "project.supportRecords.sphereHousing",
  document_recovery: "project.supportRecords.sphereDocumentRecovery",
  social_benefits: "project.supportRecords.sphereSocialBenefits",
  property_rights: "project.supportRecords.spherePropertyRights",
  employment_rights: "project.supportRecords.sphereEmploymentRights",
  family_law: "project.supportRecords.sphereFamilyLaw",
  healthcare_access: "project.supportRecords.sphereHealthcareAccess",
  education_access: "project.supportRecords.sphereEducationAccess",
  financial_aid: "project.supportRecords.sphereFinancialAid",
  psychological_support: "project.supportRecords.spherePsychologicalSupport",
  other: "project.supportRecords.sphereOther",
};

function PersonOverview() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();
  const { data: person, isLoading } = usePerson(projectId, personId);
  const { data: offices } = useOffices();

  const officeName = person?.office_id
    ? offices?.find((o) => o.id === person.office_id)?.name
    : undefined;

  const createRecord = useCreateSupportRecord(projectId);
  const toast = useToast();

  const [formOpen, setFormOpen] = useState(false);
  const [error, setError] = useState("");

  const [type, setType] = useState<SupportType>("humanitarian");
  const [sphere, setSphere] = useState("");
  const [providedAt, setProvidedAt] = useState(new Date().toISOString().slice(0, 10));
  const [notes, setNotes] = useState("");

  function resetForm() {
    setType("humanitarian");
    setSphere("");
    setProvidedAt(new Date().toISOString().slice(0, 10));
    setNotes("");
    setError("");
  }

  function handleCancel() {
    setFormOpen(false);
    resetForm();
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    try {
      await createRecord.mutateAsync({
        person_id: personId,
        type,
        sphere: sphere ? (sphere as SupportSphere) : undefined,
        provided_at: providedAt || undefined,
        notes: notes || undefined,
      });

      toast.success(t("project.people.quickSupportSaved"));
      setFormOpen(false);
      resetForm();
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setError((body as { error?: string } | null)?.error ?? err.message);
      } else {
        setError(String(err));
      }
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        {Array.from({ length: 3 }, (_, i) => (
          <div key={i} className="h-32 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
    );
  }

  if (!person) return null;

  const sexLabels: Record<string, string> = {
    male: t("project.people.sexMale"),
    female: t("project.people.sexFemale"),
    other: t("project.people.sexOther"),
    unknown: t("project.people.sexUnknown"),
  };

  const ageLabels: Record<string, string> = {
    infant: t("project.people.ageInfant"),
    toddler: t("project.people.ageToddler"),
    pre_school: t("project.people.agePreSchool"),
    middle_childhood: t("project.people.ageMiddleChildhood"),
    young_teen: t("project.people.ageYoungTeen"),
    teenager: t("project.people.ageTeenager"),
    young_adult: t("project.people.ageYoungAdult"),
    early_adult: t("project.people.ageEarlyAdult"),
    middle_aged_adult: t("project.people.ageMiddleAgedAdult"),
    old_adult: t("project.people.ageOldAdult"),
  };

  return (
    <div className="space-y-4">
      <section className="space-y-5 rounded-xl border border-border-secondary bg-bg-secondary p-5">
        <h2 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
          {t("project.people.identity")}
        </h2>
        <dl className="grid grid-cols-2 gap-x-8 gap-y-4 sm:grid-cols-3">
          <Detail label={t("project.people.firstName")} value={person.first_name} />
          <Detail label={t("project.people.lastName")} value={person.last_name} />
          <Detail label={t("project.people.patronymic")} value={person.patronymic} />
          <Detail label={t("project.people.sexLabel")} value={sexLabels[person.sex]} />
          <Detail label={t("project.people.birthDate")} value={person.birth_date} />
          <Detail
            label={t("project.people.ageGroup")}
            value={person.age_group ? ageLabels[person.age_group] : undefined}
          />
          <Detail label={t("project.people.phone")} value={person.primary_phone} />
          <Detail label={t("project.people.email")} value={person.email} />
        </dl>

        <hr className="border-border-secondary" />

        <h2 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
          {t("project.people.case")}
        </h2>
        <dl className="grid grid-cols-2 gap-x-8 gap-y-4 sm:grid-cols-3">
          <div>
            <dt className="text-xs font-medium text-fg-tertiary">
              {t("project.people.caseStatusLabel")}
            </dt>
            <dd className="mt-1">
              <StatusBadge label={person.case_status} />
            </dd>
          </div>
          <Detail label={t("project.people.externalId")} value={person.external_id} />
          <Detail label={t("project.people.office")} value={officeName} />
          <div>
            <dt className="text-xs font-medium text-fg-tertiary">
              {t("project.people.consentGiven")}
            </dt>
            <dd className="mt-0.5 text-sm text-fg">
              <StatusBadge
                label={person.consent_given ? t("admin.users.yes") : t("admin.users.no")}
                variant={person.consent_given ? "foam" : "neutral"}
                dot={false}
              />
            </dd>
          </div>
          <Detail label={t("project.people.consentDate")} value={person.consent_date} />
        </dl>
      </section>

      {!formOpen && (
        <Button
          variant="ghost"
          icon={<PlusIcon size={16} weight="bold" />}
          onClick={() => setFormOpen(true)}
        >
          {t("project.people.quickSupportAdd")}
        </Button>
      )}

      {formOpen && (
        <section className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
                <WarningIcon size={16} weight="bold" className="shrink-0" />
                {error}
              </div>
            )}

            <div className="flex items-center gap-3 pt-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                {t("project.supportRecords.recordInfo")}
              </span>
              <span className="h-px flex-1 bg-border-secondary" />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
                  {t("project.supportRecords.type")} *
                </label>
                <UISelect
                  value={type}
                  onValueChange={(v) => setType(v as SupportType)}
                  options={typeOptions.map((v) => ({
                    value: v,
                    label: t(typeKeyMap[v]),
                  }))}
                  fullWidth
                />
              </div>
              <div>
                <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
                  {t("project.supportRecords.sphere")}
                </label>
                <UISelect
                  value={sphere}
                  onValueChange={setSphere}
                  options={sphereOptions.map((v) => ({
                    value: v,
                    label: t(sphereKeyMap[v]),
                  }))}
                  placeholder={t("project.supportRecords.selectSphere")}
                  fullWidth
                />
              </div>
            </div>

            <div className="flex items-center gap-3 pt-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
                {t("project.supportRecords.notesSection")}
              </span>
              <span className="h-px flex-1 bg-border-secondary" />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
                  {t("project.supportRecords.providedAt")}
                </label>
                <DatePicker value={providedAt} onChange={setProvidedAt} />
              </div>
            </div>

            <div>
              <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
                {t("project.supportRecords.notes")}
              </label>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                rows={2}
                className={`${inputClass} h-auto py-2`}
              />
            </div>

            <div className="flex justify-end gap-2 pt-1">
              <Button variant="secondary" onClick={handleCancel}>
                {t("common.cancel")}
              </Button>
              <Button type="submit" disabled={createRecord.isPending}>
                {createRecord.isPending
                  ? t("project.supportRecords.saving")
                  : t("project.supportRecords.save")}
              </Button>
            </div>
          </form>
        </section>
      )}
    </div>
  );
}
