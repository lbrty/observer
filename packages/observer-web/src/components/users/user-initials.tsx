interface UserInitialsProps {
  firstName: string;
  lastName: string;
  size?: "sm" | "md";
}

export function UserInitials({ firstName, lastName, size = "md" }: UserInitialsProps) {
  const letters = `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase() || "?";
  const sizeClass = size === "sm" ? "size-6 text-[10px]" : "size-8 text-xs";
  return (
    <span
      className={`inline-flex shrink-0 items-center justify-center rounded-full bg-accent/10 font-semibold text-accent ${sizeClass}`}
    >
      {letters}
    </span>
  );
}
