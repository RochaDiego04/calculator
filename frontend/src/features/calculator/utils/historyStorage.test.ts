import { afterEach, describe, expect, it, vi } from "vitest";
import { makeCalculateResponse } from "../../../test/factories";
import type { HistoryEntry } from "../types/types";
import { HISTORY_STORAGE_KEY, loadHistory, saveHistory } from "./historyStorage";

const makeHistory = (overrides: Partial<HistoryEntry> = {}): HistoryEntry[] => [
  { ...makeCalculateResponse(), id: "entry-1", ...overrides },
];

afterEach(() => {
  vi.restoreAllMocks();
});

describe("loadHistory", () => {
  describe("given nothing has been saved yet", () => {
    it("then returns an empty array", () => {
      expect(loadHistory()).toEqual([]);
    });
  });

  describe("given valid history was saved", () => {
    it("then returns it unchanged", () => {
      const history = makeHistory();
      saveHistory(history);

      expect(loadHistory()).toEqual(history);
    });
  });

  describe("given the stored value is not valid JSON", () => {
    it("then returns an empty array instead of throwing", () => {
      window.localStorage.setItem(HISTORY_STORAGE_KEY, "{not json");

      expect(loadHistory()).toEqual([]);
    });
  });

  describe("given the stored value is valid JSON but the wrong shape", () => {
    it("then returns an empty array", () => {
      window.localStorage.setItem(
        HISTORY_STORAGE_KEY,
        JSON.stringify([{ id: "entry-1" }]),
      );

      expect(loadHistory()).toEqual([]);
    });
  });

  describe("given localStorage.getItem throws", () => {
    it("then returns an empty array instead of throwing", () => {
      vi.spyOn(window.localStorage, "getItem").mockImplementation(() => {
        throw new Error("storage is disabled");
      });

      expect(loadHistory()).toEqual([]);
    });
  });
});

describe("saveHistory", () => {
  describe("given localStorage.setItem throws", () => {
    it("then does not throw", () => {
      vi.spyOn(window.localStorage, "setItem").mockImplementation(() => {
        throw new Error("quota exceeded");
      });

      expect(() => saveHistory(makeHistory())).not.toThrow();
    });
  });
});
