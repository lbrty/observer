import { Button } from "@/components/button";
import { PencilSimpleIcon, TrashIcon } from "@/components/icons";

interface RowActionsProps {
  onEdit: () => void;
  onDelete: () => void;
}

export function RowActions({ onEdit, onDelete }: RowActionsProps) {
  return (
    <div className="flex gap-2">
      <Button
        variant="ghost"
        className="p-1.5"
        onClick={(e) => {
          e.stopPropagation();
          onEdit();
        }}
      >
        <PencilSimpleIcon size={16} />
      </Button>
      <Button
        variant="ghost"
        className="p-1.5 hover:text-rose"
        onClick={(e) => {
          e.stopPropagation();
          onDelete();
        }}
      >
        <TrashIcon size={16} />
      </Button>
    </div>
  );
}
