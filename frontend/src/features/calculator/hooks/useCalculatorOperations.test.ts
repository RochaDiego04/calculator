import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { getOperations } from "../../../lib/api";
import { makeOperations } from "../../../test/factories";
import { useCalculatorOperations } from "./useCalculatorOperations";

vi.mock("../../../lib/api", () => ({
  getOperations: vi.fn(),
}));

function setup() {
  return renderHook(() => useCalculatorOperations());
}

afterEach(() => {
  vi.clearAllMocks();
});

describe("useCalculatorOperations", () => {
  describe("given the discovery request is still in flight", () => {
    it("then reports loading with no operations yet", () => {
      vi.mocked(getOperations).mockReturnValue(new Promise(() => {}));

      const { result } = setup();

      expect(result.current.loading).toBe(true);
      expect(result.current.operations).toEqual([]);
    });
  });

  describe("when the discovery request succeeds", () => {
    it("then stops loading and exposes the fetched operations", async () => {
      const operations = makeOperations();
      vi.mocked(getOperations).mockResolvedValue(operations);

      const { result } = setup();
      await waitFor(() => expect(result.current.loading).toBe(false));

      expect(result.current.operations).toEqual(operations);
    });
  });

  describe("but when the discovery request fails", () => {
    it("then stops loading and falls back to an empty list", async () => {
      vi.mocked(getOperations).mockRejectedValue(new Error("network down"));

      const { result } = setup();
      await waitFor(() => expect(result.current.loading).toBe(false));

      expect(result.current.operations).toEqual([]);
    });
  });
});
