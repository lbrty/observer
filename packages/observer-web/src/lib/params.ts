/** Strips null, undefined, and empty-string values and returns a string record for use as URL search params. */
export function filterParams(params: Record<string, unknown>): Record<string, string> {
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(params)) {
    if (v != null && v !== "") {
      out[k] = String(v);
    }
  }
  return out;
}
