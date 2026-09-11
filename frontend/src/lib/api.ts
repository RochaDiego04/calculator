import {
  calculateRequest,
  calculateResponse,
  errorBody,
  operationsResponse,
  type CalculateRequest,
  type CalculateResponse,
  type OperationDescriptor,
} from "./schemas";
import { getJson, postJson } from "./http";

export { ApiErrorResponse } from "./http";

const API_BASE_URL = (import.meta.env.VITE_API_URL ?? "").replace(/\/$/, "");

export async function calculate(
  request: CalculateRequest,
  signal?: AbortSignal,
): Promise<CalculateResponse> {
  return postJson({
    url: `${API_BASE_URL}/api/v1/calculate`,
    request,
    requestSchema: calculateRequest,
    responseSchema: calculateResponse,
    errorSchema: errorBody,
    signal,
  });
}

export async function getOperations(): Promise<OperationDescriptor[]> {
  return getJson({
    url: `${API_BASE_URL}/api/v1/operations`,
    responseSchema: operationsResponse,
    errorSchema: errorBody,
  });
}
