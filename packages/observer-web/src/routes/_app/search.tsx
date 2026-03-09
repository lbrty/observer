import { useEffect, useRef, useState } from "react";

import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { FolderSimpleIcon, MagnifyingGlassIcon, PawPrintIcon, UserFocusIcon } from "@/components/icons";
import { useSearch } from "@/hooks/use-search";

export const Route = createFileRoute("/_app/search")({
  validateSearch: (search: Record<string, unknown>) => ({
    q: (search.q as string) ?? "",
  }),
  component: SearchPage,
});

const FULL_LIMIT = 50;

function SearchPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { q } = Route.useSearch();
  const [input, setInput] = useState(q);
  const [debouncedQuery, setDebouncedQuery] = useState(q);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    setInput(q);
    setDebouncedQuery(q);
  }, [q]);

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => {
      setDebouncedQuery(input);
      if (input.length >= 2) {
        navigate({ to: "/search", search: { q: input }, replace: true });
      }
    }, 300);
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [input, navigate]);

  const { data, isFetching } = useSearch(debouncedQuery, FULL_LIMIT);
  const results = data?.results ?? [];

  return (
    <div className="mx-auto w-full max-w-270 px-10 py-8">
      <div className="mb-8">
        <h1 className="font-serif text-2xl font-bold tracking-tight text-fg">
          {t("search.title")}
        </h1>
      </div>

      <div className="mb-6 flex h-10 w-full max-w-lg items-center gap-2.5 rounded-lg border border-border-secondary bg-bg-secondary px-3">
        <MagnifyingGlassIcon size={16} className="shrink-0 text-fg-tertiary" />
        <input
          autoFocus
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={t("search.placeholder")}
          className="flex-1 bg-transparent text-sm text-fg placeholder:text-fg-tertiary outline-none"
        />
        {isFetching && (
          <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-border-secondary border-t-accent" />
        )}
      </div>

      {debouncedQuery.length >= 2 && results.length === 0 && !isFetching && (
        <p className="text-sm text-fg-tertiary">{t("search.noResults")}</p>
      )}

      <div className="space-y-8">
        {results.map((group) => (
          <div key={group.project_id}>
            <h2 className="mb-3 font-serif text-base font-semibold text-fg">
              {group.project_name}
            </h2>

            {group.people.length > 0 && (
              <section className="mb-4">
                <h3 className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("search.people")}
                </h3>
                <div className="space-y-1">
                  {group.people.map((p) => (
                    <Link
                      key={p.id}
                      to="/projects/$projectId/people/$personId"
                      params={{ projectId: group.project_id, personId: p.id }}
                      className="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-fg transition-colors hover:bg-bg-tertiary"
                    >
                      <UserFocusIcon size={14} className="shrink-0 text-fg-tertiary" />
                      {p.first_name} {p.last_name}
                    </Link>
                  ))}
                </div>
              </section>
            )}

            {group.pets.length > 0 && (
              <section className="mb-4">
                <h3 className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("search.pets")}
                </h3>
                <div className="space-y-1">
                  {group.pets.map((pet) => (
                    <Link
                      key={pet.id}
                      to="/projects/$projectId/pets/$status"
                      params={{ projectId: group.project_id, status: "all" }}
                      className="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-fg transition-colors hover:bg-bg-tertiary"
                    >
                      <PawPrintIcon size={14} className="shrink-0 text-fg-tertiary" />
                      {pet.name}
                    </Link>
                  ))}
                </div>
              </section>
            )}

            {group.projects.length > 0 && (
              <section className="mb-4">
                <h3 className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                  {t("search.projects")}
                </h3>
                <div className="space-y-1">
                  {group.projects.map((proj) => (
                    <Link
                      key={proj.id}
                      to="/projects/$projectId/people"
                      params={{ projectId: proj.id }}
                      className="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-fg transition-colors hover:bg-bg-tertiary"
                    >
                      <FolderSimpleIcon size={14} className="shrink-0 text-fg-tertiary" />
                      {proj.name}
                    </Link>
                  ))}
                </div>
              </section>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
