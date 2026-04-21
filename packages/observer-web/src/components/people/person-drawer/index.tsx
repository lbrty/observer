import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/ui/alert-banner";
import { DrawerShell } from "@/components/drawer/drawer-shell";
import { SectionHeading } from "@/components/layout/section-heading";
import { TagPicker } from "@/components/tags/tag-picker";

import { CaseSection } from "./case-section";
import { IdentitySection } from "./identity-section";
import { LocationSection } from "./location-section";
import { usePersonDrawerForm } from "./use-person-drawer-form";

interface PersonDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  personId: string | null;
}

export function PersonDrawer({ open, onOpenChange, projectId, personId }: PersonDrawerProps) {
  const { t } = useTranslation();
  const {
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
  } = usePersonDrawerForm({ open, projectId, personId });

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

      <CaseSection
        form={form}
        set={(k, v) => set(k as keyof typeof form, v)}
        officeOptions={officeOptions}
      />

      <SectionHeading>{t("project.tags.title")}</SectionHeading>
      <TagPicker projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
    </DrawerShell>
  );
}
