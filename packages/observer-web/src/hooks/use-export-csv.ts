import { useState } from "react";

import { api } from "@/lib/api";
import { downloadCSV } from "@/lib/export-csv";

export function useExportCSV() {
  const [exporting, setExporting] = useState(false);

  async function exportCSV(
    apiPath: string,
    filename: string,
    searchParams: Record<string, string> = {},
  ) {
    setExporting(true);
    try {
      const text = await api.get(apiPath, { searchParams }).text();
      const date = new Date().toISOString().slice(0, 10);
      downloadCSV(text, `${filename}-${date}.csv`);
    } finally {
      setExporting(false);
    }
  }

  return { exporting, exportCSV };
}
