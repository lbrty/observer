import type { InputHTMLAttributes, ReactNode } from "react";

import { afterEach, describe, expect, it, mock } from "bun:test";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";

mock.module("@base-ui/react/field", () => ({
  Field: {
    Root: ({ children, className }: { children: ReactNode; className?: string }) => (
      <div className={className}>{children}</div>
    ),
    Label: ({ children, className }: { children: ReactNode; className?: string }) => (
      <label className={className}>{children}</label>
    ),
    Control: (props: InputHTMLAttributes<HTMLInputElement>) => <input {...props} />,
  },
}));

afterEach(cleanup);

const { FormField, FormTextarea } = await import("@/components/form-field");

describe("FormField", () => {
  it("renders label and input", () => {
    render(<FormField label="Email" value="test@example.com" onChange={() => {}} />);

    expect(screen.getByText("Email")).toBeDefined();
    expect(screen.getByDisplayValue("test@example.com")).toBeDefined();
  });

  it("calls onChange when input value changes", () => {
    const onChange = mock(() => {});
    render(<FormField label="Name" value="" onChange={onChange} />);

    const input = screen.getByDisplayValue("");
    fireEvent.change(input, { target: { value: "hello" } });

    expect(onChange).toHaveBeenCalledWith("hello");
  });
});

describe("FormTextarea", () => {
  it("renders label and textarea", () => {
    render(<FormTextarea label="Notes" value="some notes" onChange={() => {}} />);

    expect(screen.getByText("Notes")).toBeDefined();
    expect(screen.getByDisplayValue("some notes")).toBeDefined();
  });
});
