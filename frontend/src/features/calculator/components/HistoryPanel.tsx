import type { OperationDescriptor } from "../../../lib/schemas";
import type { HistoryEntry } from "../types/types";
import { formatExpression } from "../utils/expression";

type HistoryPanelProps = {
  entries: HistoryEntry[];
  operations: OperationDescriptor[];
  onClear: () => void;
};

export function HistoryPanel({
  entries,
  operations,
  onClear,
}: HistoryPanelProps) {
  return (
    <section className="glass rounded-card p-5 sm:p-6">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-base font-semibold text-ink-50">History</h2>
        <button
          type="button"
          onClick={onClear}
          disabled={entries.length === 0}
          className="rounded-control px-3 py-1.5 text-sm font-medium text-ink-400 transition-colors duration-200 hover:bg-white/5 hover:text-mint-300 focus-visible:ring-2 focus-visible:ring-mint-500/50 focus-visible:outline-none disabled:cursor-not-allowed disabled:text-ink-800 disabled:hover:bg-transparent"
        >
          Clear
        </button>
      </div>

      {entries.length === 0 ? (
        <p className="mt-5 rounded-control border border-dashed border-white/10 px-4 py-8 text-center text-sm text-ink-600">
          No calculations yet.
        </p>
      ) : (
        <ul className="mt-4 grid gap-1.5">
          {entries.map((entry) => (
            <li
              key={entry.id}
              className="rounded-control bg-white/3 px-3 py-2.5 numerals text-sm text-ink-200 transition-colors duration-200 hover:bg-white/6"
            >
              {formatExpression(entry, operations)}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
