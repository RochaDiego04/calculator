import { useState } from "react";
import { ApiErrorResponse, calculate } from "../../../lib/api";
import type {
  ApiError,
  CalculateRequest,
  CalculateResponse,
} from "../../../lib/schemas";
import type { HistoryEntry } from "../types";

const MAX_HISTORY = 10;

export function useCalculator() {
  const [result, setResult] = useState<CalculateResponse | null>(null);
  const [errors, setErrors] = useState<ApiError[]>([]);
  const [pending, setPending] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);

  async function submit(request: CalculateRequest) {
    setPending(true);
    setErrors([]);

    try {
      const response = await calculate(request);
      setResult(response);
      setHistory((previous) =>
        [{ ...response, id: crypto.randomUUID() }, ...previous].slice(
          0,
          MAX_HISTORY,
        ),
      );
    } catch (error) {
      setResult(null);
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

  function clearHistory() {
    setHistory([]);
  }

  return { submit, result, errors, pending, history, clearHistory };
}