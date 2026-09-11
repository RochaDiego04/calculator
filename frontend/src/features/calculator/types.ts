import type {
  ApiError,
  CalculateRequest,
  OperationDescriptor,
} from "../../lib/schemas";

export type CalculatorFormProps = {
  operations: OperationDescriptor[];
  onSubmit: (request: CalculateRequest) => Promise<void> | void;
  errors?: ApiError[];
  pending?: boolean;
};
