/** Converts any array of items with `id` and `name` fields into select option objects. */
export function toSelectOptions<T extends { id: string; name: string }>(
  items: T[] | undefined | null,
): Array<{ label: string; value: string }> {
  return (items ?? []).map((item) => ({ label: item.name, value: item.id }));
}
