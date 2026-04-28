import type { AuditEntry } from "@/types/audit";

export function formatAuditDetails(entry: AuditEntry): string {
  const d = entry.details;
  if (!d) return entry.summary;

  const { action, entity_type } = entry;
  const ref = entry.entity_id ?? "";

  if (entity_type === "person") {
    const name = (d.name as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created person ${name}`;
    if (action.endsWith(".update")) return `Updated person ${name}`;
    if (action.endsWith(".delete")) return `Deleted person ${name}`;
  }

  if (entity_type === "support_record") {
    const personRef =
      (d.person_name as string | undefined) ?? (d.person_id as string | undefined) ?? ref;
    const type = d.type as string | undefined;
    const sphere = d.sphere as string | undefined;
    const qualifier = [type, sphere].filter(Boolean).join("/");
    if (action.endsWith(".create"))
      return `Created ${qualifier ? qualifier + " " : ""}support record for ${personRef}`;
    if (action.endsWith(".update"))
      return `Updated ${qualifier ? qualifier + " " : ""}support record for ${personRef}`;
    if (action.endsWith(".delete")) return `Deleted support record for ${personRef}`;
  }

  if (entity_type === "migration_record") {
    const personRef = (d.person_id as string | undefined) ?? ref;
    const origin = d.origin_place_id as string | undefined;
    const dest = d.destination_place_id as string | undefined;
    const route = origin && dest ? ` (${origin} → ${dest})` : "";
    if (action.endsWith(".create")) return `Created migration record for ${personRef}${route}`;
    if (action.endsWith(".update")) return `Updated migration record for ${personRef}${route}`;
    if (action.endsWith(".delete")) return `Deleted migration record for ${personRef}`;
  }

  if (entity_type === "household") {
    const householdRef = (d.reference_number as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created household ${householdRef}`;
    if (action.endsWith(".update")) return `Updated household ${householdRef}`;
    if (action.endsWith(".delete")) return `Deleted household ${householdRef}`;
  }

  if (entity_type === "note") {
    const personRef = (d.person_id as string | undefined) ?? "";
    const forPart = personRef ? ` for ${personRef}` : "";
    if (action.endsWith(".create")) return `Created note${forPart}`;
    if (action.endsWith(".update")) return `Updated note${forPart}`;
    if (action.endsWith(".delete")) return `Deleted note${forPart}`;
  }

  if (entity_type === "document") {
    const filename = (d.filename as string | undefined) ?? ref;
    if (action === "document.upload") return `Uploaded ${filename}`;
    if (action === "document.download") return `Downloaded ${filename}`;
    if (action.endsWith(".update")) return `Updated document ${filename}`;
    if (action.endsWith(".delete")) return `Deleted ${filename}`;
  }

  if (entity_type === "pet") {
    const name = (d.name as string | undefined) ?? ref;
    const ownerId = d.owner_id as string | undefined;
    if (action.endsWith(".create"))
      return `Created pet ${name}${ownerId ? ` (owner: ${ownerId})` : ""}`;
    if (action.endsWith(".update")) return `Updated pet ${name}`;
    if (action.endsWith(".delete")) return `Deleted pet ${name}`;
  }

  if (entity_type === "tag") {
    const name = (d.name as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created tag "${name}"`;
    if (action.endsWith(".update")) return `Updated tag "${name}"`;
    if (action.endsWith(".delete")) return `Deleted tag "${name}"`;
  }

  if (entity_type === "user") {
    if (action === "admin.user.create") {
      const email = (d.email as string | undefined) ?? "";
      const role = d.role as string | undefined;
      return `Created user ${email}${role ? ` with role ${role}` : ""}`;
    }
    if (action === "user.role_change") {
      const email = (d.email as string | undefined) ?? "";
      const oldRole = d.old_role as string | undefined;
      const newRole = d.new_role as string | undefined;
      return `Changed role ${oldRole ?? ""} → ${newRole ?? ""} for ${email}`;
    }
    if (action === "user.deactivate") {
      return `Deactivated user ${(d.email as string | undefined) ?? (d.user_id as string | undefined) ?? ref}`;
    }
    if (action === "user.reactivate") {
      return `Reactivated user ${(d.email as string | undefined) ?? (d.user_id as string | undefined) ?? ref}`;
    }
  }

  if (entity_type === "permission") {
    const subjectRef =
      (d.subject_name as string | undefined) ??
      (d.subject_email as string | undefined) ??
      (d.user_id as string | undefined) ??
      ref;
    const role = d.role as string | undefined;
    if (action === "permission.grant")
      return `Granted${role ? ` ${role}` : ""} permission to ${subjectRef}`;
    if (action === "admin.permission.update")
      return `Updated permission for ${subjectRef}${role ? ` (role: ${role})` : ""}`;
    if (action === "permission.revoke") return `Revoked permission for ${subjectRef}`;
  }

  if (action === "export") {
    const count = d.count as number | undefined;
    const entityRef = (d.entity_type as string | undefined) ?? entity_type;
    return count !== undefined ? `Exported ${count} ${entityRef}` : `Exported ${entityRef}`;
  }

  return entry.summary;
}
