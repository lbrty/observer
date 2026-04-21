import { type SyntheticEvent, useEffect, useRef } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

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
import { toSelectOptions } from "@/lib/options";
import { useToast } from "@/stores/toast";
import type { HousingAtDestination, MovementReason } from "@/types/migration-record";

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

interface UseMigrationRecordFormOptions {
  open: boolean;
  projectId: string;
  personId: string;
  recordId: string | null;
}

export function useMigrationRecordForm({
  open,
  projectId,
  personId,
  recordId,
}: UseMigrationRecordFormOptions) {
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

  // Track current form values via ref to avoid stale closure in the effect below
  const formRef = useRef(form);
  formRef.current = form;

  // Resolve place -> state -> country when record loads (runs once per data change,
  // reads form via ref to avoid adding form to deps and causing a cycle)
  useEffect(() => {
    if (!isEdit || !record || !allStates?.states || !allPlaces?.places) return;

    const states = allStates.states;
    const places = allPlaces.places;
    const currentForm = formRef.current;

    if (record.from_place_id && !currentForm.from_country) {
      const place = places.find((p) => p.id === record.from_place_id);
      if (place) {
        const state = states.find((s) => s.id === place.state_id);
        if (state) {
          set("from_country", state.country_id);
          set("from_state", state.id);
        }
      }
    }

    if (record.destination_place_id && !currentForm.dest_country) {
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
            ...(form.movement_reason && {
              movement_reason: form.movement_reason as MovementReason,
            }),
            ...(form.housing_at_destination && {
              housing_at_destination: form.housing_at_destination as HousingAtDestination,
            }),
            ...(form.notes && { notes: form.notes }),
          },
        });
        await qc.invalidateQueries({ queryKey: ["migration-records", projectId, personId] });
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
        await qc.invalidateQueries({ queryKey: ["migration-records", projectId, personId] });
        setEditingId(created.id);
        toast.success(t("project.migrationRecords.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  return {
    isEdit,
    form,
    set,
    error,
    isPending: createRecord.isPending || updateRecord.isPending,
    handleSubmit,
    countryOptions: toSelectOptions(countries),
    fromStateOptions: toSelectOptions(fromStates?.states),
    fromPlaceOptions: toSelectOptions(fromPlaces?.places),
    destStateOptions: toSelectOptions(destStates?.states),
    destPlaceOptions: toSelectOptions(destPlaces?.places),
  };
}
