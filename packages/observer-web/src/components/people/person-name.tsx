import { usePerson } from "@/hooks/use-people";

interface PersonNameProps {
  projectId: string;
  personId: string;
}

export function PersonName({ projectId, personId }: PersonNameProps) {
  const { data: person } = usePerson(projectId, personId);
  if (!person) return <span className="text-fg-tertiary">…</span>;
  return <>{`${person.first_name} ${person.last_name ?? ""}`.trim()}</>;
}
