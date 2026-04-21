import { type SyntheticEvent, useEffect, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { useCountries } from "@/hooks/reference/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/reference/use-offices";
import { useCreatePerson, usePerson, useUpdatePerson } from "@/hooks/people/use-people";
import { usePlaces } from "@/hooks/reference/use-places";
import { useStates } from "@/hooks/reference/use-states";
import { usePersonTags, useReplacePersonTags } from "@/hooks/tags/use-tags";
import { handleApiError } from "@/lib/form-error";
import { toSelectOptions } from "@/lib/options";
import { useToast } from "@/stores/toast";

import type { CreatePersonInput, UpdatePersonInput } from "@/types/person";

interface UsePersonDrawerFormOptions {
  open: boolean;
  projectId: string;
  personId: string | null;
}

const emptyForm = {
  first_name: "",
  last_name: "",
  patronymic: "",
  sex: "unknown",
  birth_date: "",
  age_group: "",
  primary_phone: "",
  email: "",
  origin_place_id: "",
  current_place_id: "",
  case_status: "new",
  external_id: "",
  office_id: "",
  consent_given: false,
  consent_date: "",
};

export function usePersonDrawerForm({ open, projectId, personId }: UsePersonDrawerFormOptions) {
  const { t } = useTranslation();
  const isEdit = personId !== null;

  const { data: person } = usePerson(projectId, personId ?? "");
  const { data: personTagsData } = usePersonTags(projectId, personId ?? "");
  const qc = useQueryClient();
  const createPerson = useCreatePerson(projectId);
  const updatePerson = useUpdatePerson(projectId);
  const replacePersonTags = useReplacePersonTags(projectId);

  const [tagIds, setTagIds] = useState<string[]>([]);
  const [originPlaceLabel, setOriginPlaceLabel] = useState("");
  const [currentPlaceLabel, setCurrentPlaceLabel] = useState("");

  const toast = useToast();
  const { form, set, error, setError, editingId, setEditingId } = useDrawerForm({
    initial: emptyForm,
    open,
    isEdit,
    data: person,
    mapData: (d) => ({
      first_name: (d.first_name as string) ?? "",
      last_name: (d.last_name as string) ?? "",
      patronymic: (d.patronymic as string) ?? "",
      sex: (d.sex as string) ?? "unknown",
      birth_date: (d.birth_date as string) ?? "",
      age_group: (d.age_group as string) ?? "",
      primary_phone: (d.primary_phone as string) ?? "",
      email: (d.email as string) ?? "",
      origin_place_id: (d.origin_place_id as string) ?? "",
      current_place_id: (d.current_place_id as string) ?? "",
      case_status: (d.case_status as string) ?? "new",
      external_id: (d.external_id as string) ?? "",
      office_id: (d.office_id as string) ?? "",
      consent_given: (d.consent_given as boolean) ?? false,
      consent_date: (d.consent_date as string) ?? "",
    }),
  });

  useEffect(() => {
    if (personTagsData) setTagIds(personTagsData.tag_ids ?? []);
  }, [personTagsData]);

  useEffect(() => {
    if (!open) setTagIds([]);
  }, [open]);

  const { data: countries } = useCountries();
  const { data: statesData } = useStates();
  const { data: placesData } = usePlaces();
  const { data: offices } = useOffices();

  const isPending = createPerson.isPending || updatePerson.isPending || replacePersonTags.isPending;

  function resolvePlaceLabel(placeId: string): string {
    if (!placeId) return "";
    const place = placesData?.places.find((p) => p.id === placeId);
    if (!place) return placeId;
    const state = statesData?.states.find((s) => s.id === place.state_id);
    const country = state ? (countries ?? []).find((c) => c.id === state.country_id) : undefined;
    const parts = [place.name];
    if (state) parts.push(state.name);
    if (country) parts.push(country.name);
    return parts.join(", ");
  }

  const resolvedOriginLabel = originPlaceLabel || resolvePlaceLabel(form.origin_place_id);
  const resolvedCurrentLabel = currentPlaceLabel || resolvePlaceLabel(form.current_place_id);

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    try {
      if (isEdit && editingId) {
        const data: UpdatePersonInput = {
          first_name: form.first_name,
          last_name: form.last_name || undefined,
          patronymic: form.patronymic || undefined,
          sex: form.sex as UpdatePersonInput["sex"],
          birth_date: form.birth_date || undefined,
          age_group: (form.age_group || undefined) as UpdatePersonInput["age_group"],
          primary_phone: form.primary_phone || undefined,
          email: form.email || undefined,
          origin_place_id: form.origin_place_id || undefined,
          current_place_id: form.current_place_id || undefined,
          case_status: form.case_status as UpdatePersonInput["case_status"],
          external_id: form.external_id || undefined,
          office_id: form.office_id || undefined,
          consent_given: form.consent_given,
          consent_date: form.consent_date || undefined,
        };
        await updatePerson.mutateAsync({ personId: editingId, data });
        await replacePersonTags.mutateAsync({ personId: editingId, ids: tagIds });
        await qc.invalidateQueries({ queryKey: ["people", projectId] });
        toast.success(t("project.people.saved"));
      } else {
        const input: CreatePersonInput = {
          first_name: form.first_name,
          ...(form.last_name && { last_name: form.last_name }),
          ...(form.patronymic && { patronymic: form.patronymic }),
          ...(form.sex && { sex: form.sex as CreatePersonInput["sex"] }),
          ...(form.birth_date && { birth_date: form.birth_date }),
          ...(form.age_group && {
            age_group: form.age_group as CreatePersonInput["age_group"],
          }),
          ...(form.primary_phone && { primary_phone: form.primary_phone }),
          ...(form.email && { email: form.email }),
          ...(form.origin_place_id && {
            origin_place_id: form.origin_place_id,
          }),
          ...(form.current_place_id && {
            current_place_id: form.current_place_id,
          }),
          ...(form.case_status && {
            case_status: form.case_status as CreatePersonInput["case_status"],
          }),
          ...(form.external_id && { external_id: form.external_id }),
          ...(form.office_id && { office_id: form.office_id }),
          consent_given: form.consent_given,
          ...(form.consent_date && { consent_date: form.consent_date }),
        };
        const created = await createPerson.mutateAsync(input);
        if (tagIds.length > 0) {
          await replacePersonTags.mutateAsync({ personId: created.id, ids: tagIds });
        }
        await qc.invalidateQueries({ queryKey: ["people", projectId] });
        setEditingId(created.id);
        toast.success(t("project.people.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const officeOptions = toSelectOptions(offices);

  return {
    isEdit,
    form,
    set,
    error,
    isPending,
    handleSubmit,
    tagIds,
    setTagIds,
    officeOptions,
    resolvedOriginLabel,
    resolvedCurrentLabel,
    setOriginPlaceLabel,
    setCurrentPlaceLabel,
  };
}
