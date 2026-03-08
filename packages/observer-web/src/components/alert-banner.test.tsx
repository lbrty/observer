import { afterEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@testing-library/react";

import "@/test/mock-icons";

const { SuccessBanner, ErrorBanner } = await import("@/components/alert-banner");

afterEach(cleanup);

describe("SuccessBanner", () => {
  it("renders nothing when no message", () => {
    const { container } = render(<SuccessBanner message="" />);

    expect(container.innerHTML).toBe("");
  });

  it("renders message text", () => {
    render(<SuccessBanner message="Operation successful" />);

    expect(screen.getByText("Operation successful")).toBeDefined();
    expect(screen.getByTestId("check-icon")).toBeDefined();
  });
});

describe("ErrorBanner", () => {
  it("renders nothing when no message", () => {
    const { container } = render(<ErrorBanner message="" />);

    expect(container.innerHTML).toBe("");
  });

  it("renders message text", () => {
    render(<ErrorBanner message="Something went wrong" />);

    expect(screen.getByText("Something went wrong")).toBeDefined();
    expect(screen.getByTestId("warning-icon")).toBeDefined();
  });
});
