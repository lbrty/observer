import { type SyntheticEvent } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { FormField } from "@/components/forms/form-field";

interface MFAActiveProps {
  showDisable: boolean;
  disableCode: string;
  isPending: boolean;
  onDisableCode: (value: string) => void;
  onShowDisable: () => void;
  onCancelDisable: () => void;
  onSubmit: (e: SyntheticEvent) => void;
}

export function MFAActive({
  showDisable,
  disableCode,
  isPending,
  onDisableCode,
  onShowDisable,
  onCancelDisable,
  onSubmit,
}: MFAActiveProps) {
  const { t } = useTranslation();

  return (
    <div className="rounded-lg border border-border-secondary p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-fg">{t("profile.mfaActive")}</p>
          <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaActiveHint")}</p>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800 dark:bg-green-900/30 dark:text-green-300">
          {t("common.enabled")}
        </span>
      </div>

      {!showDisable ? (
        <Button variant="danger" className="mt-3" onClick={onShowDisable}>
          {t("profile.disableMFA")}
        </Button>
      ) : (
        <form onSubmit={onSubmit} className="mt-3 space-y-3">
          <p className="text-sm text-fg-secondary">{t("profile.disableMFAHint")}</p>
          <FormField
            label={t("auth.totpCode")}
            value={disableCode}
            onChange={onDisableCode}
            maxLength={6}
          />
          <div className="flex gap-2">
            <Button type="submit" variant="danger" disabled={isPending}>
              {t("profile.confirmDisable")}
            </Button>
            <Button type="button" variant="secondary" onClick={onCancelDisable}>
              {t("common.cancel")}
            </Button>
          </div>
        </form>
      )}
    </div>
  );
}
