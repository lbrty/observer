import { useDeferredValue, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { MagnifyingGlassIcon } from "@/components/icons";
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

interface PlaceGroup {
  country: Country;
  state: State;
  places: Place[];
}

export function PlaceCombobox({ onSelect, placeholder }: PlaceComboboxProps) {
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const [activeIndex, setActiveIndex] = useState(-1);
  const listRef = useRef<HTMLDivElement>(null);
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

  const { groups, flat } = useMemo(() => {
    const allPlaces = placesData?.places ?? [];
    const query = deferred.toLowerCase().trim();
    const filtered = query.length < 2
      ? []
      : allPlaces.filter((p) => {
          if (p.name.toLowerCase().includes(query)) return true;
          const state = stateMap.get(p.state_id);
          if (state?.name.toLowerCase().includes(query)) return true;
          const country = state ? countryMap.get(state.country_id) : undefined;
          if (country?.name.toLowerCase().includes(query)) return true;
          return false;
        });

    const groupMap = new Map<string, PlaceGroup>();
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

    const sortedGroups = Array.from(groupMap.values()).sort((a, b) =>
      a.country.name.localeCompare(b.country.name) || a.state.name.localeCompare(b.state.name),
    );

    const flatList: GroupedPlace[] = [];
    for (const g of sortedGroups) {
      for (const place of g.places) {
        flatList.push({ place, state: g.state, country: g.country });
      }
    }

    return { groups: sortedGroups, flat: flatList };
  }, [placesData, stateMap, countryMap, deferred]);

  const showDropdown = search.length >= 2;

  function select(item: GroupedPlace) {
    onSelect(item.place, item.state, item.country);
    setSearch("");
    setActiveIndex(-1);
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (!showDropdown || flat.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => (i < flat.length - 1 ? i + 1 : 0));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => (i > 0 ? i - 1 : flat.length - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (activeIndex >= 0 && activeIndex < flat.length) {
        select(flat[activeIndex]);
      }
    } else if (e.key === "Escape") {
      setSearch("");
      setActiveIndex(-1);
    }
  }

  // Track which group headers we've rendered
  let renderedIndex = 0;

  return (
    <div className="relative">
      <div className="relative">
        <MagnifyingGlassIcon
          size={16}
          className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
        />
        <input
          type="text"
          role="combobox"
          aria-expanded={showDropdown}
          aria-activedescendant={activeIndex >= 0 ? `place-option-${activeIndex}` : undefined}
          aria-autocomplete="list"
          aria-controls="place-listbox"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setActiveIndex(-1);
          }}
          onKeyDown={handleKeyDown}
          placeholder={placeholder ?? t("project.people.searchPlace")}
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </div>
      {showDropdown && (
        <div
          ref={listRef}
          id="place-listbox"
          role="listbox"
          className="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated"
        >
          {flat.length === 0 && (
            <p className="px-3 py-2 text-sm text-fg-tertiary">{t("admin.common.noData")}</p>
          )}
          {groups.map((group) => {
            const header = (
              <div
                key={`header-${group.country.id}-${group.state.id}`}
                className="sticky top-0 border-b border-border-secondary bg-bg-tertiary px-3 py-1.5 text-xs font-medium text-fg-tertiary"
              >
                {group.country.name} — {group.state.name}
              </div>
            );
            const items = group.places.map((place) => {
              const idx = renderedIndex++;
              return (
                <button
                  key={place.id}
                  id={`place-option-${idx}`}
                  role="option"
                  aria-selected={idx === activeIndex}
                  type="button"
                  onClick={() => select({ place, state: group.state, country: group.country })}
                  className={`flex w-full cursor-pointer px-3 py-2 text-left text-sm text-fg ${
                    idx === activeIndex ? "bg-bg-tertiary" : "hover:bg-bg-tertiary"
                  }`}
                >
                  {place.name}
                </button>
              );
            });
            return [header, ...items];
          })}
        </div>
      )}
    </div>
  );
}
