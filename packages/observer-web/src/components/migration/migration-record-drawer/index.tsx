import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";

import { DetailsSection } from "./details-section";
import { PlaceSection } from "./place-section";
import { useMigrationRecordForm } from "./use-migration-record-form";

interface MigrationRecordDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  personId: string;
  recordId: string | null;
}

export function MigrationRecordDrawer({
  open,
  onOpenChange,
  projectId,
  personId,
  recordId,
}: MigrationRecordDrawerProps) {
  const { t } = useTranslation();
  const {
    isEdit,
    form,
    set,
    error,
    isPending,
    handleSubmit,
    countryOptions,
    fromStateOptions,
    fromPlaceOptions,
    destStateOptions,
    destPlaceOptions,
  } = useMigrationRecordForm({ open, projectId, personId, recordId });

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={
        isEdit ? t("project.migrationRecords.editTitle") : t("project.migrationRecords.addTitle")
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
      <FormTextarea label="" value={form.notes} onChange={(v) => set("notes", v)} rows={4} />
    </DrawerShell>
  );
}
