import { afterEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@testing-library/react";

import type { Icon } from "@/components/ui/icons";
import { EmptyState } from "@/components/ui/empty-state";

function FakeIconFn({ size }: { size: number }) {
  return <svg data-testid="icon" data-size={size} />;
}
const FakeIcon = FakeIconFn as unknown as Icon;

afterEach(cleanup);

describe("EmptyState", () => {
  it("renders title and icon", () => {
    render(<EmptyState icon={FakeIcon} title="No items" />);

    expect(screen.getByText("No items")).toBeDefined();
    expect(screen.getByTestId("icon")).toBeDefined();
  });

  it("renders description when provided", () => {
    render(<EmptyState icon={FakeIcon} title="No items" description="Try adding one" />);

    expect(screen.getByText("Try adding one")).toBeDefined();
  });

  it("does not render description when omitted", () => {
    render(<EmptyState icon={FakeIcon} title="No items" />);

    expect(screen.queryByText("Try adding one")).toBeNull();
  });

  it("renders action slot when provided", () => {
    render(<EmptyState icon={FakeIcon} title="No items" action={<button>Add item</button>} />);

    expect(screen.getByText("Add item")).toBeDefined();
  });
});
