import type { FullReport } from "@/types/report";

function escapeCSV(value: string): string {
  if (value.includes(",") || value.includes('"') || value.includes("\n")) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
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

  const blob = new Blob([rows.join("\n")], {
    type: "text/csv;charset=utf-8;",
  });
  const url = URL.createObjectURL(blob);
  const date = new Date().toISOString().slice(0, 10);
  const a = document.createElement("a");
  a.href = url;
  a.download = `report-${projectId}-${date}.csv`;
  a.click();
  URL.revokeObjectURL(url);
}
