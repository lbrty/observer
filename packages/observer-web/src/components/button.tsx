import { cva, type VariantProps } from "class-variance-authority";
import {
  forwardRef,
  cloneElement,
  isValidElement,
  type ButtonHTMLAttributes,
  type ReactElement,
  type ReactNode,
} from "react";

const buttonVariants = cva(
  "inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98]",
  {
    variants: {
      variant: {
        primary: "bg-accent text-accent-fg shadow-card hover:opacity-90",
        secondary:
          "border border-border-secondary text-fg-secondary shadow-card hover:bg-bg-tertiary",
        ghost: "text-fg-secondary hover:bg-bg-tertiary hover:text-fg",
        danger: "bg-rose text-white shadow-card hover:opacity-90",
      },
      size: {
        sm: "h-8 px-3 text-xs",
        md: "h-9 px-4",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "md",
    },
  },
);

interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  loading?: boolean;
  icon?: ReactNode;
  asChild?: boolean;
}

const Spinner = (
  <span className="inline-block size-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
);

const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant, size, loading = false, icon, asChild = false, className, disabled, children, ...rest },
  ref,
) {
  const classes = buttonVariants({ variant, size, className });
  const isDisabled = disabled || loading;

  const inner = (
    <>
      {loading ? Spinner : icon}
      {children}
    </>
  );

  if (asChild && isValidElement(children)) {
    return cloneElement(children as ReactElement<Record<string, unknown>>, {
      className: buttonVariants({
        variant,
        size,
        className: [
          (children.props as { className?: string }).className,
          className,
        ]
          .filter(Boolean)
          .join(" "),
      }),
      ref,
      ...rest,
    });
  }

  return (
    <button ref={ref} type="button" className={classes} disabled={isDisabled} {...rest}>
      {inner}
    </button>
  );
});

export { Button, buttonVariants };
export type { ButtonProps };
