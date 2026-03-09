import { type SyntheticEvent, useEffect, useState } from "react";

import QRCode from "qrcode";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";
import { useDisableMFA, useEnableMFA, useMFASetup } from "@/hooks/use-mfa";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";
import { useToast } from "@/stores/toast";

export function MFASettings() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();
  const toast = useToast();

  const [showSetup, setShowSetup] = useState(false);
  const [qrDataURL, setQrDataURL] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [disableCode, setDisableCode] = useState("");
  const [showDisable, setShowDisable] = useState(false);
  const [error, setError] = useState("");

  const setup = useMFASetup(showSetup);
  const enableMFA = useEnableMFA();
  const disableMFA = useDisableMFA();

  useEffect(() => {
    if (setup.data?.otpauth_url) {
      QRCode.toDataURL(setup.data.otpauth_url, { width: 200 }).then(setQrDataURL);
    }
  }, [setup.data?.otpauth_url]);

  async function handleEnable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    if (!setup.data) return;
    try {
      await enableMFA.mutateAsync({ secret: setup.data.secret, totpCode });
      if (user) setUser({ ...user, mfa_enabled: true });
      toast.success(t("profile.mfaEnabled"));
      setShowSetup(false);
      setTotpCode("");
      setQrDataURL("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleDisable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      await disableMFA.mutateAsync(disableCode);
      if (user) setUser({ ...user, mfa_enabled: false });
      toast.success(t("profile.mfaDisabled"));
      setShowDisable(false);
      setDisableCode("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const isMFAEnabled = user?.mfa_enabled ?? false;

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.twoFactor")}</h2>

      <ErrorBanner message={error} />

      {isMFAEnabled ? (
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
            <Button variant="danger" className="mt-3" onClick={() => setShowDisable(true)}>
              {t("profile.disableMFA")}
            </Button>
          ) : (
            <form onSubmit={handleDisable} className="mt-3 space-y-3">
              <p className="text-sm text-fg-secondary">{t("profile.disableMFAHint")}</p>
              <FormField
                label={t("auth.totpCode")}
                value={disableCode}
                onChange={setDisableCode}
                maxLength={6}
              />
              <div className="flex gap-2">
                <Button type="submit" variant="danger" disabled={disableMFA.isPending}>
                  {t("profile.confirmDisable")}
                </Button>
                <Button type="button" variant="secondary" onClick={() => setShowDisable(false)}>
                  {t("common.cancel")}
                </Button>
              </div>
            </form>
          )}
        </div>
      ) : (
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
            <Button className="mt-3" onClick={() => setShowSetup(true)}>
              {t("profile.setupMFA")}
            </Button>
          ) : (
            <div className="mt-4 space-y-4">
              <p className="text-sm text-fg-secondary">{t("profile.scanQR")}</p>

              {setup.isLoading && <p className="text-sm text-fg-tertiary">{t("common.loading")}</p>}

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
                      {setup.data?.secret}
                    </code>
                  </details>
                </div>
              )}

              {setup.data && (
                <form onSubmit={handleEnable} className="space-y-3">
                  <FormField
                    label={t("profile.enterCode")}
                    value={totpCode}
                    onChange={setTotpCode}
                    maxLength={6}
                    required
                  />
                  <div className="flex gap-2">
                    <Button type="submit" disabled={enableMFA.isPending || totpCode.length !== 6}>
                      {t("profile.verifyAndEnable")}
                    </Button>
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => {
                        setShowSetup(false);
                        setTotpCode("");
                        setQrDataURL("");
                      }}
                    >
                      {t("common.cancel")}
                    </Button>
                  </div>
                </form>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
