import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { HistoryPanel } from "./HistoryPanel";

const operations = [
  { name: "add", symbol: "+", arity: 2 as const, label: "Add" },
  { name: "sqrt", symbol: "√", arity: 1 as const, label: "Square Root" },
];

describe("HistoryPanel", () => {
  it("renders each entry using the operation's symbol", () => {
    render(
      <HistoryPanel
        entries={[
          { id: "2", operation: "sqrt", a: 4, result: 2 },
          { id: "1", operation: "add", a: 1, b: 2, result: 3 },
        ]}
        operations={operations}
        onClear={vi.fn()}
      />,
    );

    const items = screen.getAllByRole("listitem");
    expect(items).toHaveLength(2);
    expect(items[0]).toHaveTextContent("√ 4 = 2");
    expect(items[1]).toHaveTextContent("1 + 2 = 3");
  });

  it("shows an empty state and disables clear when there is no history", () => {
    render(<HistoryPanel entries={[]} operations={operations} onClear={vi.fn()} />);

    expect(screen.getByText("No calculations yet.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Clear" })).toBeDisabled();
  });

  it("calls onClear when the clear button is pressed", async () => {
    const user = userEvent.setup();
    const onClear = vi.fn();
    render(
      <HistoryPanel
        entries={[{ id: "1", operation: "add", a: 1, b: 2, result: 3 }]}
        operations={operations}
        onClear={onClear}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Clear" }));
    expect(onClear).toHaveBeenCalled();
  });
});
