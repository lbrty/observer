import { XIcon } from "@/components/icons";
import type { AdminUser } from "@/types/admin";

interface SelectedUserCardProps {
  user: AdminUser;
  onClear: () => void;
}

export function SelectedUserCard({ user, onClear }: SelectedUserCardProps) {
  return (
    <div className="flex items-center justify-between rounded-lg border border-border-secondary bg-bg px-3 py-2">
      <div>
        <p className="text-sm text-fg">{user.first_name} {user.last_name}</p>
        <p className="text-xs text-fg-tertiary">{user.email}</p>
      </div>
      <button type="button" onClick={onClear} className="cursor-pointer rounded p-1 text-fg-tertiary hover:text-fg">
        <XIcon size={14} />
      </button>
    </div>
  );
}
