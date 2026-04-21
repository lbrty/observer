import { useDeferredValue, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { BaseCombobox } from "@/components/base-combobox";
import { useCountries } from "@/hooks/use-countries";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";

import type { Country, Place, State } from "@/types/reference";

interface PlaceComboboxProps {
  onSelect: (place: Place, state: State, country: Country) => void;
  placeholder?: string;
}

interface GroupedPlace {
  place: Place;
  state: State;
  country: Country;
}

export function PlaceCombobox({ onSelect, placeholder }: PlaceComboboxProps) {
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const deferred = useDeferredValue(search);

  const { data: countries } = useCountries();
  const { data: statesData } = useStates();
  const { data: placesData } = usePlaces();

  const stateMap = useMemo(() => {
    const map = new Map<string, State>();
    for (const s of statesData?.states ?? []) map.set(s.id, s);
    return map;
  }, [statesData]);

  const countryMap = useMemo(() => {
    const map = new Map<string, Country>();
    for (const c of countries ?? []) map.set(c.id, c);
    return map;
  }, [countries]);

  const flat = useMemo(() => {
    const allPlaces = placesData?.places ?? [];
    const query = deferred.toLowerCase().trim();
    const filtered =
      query.length < 2
        ? []
        : allPlaces.filter((p) => {
            if (p.name.toLowerCase().includes(query)) return true;
            const state = stateMap.get(p.state_id);
            if (state?.name.toLowerCase().includes(query)) return true;
            const country = state ? countryMap.get(state.country_id) : undefined;
            if (country?.name.toLowerCase().includes(query)) return true;
            return false;
          });

    const groupMap = new Map<string, { country: Country; state: State; places: Place[] }>();
    for (const place of filtered) {
      const state = stateMap.get(place.state_id);
      if (!state) continue;
      const country = countryMap.get(state.country_id);
      if (!country) continue;
      const key = `${country.id}:${state.id}`;
      let group = groupMap.get(key);
      if (!group) {
        group = { country, state, places: [] };
        groupMap.set(key, group);
      }
      group.places.push(place);
    }

    const sorted = Array.from(groupMap.values()).sort(
      (a, b) =>
        a.country.name.localeCompare(b.country.name) || a.state.name.localeCompare(b.state.name),
    );

    const list: GroupedPlace[] = [];
    for (const g of sorted) {
      for (const place of g.places) {
        list.push({ place, state: g.state, country: g.country });
      }
    }
    return list;
  }, [placesData, stateMap, countryMap, deferred]);

  function handleSelect(item: GroupedPlace) {
    onSelect(item.place, item.state, item.country);
  }

  function renderGroupHeader(item: GroupedPlace, index: number) {
    const prev = index > 0 ? flat[index - 1] : undefined;
    const needsHeader =
      !prev || prev.country.id !== item.country.id || prev.state.id !== item.state.id;

    if (!needsHeader) return null;

    return (
      <div
        key={`header-${item.country.id}-${item.state.id}`}
        className="sticky top-0 border-b border-border-secondary bg-bg-tertiary px-3 py-1.5 text-xs font-medium text-fg-tertiary"
      >
        {item.country.name} — {item.state.name}
      </div>
    );
  }

  return (
    <BaseCombobox
      items={flat}
      onSelect={handleSelect}
      getItemKey={(item) => item.place.id}
      search={search}
      onSearchChange={setSearch}
      placeholder={placeholder ?? t("project.people.searchPlace")}
      noDataLabel={t("admin.common.noData")}
      listboxId="place-listbox"
      optionIdPrefix="place-option"
      listClassName="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated"
      renderGroupHeader={renderGroupHeader}
      renderItem={(item) => <span className="text-sm text-fg">{item.place.name}</span>}
    />
  );
}
