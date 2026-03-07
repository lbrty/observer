import { useMemo, useRef, useState } from "react";

import { useTranslation } from "react-i18next";

import { FunnelIcon, TagIcon, XIcon } from "@/components/icons";
import { useTags } from "@/hooks/use-tags";

import type { Tag } from "@/types/tag";

interface TagFilterProps {
  projectId: string;
  selectedIds: string[];
  onChange: (ids: string[]) => void;
}

export function TagFilter({ projectId, selectedIds, onChange }: TagFilterProps) {
  const { t } = useTranslation();
  const { data } = useTags(projectId);
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  const allTags = data?.tags ?? [];
  const tagMap = useMemo(() => new Map(allTags.map((tag) => [tag.id, tag])), [allTags]);

  const selectedTags = selectedIds.map((id) => tagMap.get(id)).filter(Boolean) as Tag[];

  const filtered = allTags.filter(
    (tag) => !selectedIds.includes(tag.id) && tag.name.toLowerCase().includes(search.toLowerCase()),
  );

  function toggle(tag: Tag) {
    if (selectedIds.includes(tag.id)) {
      onChange(selectedIds.filter((v) => v !== tag.id));
    } else {
      onChange([...selectedIds, tag.id]);
    }
    setSearch("");
  }

  function remove(id: string) {
    onChange(selectedIds.filter((v) => v !== id));
  }

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg-secondary transition-colors hover:bg-bg-tertiary hover:text-fg"
      >
        <FunnelIcon size={14} />
        {t("project.tags.filterByTags")}
        {selectedIds.length > 0 && (
          <span className="ml-0.5 inline-flex size-5 items-center justify-center rounded-full bg-accent text-xs font-medium text-white">
            {selectedIds.length}
          </span>
        )}
      </button>

      {selectedTags.length > 0 && (
        <div className="ml-2 inline-flex flex-wrap gap-1.5">
          {selectedTags.map((tag) => (
            <span
              key={tag.id}
              className="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium text-white"
              style={{ backgroundColor: tag.color || "#888" }}
            >
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

      {open && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} onKeyDown={() => {}} />
          <div className="absolute top-full left-0 z-20 mt-1 w-64 rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated">
            <div className="p-2">
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder={t("project.tags.searchTags")}
                className="block h-8 w-full rounded-md border border-border-secondary bg-bg px-2.5 text-sm text-fg outline-none focus:border-accent"
                autoFocus
              />
            </div>
            <div className="max-h-48 overflow-auto px-1 pb-1">
              {filtered.length === 0 && (
                <p className="px-2 py-1.5 text-sm text-fg-tertiary">{t("admin.common.noData")}</p>
              )}
              {filtered.map((tag) => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => toggle(tag)}
                  className="flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-left hover:bg-bg-tertiary"
                >
                  <span
                    className="inline-flex size-4 shrink-0 items-center justify-center rounded"
                    style={{ backgroundColor: tag.color || "#888" }}
                  >
                    <TagIcon size={10} className="text-white" />
                  </span>
                  <span className="text-sm text-fg">{tag.name}</span>
                </button>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
