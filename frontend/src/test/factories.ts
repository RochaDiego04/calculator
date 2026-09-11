import type {
  CalculateRequest,
  CalculateResponse,
  OperationDescriptor,
} from "../lib/schemas";

export const makeOperationDescriptor = (
  overrides: Partial<OperationDescriptor> = {},
): OperationDescriptor => ({
  name: "add",
  symbol: "+",
  arity: 2,
  label: "Add",
  ...overrides,
});

export const makeOperations = (): OperationDescriptor[] => [
  makeOperationDescriptor(),
  makeOperationDescriptor({
    name: "sqrt",
    symbol: "√",
    arity: 1,
    label: "Square Root",
  }),
];

export const makeCalculateRequest = (
  overrides: Partial<CalculateRequest> = {},
): CalculateRequest => ({
  operation: "add",
  a: 1,
  b: 2,
  ...overrides,
});

export const makeCalculateResponse = (
  overrides: Partial<CalculateResponse> = {},
): CalculateResponse => ({
  operation: "add",
  a: 1,
  b: 2,
  result: 3,
  ...overrides,
});
