import type { CountResult, FullReport } from "@/types/report";

function escapeCSV(value: string): string {
  if (value.includes(",") || value.includes('"') || value.includes("\n")) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

function download(content: string, filename: string) {
  const blob = new Blob([content], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

export function exportGroupCSV(title: string, rows: CountResult[]) {
  const lines = ["Label,Count"];
  for (const row of rows) {
    lines.push(`${escapeCSV(row.label)},${row.count}`);
  }
  const slug = title.toLowerCase().replace(/[^a-z0-9]+/g, "-");
  const date = new Date().toISOString().slice(0, 10);
  download(lines.join("\n"), `${slug}-${date}.csv`);
}

export function exportReportCSV(data: FullReport, projectId: string) {
  const rows: string[] = ["Group,Label,Count"];

  const groups: [
    string,
    { rows: { label: string; count: number }[]; total: number },
  ][] = [
    ["Consultations", data.consultations],
    ["By Sex", data.by_sex],
    ["By IDP Status", data.by_idp_status],
    ["By Category", data.by_category],
    ["By Region", data.by_region],
    ["By Sphere", data.by_sphere],
    ["By Office", data.by_office],
    ["By Age Group", data.by_age_group],
    ["By Tag", data.by_tag],
    ["Family Units", data.family_units],
  ];

  for (const [name, group] of groups) {
    for (const row of group.rows) {
      rows.push(`${escapeCSV(name)},${escapeCSV(row.label)},${row.count}`);
    }
    rows.push(`${escapeCSV(name)},TOTAL,${group.total}`);
  }

  const date = new Date().toISOString().slice(0, 10);
  download(rows.join("\n"), `report-${projectId}-${date}.csv`);
}
