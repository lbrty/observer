export const SEX_VALUES = ["male", "female", "other", "unknown"] as const;

export const AGE_GROUP_VALUES = [
  "infant",
  "toddler",
  "pre_school",
  "middle_childhood",
  "young_teen",
  "teenager",
  "young_adult",
  "early_adult",
  "middle_aged_adult",
  "old_adult",
] as const;

export const AGE_RANGE_MAP: Record<string, string> = {
  infant: "0-1",
  toddler: "1-3",
  pre_school: "3-6",
  middle_childhood: "6-12",
  young_teen: "12-14",
  teenager: "14-18",
  young_adult: "18-25",
  early_adult: "25-35",
  middle_aged_adult: "35-55",
  old_adult: "55+",
};

export const CASE_STATUS_VALUES = ["new", "active", "closed", "archived"] as const;
