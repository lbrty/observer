import { DrawerPreview as Drawer } from "@base-ui/react/drawer";
import { Field } from "@base-ui/react/field";
import { useQueryClient } from "@tanstack/react-query";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { AddReferenceDialog } from "@/components/add-reference-dialog";
import { CheckIcon, PlusIcon, WarningIcon, XIcon } from "@/components/icons";
import { HTTPError } from "@/lib/api";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { useCountries, useCreateCountry } from "@/hooks/use-countries";
import { useOffices } from "@/hooks/use-offices";
import {
  useCreatePerson,
  usePerson,
  useUpdatePerson,
} from "@/hooks/use-people";
import { useCreatePlace, usePlaces } from "@/hooks/use-places";
import { useCreateState, useStates } from "@/hooks/use-states";
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

export function PersonDrawer({
  open,
  onOpenChange,
  projectId,
  personId,
}: PersonDrawerProps) {
  const { t } = useTranslation();
  const isEdit = personId !== null;

  const { data: person } = usePerson(projectId, personId ?? "");
  const qc = useQueryClient();
  const createPerson = useCreatePerson(projectId);
  const updatePerson = useUpdatePerson(projectId);

  const [form, setForm] = useState(emptyForm);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);

  const [addCountryOpen, setAddCountryOpen] = useState(false);
  const [addStateOpen, setAddStateOpen] = useState<{ open: boolean; forOrigin: boolean }>({ open: false, forOrigin: true });
  const [addPlaceOpen, setAddPlaceOpen] = useState<{ open: boolean; forOrigin: boolean }>({ open: false, forOrigin: true });

  const [newCountryName, setNewCountryName] = useState("");
  const [newCountryCode, setNewCountryCode] = useState("");
  const [newStateName, setNewStateName] = useState("");
  const [newStateConflictZone, setNewStateConflictZone] = useState("");
  const [newPlaceName, setNewPlaceName] = useState("");
  const [dialogError, setDialogError] = useState("");

  const createCountry = useCreateCountry();
  const createState = useCreateState();
  const createPlace = useCreatePlace();

  useEffect(() => {
    if (!open) {
      setForm(emptyForm);
      setSaved(false);
      setError("");
      setEditingId(null);
      return;
    }
    if (isEdit && person) {
      setForm({
        first_name: person.first_name,
        last_name: person.last_name ?? "",
        patronymic: person.patronymic ?? "",
        sex: person.sex,
        birth_date: person.birth_date ?? "",
        age_group: person.age_group ?? "",
        primary_phone: person.primary_phone ?? "",
        email: person.email ?? "",
        origin_country: "",
        origin_state: "",
        origin_place_id: person.origin_place_id ?? "",
        current_country: "",
        current_state: "",
        current_place_id: person.current_place_id ?? "",
        case_status: person.case_status,
        external_id: person.external_id ?? "",
        office_id: person.office_id ?? "",
        consent_given: person.consent_given,
        consent_date: person.consent_date ?? "",
      });
      setEditingId(person.id);
    }
  }, [open, isEdit, person]);

  const { data: countries } = useCountries();
  const { data: originStates } = useStates(form.origin_country || undefined);
  const { data: originPlaces } = usePlaces(form.origin_state || undefined);
  const { data: currentStates } = useStates(form.current_country || undefined);
  const { data: currentPlaces } = usePlaces(form.current_state || undefined);
  const { data: offices } = useOffices();

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }));
    setSaved(false);
    setError("");
  }

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
        data: { name: newStateName, ...(newStateConflictZone && { conflict_zone: newStateConflictZone }) },
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

  async function handleSubmit(e: FormEvent) {
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
          age_group: (form.age_group ||
            undefined) as UpdatePersonInput["age_group"],
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
        setSaved(true);
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
        setSaved(true);
      }
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
    }
  }

  const countryOptions = (countries ?? []).map((c) => ({
    label: c.name,
    value: c.id,
  }));
  const originStateOptions = (originStates?.states ?? []).map((s) => ({
    label: s.name,
    value: s.id,
  }));
  const originPlaceOptions = (originPlaces?.places ?? []).map((p) => ({
    label: p.name,
    value: p.id,
  }));
  const currentStateOptions = (currentStates?.states ?? []).map((s) => ({
    label: s.name,
    value: s.id,
  }));
  const currentPlaceOptions = (currentPlaces?.places ?? []).map((p) => ({
    label: p.name,
    value: p.id,
  }));
  const officeOptions = (offices ?? []).map((o) => ({
    label: o.name,
    value: o.id,
  }));

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
    {
      label: t("project.people.ageMiddleChildhood"),
      value: "middle_childhood",
    },
    { label: t("project.people.ageYoungTeen"), value: "young_teen" },
    { label: t("project.people.ageTeenager"), value: "teenager" },
    { label: t("project.people.ageYoungAdult"), value: "young_adult" },
    { label: t("project.people.ageEarlyAdult"), value: "early_adult" },
    { label: t("project.people.ageMiddleAged"), value: "middle_aged_adult" },
    { label: t("project.people.ageOlderAdult"), value: "old_adult" },
    { label: t("project.people.sexUnknown"), value: "unknown" },
  ];

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";

  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs transition-opacity duration-200 data-ending-style:opacity-0 data-starting-style:opacity-0" />
        <Drawer.Viewport className="fixed inset-0 z-50">
          <Drawer.Popup className="fixed top-0 right-0 flex h-dvh w-full max-w-[840px] flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full">
            <div className="flex shrink-0 items-center justify-between border-b border-border-secondary px-6 py-4">
              <Drawer.Title className="font-serif text-lg font-semibold text-fg">
                {isEdit
                  ? t("project.people.editTitle")
                  : t("project.people.formTitle")}
              </Drawer.Title>
              <Drawer.Close className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg">
                <XIcon size={18} />
              </Drawer.Close>
            </div>

            <form
              onSubmit={handleSubmit}
              className="flex min-h-0 flex-1 flex-col"
            >
              <div className="flex-1 space-y-5 overflow-y-auto px-6 py-5">
                {saved && (
                  <div className="flex items-center gap-2 rounded-lg border border-foam/20 bg-foam/8 px-3 py-2.5 text-sm font-medium text-foam">
                    <CheckIcon size={16} weight="bold" className="shrink-0" />
                    {t("project.people.saved")}
                  </div>
                )}
                {error && (
                  <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
                    <WarningIcon size={16} weight="bold" className="shrink-0" />
                    {error}
                  </div>
                )}

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">{t("project.people.identity")}</h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.firstName")} *
                    </Field.Label>
                    <Field.Control
                      required
                      value={form.first_name}
                      onChange={(e) => set("first_name", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.lastName")}
                    </Field.Label>
                    <Field.Control
                      value={form.last_name}
                      onChange={(e) => set("last_name", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.patronymic")}
                    </Field.Label>
                    <Field.Control
                      value={form.patronymic}
                      onChange={(e) => set("patronymic", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

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

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.birthDate")}
                    </Field.Label>
                    <Field.Control
                      type="date"
                      value={form.birth_date}
                      onChange={(e) => set("birth_date", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

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

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.phone")}
                    </Field.Label>
                    <Field.Control
                      value={form.primary_phone}
                      onChange={(e) => set("primary_phone", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.email")}
                    </Field.Label>
                    <Field.Control
                      type="email"
                      value={form.email}
                      onChange={(e) => set("email", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>
                </div>

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">{t("project.people.location")}</h3>
                <p className="text-xs font-medium text-fg-tertiary">
                  {t("project.people.originPlace")}
                </p>
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
                      onClick={() => { setDialogError(""); setNewCountryName(""); setNewCountryCode(""); setAddCountryOpen(true); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg"
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
                      onClick={() => { setDialogError(""); setNewStateName(""); setNewStateConflictZone(""); setAddStateOpen({ open: true, forOrigin: true }); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg disabled:opacity-30"
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
                      onClick={() => { setDialogError(""); setNewPlaceName(""); setAddPlaceOpen({ open: true, forOrigin: true }); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg disabled:opacity-30"
                      title={t("project.people.addPlace")}
                    >
                      <PlusIcon size={16} />
                    </button>
                  </div>
                </div>

                <p className="text-xs font-medium text-fg-tertiary">
                  {t("project.people.currentPlace")}
                </p>
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
                      onClick={() => { setDialogError(""); setNewCountryName(""); setNewCountryCode(""); setAddCountryOpen(true); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg"
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
                      onClick={() => { setDialogError(""); setNewStateName(""); setNewStateConflictZone(""); setAddStateOpen({ open: true, forOrigin: false }); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg disabled:opacity-30"
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
                      onClick={() => { setDialogError(""); setNewPlaceName(""); setAddPlaceOpen({ open: true, forOrigin: false }); }}
                      className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg disabled:opacity-30"
                      title={t("project.people.addPlace")}
                    >
                      <PlusIcon size={16} />
                    </button>
                  </div>
                </div>

                <h3 className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">{t("project.people.case")}</h3>
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

                  <Field.Root>
                    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                      {t("project.people.externalId")}
                    </Field.Label>
                    <Field.Control
                      value={form.external_id}
                      onChange={(e) => set("external_id", e.target.value)}
                      className={inputClass}
                    />
                  </Field.Root>

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
                      <Field.Root>
                        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                          {t("project.people.consentDate")}
                        </Field.Label>
                        <Field.Control
                          type="date"
                          value={form.consent_date}
                          onChange={(e) => set("consent_date", e.target.value)}
                          className={inputClass}
                        />
                      </Field.Root>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex shrink-0 items-center justify-end gap-2 border-t border-border-secondary px-6 py-4">
                <Drawer.Close className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary shadow-card hover:bg-bg-tertiary">
                  {t("admin.common.cancel")}
                </Drawer.Close>
                <button
                  type="submit"
                  disabled={isPending}
                  className="cursor-pointer rounded-lg bg-accent px-5 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
                >
                  {isPending
                    ? t("project.people.saving")
                    : t("project.people.save")}
                </button>
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
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
