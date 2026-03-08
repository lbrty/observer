import { type SyntheticEvent, useEffect } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";
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
import { useToast } from "@/stores/toast";

import type { HousingAtDestination, MovementReason } from "@/types/migration-record";

import { DetailsSection } from "./details-section";
import { PlaceSection } from "./place-section";

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

  const toast = useToast();
  const { form, set, error, setError, editingId, setEditingId } = useDrawerForm({
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
  const { data: allStates } = useStates();
  const { data: allPlaces } = usePlaces();
  const { data: fromStates } = useStates(form.from_country || undefined);
  const { data: fromPlaces } = usePlaces(form.from_state || undefined);
  const { data: destStates } = useStates(form.dest_country || undefined);
  const { data: destPlaces } = usePlaces(form.dest_state || undefined);

  // resolve place -> state -> country when record loads
  useEffect(() => {
    if (!isEdit || !record || !allStates?.states || !allPlaces?.places) return;

    const states = allStates.states;
    const places = allPlaces.places;

    if (record.from_place_id && !form.from_country) {
      const place = places.find((p) => p.id === record.from_place_id);
      if (place) {
        const state = states.find((s) => s.id === place.state_id);
        if (state) {
          set("from_country", state.country_id);
          set("from_state", state.id);
        }
      }
    }

    if (record.destination_place_id && !form.dest_country) {
      const place = places.find((p) => p.id === record.destination_place_id);
      if (place) {
        const state = states.find((s) => s.id === place.state_id);
        if (state) {
          set("dest_country", state.country_id);
          set("dest_state", state.id);
        }
      }
    }
  }, [record, allStates, allPlaces]);

  const isPending = createRecord.isPending || updateRecord.isPending;

  async function handleSubmit(e: SyntheticEvent) {
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
        toast.success(t("project.migrationRecords.saved"));
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
        toast.success(t("project.migrationRecords.saved"));
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
      {error && <ErrorBanner message={error} />}

      <PlaceSection
        title={t("project.migrationRecords.from")}
        country={form.from_country}
        state={form.from_state}
        place={form.from_place_id}
        countryOptions={countryOptions}
        stateOptions={fromStateOptions}
        placeOptions={fromPlaceOptions}
        countryPlaceholder={t("project.people.selectCountry")}
        statePlaceholder={t("project.people.selectState")}
        placePlaceholder={t("project.people.selectPlace")}
        onCountryChange={(v) => {
          set("from_country", v);
          set("from_state", "");
          set("from_place_id", "");
        }}
        onStateChange={(v) => {
          set("from_state", v);
          set("from_place_id", "");
        }}
        onPlaceChange={(v) => set("from_place_id", v)}
      />

      <PlaceSection
        title={t("project.migrationRecords.to")}
        country={form.dest_country}
        state={form.dest_state}
        place={form.destination_place_id}
        countryOptions={countryOptions}
        stateOptions={destStateOptions}
        placeOptions={destPlaceOptions}
        countryPlaceholder={t("project.people.selectCountry")}
        statePlaceholder={t("project.people.selectState")}
        placePlaceholder={t("project.people.selectPlace")}
        onCountryChange={(v) => {
          set("dest_country", v);
          set("dest_state", "");
          set("destination_place_id", "");
        }}
        onStateChange={(v) => {
          set("dest_state", v);
          set("destination_place_id", "");
        }}
        onPlaceChange={(v) => set("destination_place_id", v)}
      />

      <DetailsSection
        migrationDate={form.migration_date}
        movementReason={form.movement_reason}
        housingAtDestination={form.housing_at_destination}
        onDateChange={(v) => set("migration_date", v)}
        onReasonChange={(v) => set("movement_reason", v)}
        onHousingChange={(v) => set("housing_at_destination", v)}
      />

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
