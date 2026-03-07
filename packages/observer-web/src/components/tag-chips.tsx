import { useMemo } from "react";

import { useTags } from "@/hooks/use-tags";

export function TagChips({ projectId, tagIds }: { projectId: string; tagIds: string[] }) {
  const { data } = useTags(projectId);

  const tags = useMemo(() => {
    if (!data?.tags || tagIds.length === 0) return [];
    const map = new Map(data.tags.map((t) => [t.id, t]));
    return tagIds.map((id) => map.get(id)).filter(Boolean) as typeof data.tags;
  }, [data?.tags, tagIds]);

  if (tags.length === 0) return null;

  return (
    <div className="flex flex-wrap gap-1">
      {tags.map((tag) => (
        <span
          key={tag.id}
          className="inline-flex items-center rounded px-1.5 py-0.5 text-[11px] font-medium text-white"
          style={{ backgroundColor: tag.color || "#888" }}
        >
          {tag.name}
        </span>
      ))}
    </div>
  );
}
