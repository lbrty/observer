import { CheckIcon, WarningIcon, XIcon } from "@/components/icons";
import { useToasts, useToast } from "@/stores/toast";

const variantStyles = {
  success: "border-foam/20 bg-foam/10 text-foam",
  error: "border-rose/20 bg-rose/10 text-rose",
  info: "border-accent/20 bg-accent/10 text-accent",
};

const variantIcons = {
  success: CheckIcon,
  error: WarningIcon,
  info: CheckIcon,
};

export function ToastContainer() {
  const toasts = useToasts();
  const { dismiss } = useToast();

  if (toasts.length === 0) return null;

  return (
    <div className="fixed right-4 bottom-4 z-200 flex flex-col gap-2">
      {toasts.map((t) => {
        const Icon = variantIcons[t.variant];
        return (
          <div
            key={t.id}
            role="status"
            aria-live="polite"
            className={`flex items-center gap-2 rounded-lg border px-4 py-3 text-sm font-medium shadow-elevated backdrop-blur-sm animate-slide-in ${variantStyles[t.variant]}`}
          >
            <Icon size={16} weight="bold" className="shrink-0" />
            <span className="flex-1">{t.message}</span>
            <button
              type="button"
              onClick={() => dismiss(t.id)}
              className="shrink-0 cursor-pointer rounded p-0.5 opacity-60 hover:opacity-100"
            >
              <XIcon size={14} />
            </button>
          </div>
        );
      })}
    </div>
  );
}
