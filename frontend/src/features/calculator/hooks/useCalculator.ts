import { useState } from "react";
import { ApiErrorResponse, calculate } from "../../../lib/api";
import type {
  ApiError,
  CalculateRequest,
  CalculateResponse,
} from "../../../lib/schemas";

export function useCalculator() {
  const [result, setResult] = useState<CalculateResponse | null>(null);
  const [errors, setErrors] = useState<ApiError[]>([]);
  const [pending, setPending] = useState(false);

  async function submit(request: CalculateRequest) {
    setPending(true);
    setErrors([]);

    try {
      const response = await calculate(request);
      setResult(response);
    } catch (error) {
      if (error instanceof ApiErrorResponse) {
        setErrors(error.errors);
      } else {
        setErrors([
          {
            code: "NETWORK_ERROR",
            message: "the calculation could not be completed",
          },
        ]);
      }
    } finally {
      setPending(false);
    }
  }

  return { submit, result, errors, pending };
}