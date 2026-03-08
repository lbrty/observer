import { useEffect, useState } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { SectionHeading } from "@/components/section-heading";
import { TagPicker } from "@/components/tag-picker";
import { useCountries } from "@/hooks/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/use-offices";
import { useCreatePerson, usePerson, useUpdatePerson } from "@/hooks/use-people";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";
import { usePersonTags, useReplacePersonTags } from "@/hooks/use-tags";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";

import type { CreatePersonInput, UpdatePersonInput } from "@/types/person";

import { CaseSection } from "./case-section";
import { IdentitySection } from "./identity-section";
import { LocationSection } from "./location-section";

interface PersonDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
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

export function PersonDrawer({ open, onOpenChange, projectId, personId }: PersonDrawerProps) {
  const { t } = useTranslation();
  const isEdit = personId !== null;

  const { data: person } = usePerson(projectId, personId ?? "");
  const { data: personTagsData } = usePersonTags(projectId, personId ?? "");
  const qc = useQueryClient();
  const createPerson = useCreatePerson(projectId);
  const updatePerson = useUpdatePerson(projectId);
  const replacePersonTags = useReplacePersonTags(projectId);

  const [tagIds, setTagIds] = useState<string[]>([]);

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

  const [originPlaceLabel, setOriginPlaceLabel] = useState("");
  const [currentPlaceLabel, setCurrentPlaceLabel] = useState("");

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

  async function handleSubmit(e: React.SyntheticEvent) {
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

  const officeOptions = (offices ?? []).map((o) => ({ label: o.name, value: o.id }));

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? t("project.people.editTitle") : t("project.people.formTitle")}
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.people.save")}
      savingLabel={t("project.people.saving")}
    >
      <ErrorBanner message={error} />

      <IdentitySection form={form} set={(k, v) => set(k as keyof typeof form, v)} />

      <LocationSection
        originPlaceId={form.origin_place_id}
        currentPlaceId={form.current_place_id}
        originPlaceLabel={resolvedOriginLabel}
        currentPlaceLabel={resolvedCurrentLabel}
        onSelectOrigin={(place, state, country) => {
          set("origin_place_id", place.id);
          setOriginPlaceLabel(`${place.name}, ${state.name}, ${country.name}`);
        }}
        onClearOrigin={() => {
          set("origin_place_id", "");
          setOriginPlaceLabel("");
        }}
        onSelectCurrent={(place, state, country) => {
          set("current_place_id", place.id);
          setCurrentPlaceLabel(`${place.name}, ${state.name}, ${country.name}`);
        }}
        onClearCurrent={() => {
          set("current_place_id", "");
          setCurrentPlaceLabel("");
        }}
      />

      <CaseSection form={form} set={(k, v) => set(k as keyof typeof form, v)} officeOptions={officeOptions} />

      <SectionHeading>{t("project.tags.title")}</SectionHeading>
      <TagPicker projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
    </DrawerShell>
  );
}
