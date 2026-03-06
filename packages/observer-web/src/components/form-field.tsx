import { Field } from "@base-ui/react/field";

const inputClass =
  "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg";

const textareaClass =
  "block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg";

export { inputClass, textareaClass };

interface FormFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  required?: boolean;
  disabled?: boolean;
  type?: string;
  maxLength?: number;
  className?: string;
  error?: string;
}

export function FormField({
  label,
  value,
  onChange,
  required,
  disabled,
  type,
  maxLength,
  className,
  error,
}: FormFieldProps) {
  return (
    <Field.Root className={className} invalid={!!error}>
      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
        {label}
        {required && " *"}
      </Field.Label>
      <Field.Control
        required={required}
        disabled={disabled}
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        maxLength={maxLength}
        className={`${inputClass} ${error ? "border-rose!" : ""}`}
      />
      {error && <p className="mt-1 text-xs text-rose">{error}</p>}
    </Field.Root>
  );
}

interface FormTextareaProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  rows?: number;
  className?: string;
  error?: string;
}

export function FormTextarea({
  label,
  value,
  onChange,
  rows = 4,
  className,
  error,
}: FormTextareaProps) {
  return (
    <Field.Root className={className} invalid={!!error}>
      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
        {label}
      </Field.Label>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        rows={rows}
        className={`${textareaClass} ${error ? "border-rose!" : ""}`}
      />
      {error && <p className="mt-1 text-xs text-rose">{error}</p>}
    </Field.Root>
  );
}
