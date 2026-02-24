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
      <h1 className="mb-6 text-center text-2xl font-semibold">
        {t("auth.loginTitle")}
      </h1>

      {error && (
        <div className="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <Field.Root name="email">
          <Field.Label className="mb-1 block text-sm font-medium text-zinc-700">
            {t("common.email")}
          </Field.Label>
          <Field.Control
            type="email"
            required
            autoComplete="email"
            className="block w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
          />
        </Field.Root>

        <Field.Root name="password">
          <Field.Label className="mb-1 block text-sm font-medium text-zinc-700">
            {t("auth.password")}
          </Field.Label>
          <Field.Control
            type="password"
            required
            autoComplete="current-password"
            className="block w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
          />
        </Field.Root>

        <button
          type="submit"
          disabled={submitting}
          className="w-full cursor-pointer rounded-md bg-zinc-900 px-3 py-2 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {submitting ? t("auth.loggingIn") : t("auth.login")}
        </button>
      </form>

      <p className="mt-4 text-center text-sm text-zinc-500">
        {t("auth.noAccount")}{" "}
        <Link to="/register" className="text-zinc-900 underline">
          {t("auth.register")}
        </Link>
      </p>
    </>
  );
}
