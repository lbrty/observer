import { Field } from "@base-ui/react/field";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { HTTPError } from "@/lib/api";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_auth/login")({
  component: LoginPage,
});

function LoginPage() {
  const { t } = useTranslation();
  const { login } = useAuth();
  const navigate = useNavigate();

  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setSubmitting(true);

    const form = new FormData(e.currentTarget);
    const email = form.get("email") as string;
    const password = form.get("password") as string;

    try {
      const result = await login({ email, password });

      if (result.requires_mfa) {
        setError(t("auth.mfaRequired"));
        return;
      }

      navigate({ to: "/" });
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

  return (
    <>
      <div className="mb-8 flex flex-col items-center">
        <span className="brand-icon mb-4 inline-flex size-14 items-center justify-center rounded-2xl text-xl font-bold text-white">
          O
        </span>
        <h1 className="font-serif text-xl font-semibold text-fg">
          {t("auth.loginTitle")}
        </h1>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">
          {error}
        </div>
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
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
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
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <button
          type="submit"
          disabled={submitting}
          className="w-full cursor-pointer rounded-lg bg-accent px-3 py-2.5 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {submitting ? t("auth.loggingIn") : t("auth.login")}
        </button>
      </form>

      <p className="mt-5 text-center text-sm text-fg-tertiary">
        {t("auth.noAccount")}{" "}
        <Link
          to="/register"
          className="font-medium text-accent hover:underline"
        >
          {t("auth.register")}
        </Link>
      </p>
    </>
  );
}
