import { useEffect, useState } from "react";

interface UseDrawerFormOptions<T extends Record<string, unknown>, D> {
  initial: T;
  open: boolean;
  isEdit: boolean;
  data?: D | null;
  mapData: (data: D) => T;
}

export function useDrawerForm<T extends Record<string, unknown>, D extends { id: string }>({
  initial,
  open,
  isEdit,
  data,
  mapData,
}: UseDrawerFormOptions<T, D>) {
  const [form, setForm] = useState<T>(initial);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setForm(initial);
      setSaved(false);
      setError("");
      setEditingId(null);
      return;
    }
    if (isEdit && data) {
      setForm(mapData(data));
      setEditingId(data.id);
    }
  }, [open, isEdit, data]);

  function set<K extends keyof T>(key: K, value: T[K]) {
    setForm((f) => ({ ...f, [key]: value }));
    setSaved(false);
    setError("");
  }

  return { form, setForm, set, saved, setSaved, error, setError, editingId, setEditingId };
}
