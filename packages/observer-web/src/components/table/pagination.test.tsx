import { afterEach, describe, expect, it, mock } from "bun:test";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";

import "@/test/mock-icons";

mock.module("react-i18next", () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

const { Pagination } = await import("@/components/pagination");

afterEach(cleanup);

describe("Pagination", () => {
  it("renders current page info", () => {
    render(<Pagination page={1} perPage={10} total={50} onChange={() => {}} />);

    expect(screen.getByText("admin.common.paginationRange")).toBeDefined();
  });

  it("calls onChange when next page button clicked", () => {
    const onChange = mock(() => {});
    render(<Pagination page={1} perPage={10} total={50} onChange={onChange} />);

    const nextBtn = screen.getByTestId("caret-right").closest("button")!;
    fireEvent.click(nextBtn);

    expect(onChange).toHaveBeenCalledWith(2);
  });

  it("calls onChange when prev page button clicked", () => {
    const onChange = mock(() => {});
    render(<Pagination page={3} perPage={10} total={50} onChange={onChange} />);

    const prevBtn = screen.getByTestId("caret-left").closest("button")!;
    fireEvent.click(prevBtn);

    expect(onChange).toHaveBeenCalledWith(2);
  });

  it("disables prev on first page", () => {
    render(<Pagination page={1} perPage={10} total={50} onChange={() => {}} />);

    const prevBtn = screen.getByTestId("caret-left").closest("button")!;
    expect(prevBtn.disabled).toBe(true);
  });

  it("disables next on last page", () => {
    render(<Pagination page={5} perPage={10} total={50} onChange={() => {}} />);

    const nextBtn = screen.getByTestId("caret-right").closest("button")!;
    expect(nextBtn.disabled).toBe(true);
  });
});
