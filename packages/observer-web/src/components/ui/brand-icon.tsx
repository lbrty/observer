interface BrandIconProps {
  size?: "sm" | "lg";
}

export function BrandIcon({ size = "sm" }: BrandIconProps) {
  const dim = size === "lg" ? "size-14 rounded-2xl" : "size-7 rounded-lg";
  const svg = size === "lg" ? 28 : 14;
  return (
    <span className={`brand-icon inline-flex items-center justify-center ${dim}`}>
      <svg width={svg} height={svg} viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="1.5" fill="none" />
        <circle cx="12" cy="12" r="3" fill="currentColor" />
      </svg>
    </span>
  );
}
