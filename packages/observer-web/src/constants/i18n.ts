export const sexKeys: Record<string, string> = {
  male: "project.people.sexMale",
  female: "project.people.sexFemale",
  other: "project.people.sexOther",
  unknown: "project.people.sexUnknown",
};

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

export const caseStatusKeys: Record<string, string> = {
  new: "project.people.new",
  active: "project.people.active",
  closed: "project.people.closed",
  archived: "project.people.archived",
};

export const sphereKeys: Record<string, string> = {
  housing_assistance: "project.supportRecords.sphereHousing",
  document_recovery: "project.supportRecords.sphereDocumentRecovery",
  social_benefits: "project.supportRecords.sphereSocialBenefits",
  property_rights: "project.supportRecords.spherePropertyRights",
  employment_rights: "project.supportRecords.sphereEmploymentRights",
  family_law: "project.supportRecords.sphereFamilyLaw",
  healthcare_access: "project.supportRecords.sphereHealthcareAccess",
  education_access: "project.supportRecords.sphereEducationAccess",
  financial_aid: "project.supportRecords.sphereFinancialAid",
  psychological_support: "project.supportRecords.spherePsychologicalSupport",
  other: "project.supportRecords.sphereOther",
};

export const typeKeys: Record<string, string> = {
  humanitarian: "project.supportRecords.typeHumanitarian",
  legal: "project.supportRecords.typeLegal",
  social: "project.supportRecords.typeSocial",
  psychological: "project.supportRecords.typePsychological",
  medical: "project.supportRecords.typeMedical",
  general: "project.supportRecords.typeGeneral",
};

export const referralKeys: Record<string, string> = {
  pending: "project.supportRecords.referralPending",
  accepted: "project.supportRecords.referralAccepted",
  completed: "project.supportRecords.referralCompleted",
  declined: "project.supportRecords.referralDeclined",
  no_response: "project.supportRecords.referralNoResponse",
};

export const roleKeys: Record<string, string> = {
  admin: "admin.users.roleAdmin",
  staff: "admin.users.roleStaff",
  consultant: "admin.users.roleConsultant",
  guest: "admin.users.roleGuest",
};

export const projectRoleKeys: Record<string, string> = {
  owner: "admin.permissions.roleOwner",
  manager: "admin.permissions.roleManager",
  consultant: "admin.permissions.roleConsultant",
  viewer: "admin.permissions.roleViewer",
};

export const relationshipKeys: Record<string, string> = {
  head: "project.households.relationshipHead",
  spouse: "project.households.relationshipSpouse",
  child: "project.households.relationshipChild",
  parent: "project.households.relationshipParent",
  sibling: "project.households.relationshipSibling",
  grandchild: "project.households.relationshipGrandchild",
  grandparent: "project.households.relationshipGrandparent",
  other_relative: "project.households.relationshipOtherRelative",
  non_relative: "project.households.relationshipNonRelative",
};

export const reasonKeys: Record<string, string> = {
  conflict: "project.migrationRecords.reasonConflict",
  security: "project.migrationRecords.reasonSecurity",
  service_access: "project.migrationRecords.reasonService",
  return: "project.migrationRecords.reasonReturn",
  relocation_program: "project.migrationRecords.reasonRelocation",
  economic: "project.migrationRecords.reasonEconomic",
  other: "project.migrationRecords.reasonOther",
};

export const housingKeys: Record<string, string> = {
  own_property: "project.migrationRecords.housingOwn",
  renting: "project.migrationRecords.housingRenting",
  with_relatives: "project.migrationRecords.housingRelatives",
  collective_site: "project.migrationRecords.housingCollective",
  hotel: "project.migrationRecords.housingHotel",
  other: "project.migrationRecords.housingOther",
  unknown: "project.migrationRecords.housingUnknown",
};

export const petStatusKeys: Record<string, string> = {
  registered: "project.pets.statusRegistered",
  adopted: "project.pets.statusAdopted",
  owner_found: "project.pets.statusOwnerFound",
  needs_shelter: "project.pets.statusNeedsShelter",
  unknown: "project.pets.statusUnknown",
};

export const petOwnershipKeys: Record<string, string> = {
  with_owner: "project.petReports.withOwner",
  without_owner: "project.petReports.withoutOwner",
};
