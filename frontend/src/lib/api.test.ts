import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiErrorResponse, calculate } from "./api";

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("calculate", () => {
  it("parses a successful response", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({ operation: "divide", a: 12, b: 3, result: 4 }),
          { status: 200 },
        ),
      ),
    );

    await expect(
      calculate({ operation: "divide", a: 12, b: 3 }),
    ).resolves.toEqual({ operation: "divide", a: 12, b: 3, result: 4 });
  });

  it("preserves the server error envelope", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            errors: [
              {
                code: "DIVISION_BY_ZERO",
                message: "cannot divide by zero",
                field: "b",
              },
            ],
          }),
          { status: 422 },
        ),
      ),
    );

    const result = calculate({ operation: "divide", a: 12, b: 0 });

    await expect(result).rejects.toMatchObject({
      status: 422,
      errors: [
        {
          code: "DIVISION_BY_ZERO",
          message: "cannot divide by zero",
          field: "b",
        },
      ],
    });
    await expect(result).rejects.toBeInstanceOf(ApiErrorResponse);
  });

  it("turns a malformed success into a rendered API error", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ result: "not a number" }), {
          status: 200,
        }),
      ),
    );

    await expect(
      calculate({ operation: "add", a: 1, b: 2 }),
    ).rejects.toMatchObject({
      status: 200,
      errors: [{ code: "MALFORMED_RESPONSE" }],
    });
  });
});