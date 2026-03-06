import { useTranslation } from "react-i18next";

import { PathIcon } from "@/components/icons";
import { StatusBadge } from "@/components/status-badge";
import type { MigrationRecord } from "@/types/migration-record";
import type { Place } from "@/types/reference";

interface JourneyTimelineProps {
  records: MigrationRecord[];
  places: Place[];
  onEdit?: (id: string) => void;
}

function placeName(id: string | undefined, places: Place[]): string | null {
  if (!id) return null;
  return places.find((p) => p.id === id)?.name ?? null;
}

export function JourneyTimeline({ records, places, onEdit }: JourneyTimelineProps) {
  const { t } = useTranslation();

  const sorted = [...records].sort((a, b) => {
    const da = a.migration_date ?? a.created_at;
    const db = b.migration_date ?? b.created_at;
    return da.localeCompare(db);
  });

  if (sorted.length === 0) return null;

  return (
    <div className="relative">
      {/* vertical line */}
      <div className="absolute top-3 bottom-3 left-4 w-px bg-border-secondary" />

      <ol className="space-y-0">
        {sorted.map((record, i) => {
          const from = placeName(record.from_place_id, places);
          const to = placeName(record.destination_place_id, places);
          const date = record.migration_date
            ? new Date(record.migration_date).toLocaleDateString("en-CA")
            : null;
          const isFirst = i === 0;
          const isLast = i === sorted.length - 1;

          return (
            <li key={record.id} className="relative pl-11">
              {/* node */}
              <div
                className={`absolute left-2 top-3 size-5 rounded-full border-2 ${
                  isLast
                    ? "border-accent bg-accent"
                    : isFirst
                      ? "border-fg-tertiary bg-bg-secondary"
                      : "border-border-primary bg-bg-secondary"
                }`}
              >
                {isLast && (
                  <PathIcon
                    size={11}
                    weight="bold"
                    className="absolute top-0.5 left-0.5 text-accent-fg"
                  />
                )}
              </div>

              <button
                type="button"
                onClick={() => onEdit?.(record.id)}
                className={`group w-full cursor-pointer rounded-lg border border-transparent p-3 text-left transition-colors hover:border-border-secondary hover:bg-bg-secondary ${
                  isLast ? "bg-bg-secondary" : ""
                }`}
              >
                {/* date */}
                {date && (
                  <span className="mb-1 block font-mono text-xs tabular-nums text-fg-tertiary">
                    {date}
                  </span>
                )}

                {/* from → to */}
                <div className="flex items-center gap-2 text-sm">
                  {from ? (
                    <>
                      <span className="font-medium text-fg">{from}</span>
                      <span className="text-fg-tertiary">→</span>
                      <span className="font-medium text-fg">{to ?? "—"}</span>
                    </>
                  ) : to ? (
                    <>
                      <span className="text-fg-tertiary">{t("project.migrationRecords.arrivedAt")}</span>
                      <span className="font-medium text-fg">{to}</span>
                    </>
                  ) : (
                    <span className="text-fg-tertiary">{t("project.migrationRecords.movement")}</span>
                  )}
                </div>

                {/* badges row */}
                <div className="mt-2 flex flex-wrap gap-2">
                  {record.movement_reason && (
                    <StatusBadge label={record.movement_reason} />
                  )}
                  {record.housing_at_destination && (
                    <span className="inline-flex items-center rounded-full bg-bg-tertiary px-2.5 py-0.5 text-xs font-medium text-fg-secondary">
                      {record.housing_at_destination.replaceAll("_", " ")}
                    </span>
                  )}
                </div>

                {/* notes */}
                {record.notes && (
                  <p className="mt-1.5 text-xs text-fg-tertiary line-clamp-2">{record.notes}</p>
                )}
              </button>
            </li>
          );
        })}
      </ol>
    </div>
  );
}
