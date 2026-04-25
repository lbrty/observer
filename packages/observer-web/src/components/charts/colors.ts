export const SEX_COLORS: Record<string, string> = {
  female: "#e05a8a",
  male: "#5b8af8",
  other: "#8b5cf6",
  unknown: "#8b909e",
};

export const SUPPORT_TYPE_COLORS: Record<string, string> = {
  humanitarian: "#e08c1a",
  legal: "#5b8af8",
  social: "#30a46c",
  psychological: "#8b5cf6",
  medical: "#e5534b",
  general: "#60636c",
};

export const SPHERE_COLORS: Record<string, string> = {
  housing_assistance: "#38bdf8",
  document_recovery: "#5b8af8",
  social_benefits: "#30a46c",
  property_rights: "#f59e0b",
  employment_rights: "#ec4899",
  family_law: "#8b5cf6",
  healthcare_access: "#e5534b",
  education_access: "#14b8a6",
  financial_aid: "#e08c1a",
  psychological_support: "#a78bfa",
  other: "#8b909e",
  unspecified: "#8b909e",
};

export const CASE_STATUS_COLORS: Record<string, string> = {
  new: "#5b8af8",
  active: "#30a46c",
  closed: "#8b909e",
  archived: "#60636c",
};

export const IDP_STATUS_COLORS: Record<string, string> = {
  idp: "#e5534b",
  non_idp: "#30a46c",
  unknown: "#8b909e",
};

export const AGE_GROUP_COLORS: Record<string, string> = {
  infant: "#38bdf8",
  toddler: "#0ea5e9",
  pre_school: "#0284c7",
  middle_childhood: "#0369a1",
  young_teen: "#3b82f6",
  teenager: "#2563eb",
  young_adult: "#1d4ed8",
  early_adult: "#6366f1",
  middle_aged_adult: "#4f46e5",
  old_adult: "#4338ca",
};

export const PET_STATUS_COLORS: Record<string, string> = {
  registered: "#5b8af8",
  adopted: "#30a46c",
  owner_found: "#38bdf8",
  needs_shelter: "#e5534b",
  unknown: "#8b909e",
};

export const PET_OWNERSHIP_COLORS: Record<string, string> = {
  with_owner: "#30a46c",
  without_owner: "#f59e0b",
};

export const FALLBACK_PALETTE = [
  "#5b8af8",
  "#30a46c",
  "#e08c1a",
  "#e5534b",
  "#8b5cf6",
  "#ec4899",
  "#14b8a6",
  "#f97316",
  "#38bdf8",
  "#84cc16",
  "#e879f9",
  "#a78bfa",
];

export function getColor(label: string, colorMap?: Record<string, string>, index?: number): string {
  if (colorMap?.[label]) return colorMap[label];
  return FALLBACK_PALETTE[(index ?? 0) % FALLBACK_PALETTE.length];
}

export function getContrastColor(hex: string): string {
  const r = parseInt(hex.slice(1, 3), 16) / 255;
  const g = parseInt(hex.slice(3, 5), 16) / 255;
  const b = parseInt(hex.slice(5, 7), 16) / 255;
  const toLinear = (c: number) => (c <= 0.04045 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4);
  const L = 0.2126 * toLinear(r) + 0.7152 * toLinear(g) + 0.0722 * toLinear(b);
  return L > 0.179 ? "#1c1d22" : "#ffffff";
}
