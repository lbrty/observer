import type { TFunction } from "i18next";

import { HTTPError } from "@/lib/api";

export async function handleApiError(err: unknown, t: TFunction): Promise<string> {
  if (err instanceof HTTPError) {
    const body = await err.response.json().catch(() => null);
    const code = body?.code;
    const translated = code ? t(code, { defaultValue: "" }) : "";
    return translated || body?.error || err.message;
  }
  return t("common.unexpectedError");
}
