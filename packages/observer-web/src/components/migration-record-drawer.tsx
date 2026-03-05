import { type FormEvent } from "react";

import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner, SuccessBanner } from "@/components/alert-banner";
import { DatePicker } from "@/components/date-picker";
import { DrawerShell } from "@/components/drawer-shell";
import { FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { useCountries } from "@/hooks/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import {
  useCreateMigrationRecord,
  useMigrationRecord,
  useUpdateMigrationRecord,
} from "@/hooks/use-migration-records";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";
import { handleApiError } from "@/lib/form-error";

import type { HousingAtDestination, MovementReason } from "@/types/migration-record";

interface MigrationRecordDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  personId: string;
  recordId: string | null;
}

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

export function MigrationRecordDrawer({
  open,
  onOpenChange,
  projectId,
  personId,
  recordId,
}: MigrationRecordDrawerProps) {
  const { t } = useTranslation();
  const isEdit = recordId !== null;

  const { data: record } = useMigrationRecord(projectId, personId, recordId ?? "");
  const qc = useQueryClient();
  const createRecord = useCreateMigrationRecord(projectId, personId);
  const updateRecord = useUpdateMigrationRecord(projectId, personId);

  const { form, set, saved, setSaved, error, setError, editingId, setEditingId } = useDrawerForm({
    initial: emptyForm,
    open,
    isEdit,
    data: record,
    mapData: (d) => ({
      from_country: "",
      from_state: "",
      from_place_id: (d.from_place_id as string) ?? "",
      dest_country: "",
      dest_state: "",
      destination_place_id: (d.destination_place_id as string) ?? "",
      migration_date: (d.migration_date as string) ?? "",
      movement_reason: (d.movement_reason as string) ?? "",
      housing_at_destination: (d.housing_at_destination as string) ?? "",
      notes: (d.notes as string) ?? "",
    }),
  });

  const { data: countries } = useCountries();
  const { data: fromStates } = useStates(form.from_country || undefined);
  const { data: fromPlaces } = usePlaces(form.from_state || undefined);
  const { data: destStates } = useStates(form.dest_country || undefined);
  const { data: destPlaces } = usePlaces(form.dest_state || undefined);

  const isPending = createRecord.isPending || updateRecord.isPending;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    try {
      if (isEdit && editingId) {
        await updateRecord.mutateAsync({
          id: editingId,
          data: {
            ...(form.from_place_id && { from_place_id: form.from_place_id }),
            ...(form.destination_place_id && { destination_place_id: form.destination_place_id }),
            ...(form.migration_date && { migration_date: form.migration_date }),
            ...(form.movement_reason && { movement_reason: form.movement_reason }),
            ...(form.housing_at_destination && {
              housing_at_destination: form.housing_at_destination,
            }),
            ...(form.notes && { notes: form.notes }),
          },
        });
        await qc.invalidateQueries({
          queryKey: ["migration-records", projectId, personId],
        });
        setSaved(true);
      } else {
        const created = await createRecord.mutateAsync({
          ...(form.from_place_id && { from_place_id: form.from_place_id }),
          ...(form.destination_place_id && { destination_place_id: form.destination_place_id }),
          ...(form.migration_date && { migration_date: form.migration_date }),
          ...(form.movement_reason && {
            movement_reason: form.movement_reason as MovementReason,
          }),
          ...(form.housing_at_destination && {
            housing_at_destination: form.housing_at_destination as HousingAtDestination,
          }),
          ...(form.notes && { notes: form.notes }),
        });
        await qc.invalidateQueries({
          queryKey: ["migration-records", projectId, personId],
        });
        setEditingId(created.id);
        setSaved(true);
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
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
    { label: t("project.migrationRecords.reasonService"), value: "service_access" },
    { label: t("project.migrationRecords.reasonReturn"), value: "return" },
    { label: t("project.migrationRecords.reasonRelocation"), value: "relocation_program" },
    { label: t("project.migrationRecords.reasonEconomic"), value: "economic" },
    { label: t("project.migrationRecords.reasonOther"), value: "other" },
  ];

  const housingOptions = [
    { label: t("project.migrationRecords.housingOwn"), value: "own_property" },
    { label: t("project.migrationRecords.housingRenting"), value: "renting" },
    { label: t("project.migrationRecords.housingRelatives"), value: "with_relatives" },
    { label: t("project.migrationRecords.housingCollective"), value: "collective_site" },
    { label: t("project.migrationRecords.housingHotel"), value: "hotel" },
    { label: t("project.migrationRecords.housingOther"), value: "other" },
    { label: t("project.migrationRecords.housingUnknown"), value: "unknown" },
  ];

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={
        isEdit
          ? t("project.migrationRecords.editTitle")
          : t("project.migrationRecords.addTitle")
      }
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.migrationRecords.save")}
      savingLabel={t("project.migrationRecords.saving")}
    >
      {saved && <SuccessBanner message={t("project.migrationRecords.saved")} />}
      {error && <ErrorBanner message={error} />}

      <SectionHeading>{t("project.migrationRecords.from")}</SectionHeading>
      <div className="grid grid-cols-3 gap-3">
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

      <SectionHeading>{t("project.migrationRecords.to")}</SectionHeading>
      <div className="grid grid-cols-3 gap-3">
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

      <SectionHeading>{t("project.migrationRecords.details")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.migrationRecords.date")}
          </span>
          <DatePicker value={form.migration_date} onChange={(v) => set("migration_date", v)} />
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
      </div>

      <SectionHeading>{t("project.migrationRecords.notes")}</SectionHeading>
      <FormTextarea
        label=""
        value={form.notes}
        onChange={(v) => set("notes", v)}
        rows={4}
      />
    </DrawerShell>
  );
}
