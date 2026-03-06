import { useState } from "react";

import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { AddReferenceDialog } from "@/components/add-reference-dialog";
import { ErrorBanner } from "@/components/alert-banner";
import { DatePicker } from "@/components/date-picker";
import { DrawerShell } from "@/components/drawer-shell";
import { inputClass } from "@/components/form-field";
import { FormField } from "@/components/form-field";
import { PlusIcon } from "@/components/icons";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { useCountries, useCreateCountry } from "@/hooks/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import { useOffices } from "@/hooks/use-offices";
import { useCreatePerson, usePerson, useUpdatePerson } from "@/hooks/use-people";
import { useCreatePlace, usePlaces } from "@/hooks/use-places";
import { useCreateState, useStates } from "@/hooks/use-states";
import { HTTPError } from "@/lib/api";
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
  origin_country: "",
  origin_state: "",
  origin_place_id: "",
  current_country: "",
  current_state: "",
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
      origin_country: "",
      origin_state: "",
      origin_place_id: (d.origin_place_id as string) ?? "",
      current_country: "",
      current_state: "",
      current_place_id: (d.current_place_id as string) ?? "",
      case_status: (d.case_status as string) ?? "new",
      external_id: (d.external_id as string) ?? "",
      office_id: (d.office_id as string) ?? "",
      consent_given: (d.consent_given as boolean) ?? false,
      consent_date: (d.consent_date as string) ?? "",
    }),
  });

  const [addCountryOpen, setAddCountryOpen] = useState(false);
  const [addStateOpen, setAddStateOpen] = useState<{ open: boolean; forOrigin: boolean }>({
    open: false,
    forOrigin: true,
  });
  const [addPlaceOpen, setAddPlaceOpen] = useState<{ open: boolean; forOrigin: boolean }>({
    open: false,
    forOrigin: true,
  });

  const [newCountryName, setNewCountryName] = useState("");
  const [newCountryCode, setNewCountryCode] = useState("");
  const [newStateName, setNewStateName] = useState("");
  const [newStateConflictZone, setNewStateConflictZone] = useState("");
  const [newPlaceName, setNewPlaceName] = useState("");
  const [dialogError, setDialogError] = useState("");

  const createCountry = useCreateCountry();
  const createState = useCreateState();
  const createPlace = useCreatePlace();

  const { data: countries } = useCountries();
  const { data: originStates } = useStates(form.origin_country || undefined);
  const { data: originPlaces } = usePlaces(form.origin_state || undefined);
  const { data: currentStates } = useStates(form.current_country || undefined);
  const { data: currentPlaces } = usePlaces(form.current_state || undefined);
  const { data: offices } = useOffices();

  const isPending = createPerson.isPending || updatePerson.isPending;

  async function handleAddCountry() {
    setDialogError("");
    try {
      await createCountry.mutateAsync({ name: newCountryName, code: newCountryCode });
      setAddCountryOpen(false);
      setNewCountryName("");
      setNewCountryCode("");
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setDialogError((body as { error?: string } | null)?.error ?? err.message);
      }
    }
  }

  async function handleAddState() {
    setDialogError("");
    const countryId = addStateOpen.forOrigin ? form.origin_country : form.current_country;
    try {
      const created = await createState.mutateAsync({
        countryId,
        data: {
          name: newStateName,
          ...(newStateConflictZone && { conflict_zone: newStateConflictZone }),
        },
      });
      if (addStateOpen.forOrigin) {
        set("origin_state", created.id);
        set("origin_place_id", "");
      } else {
        set("current_state", created.id);
        set("current_place_id", "");
      }
      setAddStateOpen({ open: false, forOrigin: true });
      setNewStateName("");
      setNewStateConflictZone("");
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setDialogError((body as { error?: string } | null)?.error ?? err.message);
      }
    }
  }

  async function handleAddPlace() {
    setDialogError("");
    const stateId = addPlaceOpen.forOrigin ? form.origin_state : form.current_state;
    try {
      const created = await createPlace.mutateAsync({
        stateId,
        data: { name: newPlaceName },
      });
      if (addPlaceOpen.forOrigin) {
        set("origin_place_id", created.id);
      } else {
        set("current_place_id", created.id);
      }
      setAddPlaceOpen({ open: false, forOrigin: true });
      setNewPlaceName("");
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setDialogError((body as { error?: string } | null)?.error ?? err.message);
      }
    }
  }

  async function handleSubmit(e: React.FormEvent) {
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

  const countryOptions = (countries ?? []).map((c) => ({ label: c.name, value: c.id }));
  const originStateOptions = (originStates?.states ?? []).map((s) => ({ label: s.name, value: s.id }));
  const originPlaceOptions = (originPlaces?.places ?? []).map((p) => ({ label: p.name, value: p.id }));
  const currentStateOptions = (currentStates?.states ?? []).map((s) => ({ label: s.name, value: s.id }));
  const currentPlaceOptions = (currentPlaces?.places ?? []).map((p) => ({ label: p.name, value: p.id }));
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
    { label: t("project.people.ageMiddleAged"), value: "middle_aged_adult" },
    { label: t("project.people.ageOlderAdult"), value: "old_adult" },
    { label: t("project.people.sexUnknown"), value: "unknown" },
  ];

  const addBtnClass =
    "inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg disabled:opacity-30";

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
      <p className="text-xs font-medium text-fg-tertiary">{t("project.people.originPlace")}</p>
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.origin_country}
              onValueChange={(v) => {
                set("origin_country", v);
                set("origin_state", "");
                set("origin_place_id", "");
              }}
              options={countryOptions}
              placeholder={t("project.people.selectCountry")}
              fullWidth
            />
          </div>
          <button
            type="button"
            onClick={() => {
              setDialogError("");
              setNewCountryName("");
              setNewCountryCode("");
              setAddCountryOpen(true);
            }}
            className={addBtnClass}
            title={t("project.people.addCountry")}
          >
            <PlusIcon size={16} />
          </button>
        </div>
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.origin_state}
              onValueChange={(v) => {
                set("origin_state", v);
                set("origin_place_id", "");
              }}
              options={originStateOptions}
              placeholder={t("project.people.selectState")}
              disabled={!form.origin_country}
              fullWidth
            />
          </div>
          <button
            type="button"
            disabled={!form.origin_country}
            onClick={() => {
              setDialogError("");
              setNewStateName("");
              setNewStateConflictZone("");
              setAddStateOpen({ open: true, forOrigin: true });
            }}
            className={addBtnClass}
            title={t("project.people.addState")}
          >
            <PlusIcon size={16} />
          </button>
        </div>
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.origin_place_id}
              onValueChange={(v) => set("origin_place_id", v)}
              options={originPlaceOptions}
              placeholder={t("project.people.selectPlace")}
              disabled={!form.origin_state}
              fullWidth
            />
          </div>
          <button
            type="button"
            disabled={!form.origin_state}
            onClick={() => {
              setDialogError("");
              setNewPlaceName("");
              setAddPlaceOpen({ open: true, forOrigin: true });
            }}
            className={addBtnClass}
            title={t("project.people.addPlace")}
          >
            <PlusIcon size={16} />
          </button>
        </div>
      </div>

      <p className="text-xs font-medium text-fg-tertiary">{t("project.people.currentPlace")}</p>
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.current_country}
              onValueChange={(v) => {
                set("current_country", v);
                set("current_state", "");
                set("current_place_id", "");
              }}
              options={countryOptions}
              placeholder={t("project.people.selectCountry")}
              fullWidth
            />
          </div>
          <button
            type="button"
            onClick={() => {
              setDialogError("");
              setNewCountryName("");
              setNewCountryCode("");
              setAddCountryOpen(true);
            }}
            className={addBtnClass}
            title={t("project.people.addCountry")}
          >
            <PlusIcon size={16} />
          </button>
        </div>
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.current_state}
              onValueChange={(v) => {
                set("current_state", v);
                set("current_place_id", "");
              }}
              options={currentStateOptions}
              placeholder={t("project.people.selectState")}
              disabled={!form.current_country}
              fullWidth
            />
          </div>
          <button
            type="button"
            disabled={!form.current_country}
            onClick={() => {
              setDialogError("");
              setNewStateName("");
              setNewStateConflictZone("");
              setAddStateOpen({ open: true, forOrigin: false });
            }}
            className={addBtnClass}
            title={t("project.people.addState")}
          >
            <PlusIcon size={16} />
          </button>
        </div>
        <div className="flex gap-1.5">
          <div className="flex-1">
            <UISelect
              value={form.current_place_id}
              onValueChange={(v) => set("current_place_id", v)}
              options={currentPlaceOptions}
              placeholder={t("project.people.selectPlace")}
              disabled={!form.current_state}
              fullWidth
            />
          </div>
          <button
            type="button"
            disabled={!form.current_state}
            onClick={() => {
              setDialogError("");
              setNewPlaceName("");
              setAddPlaceOpen({ open: true, forOrigin: false });
            }}
            className={addBtnClass}
            title={t("project.people.addPlace")}
          >
            <PlusIcon size={16} />
          </button>
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

      <AddReferenceDialog
        open={addCountryOpen}
        onOpenChange={setAddCountryOpen}
        title={t("project.people.addCountry")}
        onSubmit={handleAddCountry}
        isPending={createCountry.isPending}
        error={dialogError}
      >
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.countries.name")} *
          </Field.Label>
          <Field.Control
            required
            value={newCountryName}
            onChange={(e) => setNewCountryName(e.target.value)}
            className={inputClass}
          />
        </Field.Root>
        <Field.Root className="mt-3">
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.countries.code")}
          </Field.Label>
          <Field.Control
            value={newCountryCode}
            onChange={(e) => setNewCountryCode(e.target.value)}
            className={inputClass}
            maxLength={3}
          />
        </Field.Root>
      </AddReferenceDialog>

      <AddReferenceDialog
        open={addStateOpen.open}
        onOpenChange={(v) => setAddStateOpen((s) => ({ ...s, open: v }))}
        title={t("project.people.addState")}
        onSubmit={handleAddState}
        isPending={createState.isPending}
        error={dialogError}
      >
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.states.name")} *
          </Field.Label>
          <Field.Control
            required
            value={newStateName}
            onChange={(e) => setNewStateName(e.target.value)}
            className={inputClass}
          />
        </Field.Root>
        <Field.Root className="mt-3">
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.states.conflictZone")}
          </Field.Label>
          <Field.Control
            value={newStateConflictZone}
            onChange={(e) => setNewStateConflictZone(e.target.value)}
            className={inputClass}
          />
        </Field.Root>
      </AddReferenceDialog>

      <AddReferenceDialog
        open={addPlaceOpen.open}
        onOpenChange={(v) => setAddPlaceOpen((s) => ({ ...s, open: v }))}
        title={t("project.people.addPlace")}
        onSubmit={handleAddPlace}
        isPending={createPlace.isPending}
        error={dialogError}
      >
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.places.name")} *
          </Field.Label>
          <Field.Control
            required
            value={newPlaceName}
            onChange={(e) => setNewPlaceName(e.target.value)}
            className={inputClass}
          />
        </Field.Root>
      </AddReferenceDialog>
    </DrawerShell>
  );
}
