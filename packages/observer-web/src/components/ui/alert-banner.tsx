import { CheckIcon, WarningIcon } from "@/components/icons";

interface AlertBannerProps {
  message: string;
}

export function SuccessBanner({ message }: AlertBannerProps) {
  if (!message) return null;
  return (
    <div className="flex items-center gap-2 rounded-lg border border-foam/20 bg-foam/8 px-3 py-2.5 text-sm font-medium text-foam">
      <CheckIcon size={16} weight="bold" className="shrink-0" />
      {message}
    </div>
  );
}

export function ErrorBanner({ message }: AlertBannerProps) {
  if (!message) return null;
  return (
    <div className="flex items-center gap-2 rounded-lg border border-rose/20 bg-rose/8 px-3 py-2.5 text-sm font-medium text-rose">
      <WarningIcon size={16} weight="bold" className="shrink-0" />
      {message}
    </div>
  );
}
