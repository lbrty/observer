export const SEX_COLORS: Record<string, string> = {
  female: "#e05a8a",
  male: "#5a8ae0",
  other: "#8b5cf6",
  unknown: "#94a3b8",
};

export const SUPPORT_TYPE_COLORS: Record<string, string> = {
  humanitarian: "#d97706",
  legal: "#3b82f6",
  social: "#10b981",
  psychological: "#8b5cf6",
  medical: "#ef4444",
  general: "#64748b",
};

export const SPHERE_COLORS: Record<string, string> = {
  housing_assistance: "#0ea5e9",
  document_recovery: "#6366f1",
  social_benefits: "#10b981",
  property_rights: "#f59e0b",
  employment_rights: "#ec4899",
  family_law: "#8b5cf6",
  healthcare_access: "#ef4444",
  education_access: "#14b8a6",
  financial_aid: "#d97706",
  psychological_support: "#a78bfa",
  other: "#94a3b8",
  unspecified: "#94a3b8",
};

export const CASE_STATUS_COLORS: Record<string, string> = {
  new: "#6366f1",
  active: "#10b981",
  closed: "#94a3b8",
  archived: "#64748b",
};

export const IDP_STATUS_COLORS: Record<string, string> = {
  idp: "#ef4444",
  non_idp: "#10b981",
  unknown: "#94a3b8",
};

export const AGE_GROUP_COLORS: Record<string, string> = {
  infant: "#fef3c7",
  toddler: "#fde68a",
  pre_school: "#fcd34d",
  middle_childhood: "#fbbf24",
  young_teen: "#f59e0b",
  teenager: "#d97706",
  young_adult: "#b45309",
  early_adult: "#92400e",
  middle_aged_adult: "#78350f",
  old_adult: "#451a03",
};

export const PET_STATUS_COLORS: Record<string, string> = {
  registered: "#6366f1",
  adopted: "#10b981",
  owner_found: "#3b82f6",
  needs_shelter: "#ef4444",
  unknown: "#94a3b8",
};

export const PET_OWNERSHIP_COLORS: Record<string, string> = {
  with_owner: "#10b981",
  without_owner: "#f59e0b",
};

export const FALLBACK_PALETTE = [
  "#6366f1", "#f59e0b", "#10b981", "#ef4444",
  "#8b5cf6", "#ec4899", "#14b8a6", "#f97316",
  "#3b82f6", "#84cc16", "#e879f9", "#06b6d4",
];

export function getColor(
  label: string,
  colorMap?: Record<string, string>,
  index?: number,
): string {
  if (colorMap?.[label]) return colorMap[label];
  return FALLBACK_PALETTE[(index ?? 0) % FALLBACK_PALETTE.length];
}
