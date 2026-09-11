import { useState } from "react";
import type {
  ApiError,
  CalculateRequest,
  OperationDescriptor,
} from "../../../lib/schemas";
import { validateOperand } from "../utils/validation";

export function useCalculatorForm(
  operations: OperationDescriptor[],
  onSubmit: (request: CalculateRequest) => Promise<void> | void,
) {
  const [selectedName, setSelectedName] = useState<string>();
  const [validationErrors, setValidationErrors] = useState<ApiError[]>([]);

  const selectedOperation =
    operations.find((item) => item.name === selectedName) ?? operations[0];
  const isBinary = selectedOperation?.arity === 2;

  function handleOperationChange(value: string) {
    setSelectedName(value);
    setValidationErrors([]);
  }

  function handleSubmit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedOperation) return;

    const formData = new FormData(event.currentTarget);
    const a = String(formData.get("a") ?? "");
    const b = String(formData.get("b") ?? "");
    const nextErrors = [validateOperand(a, "a")];
    if (isBinary) nextErrors.push(validateOperand(b, "b"));
    const errors = nextErrors.filter(
      (error): error is ApiError => error !== undefined,
    );
    setValidationErrors(errors);
    if (errors.length > 0) return;

    void onSubmit({
      operation: selectedOperation.name,
      a: Number(a),
      ...(isBinary ? { b: Number(b) } : {}),
    });
  }

  return {
    selectedOperation,
    isBinary,
    validationErrors,
    handleOperationChange,
    handleSubmit,
  };
}
