import { DatePicker } from "@/components/date-picker";
import { PlusIcon } from "@/components/icons";
import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { UISelect } from "@/components/ui-select";
import {
  useCreateMigrationRecord,
  useMigrationRecords,
} from "@/hooks/use-migration-records";
import { useCountries } from "@/hooks/use-countries";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";
import type { MigrationRecord } from "@/types/migration-record";
import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute(
  "/_app/projects/$projectId/people/$personId/migration-records",
)({
  component: PersonMigrationRecords,
});

const emptyForm = {
  from_country: "",
  from_state: "",
  from_place_id: "",
  dest_country: "",
  dest_state: "",
  destination_place_id: "",
  migration_date: "",
  movement_reason: "",
  housing_at_destination: "",
  notes: "",
};

function PersonMigrationRecords() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();
  const { data, isLoading } = useMigrationRecords(projectId, personId);
  const createMutation = useCreateMigrationRecord(projectId, personId);

  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState(emptyForm);

  const { data: countries } = useCountries();
  const { data: fromStates } = useStates(form.from_country || undefined);
  const { data: fromPlaces } = usePlaces(form.from_state || undefined);
  const { data: destStates } = useStates(form.dest_country || undefined);
  const { data: destPlaces } = usePlaces(form.dest_state || undefined);

  function set<K extends keyof typeof form>(key: K, value: string) {
    setForm((f) => ({ ...f, [key]: value }));
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    await createMutation.mutateAsync({
      ...(form.from_place_id && { from_place_id: form.from_place_id }),
      ...(form.destination_place_id && {
        destination_place_id: form.destination_place_id,
      }),
      ...(form.migration_date && { migration_date: form.migration_date }),
      ...(form.movement_reason && {
        movement_reason:
          form.movement_reason as MigrationRecord["movement_reason"],
      }),
      ...(form.housing_at_destination && {
        housing_at_destination:
          form.housing_at_destination as MigrationRecord["housing_at_destination"],
      }),
      ...(form.notes && { notes: form.notes }),
    });
    setForm(emptyForm);
    setDialogOpen(false);
  }

  const countryOptions = (countries ?? []).map((c) => ({
    label: c.name,
    value: c.id,
  }));
  const fromStateOptions = (fromStates?.states ?? []).map((s) => ({
    label: s.name,
    value: s.id,
  }));
  const fromPlaceOptions = (fromPlaces?.places ?? []).map((p) => ({
    label: p.name,
    value: p.id,
  }));
  const destStateOptions = (destStates?.states ?? []).map((s) => ({
    label: s.name,
    value: s.id,
  }));
  const destPlaceOptions = (destPlaces?.places ?? []).map((p) => ({
    label: p.name,
    value: p.id,
  }));

  const reasonOptions = [
    { label: t("project.migrationRecords.reasonConflict"), value: "conflict" },
    { label: t("project.migrationRecords.reasonSecurity"), value: "security" },
    {
      label: t("project.migrationRecords.reasonService"),
      value: "service_access",
    },
    { label: t("project.migrationRecords.reasonReturn"), value: "return" },
    {
      label: t("project.migrationRecords.reasonRelocation"),
      value: "relocation_program",
    },
    { label: t("project.migrationRecords.reasonEconomic"), value: "economic" },
    { label: t("project.migrationRecords.reasonOther"), value: "other" },
  ];

  const housingOptions = [
    {
      label: t("project.migrationRecords.housingOwn"),
      value: "own_property",
    },
    {
      label: t("project.migrationRecords.housingRenting"),
      value: "renting",
    },
    {
      label: t("project.migrationRecords.housingRelatives"),
      value: "with_relatives",
    },
    {
      label: t("project.migrationRecords.housingCollective"),
      value: "collective_site",
    },
    { label: t("project.migrationRecords.housingHotel"), value: "hotel" },
    { label: t("project.migrationRecords.housingOther"), value: "other" },
    {
      label: t("project.migrationRecords.housingUnknown"),
      value: "unknown",
    },
  ];

  const columns: Column<MigrationRecord>[] = [
    {
      key: "migration_date",
      header: t("project.migrationRecords.date"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {r.migration_date
            ? new Date(r.migration_date).toLocaleDateString("en-CA")
            : "—"}
        </span>
      ),
    },
    {
      key: "movement_reason",
      header: t("project.migrationRecords.reason"),
      render: (r) => (
        <span className="text-fg-secondary">{r.movement_reason ?? "—"}</span>
      ),
    },
    {
      key: "housing",
      header: t("project.migrationRecords.housing"),
      render: (r) => (
        <span className="text-fg-secondary">
          {r.housing_at_destination ?? "—"}
        </span>
      ),
    },
    {
      key: "created_at",
      header: t("project.migrationRecords.created"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(r.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
  ];

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-fg">
          {t("project.migrationRecords.title")}
        </h2>
        <button
          type="button"
          onClick={() => setDialogOpen(true)}
          className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
        >
          <PlusIcon size={14} weight="bold" />
          {t("admin.common.add")}
        </button>
      </div>

      <DataTable
        columns={columns}
        data={data?.migration_records ?? []}
        keyExtractor={(r) => r.id}
        isLoading={isLoading}
      />

      <FormDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        title={t("project.migrationRecords.addTitle")}
        loading={createMutation.isPending}
        onSubmit={handleCreate}
        maxWidth="md"
      >
        <div className="space-y-3">
          <p className="text-xs font-medium text-fg-tertiary">
            {t("project.migrationRecords.from")}
          </p>
          <div className="grid grid-cols-3 gap-2">
            <UISelect
              value={form.from_country}
              onValueChange={(v) => {
                set("from_country", v);
                set("from_state", "");
                set("from_place_id", "");
              }}
              options={countryOptions}
              placeholder={t("project.people.selectCountry")}
              fullWidth
            />
            <UISelect
              value={form.from_state}
              onValueChange={(v) => {
                set("from_state", v);
                set("from_place_id", "");
              }}
              options={fromStateOptions}
              placeholder={t("project.people.selectState")}
              disabled={!form.from_country}
              fullWidth
            />
            <UISelect
              value={form.from_place_id}
              onValueChange={(v) => set("from_place_id", v)}
              options={fromPlaceOptions}
              placeholder={t("project.people.selectPlace")}
              disabled={!form.from_state}
              fullWidth
            />
          </div>

          <p className="text-xs font-medium text-fg-tertiary">
            {t("project.migrationRecords.to")}
          </p>
          <div className="grid grid-cols-3 gap-2">
            <UISelect
              value={form.dest_country}
              onValueChange={(v) => {
                set("dest_country", v);
                set("dest_state", "");
                set("destination_place_id", "");
              }}
              options={countryOptions}
              placeholder={t("project.people.selectCountry")}
              fullWidth
            />
            <UISelect
              value={form.dest_state}
              onValueChange={(v) => {
                set("dest_state", v);
                set("destination_place_id", "");
              }}
              options={destStateOptions}
              placeholder={t("project.people.selectState")}
              disabled={!form.dest_country}
              fullWidth
            />
            <UISelect
              value={form.destination_place_id}
              onValueChange={(v) => set("destination_place_id", v)}
              options={destPlaceOptions}
              placeholder={t("project.people.selectPlace")}
              disabled={!form.dest_state}
              fullWidth
            />
          </div>

          <div>
            <span className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.migrationRecords.date")}
            </span>
            <DatePicker
              value={form.migration_date}
              onChange={(v) => set("migration_date", v)}
            />
          </div>

          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.migrationRecords.reason")}
            </Field.Label>
            <UISelect
              value={form.movement_reason}
              onValueChange={(v) => set("movement_reason", v)}
              options={reasonOptions}
              fullWidth
            />
          </Field.Root>

          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.migrationRecords.housing")}
            </Field.Label>
            <UISelect
              value={form.housing_at_destination}
              onValueChange={(v) => set("housing_at_destination", v)}
              options={housingOptions}
              fullWidth
            />
          </Field.Root>

          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.migrationRecords.notes")}
            </Field.Label>
            <textarea
              value={form.notes}
              onChange={(e) => set("notes", e.target.value)}
              rows={3}
              className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
            />
          </Field.Root>
        </div>
      </FormDialog>
    </div>
  );
}
