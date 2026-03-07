import { useMemo, useState } from "react";

import { useTranslation } from "react-i18next";

import { TagIcon, XIcon } from "@/components/icons";
import { useTags } from "@/hooks/use-tags";

import type { Tag } from "@/types/tag";

interface TagPickerProps {
  projectId: string;
  selectedIds: string[];
  onChange: (ids: string[]) => void;
}

export function TagPicker({ projectId, selectedIds, onChange }: TagPickerProps) {
  const { t } = useTranslation();
  const { data } = useTags(projectId);
  const [search, setSearch] = useState("");

  const allTags = data?.tags ?? [];
  const tagMap = useMemo(() => new Map(allTags.map((tag) => [tag.id, tag])), [allTags]);

  const selectedTags = selectedIds.map((id) => tagMap.get(id)).filter(Boolean) as Tag[];

  const filtered = allTags.filter(
    (tag) =>
      !selectedIds.includes(tag.id) && tag.name.toLowerCase().includes(search.toLowerCase()),
  );

  function add(tag: Tag) {
    onChange([...selectedIds, tag.id]);
    setSearch("");
  }

  function remove(id: string) {
    onChange(selectedIds.filter((v) => v !== id));
  }

  return (
    <div className="space-y-2">
      {selectedTags.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {selectedTags.map((tag) => (
            <span
              key={tag.id}
              className="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium text-white"
              style={{ backgroundColor: tag.color || "#888" }}
            >
              <TagIcon size={12} />
              {tag.name}
              <button
                type="button"
                onClick={() => remove(tag.id)}
                className="ml-0.5 cursor-pointer rounded-sm p-0.5 hover:bg-black/20"
              >
                <XIcon size={10} />
              </button>
            </span>
          ))}
        </div>
      )}

      <div className="relative">
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t("project.tags.searchTags")}
          className="block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />

        {search.length > 0 && filtered.length > 0 && (
          <div className="absolute z-10 mt-1 max-h-40 w-full overflow-auto rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated">
            {filtered.map((tag) => (
              <button
                key={tag.id}
                type="button"
                onClick={() => add(tag)}
                className="flex w-full cursor-pointer items-center gap-2 px-3 py-2 text-left hover:bg-bg-tertiary"
              >
                <span
                  className="inline-flex size-5 shrink-0 items-center justify-center rounded"
                  style={{ backgroundColor: tag.color || "#888" }}
                >
                  <TagIcon size={10} className="text-white" />
                </span>
                <span className="text-sm text-fg">{tag.name}</span>
              </button>
            ))}
          </div>
        )}

        {search.length > 0 && filtered.length === 0 && (
          <div className="absolute z-10 mt-1 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 shadow-elevated">
            <p className="text-sm text-fg-tertiary">{t("admin.common.noData")}</p>
          </div>
        )}
      </div>
    </div>
  );
}
