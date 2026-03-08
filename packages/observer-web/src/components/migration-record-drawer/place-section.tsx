import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";

interface SelectOption {
  label: string;
  value: string;
}

interface PlaceSectionProps {
  title: string;
  country: string;
  state: string;
  place: string;
  countryOptions: SelectOption[];
  stateOptions: SelectOption[];
  placeOptions: SelectOption[];
  countryPlaceholder: string;
  statePlaceholder: string;
  placePlaceholder: string;
  onCountryChange: (v: string) => void;
  onStateChange: (v: string) => void;
  onPlaceChange: (v: string) => void;
}

export function PlaceSection({
  title,
  country,
  state,
  place,
  countryOptions,
  stateOptions,
  placeOptions,
  countryPlaceholder,
  statePlaceholder,
  placePlaceholder,
  onCountryChange,
  onStateChange,
  onPlaceChange,
}: PlaceSectionProps) {
  return (
    <>
      <SectionHeading>{title}</SectionHeading>
      <div className="grid grid-cols-3 gap-3">
        <UISelect
          value={country}
          onValueChange={onCountryChange}
          options={countryOptions}
          placeholder={countryPlaceholder}
          fullWidth
        />
        <UISelect
          value={state}
          onValueChange={onStateChange}
          options={stateOptions}
          placeholder={statePlaceholder}
          disabled={!country}
          fullWidth
        />
        <UISelect
          value={place}
          onValueChange={onPlaceChange}
          options={placeOptions}
          placeholder={placePlaceholder}
          disabled={!state}
          fullWidth
        />
      </div>
    </>
  );
}
