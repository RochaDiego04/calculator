import { formatResult } from "../../../lib/format";
import type { OperationDescriptor } from "../../../lib/schemas";
import type { HistoryEntry } from "../types";

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
    <section>
      <div>
        <h2>History</h2>
        <button
          type="button"
          onClick={onClear}
          disabled={entries.length === 0}
        >
          Clear
        </button>
      </div>

      {entries.length === 0 ? (
        <p>No calculations yet.</p>
      ) : (
        <ul>
          {entries.map((entry) => {
            const descriptor = operations.find(
              (item) => item.name === entry.operation,
            );
            const symbol = descriptor?.symbol ?? entry.operation;
            const isBinary = descriptor
              ? descriptor.arity === 2
              : entry.b !== undefined;

            return (
              <li key={entry.id}>
                {isBinary
                  ? `${formatResult(entry.a)} ${symbol} ${formatResult(entry.b ?? 0)} = ${formatResult(entry.result)}`
                  : `${symbol} ${formatResult(entry.a)} = ${formatResult(entry.result)}`}
              </li>
            );
          })}
        </ul>
      )}
    </section>
  );
}
