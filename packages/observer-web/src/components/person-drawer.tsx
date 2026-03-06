import { useState } from "react";

import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DatePicker } from "@/components/date-picker";
import { DrawerShell } from "@/components/drawer-shell";
import { FormField } from "@/components/form-field";
import { PlaceCombobox } from "@/components/place-combobox";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { useCountries } from "@/hooks/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/use-offices";
import { useCreatePerson, usePerson, useUpdatePerson } from "@/hooks/use-people";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";

import type { CreatePersonInput, UpdatePersonInput } from "@/types/person";

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
  const qc = useQueryClient();
  const createPerson = useCreatePerson(projectId);
  const updatePerson = useUpdatePerson(projectId);

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

  // Labels for selected places (resolved from reference data or set on select)
  const [originPlaceLabel, setOriginPlaceLabel] = useState("");
  const [currentPlaceLabel, setCurrentPlaceLabel] = useState("");

  // Reference data for resolving place IDs to labels in edit mode
  const { data: countries } = useCountries();
  const { data: statesData } = useStates();
  const { data: placesData } = usePlaces();
  const { data: offices } = useOffices();

  const isPending = createPerson.isPending || updatePerson.isPending;

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
        await qc.invalidateQueries({ queryKey: ["people", projectId] });
        setEditingId(created.id);
        toast.success(t("project.people.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const officeOptions = (offices ?? []).map((o) => ({ label: o.name, value: o.id }));

  const sexOptions = [
    { label: t("project.people.sexMale"), value: "male" },
    { label: t("project.people.sexFemale"), value: "female" },
    { label: t("project.people.sexOther"), value: "other" },
    { label: t("project.people.sexUnknown"), value: "unknown" },
  ];
  const caseStatusOptions = [
    { label: t("project.people.new"), value: "new" },
    { label: t("project.people.active"), value: "active" },
    { label: t("project.people.closed"), value: "closed" },
    { label: t("project.people.archived"), value: "archived" },
  ];
  const ageGroupOptions = [
    { label: t("project.people.ageInfant"), value: "infant" },
    { label: t("project.people.ageToddler"), value: "toddler" },
    { label: t("project.people.agePreSchool"), value: "pre_school" },
    { label: t("project.people.ageMiddleChildhood"), value: "middle_childhood" },
    { label: t("project.people.ageYoungTeen"), value: "young_teen" },
    { label: t("project.people.ageTeenager"), value: "teenager" },
    { label: t("project.people.ageYoungAdult"), value: "young_adult" },
    { label: t("project.people.ageEarlyAdult"), value: "early_adult" },
    { label: t("project.people.ageMiddleAgedAdult"), value: "middle_aged_adult" },
    { label: t("project.people.ageOldAdult"), value: "old_adult" },
    { label: t("project.people.sexUnknown"), value: "unknown" },
  ];

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

      <SectionHeading>{t("project.people.identity")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField
          label={t("project.people.firstName")}
          value={form.first_name}
          onChange={(v) => set("first_name", v)}
          required
        />
        <FormField
          label={t("project.people.lastName")}
          value={form.last_name}
          onChange={(v) => set("last_name", v)}
        />
        <FormField
          label={t("project.people.patronymic")}
          value={form.patronymic}
          onChange={(v) => set("patronymic", v)}
        />

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.sexLabel")}
          </Field.Label>
          <UISelect
            value={form.sex}
            onValueChange={(v) => set("sex", v)}
            options={sexOptions}
            fullWidth
          />
        </Field.Root>

        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.birthDate")}
          </span>
          <DatePicker value={form.birth_date} onChange={(v) => set("birth_date", v)} />
        </div>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.ageGroup")}
          </Field.Label>
          <UISelect
            value={form.age_group}
            onValueChange={(v) => set("age_group", v)}
            options={ageGroupOptions}
            fullWidth
          />
        </Field.Root>

        <FormField
          label={t("project.people.phone")}
          value={form.primary_phone}
          onChange={(v) => set("primary_phone", v)}
        />
        <FormField
          label={t("project.people.email")}
          value={form.email}
          onChange={(v) => set("email", v)}
          type="email"
        />
      </div>

      <SectionHeading>{t("project.people.location")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.originPlace")}
          </span>
          {form.origin_place_id ? (
            <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
              <span className="flex-1 truncate text-sm text-fg">
                {resolvedOriginLabel || form.origin_place_id}
              </span>
              <button
                type="button"
                onClick={() => {
                  set("origin_place_id", "");
                  setOriginPlaceLabel("");
                }}
                className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
              >
                ×
              </button>
            </div>
          ) : (
            <PlaceCombobox
              onSelect={(place, state, country) => {
                set("origin_place_id", place.id);
                setOriginPlaceLabel(`${place.name}, ${state.name}, ${country.name}`);
              }}
            />
          )}
        </div>

        <div>
          <span className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.currentPlace")}
          </span>
          {form.current_place_id ? (
            <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
              <span className="flex-1 truncate text-sm text-fg">
                {resolvedCurrentLabel || form.current_place_id}
              </span>
              <button
                type="button"
                onClick={() => {
                  set("current_place_id", "");
                  setCurrentPlaceLabel("");
                }}
                className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
              >
                ×
              </button>
            </div>
          ) : (
            <PlaceCombobox
              onSelect={(place, state, country) => {
                set("current_place_id", place.id);
                setCurrentPlaceLabel(`${place.name}, ${state.name}, ${country.name}`);
              }}
            />
          )}
        </div>
      </div>

      <SectionHeading>{t("project.people.case")}</SectionHeading>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.caseStatusLabel")}
          </Field.Label>
          <UISelect
            value={form.case_status}
            onValueChange={(v) => set("case_status", v)}
            options={caseStatusOptions}
            fullWidth
          />
        </Field.Root>

        <FormField
          label={t("project.people.externalId")}
          value={form.external_id}
          onChange={(v) => set("external_id", v)}
        />

        {officeOptions.length > 0 && (
          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.people.office")}
            </Field.Label>
            <UISelect
              value={form.office_id}
              onValueChange={(v) => set("office_id", v)}
              options={officeOptions}
              fullWidth
            />
          </Field.Root>
        )}

        <div className="col-span-full space-y-4">
          <UISwitch
            checked={form.consent_given}
            onCheckedChange={(v) => set("consent_given", v)}
            label={t("project.people.consentGiven")}
          />

          {form.consent_given && (
            <div>
              <span className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("project.people.consentDate")}
              </span>
              <DatePicker value={form.consent_date} onChange={(v) => set("consent_date", v)} />
            </div>
          )}
        </div>
      </div>
    </DrawerShell>
  );
}
