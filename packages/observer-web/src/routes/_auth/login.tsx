import { Field } from "@base-ui/react/field";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { type SyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { HTTPError } from "@/lib/api";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_auth/login")({
  component: LoginPage,
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
async function parseLoginError(err: unknown, t: (k: string, opts?: any) => string): Promise<string> {
  if (err instanceof HTTPError) {
    const body = await err.response.json().catch(() => null);
    if (body?.code === "errors.user.notActive") return t("auth.pendingApproval");
    const translated = body?.code ? t(body.code, { defaultValue: "" }) : "";
    return translated || body?.error || (err as Error).message;
  }
  return t("common.unexpectedError");
}

function LoginPage() {
  const { t } = useTranslation();
  const { login, verifyMFA } = useAuth();
  const navigate = useNavigate();

  const [step, setStep] = useState<"credentials" | "mfa">("credentials");
  const [mfaToken, setMfaToken] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleCredentials(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setSubmitting(true);
    const form = new FormData(e.currentTarget);
    try {
      const result = await login({
        email: form.get("email") as string,
        password: form.get("password") as string,
      });
      if (result.requires_mfa && result.mfa_token) {
        setMfaToken(result.mfa_token);
        setStep("mfa");
      } else {
        navigate({ to: "/" });
      }
    } catch (err) {
      setError(await parseLoginError(err, t));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleMFA(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setSubmitting(true);
    const form = new FormData(e.currentTarget);
    try {
      await verifyMFA(mfaToken, form.get("totp_code") as string);
      navigate({ to: "/" });
    } catch (err) {
      setError(await parseLoginError(err, t));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <>
      <div className="mb-8 flex flex-col items-center">
        <span className="brand-icon mb-4 inline-flex size-14 items-center justify-center rounded-2xl text-xl font-bold text-white">
          O
        </span>
        <h1 className="font-serif text-xl font-semibold text-fg">
          {step === "mfa" ? t("auth.mfaTitle") : t("auth.loginTitle")}
        </h1>
        {step === "mfa" && (
          <p className="mt-1 text-sm text-fg-secondary">{t("auth.mfaHint")}</p>
        )}
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">{error}</div>
      )}

      {step === "credentials" ? (
        <form onSubmit={handleCredentials} className="space-y-4">
          <Field.Root name="email">
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("common.email")}
            </Field.Label>
            <Field.Control
              type="email"
              required
              autoComplete="email"
              className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
            />
          </Field.Root>

          <Field.Root name="password">
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("auth.password")}
            </Field.Label>
            <Field.Control
              type="password"
              required
              autoComplete="current-password"
              className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
            />
          </Field.Root>

          <Button type="submit" disabled={submitting} className="w-full">
            {submitting ? t("auth.loggingIn") : t("auth.login")}
          </Button>
        </form>
      ) : (
        <form onSubmit={handleMFA} className="space-y-4">
          <Field.Root name="totp_code">
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("auth.totpCode")}
            </Field.Label>
            <Field.Control
              type="text"
              inputMode="numeric"
              pattern="[0-9]{6}"
              maxLength={6}
              required
              autoComplete="one-time-code"
              autoFocus
              className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-center text-lg tracking-widest text-fg outline-none focus:border-accent"
            />
          </Field.Root>

          <Button type="submit" disabled={submitting} className="w-full">
            {submitting ? t("auth.verifying") : t("auth.verifyMFA")}
          </Button>

          <button
            type="button"
            className="w-full text-sm text-fg-tertiary hover:text-fg"
            onClick={() => {
              setStep("credentials");
              setError("");
            }}
          >
            {t("auth.backToLogin")}
          </button>
        </form>
      )}

      {step === "credentials" && (
        <p className="mt-5 text-center text-sm text-fg-tertiary">
          {t("auth.noAccount")}{" "}
          <Link to="/register" className="font-medium text-accent hover:underline">
            {t("auth.register")}
          </Link>
        </p>
      )}
    </>
  );
}
