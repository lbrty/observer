import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { StatusBadge } from "@/components/status-badge";
import { usePerson } from "@/hooks/use-people";

export const Route = createFileRoute(
  "/_app/projects/$projectId/people/$personId/",
)({
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

function PersonOverview() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();
  const { data: person, isLoading } = usePerson(projectId, personId);

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
    middle_aged_adult: t("project.people.ageMiddleAged"),
    old_adult: t("project.people.ageOlderAdult"),
  };

  return (
    <section className="space-y-5 rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <h2 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
        {t("project.people.identity")}
      </h2>
      <dl className="grid grid-cols-2 gap-x-8 gap-y-4 sm:grid-cols-3">
        <Detail
          label={t("project.people.firstName")}
          value={person.first_name}
        />
        <Detail
          label={t("project.people.lastName")}
          value={person.last_name}
        />
        <Detail
          label={t("project.people.patronymic")}
          value={person.patronymic}
        />
        <Detail
          label={t("project.people.sexLabel")}
          value={sexLabels[person.sex]}
        />
        <Detail
          label={t("project.people.birthDate")}
          value={person.birth_date}
        />
        <Detail
          label={t("project.people.ageGroup")}
          value={person.age_group ? ageLabels[person.age_group] : undefined}
        />
        <Detail
          label={t("project.people.phone")}
          value={person.primary_phone}
        />
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
        <Detail
          label={t("project.people.externalId")}
          value={person.external_id}
        />
        <Detail
          label={t("project.people.office")}
          value={person.office_id}
        />
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
        <Detail
          label={t("project.people.consentDate")}
          value={person.consent_date}
        />
      </dl>
    </section>
  );
}
