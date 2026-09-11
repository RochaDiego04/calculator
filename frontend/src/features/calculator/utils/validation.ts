import type { ApiError } from "../../../lib/schemas";

export function validateOperand(
  value: string,
  field: string,
): ApiError | undefined {
  if (value.trim() === "") {
    return { code: "MISSING_OPERAND", message: `${field} is required`, field };
  }
  if (!Number.isFinite(Number(value))) {
    return {
      code: "INVALID_OPERAND",
      message: `${field} must be a finite number`,
      field,
    };
  }
  return undefined;
}
