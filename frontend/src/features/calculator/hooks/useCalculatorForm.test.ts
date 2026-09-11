import { act, renderHook, type RenderHookResult } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { makeOperations } from "../../../test/factories";
import { useCalculatorForm } from "./useCalculatorForm";

type FormResult = RenderHookResult<
  ReturnType<typeof useCalculatorForm>,
  unknown
>["result"];

function makeForm(values: Record<string, string>) {
  const form = document.createElement("form");
  for (const [name, value] of Object.entries(values)) {
    const input = document.createElement("input");
    input.name = name;
    input.value = value;
    form.appendChild(input);
  }
  return form;
}

function setup(onSubmit = vi.fn()) {
  const { result } = renderHook(() =>
    useCalculatorForm(makeOperations(), onSubmit),
  );
  return { result, onSubmit };
}

function submit(result: FormResult, values: Record<string, string>) {
  const event = {
    preventDefault: vi.fn(),
    currentTarget: makeForm(values),
  } as unknown as Parameters<typeof result.current.handleSubmit>[0];
  act(() => result.current.handleSubmit(event));
}

describe("useCalculatorForm", () => {
  it("then defaults to the first operation and its arity", () => {
    const { result } = setup();

    expect(result.current.selectedOperation?.name).toBe("add");
    expect(result.current.isBinary).toBe(true);
  });

  describe("when the operation changes", () => {
    it("then updates the selected operation and its arity", () => {
      const { result } = setup();

      act(() => result.current.handleOperationChange("sqrt"));

      expect(result.current.selectedOperation?.name).toBe("sqrt");
      expect(result.current.isBinary).toBe(false);
    });

    it("and clears any pending validation errors", () => {
      const { result } = setup();
      submit(result, { a: "", b: "" });
      expect(result.current.validationErrors).not.toEqual([]);

      act(() => result.current.handleOperationChange("sqrt"));

      expect(result.current.validationErrors).toEqual([]);
    });
  });

  describe("given both operands are blank", () => {
    it("then blocks submit and reports a validation error per field", () => {
      const { result, onSubmit } = setup();

      submit(result, { a: "", b: "" });

      expect(onSubmit).not.toHaveBeenCalled();
      expect(result.current.validationErrors).toEqual([
        { code: "MISSING_OPERAND", message: "a is required", field: "a" },
        { code: "MISSING_OPERAND", message: "b is required", field: "b" },
      ]);
    });
  });

  describe("given a non-numeric operand", () => {
    it("then reports it as invalid rather than missing", () => {
      const { result, onSubmit } = setup();

      submit(result, { a: "abc", b: "3" });

      expect(onSubmit).not.toHaveBeenCalled();
      expect(result.current.validationErrors).toEqual([
        {
          code: "INVALID_OPERAND",
          message: "a must be a finite number",
          field: "a",
        },
      ]);
    });
  });

  describe("given a valid binary request", () => {
    it("then submits both operands and clears validation errors", () => {
      const { result, onSubmit } = setup();

      submit(result, { a: "12", b: "3" });

      expect(onSubmit).toHaveBeenCalledWith({ operation: "add", a: 12, b: 3 });
      expect(result.current.validationErrors).toEqual([]);
    });
  });

  describe("given a unary operation", () => {
    it("then submits without a b operand", () => {
      const { result, onSubmit } = setup();
      act(() => result.current.handleOperationChange("sqrt"));

      submit(result, { a: "4" });

      expect(onSubmit).toHaveBeenCalledWith({ operation: "sqrt", a: 4 });
    });
  });
});
