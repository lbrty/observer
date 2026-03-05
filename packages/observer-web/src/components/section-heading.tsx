interface SectionHeadingProps {
  children: React.ReactNode;
  className?: string;
}

export function SectionHeading({ children, className }: SectionHeadingProps) {
  return (
    <h3 className={`text-xs font-semibold uppercase tracking-wide text-fg-tertiary ${className ?? ""}`}>
      {children}
    </h3>
  );
}
