import { z } from "zod";
import { calculateResponse } from "../../../lib/schemas";
import type {
  ApiError,
  CalculateRequest,
  OperationDescriptor,
} from "../../../lib/schemas";

export type CalculatorFormProps = {
  operations: OperationDescriptor[];
  onSubmit: (request: CalculateRequest) => Promise<void> | void;
  errors?: ApiError[];
  pending?: boolean;
};

export const historyEntry = calculateResponse.extend({ id: z.string() });
export type HistoryEntry = z.infer<typeof historyEntry>;
