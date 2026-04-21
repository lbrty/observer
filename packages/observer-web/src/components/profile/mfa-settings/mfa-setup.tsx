import { type SyntheticEvent } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { FormField } from "@/components/forms/form-field";

interface MFASetupProps {
  showSetup: boolean;
  qrDataURL: string;
  secret: string | undefined;
  totpCode: string;
  isLoading: boolean;
  hasData: boolean;
  isPending: boolean;
  onTotpCode: (value: string) => void;
  onShowSetup: () => void;
  onCancelSetup: () => void;
  onSubmit: (e: SyntheticEvent) => void;
}

export function MFASetup({
  showSetup,
  qrDataURL,
  secret,
  totpCode,
  isLoading,
  hasData,
  isPending,
  onTotpCode,
  onShowSetup,
  onCancelSetup,
  onSubmit,
}: MFASetupProps) {
  const { t } = useTranslation();

  return (
    <div className="rounded-lg border border-border-secondary p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-fg">{t("profile.mfaInactive")}</p>
          <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaInactiveHint")}</p>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-bg-secondary px-2.5 py-0.5 text-xs font-medium text-fg-tertiary">
          {t("common.disabled")}
        </span>
      </div>

      {!showSetup ? (
        <Button className="mt-3" onClick={onShowSetup}>
          {t("profile.setupMFA")}
        </Button>
      ) : (
        <div className="mt-4 space-y-4">
          <p className="text-sm text-fg-secondary">{t("profile.scanQR")}</p>

          {isLoading && <p className="text-sm text-fg-tertiary">{t("common.loading")}</p>}

          {qrDataURL && (
            <div className="flex flex-col items-start gap-3">
              <img
                src={qrDataURL}
                alt={t("profile.totpQrCodeAlt")}
                className="rounded-lg border border-border-secondary"
                width={180}
                height={180}
              />
              <details className="text-xs text-fg-tertiary">
                <summary className="cursor-pointer">{t("profile.cantScanQR")}</summary>
                <code className="mt-1 block break-all font-mono text-xs text-fg-secondary select-all">
                  {secret}
                </code>
              </details>
            </div>
          )}

          {hasData && (
            <form onSubmit={onSubmit} className="space-y-3">
              <FormField
                label={t("profile.enterCode")}
                value={totpCode}
                onChange={onTotpCode}
                maxLength={6}
                required
              />
              <div className="flex gap-2">
                <Button type="submit" disabled={isPending || totpCode.length !== 6}>
                  {t("profile.verifyAndEnable")}
                </Button>
                <Button type="button" variant="secondary" onClick={onCancelSetup}>
                  {t("common.cancel")}
                </Button>
              </div>
            </form>
          )}
        </div>
      )}
    </div>
  );
}
