import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiErrorResponse, calculate } from "./lib/api";
import App from "./App";

vi.mock("./lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("./lib/api")>();
  return {
    ...actual,
    calculate: vi.fn(),
    getOperations: vi.fn().mockResolvedValue([
      { name: "add", symbol: "+", arity: 2, label: "Add" },
      { name: "divide", symbol: "÷", arity: 2, label: "Divide" },
    ]),
  };
});

function setup() {
  const user = userEvent.setup();
  return { user, ...render(<App />) };
}

afterEach(() => {
  vi.clearAllMocks();
});

describe("App", () => {
  it("renders the calculator heading", () => {
    setup();

    expect(
      screen.getByRole("heading", { name: "Calculator" }),
    ).toBeInTheDocument();
  });

  it("then renders the formatted result after a successful calculation", async () => {
    vi.mocked(calculate).mockResolvedValue({
      operation: "add",
      a: 0.1,
      b: 0.2,
      result: 0.30000000000000004,
    });
    const { user } = setup();

    await user.type(screen.getByLabelText("A"), "0.1");
    await user.type(await screen.findByLabelText("B"), "0.2");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("status")).toHaveTextContent("0.3");
    expect(calculate).toHaveBeenCalledWith({
      operation: "add",
      a: 0.1,
      b: 0.2,
    });
  });

  it("then renders the server error message", async () => {
    vi.mocked(calculate).mockRejectedValue(
      new ApiErrorResponse(422, [
        {
          code: "DIVISION_BY_ZERO",
          message: "cannot divide by zero",
          field: "b",
        },
      ]),
    );
    const { user } = setup();

    await user.type(screen.getByLabelText("A"), "12");
    await user.type(await screen.findByLabelText("B"), "0");
    await user.selectOptions(screen.getByLabelText("Operation"), "divide");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "cannot divide by zero",
    );
  });
});
