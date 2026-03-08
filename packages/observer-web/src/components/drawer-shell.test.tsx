import { describe, it, expect, mock } from "bun:test";
import { render, screen, fireEvent } from "@testing-library/react";

mock.module("@base-ui/react/drawer", () => {
  const Passthrough = ({ children, ...props }: any) => (
    <div {...props}>{children}</div>
  );
  const Close = ({ children, ...props }: any) => (
    <button type="button" {...props}>
      {children}
    </button>
  );
  return {
    DrawerPreview: {
      Root: Passthrough,
      Portal: Passthrough,
      Backdrop: Passthrough,
      Viewport: Passthrough,
      Popup: Passthrough,
      Title: ({ children }: any) => <h2>{children}</h2>,
      Close,
    },
  };
});

mock.module("@/components/button", () => ({
  Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
}));

mock.module("@/components/icons", () => ({
  XIcon: () => <span>X</span>,
}));

mock.module("@/components/tooltip", () => ({
  Tooltip: ({ children }: any) => <>{children}</>,
}));

mock.module("react-i18next", () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

const { DrawerShell } = await import("@/components/drawer-shell");

describe("DrawerShell", () => {
  const defaultProps = {
    open: true,
    onOpenChange: mock(() => {}),
    title: "Test Drawer",
    onSubmit: mock(() => {}),
  };

  it("renders title and children when open", () => {
    render(
      <DrawerShell {...defaultProps}>
        <p>Drawer content</p>
      </DrawerShell>,
    );

    expect(screen.getByText("Test Drawer")).toBeDefined();
    expect(screen.getByText("Drawer content")).toBeDefined();
  });

  it("calls onSubmit when form submitted", () => {
    const onSubmit = mock((e: any) => e.preventDefault());

    render(
      <DrawerShell {...defaultProps} onSubmit={onSubmit}>
        <p>Content</p>
      </DrawerShell>,
    );

    fireEvent.submit(screen.getByText("Content").closest("form")!);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it("shows saving text when isPending", () => {
    render(
      <DrawerShell {...defaultProps} isPending>
        <p>Content</p>
      </DrawerShell>,
    );

    expect(screen.getByText("admin.common.saving")).toBeDefined();
  });
});
