import { Field } from "@base-ui/react/field";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { HTTPError } from "@/lib/api";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_auth/register")({
  component: RegisterPage,
});

function RegisterPage() {
  const { t } = useTranslation();
  const { register } = useAuth();
  const navigate = useNavigate();

  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setSubmitting(true);

    const form = new FormData(e.currentTarget);
    const email = form.get("email") as string;
    const phone = form.get("phone") as string;
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
      await register({ email, phone, password, role: "staff" });
      navigate({ to: "/login" });
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
        {t("auth.registerTitle")}
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

        <Field.Root name="phone">
          <Field.Label className="mb-1 block text-sm font-medium text-zinc-700">
            {t("auth.phone")}
          </Field.Label>
          <Field.Control
            type="tel"
            required
            autoComplete="tel"
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
            minLength={8}
            autoComplete="new-password"
            className="block w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
          />
        </Field.Root>

        <Field.Root name="confirm_password">
          <Field.Label className="mb-1 block text-sm font-medium text-zinc-700">
            {t("auth.confirmPassword")}
          </Field.Label>
          <Field.Control
            type="password"
            required
            minLength={8}
            autoComplete="new-password"
            className="block w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
          />
        </Field.Root>

        <button
          type="submit"
          disabled={submitting}
          className="w-full cursor-pointer rounded-md bg-zinc-900 px-3 py-2 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {submitting ? t("auth.registering") : t("auth.register")}
        </button>
      </form>

      <p className="mt-4 text-center text-sm text-zinc-500">
        {t("auth.hasAccount")}{" "}
        <Link to="/login" className="text-zinc-900 underline">
          {t("auth.login")}
        </Link>
      </p>
    </>
  );
}
