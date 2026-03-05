export type MovementReason =
  | "conflict"
  | "security"
  | "service_access"
  | "return"
  | "relocation_program"
  | "economic"
  | "other";

export type HousingAtDestination =
  | "own_property"
  | "renting"
  | "with_relatives"
  | "collective_site"
  | "hotel"
  | "other"
  | "unknown";

export interface MigrationRecord {
  id: string;
  person_id: string;
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
  created_at: string;
}

export interface CreateMigrationRecordInput {
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
}

export interface ListMigrationRecordsOutput {
  migration_records: MigrationRecord[];
  total: number;
}
