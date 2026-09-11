import type {
  ApiError,
  CalculateRequest,
  CalculateResponse,
  OperationDescriptor,
} from "../../lib/schemas";

export type CalculatorFormProps = {
  operations: OperationDescriptor[];
  onSubmit: (request: CalculateRequest) => Promise<void> | void;
  errors?: ApiError[];
  pending?: boolean;
};

export type HistoryEntry = CalculateResponse & { id: string };
