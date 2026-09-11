import { z } from "zod";

export const operationName = z.string().min(1);

export const operationDescriptor = z.object({
  name: operationName,
  symbol: z.string(),
  arity: z.union([z.literal(1), z.literal(2)]),
  label: z.string(),
});

export const operationsResponse = z.array(operationDescriptor);

export const calculateRequest = z.object({
  operation: operationName,
  a: z.number(),
  b: z.number().optional(),
});

export const calculateResponse = z.object({
  operation: operationName,
  a: z.number(),
  b: z.number().optional(),
  result: z.number(),
});

export const apiError = z.object({
  code: z.string(),
  message: z.string(),
  field: z.string().optional(),
});

export const errorBody = z.object({
  errors: z.array(apiError),
});

export type CalculateRequest = z.infer<typeof calculateRequest>;
export type CalculateResponse = z.infer<typeof calculateResponse>;
export type ApiError = z.infer<typeof apiError>;
export type OperationDescriptor = z.infer<typeof operationDescriptor>;
