export const sexKeys: Record<string, string> = {
  male: "project.people.sexMale",
  female: "project.people.sexFemale",
  other: "project.people.sexOther",
  unknown: "project.people.sexUnknown",
};

export const SEX_VALUES = ["male", "female", "other", "unknown"] as const;

export const ageGroupKeys: Record<string, string> = {
  infant: "project.people.ageInfant",
  toddler: "project.people.ageToddler",
  pre_school: "project.people.agePreSchool",
  middle_childhood: "project.people.ageMiddleChildhood",
  young_teen: "project.people.ageYoungTeen",
  teenager: "project.people.ageTeenager",
  young_adult: "project.people.ageYoungAdult",
  early_adult: "project.people.ageEarlyAdult",
  middle_aged_adult: "project.people.ageMiddleAgedAdult",
  old_adult: "project.people.ageOldAdult",
};

export const AGE_GROUP_VALUES = [
  "infant", "toddler", "pre_school", "middle_childhood", "young_teen",
  "teenager", "young_adult", "early_adult", "middle_aged_adult", "old_adult",
] as const;

export const AGE_RANGE_MAP: Record<string, string> = {
  infant: "0-1", toddler: "1-3", pre_school: "3-6", middle_childhood: "6-12",
  young_teen: "12-14", teenager: "14-18", young_adult: "18-25",
  early_adult: "25-35", middle_aged_adult: "35-55", old_adult: "55+",
};

export const caseStatusKeys: Record<string, string> = {
  new: "project.people.new",
  active: "project.people.active",
  closed: "project.people.closed",
  archived: "project.people.archived",
};

export const CASE_STATUS_VALUES = ["new", "active", "closed", "archived"] as const;
