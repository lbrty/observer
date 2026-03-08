import type { ReactNode } from "react";

import { afterEach, describe, expect, it, mock } from "bun:test";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";

mock.module("react-i18next", () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

mock.module("@base-ui/react/dialog", () => ({
  Dialog: {
    Root: ({ children }: { children: ReactNode }) => <div>{children}</div>,
    Portal: ({ children }: { children: ReactNode }) => <div>{children}</div>,
    Backdrop: ({ className }: { className: string }) => <div className={className} />,
    Popup: ({ children, className }: { children: ReactNode; className: string }) => (
      <div className={className}>{children}</div>
    ),
    Title: ({ children, className }: { children: ReactNode; className: string }) => (
      <h2 className={className}>{children}</h2>
    ),
    Description: ({ children, className }: { children: ReactNode; className: string }) => (
      <p className={className}>{children}</p>
    ),
    Close: ({ children }: { children: ReactNode }) => <button type="button">{children}</button>,
  },
}));

mock.module("@/components/button", () => ({
  Button: ({
    children,
    onClick,
    disabled,
    ...rest
  }: {
    children: ReactNode;
    onClick?: () => void;
    disabled?: boolean;
    asChild?: boolean;
    variant?: string;
  }) => {
    if (rest.asChild) return <>{children}</>;
    return (
      <button type="button" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    );
  },
}));

afterEach(cleanup);

const { ConfirmDialog } = await import("@/components/confirm-dialog");

describe("ConfirmDialog", () => {
  it("renders title and description when open", () => {
    render(
      <ConfirmDialog
        open
        onOpenChange={() => {}}
        title="Delete item?"
        description="This cannot be undone."
        onConfirm={() => {}}
      />,
    );

    expect(screen.getByText("Delete item?")).toBeDefined();
    expect(screen.getByText("This cannot be undone.")).toBeDefined();
  });

  it("calls onConfirm when confirm button clicked", () => {
    const onConfirm = mock(() => {});

    render(
      <ConfirmDialog
        open
        onOpenChange={() => {}}
        title="Delete item?"
        description="This cannot be undone."
        onConfirm={onConfirm}
      />,
    );

    const deleteBtn = screen.getByText("admin.common.delete");
    fireEvent.click(deleteBtn);

    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it("disables confirm button when loading", () => {
    render(
      <ConfirmDialog
        open
        onOpenChange={() => {}}
        title="Delete item?"
        description="Deleting..."
        onConfirm={() => {}}
        loading
      />,
    );

    const deleteBtn = screen.getByText("admin.common.delete");
    expect(deleteBtn.closest("button")!.disabled).toBe(true);
  });
});
