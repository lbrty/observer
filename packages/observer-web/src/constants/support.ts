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

export const SUPPORT_TYPE_VALUES = [
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

export const SUPPORT_SPHERE_VALUES = [
  "housing_assistance",
  "document_recovery",
  "social_benefits",
  "property_rights",
  "employment_rights",
  "family_law",
  "healthcare_access",
  "education_access",
  "financial_aid",
  "psychological_support",
  "other",
] as const;

export const referralKeys: Record<string, string> = {
  pending: "project.supportRecords.referralPending",
  accepted: "project.supportRecords.referralAccepted",
  completed: "project.supportRecords.referralCompleted",
  declined: "project.supportRecords.referralDeclined",
  no_response: "project.supportRecords.referralNoResponse",
};
