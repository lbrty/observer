import { type SyntheticEvent, useEffect, useState } from "react";

import QRCode from "qrcode";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/ui/alert-banner";
import { useDisableMFA, useEnableMFA, useMFASetup } from "@/hooks/users/use-mfa";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";
import { useToast } from "@/stores/toast";

import { MFAActive } from "./mfa-active";
import { MFASetup } from "./mfa-setup";

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
        <MFAActive
          showDisable={showDisable}
          disableCode={disableCode}
          isPending={disableMFA.isPending}
          onDisableCode={setDisableCode}
          onShowDisable={() => setShowDisable(true)}
          onCancelDisable={() => setShowDisable(false)}
          onSubmit={handleDisable}
        />
      ) : (
        <MFASetup
          showSetup={showSetup}
          qrDataURL={qrDataURL}
          secret={setup.data?.secret}
          totpCode={totpCode}
          isLoading={setup.isLoading}
          hasData={!!setup.data}
          isPending={enableMFA.isPending}
          onTotpCode={setTotpCode}
          onShowSetup={() => setShowSetup(true)}
          onCancelSetup={() => {
            setShowSetup(false);
            setTotpCode("");
            setQrDataURL("");
          }}
          onSubmit={handleEnable}
        />
      )}
    </div>
  );
}
