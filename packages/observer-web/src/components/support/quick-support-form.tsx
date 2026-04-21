import { type SyntheticEvent, useState } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { DatePicker } from "@/components/date-picker";
import { WarningIcon } from "@/components/ui/icons";
import { UISelect } from "@/components/ui/ui-select";
import {
  sphereKeys,
  typeKeys,
  SUPPORT_TYPE_VALUES,
  SUPPORT_SPHERE_VALUES,
} from "@/constants/support";
import { useCreateSupportRecord } from "@/hooks/support/use-support-records";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";

import type { SupportSphere, SupportType } from "@/types/support-record";

interface QuickSupportFormProps {
  projectId: string;
  personId: string;
  onSaved: () => void;
  onCancel: () => void;
}

const inputClass =
  "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg";

const typeKeyMap = typeKeys as Record<SupportType, string>;
const sphereKeyMap = sphereKeys as Record<SupportSphere, string>;

export function QuickSupportForm({
  projectId,
  personId,
  onSaved,
  onCancel,
}: QuickSupportFormProps) {
  const { t } = useTranslation();
  const createRecord = useCreateSupportRecord(projectId);
  const toast = useToast();

  const [type, setType] = useState<SupportType>("humanitarian");
  const [sphere, setSphere] = useState("");
  const [providedAt, setProvidedAt] = useState(new Date().toISOString().slice(0, 10));
  const [notes, setNotes] = useState("");
  const [error, setError] = useState("");

  function resetForm() {
    setType("humanitarian");
    setSphere("");
    setProvidedAt(new Date().toISOString().slice(0, 10));
    setNotes("");
    setError("");
  }

  function handleCancel() {
    resetForm();
    onCancel();
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    try {
      await createRecord.mutateAsync({
        person_id: personId,
        type,
        sphere: sphere ? (sphere as SupportSphere) : undefined,
        provided_at: providedAt || undefined,
        notes: notes || undefined,
      });

      toast.success(t("project.people.quickSupportSaved"));
      resetForm();
      onSaved();
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  return (
    <section className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
            <WarningIcon size={16} weight="bold" className="shrink-0" />
            {error}
          </div>
        )}

        <div className="flex items-center gap-3 pt-2">
          <span className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
            {t("project.supportRecords.recordInfo")}
          </span>
          <span className="h-px flex-1 bg-border-secondary" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
              {t("project.supportRecords.type")} *
            </label>
            <UISelect
              value={type}
              onValueChange={(v) => setType(v as SupportType)}
              options={SUPPORT_TYPE_VALUES.map((v) => ({
                value: v,
                label: t(typeKeyMap[v]),
              }))}
              fullWidth
            />
          </div>
          <div>
            <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
              {t("project.supportRecords.sphere")}
            </label>
            <UISelect
              value={sphere}
              onValueChange={setSphere}
              options={SUPPORT_SPHERE_VALUES.map((v) => ({
                value: v,
                label: t(sphereKeyMap[v]),
              }))}
              placeholder={t("project.supportRecords.selectSphere")}
              fullWidth
            />
          </div>
        </div>

        <div className="flex items-center gap-3 pt-2">
          <span className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
            {t("project.supportRecords.notesSection")}
          </span>
          <span className="h-px flex-1 bg-border-secondary" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
              {t("project.supportRecords.providedAt")}
            </label>
            <DatePicker value={providedAt} onChange={setProvidedAt} />
          </div>
        </div>

        <div>
          <label className="mb-1.5 block text-xs font-medium text-fg-secondary">
            {t("project.supportRecords.notes")}
          </label>
          <textarea
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            rows={2}
            className={`${inputClass} h-auto py-2`}
          />
        </div>

        <div className="flex justify-end gap-2 pt-1">
          <Button variant="secondary" onClick={handleCancel}>
            {t("common.cancel")}
          </Button>
          <Button type="submit" disabled={createRecord.isPending}>
            {createRecord.isPending
              ? t("project.supportRecords.saving")
              : t("project.supportRecords.save")}
          </Button>
        </div>
      </form>
    </section>
  );
}
