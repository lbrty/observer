import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";

import { Button } from "@/components/button";
import { TrashIcon } from "@/components/icons";
import { PersonCombobox } from "@/components/person-combobox";
import { PersonName } from "@/components/person-name";
import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";
import { relationshipKeys } from "@/constants/household";

import type { Household, Relationship } from "@/types/household";

interface MemberForm {
  person_id: string;
  relationship: Relationship;
}

interface MembersSectionProps {
  editingId: string;
  household: Household | undefined;
  projectId: string;
  memberForm: MemberForm;
  memberPersonName: string;
  addMemberPending: boolean;
  onMemberFormChange: (form: MemberForm) => void;
  onMemberPersonNameChange: (name: string) => void;
  onAddMember: () => void;
  onRemoveMember: (personId: string) => void;
  onCloseDrawer: () => void;
}

export function MembersSection({
  editingId,
  household,
  projectId,
  memberForm,
  memberPersonName,
  addMemberPending,
  onMemberFormChange,
  onMemberPersonNameChange,
  onAddMember,
  onRemoveMember,
  onCloseDrawer,
}: MembersSectionProps) {
  const { t } = useTranslation();

  const relationshipOptions = Object.entries(relationshipKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  const memberPersonLabel = memberPersonName || memberForm.person_id;

  return (
    <>
      <SectionHeading>{t("project.households.members")}</SectionHeading>
      {household?.members && household.members.length > 0 ? (
        <div className="overflow-hidden rounded-lg border border-border-secondary">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-secondary bg-bg-secondary">
                <th className="px-3 py-2 text-left font-medium text-fg-secondary">
                  {t("project.households.member")}
                </th>
                <th className="px-3 py-2 text-left font-medium text-fg-secondary">
                  {t("project.households.relationship")}
                </th>
                <th className="w-10 px-3 py-2" />
              </tr>
            </thead>
            <tbody>
              {household.members.map((m) => (
                <tr
                  key={m.person_id}
                  className="border-b border-border-secondary last:border-b-0"
                >
                  <td className="px-3 py-2 text-sm">
                    <Link
                      to="/projects/$projectId/people/$personId"
                      params={{ projectId, personId: m.person_id }}
                      className="text-fg underline-offset-2 hover:underline"
                      onClick={onCloseDrawer}
                    >
                      <PersonName projectId={projectId} personId={m.person_id} />
                    </Link>
                  </td>
                  <td className="px-3 py-2 text-fg-secondary">
                    {t(
                      `project.households.relationship${m.relationship.charAt(0).toUpperCase()}${m.relationship.slice(1).replace(/_(\w)/g, (_, c: string) => c.toUpperCase())}`,
                    )}
                  </td>
                  <td className="px-3 py-2">
                    <Button
                      variant="ghost"
                      className="p-1 hover:text-rose"
                      onClick={() => onRemoveMember(m.person_id)}
                    >
                      <TrashIcon size={14} />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p className="text-sm text-fg-tertiary">
          {t("project.households.noMembers")}
        </p>
      )}

      <p className="text-xs font-medium text-fg-tertiary">
        {t("project.households.addMember")}
      </p>
      <div className="flex items-end gap-3">
        <div className="flex-1">
          <label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.households.member")}
          </label>
          {memberForm.person_id ? (
            <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
              <span className="flex-1 truncate text-sm text-fg">
                {memberPersonLabel}
              </span>
              <button
                type="button"
                onClick={() => {
                  onMemberFormChange({ person_id: "", relationship: "head" });
                  onMemberPersonNameChange("");
                }}
                className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
              >
                ×
              </button>
            </div>
          ) : (
            <PersonCombobox
              projectId={projectId}
              excludeIds={household?.members?.map((m) => m.person_id) ?? []}
              onSelect={(p) => {
                onMemberFormChange({ ...memberForm, person_id: p.id });
                onMemberPersonNameChange(`${p.first_name} ${p.last_name ?? ""}`.trim());
              }}
            />
          )}
        </div>

        <div className="flex-1">
          <label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.households.relationship")}
          </label>
          <UISelect
            value={memberForm.relationship}
            onValueChange={(v) =>
              onMemberFormChange({
                ...memberForm,
                relationship: v as Relationship,
              })
            }
            options={relationshipOptions}
            fullWidth
          />
        </div>

        <Button
          onClick={() => {
            onAddMember();
            onMemberPersonNameChange("");
          }}
          disabled={addMemberPending || !memberForm.person_id}
        >
          {t("project.households.addMember")}
        </Button>
      </div>
    </>
  );
}
