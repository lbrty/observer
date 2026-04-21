import { Command } from "cmdk";

import { PawPrintIcon, UserFocusIcon, FolderSimpleIcon } from "@/components/ui/icons";
import type { ProjectGroup } from "@/types/search";

export interface ProjectGroupSectionProps {
  group: ProjectGroup;
  onPerson: (projectId: string, personId: string) => void;
  onPet: (projectId: string, petId: string) => void;
  onProject: (projectId: string) => void;
  t: (key: string) => string;
}

export function ProjectGroupSection({
  group,
  onPerson,
  onPet,
  onProject,
  t,
}: ProjectGroupSectionProps) {
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
