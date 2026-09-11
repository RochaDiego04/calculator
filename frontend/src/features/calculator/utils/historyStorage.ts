import { z } from "zod";
import { historyEntry } from "../types/types";
import type { HistoryEntry } from "../types/types";

export const HISTORY_STORAGE_KEY = "calculator:history";

const storedHistory = z.array(historyEntry);

export function loadHistory(): HistoryEntry[] {
  try {
    const raw = window.localStorage.getItem(HISTORY_STORAGE_KEY);
    if (!raw) return [];
    const parsed = storedHistory.safeParse(JSON.parse(raw));
    return parsed.success ? parsed.data : [];
  } catch {
    return [];
  }
}

export function saveHistory(history: HistoryEntry[]): void {
  try {
    window.localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(history));
  } catch {
    // Storage can be full, disabled, or unavailable (private browsing) -
    // history still works for the session, persistence is just best-effort.
  }
}
