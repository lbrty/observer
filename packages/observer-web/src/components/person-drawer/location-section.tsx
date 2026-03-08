import { useTranslation } from "react-i18next";

import { FormSection } from "@/components/form-section";
import { PlaceCombobox } from "@/components/place-combobox";

import type { Country, Place, State } from "@/types/reference";

interface LocationSectionProps {
  originPlaceId: string;
  currentPlaceId: string;
  originPlaceLabel: string;
  currentPlaceLabel: string;
  onSelectOrigin: (place: Place, state: State, country: Country) => void;
  onClearOrigin: () => void;
  onSelectCurrent: (place: Place, state: State, country: Country) => void;
  onClearCurrent: () => void;
}

export function LocationSection({
  originPlaceId,
  currentPlaceId,
  originPlaceLabel,
  currentPlaceLabel,
  onSelectOrigin,
  onClearOrigin,
  onSelectCurrent,
  onClearCurrent,
}: LocationSectionProps) {
  const { t } = useTranslation();

  return (
    <FormSection title={t("project.people.location")}>
      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.originPlace")}
        </span>
        {originPlaceId ? (
          <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
            <span className="flex-1 truncate text-sm text-fg">
              {originPlaceLabel || originPlaceId}
            </span>
            <button
              type="button"
              onClick={onClearOrigin}
              className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
            >
              ×
            </button>
          </div>
        ) : (
          <PlaceCombobox onSelect={onSelectOrigin} />
        )}
      </div>

      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.currentPlace")}
        </span>
        {currentPlaceId ? (
          <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
            <span className="flex-1 truncate text-sm text-fg">
              {currentPlaceLabel || currentPlaceId}
            </span>
            <button
              type="button"
              onClick={onClearCurrent}
              className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
            >
              ×
            </button>
          </div>
        ) : (
          <PlaceCombobox onSelect={onSelectCurrent} />
        )}
      </div>
    </FormSection>
  );
}
