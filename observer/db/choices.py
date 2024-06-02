from observer.common.types import AgeGroup, DisplacedPersonStatus, PetStatus, Role, Sex
from observer.db.utils import choices_from_enum


age_groups = choices_from_enum(AgeGroup)
person_statuses = choices_from_enum(DisplacedPersonStatus)
pet_statuses = choices_from_enum(PetStatus)
sex = choices_from_enum(Sex)
user_roles = choices_from_enum(Role)
