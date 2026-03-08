import { describe, it, expect, mock } from "bun:test";
import { render, screen, fireEvent } from "@testing-library/react";

mock.module("react-i18next", () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

const { DataTable } = await import("@/components/data-table");

type Item = { id: string; name: string; email: string };

const columns = [
  { key: "name", header: "Name", render: (item: Item) => item.name },
  { key: "email", header: "Email", render: (item: Item) => item.email },
];

const data: Item[] = [
  { id: "1", name: "Alice", email: "alice@example.com" },
  { id: "2", name: "Bob", email: "bob@example.com" },
];

const keyExtractor = (item: Item) => item.id;

describe("DataTable", () => {
  it("renders column headers and data rows", () => {
    render(
      <DataTable columns={columns} data={data} keyExtractor={keyExtractor} />,
    );

    expect(screen.getByText("Name")).toBeDefined();
    expect(screen.getByText("Email")).toBeDefined();
    expect(screen.getByText("Alice")).toBeDefined();
    expect(screen.getByText("bob@example.com")).toBeDefined();
  });

  it("renders skeleton rows when isLoading is true", () => {
    const { container } = render(
      <DataTable
        columns={columns}
        data={[]}
        keyExtractor={keyExtractor}
        isLoading
      />,
    );

    const skeletonCells = container.querySelectorAll(".animate-pulse");
    expect(skeletonCells.length).toBe(10); // 5 rows * 2 columns
  });

  it("renders default empty state when data is empty", () => {
    render(
      <DataTable columns={columns} data={[]} keyExtractor={keyExtractor} />,
    );

    expect(screen.getByText("admin.common.noData")).toBeDefined();
  });

  it("renders custom empty state node", () => {
    render(
      <DataTable
        columns={columns}
        data={[]}
        keyExtractor={keyExtractor}
        emptyState={<span>Nothing here</span>}
      />,
    );

    expect(screen.getByText("Nothing here")).toBeDefined();
  });

  it("calls onRowClick when row is clicked", () => {
    const handleClick = mock(() => {});

    const { container } = render(
      <DataTable
        columns={columns}
        data={data}
        keyExtractor={keyExtractor}
        onRowClick={handleClick}
      />,
    );

    const dataRows = container.querySelectorAll("tbody tr");
    fireEvent.click(dataRows[0]);
    expect(handleClick).toHaveBeenCalledTimes(1);
    expect(handleClick.mock.calls[0][0]).toEqual(data[0]);
  });
});
