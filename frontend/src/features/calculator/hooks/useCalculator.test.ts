import { act, renderHook } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiErrorResponse, calculate } from "../../../lib/api";
import type { CalculateResponse } from "../../../lib/schemas";
import { makeCalculateRequest, makeCalculateResponse } from "../../../test/factories";
import {
  HISTORY_STORAGE_KEY,
  loadHistory,
  saveHistory,
} from "../utils/historyStorage";
import { useCalculator } from "./useCalculator";

vi.mock("../../../lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../../../lib/api")>();
  return { ...actual, calculate: vi.fn() };
});

function setup() {
  return renderHook(() => useCalculator());
}

afterEach(() => {
  vi.clearAllMocks();
});

describe("useCalculator", () => {
  it("then starts idle with no result, errors or history", () => {
    const { result } = setup();

    expect(result.current.result).toBeNull();
    expect(result.current.errors).toEqual([]);
    expect(result.current.pending).toBe(false);
    expect(result.current.history).toEqual([]);
  });

  describe("given a calculation is in flight", () => {
    it("then reports pending until it settles", async () => {
      let resolveCalculate: (value: CalculateResponse) => void = () => {};
      vi.mocked(calculate).mockReturnValue(
        new Promise((resolve) => {
          resolveCalculate = resolve;
        }),
      );
      const { result } = setup();

      act(() => {
        void result.current.submit(makeCalculateRequest());
      });
      expect(result.current.pending).toBe(true);

      await act(async () => resolveCalculate(makeCalculateResponse()));
      expect(result.current.pending).toBe(false);
    });
  });

  describe("when the calculation succeeds", () => {
    it("then stores the result and prepends it to history", async () => {
      vi.mocked(calculate).mockResolvedValue(makeCalculateResponse());
      const { result } = setup();

      await act(async () => result.current.submit(makeCalculateRequest()));

      expect(result.current.result).toEqual(makeCalculateResponse());
      expect(result.current.history).toHaveLength(1);
      expect(result.current.history[0]).toMatchObject(makeCalculateResponse());
      expect(result.current.history[0].id).toEqual(expect.any(String));
    });

    it("and caps history at 10 entries, newest first", async () => {
      const { result } = setup();

      for (let i = 0; i < 11; i++) {
        vi.mocked(calculate).mockResolvedValueOnce(
          makeCalculateResponse({ a: i, result: i + 2 }),
        );
        await act(async () =>
          result.current.submit(makeCalculateRequest({ a: i })),
        );
      }

      expect(result.current.history).toHaveLength(10);
      expect(result.current.history[0]).toMatchObject({ a: 10, result: 12 });
    });
  });

  describe("but when the calculation fails with a server error", () => {
    it("then clears any previous result and reports the server's errors", async () => {
      vi.mocked(calculate).mockResolvedValueOnce(makeCalculateResponse());
      const { result } = setup();
      await act(async () => result.current.submit(makeCalculateRequest()));
      expect(result.current.result).not.toBeNull();

      vi.mocked(calculate).mockRejectedValueOnce(
        new ApiErrorResponse(422, [
          {
            code: "DIVISION_BY_ZERO",
            message: "cannot divide by zero",
            field: "b",
          },
        ]),
      );
      await act(async () =>
        result.current.submit(
          makeCalculateRequest({ operation: "divide", a: 12, b: 0 }),
        ),
      );

      expect(result.current.result).toBeNull();
      expect(result.current.errors).toEqual([
        {
          code: "DIVISION_BY_ZERO",
          message: "cannot divide by zero",
          field: "b",
        },
      ]);
    });
  });

  describe("but when the calculation fails without a server response", () => {
    it("then reports a generic network error", async () => {
      vi.mocked(calculate).mockRejectedValue(new TypeError("Failed to fetch"));
      const { result } = setup();

      await act(async () => result.current.submit(makeCalculateRequest()));

      expect(result.current.errors).toEqual([
        {
          code: "NETWORK_ERROR",
          message: "the calculation could not be completed",
        },
      ]);
    });
  });

  it("then clears history on demand", async () => {
    vi.mocked(calculate).mockResolvedValue(makeCalculateResponse());
    const { result } = setup();
    await act(async () => result.current.submit(makeCalculateRequest()));
    expect(result.current.history).toHaveLength(1);

    act(() => result.current.clearHistory());

    expect(result.current.history).toEqual([]);
  });

  describe("persistence", () => {
    describe("given history was already saved from a previous session", () => {
      it("then loads it on mount", () => {
        saveHistory([{ ...makeCalculateResponse(), id: "saved-1" }]);

        const { result } = setup();

        expect(result.current.history).toEqual([
          { ...makeCalculateResponse(), id: "saved-1" },
        ]);
      });
    });

    describe("given the stored history is corrupted", () => {
      it("then starts empty instead of throwing", () => {
        window.localStorage.setItem(HISTORY_STORAGE_KEY, "not valid json");

        const { result } = setup();

        expect(result.current.history).toEqual([]);
      });
    });

    describe("given the stored history has the wrong shape", () => {
      it("then discards it and starts empty", () => {
        window.localStorage.setItem(
          HISTORY_STORAGE_KEY,
          JSON.stringify([{ operation: "add", result: "not a number" }]),
        );

        const { result } = setup();

        expect(result.current.history).toEqual([]);
      });
    });

    describe("when a calculation succeeds", () => {
      it("then persists the updated history", async () => {
        vi.mocked(calculate).mockResolvedValue(makeCalculateResponse());
        const { result } = setup();

        await act(async () => result.current.submit(makeCalculateRequest()));

        expect(loadHistory()).toEqual(result.current.history);
      });
    });

    describe("when history is cleared", () => {
      it("then persists the empty history", async () => {
        vi.mocked(calculate).mockResolvedValue(makeCalculateResponse());
        const { result } = setup();
        await act(async () => result.current.submit(makeCalculateRequest()));

        act(() => result.current.clearHistory());

        expect(loadHistory()).toEqual([]);
      });
    });
  });
});
