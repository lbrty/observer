import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { PieChart } from "@/components/charts/pie-chart";
import { SUPPORT_TYPE_COLORS } from "@/components/charts/colors";
import { useDocuments } from "@/hooks/use-documents";
import { useMigrationRecords } from "@/hooks/use-migration-records";
import { useMyProjects } from "@/hooks/use-my-projects";
import { useNotes } from "@/hooks/use-notes";
import { useSupportRecords } from "@/hooks/use-support-records";

import type { CountResult } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/people/$personId/stats")({
  component: PersonStats,
});

const typeKeyMap: Record<string, string> = {
  humanitarian: "project.supportRecords.typeHumanitarian",
  legal: "project.supportRecords.typeLegal",
  social: "project.supportRecords.typeSocial",
  psychological: "project.supportRecords.typePsychological",
  medical: "project.supportRecords.typeMedical",
  general: "project.supportRecords.typeGeneral",
};

function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
  );
}

function PersonStats() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const { data: projectsData } = useMyProjects();
  const project = projectsData?.projects.find((p) => p.id === projectId);
  const canViewDocuments = project?.can_view_documents ?? false;

  const { data: supportData, isLoading: loadingSupport } = useSupportRecords(projectId, {
    person_id: personId,
    per_page: 1000,
  });
  const { data: notesData, isLoading: loadingNotes } = useNotes(projectId, personId);
  const { data: migrationData, isLoading: loadingMigration } = useMigrationRecords(
    projectId,
    personId,
  );
  const { data: documentsData, isLoading: loadingDocs } = useDocuments(
    canViewDocuments ? projectId : "",
    personId,
  );

  const isLoading = loadingSupport || loadingNotes || loadingMigration || (canViewDocuments && loadingDocs);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {Array.from({ length: 4 }, (_, i) => (
            <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
          ))}
        </div>
        <div className="grid gap-6 lg:grid-cols-2">
          {Array.from({ length: 2 }, (_, i) => (
            <div key={i} className="h-60 animate-pulse rounded-xl bg-bg-tertiary" />
          ))}
        </div>
      </div>
    );
  }

  const supportRecords = supportData?.records ?? [];
  const totalSupport = supportData?.total ?? 0;
  const totalNotes = notesData?.notes.length ?? 0;
  const totalMigration = migrationData?.records.length ?? 0;
  const totalDocuments = canViewDocuments ? (documentsData?.documents.length ?? 0) : 0;

  const typeCounts = new Map<string, number>();
  for (const r of supportRecords) {
    typeCounts.set(r.type, (typeCounts.get(r.type) ?? 0) + 1);
  }
  const byType: CountResult[] = Array.from(typeCounts, ([label, count]) => ({
    label: t(typeKeyMap[label] ?? label),
    count,
  }));

  const sphereCounts = new Map<string, number>();
  for (const r of supportRecords) {
    if (r.sphere) {
      sphereCounts.set(r.sphere, (sphereCounts.get(r.sphere) ?? 0) + 1);
    }
  }

  const sphereKeyMap: Record<string, string> = {
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

  const bySphere: CountResult[] = Array.from(sphereCounts, ([label, count]) => ({
    label: t(sphereKeyMap[label] ?? label),
    count,
  }));

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <KpiCard label={t("project.people.statsSupportRecords")} value={totalSupport} />
        <KpiCard label={t("project.people.statsNotes")} value={totalNotes} />
        <KpiCard label={t("project.people.statsMigrationRecords")} value={totalMigration} />
        {canViewDocuments && (
          <KpiCard label={t("project.people.statsDocuments")} value={totalDocuments} />
        )}
      </div>

      {byType.length > 0 && (
        <div className="grid gap-6 lg:grid-cols-2">
          <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
            <h3 className="mb-3 text-sm font-semibold text-fg">
              {t("project.people.statsByType")}
            </h3>
            <PieChart data={byType} colorMap={SUPPORT_TYPE_COLORS} />
          </div>

          {bySphere.length > 0 && (
            <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
              <h3 className="mb-3 text-sm font-semibold text-fg">
                {t("project.people.statsBySphere")}
              </h3>
              <PieChart data={bySphere} />
            </div>
          )}
        </div>
      )}

      {byType.length === 0 && (
        <p className="py-12 text-center text-sm text-fg-tertiary">
          {t("project.people.statsEmpty")}
        </p>
      )}
    </div>
  );
}
