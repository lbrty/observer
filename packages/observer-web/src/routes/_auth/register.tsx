import { Field } from "@base-ui/react/field";
import { createFileRoute, Link } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { HTTPError } from "@/lib/api";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_auth/register")({
  component: RegisterPage,
});

function RegisterPage() {
  const { t } = useTranslation();
  const { register } = useAuth();

  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setSubmitting(true);

    const form = new FormData(e.currentTarget);
    const email = form.get("email") as string;
    const password = form.get("password") as string;
    const confirmPassword = form.get("confirm_password") as string;

    if (password !== confirmPassword) {
      setError(t("auth.passwordsMismatch"));
      setSubmitting(false);
      return;
    }

    if (password.length < 8) {
      setError(t("auth.passwordTooShort"));
      setSubmitting(false);
      return;
    }

    try {
      await register({ email, password, role: "staff" });
      setSuccess(true);
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setError(body?.error ?? err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
    } finally {
      setSubmitting(false);
    }
  }

  if (success) {
    return (
      <>
        <div className="mb-8 flex flex-col items-center">
          <span className="brand-icon mb-4 inline-flex size-14 items-center justify-center rounded-2xl text-xl font-bold text-white">
            O
          </span>
          <h1 className="font-serif text-xl font-semibold text-fg">{t("auth.registerTitle")}</h1>
        </div>
        <div className="mb-4 rounded-lg bg-foam/10 px-3 py-2 text-sm text-foam">
          {t("auth.pendingApproval")}
        </div>
        <Button asChild className="w-full">
          <Link to="/login">{t("auth.login")}</Link>
        </Button>
      </>
    );
  }

  return (
    <>
      <div className="mb-8 flex flex-col items-center">
        <span className="brand-icon mb-4 inline-flex size-14 items-center justify-center rounded-2xl text-xl font-bold text-white">
          O
        </span>
        <h1 className="font-serif text-xl font-semibold text-fg">{t("auth.registerTitle")}</h1>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">{error}</div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
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
            minLength={8}
            autoComplete="new-password"
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
          />
        </Field.Root>

        <Field.Root name="confirm_password">
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("auth.confirmPassword")}
          </Field.Label>
          <Field.Control
            type="password"
            required
            minLength={8}
            autoComplete="new-password"
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
          />
        </Field.Root>

        <Button type="submit" disabled={submitting} className="w-full">
          {submitting ? t("auth.registering") : t("auth.register")}
        </Button>
      </form>

      <p className="mt-5 text-center text-sm text-fg-tertiary">
        {t("auth.hasAccount")}{" "}
        <Link to="/login" className="font-medium text-accent hover:underline">
          {t("auth.login")}
        </Link>
      </p>
    </>
  );
}
