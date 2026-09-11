import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { CalculatorForm } from "./CalculatorForm";

const operations = [
  { name: "add", symbol: "+", arity: 2 as const, label: "Add" },
  { name: "divide", symbol: "÷", arity: 2 as const, label: "Divide" },
  { name: "sqrt", symbol: "√", arity: 1 as const, label: "Square Root" },
];

describe("CalculatorForm", () => {
  it("renders controls by their accessible labels", () => {
    render(<CalculatorForm operations={operations} onSubmit={vi.fn()} />);

    expect(screen.getByLabelText("Operation")).toBeInTheDocument();
    expect(screen.getByLabelText("A")).toBeInTheDocument();
    expect(screen.getByLabelText("B")).toBeInTheDocument();
  });

  it("blocks an empty submission without calling the API handler", async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<CalculatorForm operations={operations} onSubmit={onSubmit} />);

    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(onSubmit).not.toHaveBeenCalled();
    expect(screen.getByRole("alert")).toHaveTextContent("a is required");
    expect(screen.getByRole("alert")).toHaveTextContent("b is required");
  });

  it("hides operand B for a unary operation", async () => {
    const user = userEvent.setup();
    render(<CalculatorForm operations={operations} onSubmit={vi.fn()} />);

    await user.selectOptions(screen.getByLabelText("Operation"), "sqrt");

    expect(screen.queryByLabelText("B")).not.toBeInTheDocument();
  });

  it("replaces a stale server error with a new validation error instead of stacking them", async () => {
    const user = userEvent.setup();
    render(
      <CalculatorForm
        operations={operations}
        onSubmit={vi.fn()}
        errors={[
          {
            code: "DIVISION_BY_ZERO",
            message: "cannot divide by zero",
            field: "b",
          },
        ]}
      />,
    );

    expect(screen.getByRole("alert")).toHaveTextContent(
      "cannot divide by zero",
    );

    await user.click(screen.getByRole("button", { name: "Calculate" }));

    const alert = screen.getByRole("alert");
    expect(alert).toHaveTextContent("a is required");
    expect(alert).not.toHaveTextContent("cannot divide by zero");
  });

  it("rejects a non-numeric operand without calling the API handler", async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<CalculatorForm operations={operations} onSubmit={onSubmit} />);

    await user.type(screen.getByLabelText("A"), "abc");
    await user.type(screen.getByLabelText("B"), "3");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(onSubmit).not.toHaveBeenCalled();
    expect(screen.getByRole("alert")).toHaveTextContent(
      "a must be a finite number",
    );
  });
});
