import { useEffect, useRef, useState } from "react";

import { Command } from "cmdk";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { MagnifyingGlassIcon, PawPrintIcon, UserFocusIcon, FolderSimpleIcon } from "@/components/icons";
import { useSearch } from "@/hooks/use-search";
import type { ProjectGroup } from "@/types/search";

interface SearchPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const PALETTE_LIMIT = 5;

export function SearchPalette({ open, onOpenChange }: SearchPaletteProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [input, setInput] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => setDebouncedQuery(input), 300);
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [input]);

  useEffect(() => {
    if (!open) {
      setInput("");
      setDebouncedQuery("");
    }
  }, [open]);

  const { data } = useSearch(debouncedQuery, PALETTE_LIMIT);
  const results = debouncedQuery.length >= 2 ? (data?.results ?? []) : [];
  const hasResults = results.length > 0;
  const showViewAll = data && (
    results.some(g => g.people.length >= PALETTE_LIMIT || g.pets.length >= PALETTE_LIMIT || g.projects.length >= PALETTE_LIMIT)
  );

  function close() {
    onOpenChange(false);
  }

  function goToViewAll() {
    if (debouncedQuery.length >= 2) {
      navigate({ to: "/search", search: { q: debouncedQuery } });
    }
    close();
  }

  function goToPerson(projectId: string, personId: string) {
    navigate({ to: "/projects/$projectId/people/$personId", params: { projectId, personId } });
    close();
  }

  function goToPet(projectId: string, _petId: string) {
    navigate({ to: "/projects/$projectId/pets/$status", params: { projectId, status: "all" } });
    close();
  }

  function goToProject(projectId: string) {
    navigate({ to: "/projects/$projectId/people", params: { projectId } });
    close();
  }

  return (
    <Command.Dialog
      open={open}
      onOpenChange={onOpenChange}
      label={t("search.title")}
      shouldFilter={false}
      loop
      overlayClassName="fixed inset-0 z-[199] bg-black/10"
      contentClassName="fixed left-1/2 top-[15vh] -translate-x-1/2 z-[200] w-[calc(100%-2rem)] max-w-lg outline-none"
    >
      <div className="overflow-hidden rounded-xl border border-border-secondary bg-bg shadow-elevated">
        <div className="flex items-center gap-3 border-b border-border-secondary px-4 py-3">
          <MagnifyingGlassIcon size={16} className="shrink-0 text-fg-tertiary" />
          <Command.Input
            value={input}
            onValueChange={setInput}
            placeholder={t("search.placeholder")}
            className="flex-1 bg-transparent text-sm text-fg placeholder:text-fg-tertiary outline-none"
          />
          <kbd className="hidden shrink-0 rounded border border-border-secondary px-1.5 py-0.5 font-mono text-[11px] text-fg-tertiary sm:block">
            Esc
          </kbd>
        </div>

        {debouncedQuery.length >= 2 ? (
          <Command.List className="max-h-[60vh] overflow-y-auto p-2">
            {!hasResults && (
              <Command.Empty className="py-8 text-center text-sm text-fg-tertiary">
                {t("search.noResults")}
              </Command.Empty>
            )}

            {results.map((group) => (
              <ProjectGroupSection
                key={group.project_id}
                group={group}
                onPerson={goToPerson}
                onPet={goToPet}
                onProject={goToProject}
                t={t}
              />
            ))}

            {showViewAll && (
              <Command.Item
                onSelect={goToViewAll}
                className="mt-1 flex cursor-pointer items-center gap-2 rounded-lg border border-border-secondary px-3 py-2 text-sm text-accent outline-none aria-selected:bg-bg-tertiary"
              >
                <MagnifyingGlassIcon size={14} />
                {t("search.viewAll")} &ldquo;{debouncedQuery}&rdquo;
              </Command.Item>
            )}
          </Command.List>
        ) : (
          <div className="px-4 py-8 text-center text-sm text-fg-tertiary">
            {t("search.hint")}
          </div>
        )}
      </div>
    </Command.Dialog>
  );
}

interface ProjectGroupSectionProps {
  group: ProjectGroup;
  onPerson: (projectId: string, personId: string) => void;
  onPet: (projectId: string, petId: string) => void;
  onProject: (projectId: string) => void;
  t: (key: string) => string;
}

function ProjectGroupSection({ group, onPerson, onPet, onProject, t }: ProjectGroupSectionProps) {
  return (
    <Command.Group
      heading={group.project_name}
      className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1 [&_[cmdk-group-heading]]:text-[11px] [&_[cmdk-group-heading]]:font-semibold [&_[cmdk-group-heading]]:uppercase [&_[cmdk-group-heading]]:tracking-wide [&_[cmdk-group-heading]]:text-fg-tertiary"
    >
      {group.people.length > 0 && (
        <Command.Group
          heading={t("search.people")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.people.map((p) => (
            <Command.Item
              key={p.id}
              value={`person-${p.id}-${p.first_name} ${p.last_name}`}
              onSelect={() => onPerson(group.project_id, p.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <UserFocusIcon size={14} className="shrink-0 text-fg-tertiary" />
              {p.first_name} {p.last_name}
            </Command.Item>
          ))}
        </Command.Group>
      )}

      {group.pets.length > 0 && (
        <Command.Group
          heading={t("search.pets")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.pets.map((pet) => (
            <Command.Item
              key={pet.id}
              value={`pet-${pet.id}-${pet.name}`}
              onSelect={() => onPet(group.project_id, pet.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <PawPrintIcon size={14} className="shrink-0 text-fg-tertiary" />
              {pet.name}
            </Command.Item>
          ))}
        </Command.Group>
      )}

      {group.projects.length > 0 && (
        <Command.Group
          heading={t("search.projects")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.projects.map((proj) => (
            <Command.Item
              key={proj.id}
              value={`project-${proj.id}-${proj.name}`}
              onSelect={() => onProject(proj.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <FolderSimpleIcon size={14} className="shrink-0 text-fg-tertiary" />
              {proj.name}
            </Command.Item>
          ))}
        </Command.Group>
      )}
    </Command.Group>
  );
}
