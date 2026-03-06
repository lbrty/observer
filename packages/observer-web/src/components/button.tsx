import {
  forwardRef,
  cloneElement,
  isValidElement,
  type ButtonHTMLAttributes,
  type ReactElement,
  type ReactNode,
} from "react";

type Variant = "primary" | "secondary" | "ghost" | "danger";
type Size = "sm" | "md";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: Size;
  loading?: boolean;
  icon?: ReactNode;
  asChild?: boolean;
}

const base =
  "inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98]";

const variantStyles: Record<Variant, string> = {
  primary: "bg-accent text-accent-fg shadow-card hover:opacity-90",
  secondary:
    "border border-border-secondary text-fg-secondary shadow-card hover:bg-bg-tertiary",
  ghost: "text-fg-secondary hover:bg-bg-tertiary hover:text-fg",
  danger: "bg-rose text-white shadow-card hover:opacity-90",
};

const sizeStyles: Record<Size, string> = {
  sm: "h-8 px-3 text-xs",
  md: "h-9 px-4",
};

function cx(...classes: (string | undefined | false)[]): string {
  return classes.filter(Boolean).join(" ");
}

const Spinner = (
  <span className="inline-block size-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
);

const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  {
    variant = "primary",
    size = "md",
    loading = false,
    icon,
    asChild = false,
    className,
    disabled,
    children,
    ...rest
  },
  ref,
) {
  const classes = cx(base, variantStyles[variant], sizeStyles[size], className);
  const isDisabled = disabled || loading;

  const inner = (
    <>
      {loading ? Spinner : icon}
      {children}
    </>
  );

  if (asChild && isValidElement(children)) {
    return cloneElement(children as ReactElement<Record<string, unknown>>, {
      className: cx(
        base,
        variantStyles[variant],
        sizeStyles[size],
        (children.props as { className?: string }).className,
        className,
      ),
      ref,
      ...rest,
    });
  }

  return (
    <button
      ref={ref}
      type="button"
      className={classes}
      disabled={isDisabled}
      {...rest}
    >
      {inner}
    </button>
  );
});

export { Button };
export type { ButtonProps };
