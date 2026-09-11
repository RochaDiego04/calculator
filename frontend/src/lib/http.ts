import { z } from "zod";

type ApiError = {
  code: string;
  message: string;
  field?: string;
};

export class ApiErrorResponse extends Error {
  readonly status: number;
  readonly errors: ApiError[];

  constructor(status: number, errors: ApiError[]) {
    super(errors.map((error) => error.message).join("; "));
    this.name = "ApiErrorResponse";
    this.status = status;
    this.errors = errors;
  }
}

type JsonRequestOptions<TRequest, TResponse> = {
  url: string;
  request: TRequest;
  requestSchema: z.ZodType<TRequest>;
  responseSchema: z.ZodType<TResponse>;
  errorSchema: z.ZodType<{ errors: ApiError[] }>;
  signal?: AbortSignal;
};

type JsonResponseOptions<TResponse> = {
  url: string;
  responseSchema: z.ZodType<TResponse>;
  errorSchema: z.ZodType<{ errors: ApiError[] }>;
  signal?: AbortSignal;
};

const malformedResponse = (status: number) =>
  new ApiErrorResponse(status, [
    {
      code: "MALFORMED_RESPONSE",
      message: "the server returned an invalid response",
    },
  ]);

async function readJson(response: Response): Promise<unknown> {
  try {
    return await response.json();
  } catch {
    throw malformedResponse(response.status);
  }
}

export async function postJson<TRequest, TResponse>({
  url,
  request,
  requestSchema,
  responseSchema,
  errorSchema,
  signal,
}: JsonRequestOptions<TRequest, TResponse>): Promise<TResponse> {
  const parsedRequest = requestSchema.safeParse(request);
  if (!parsedRequest.success) {
    throw new ApiErrorResponse(0, [
      { code: "INVALID_REQUEST", message: "the request is invalid" },
    ]);
  }

  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(parsedRequest.data),
    signal,
  });
  const body = await readJson(response);

  if (!response.ok) {
    const parsedError = errorSchema.safeParse(body);
    throw new ApiErrorResponse(
      response.status,
      parsedError.success
        ? parsedError.data.errors
        : malformedResponse(response.status).errors,
    );
  }

  const parsedResponse = responseSchema.safeParse(body);
  if (!parsedResponse.success) throw malformedResponse(response.status);
  return parsedResponse.data;
}

export async function getJson<TResponse>({
  url,
  responseSchema,
  errorSchema,
  signal,
}: JsonResponseOptions<TResponse>): Promise<TResponse> {
  const response = await fetch(url, { signal });
  const body = await readJson(response);

  if (!response.ok) {
    const parsedError = errorSchema.safeParse(body);
    throw new ApiErrorResponse(
      response.status,
      parsedError.success
        ? parsedError.data.errors
        : malformedResponse(response.status).errors,
    );
  }

  const parsedResponse = responseSchema.safeParse(body);
  if (!parsedResponse.success) throw malformedResponse(response.status);
  return parsedResponse.data;
}
