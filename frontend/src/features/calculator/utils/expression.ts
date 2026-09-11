import { formatResult } from "../../../lib/format";
import type {
  CalculateResponse,
  OperationDescriptor,
} from "../../../lib/schemas";

type Calculation = Pick<CalculateResponse, "operation" | "a" | "b" | "result">;

export function formatOperands(
  calculation: Calculation,
  operations: OperationDescriptor[],
): string {
  const descriptor = operations.find(
    (item) => item.name === calculation.operation,
  );
  const symbol = descriptor?.symbol ?? calculation.operation;
  const isBinary = descriptor
    ? descriptor.arity === 2
    : calculation.b !== undefined;

  return isBinary
    ? `${formatResult(calculation.a)} ${symbol} ${formatResult(calculation.b ?? 0)}`
    : `${symbol} ${formatResult(calculation.a)}`;
}

export function formatExpression(
  calculation: Calculation,
  operations: OperationDescriptor[],
): string {
  return `${formatOperands(calculation, operations)} = ${formatResult(calculation.result)}`;
}
