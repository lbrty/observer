import { Button } from "@/components/button";
import { PencilSimpleIcon, TrashIcon } from "@/components/icons";
import { Tooltip } from "@/components/tooltip";
import { useTranslation } from "react-i18next";

interface RowActionsProps {
  onEdit: () => void;
  onDelete: () => void;
}

export function RowActions({ onEdit, onDelete }: RowActionsProps) {
  const { t } = useTranslation();

  return (
    <div className="flex gap-1">
      <Tooltip label={t("admin.common.edit")}>
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
      </Tooltip>
      <Tooltip label={t("admin.common.delete")}>
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
      </Tooltip>
    </div>
  );
}
