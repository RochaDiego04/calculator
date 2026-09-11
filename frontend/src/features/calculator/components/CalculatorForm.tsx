import { useCalculatorForm } from "../hooks/useCalculatorForm";
import type { CalculatorFormProps } from "../types/types";

const labelClass = "mb-2 block text-sm font-medium text-ink-400";

const controlClass =
  "w-full rounded-control border border-white/10 bg-ink-950/60 px-4 py-3 text-base text-ink-50 transition-colors duration-200 focus:border-mint-500 focus:outline-none focus:ring-2 focus:ring-mint-500/40 aria-invalid:border-danger-500 aria-invalid:ring-2 aria-invalid:ring-danger-500/25";

export function CalculatorForm({
  operations,
  onSubmit,
  errors = [],
  pending = false,
}: CalculatorFormProps) {
  const form = useCalculatorForm(operations, onSubmit);
  const visibleErrors =
    form.validationErrors.length > 0 ? form.validationErrors : errors;
  const hasError = (field: string) =>
    visibleErrors.some((error) => error.field === field);
  const describedBy = visibleErrors.length > 0 ? "calculator-errors" : undefined;
  const awaitingOperations = operations.length === 0;

  return (
    <form onSubmit={form.handleSubmit} noValidate className="grid gap-5">
      <div>
        <label htmlFor="operation" className={labelClass}>
          Operation
        </label>
        <select
          id="operation"
          value={form.selectedOperation?.name ?? ""}
          onChange={(event) => form.handleOperationChange(event.target.value)}
          aria-invalid={hasError("operation")}
          aria-describedby={describedBy}
          className={`${controlClass} select-chevron appearance-none pr-11`}
        >
          {operations.map((item) => (
            <option key={item.name} value={item.name}>
              {item.label}
            </option>
          ))}
        </select>
      </div>

      <div className={form.isBinary ? "grid gap-4 sm:grid-cols-2" : "grid gap-4"}>
        <div>
          <label htmlFor="operand-a" className={labelClass}>
            A
          </label>
          <input
            id="operand-a"
            name="a"
            inputMode="decimal"
            autoComplete="off"
            aria-invalid={hasError("a")}
            aria-describedby={describedBy}
            className={`${controlClass} numerals`}
          />
        </div>

        {form.isBinary && (
          <div>
            <label htmlFor="operand-b" className={labelClass}>
              B
            </label>
            <input
              id="operand-b"
              name="b"
              inputMode="decimal"
              autoComplete="off"
              aria-invalid={hasError("b")}
              aria-describedby={describedBy}
              className={`${controlClass} numerals`}
            />
          </div>
        )}
      </div>

      {visibleErrors.length > 0 && (
        <div
          role="alert"
          id="calculator-errors"
          className="rounded-control border border-danger-500/30 bg-danger-500/10 px-4 py-3"
        >
          {visibleErrors.map((error, index) => (
            <p
              key={`${error.code}-${error.field ?? "general"}-${index}`}
              className="text-sm leading-relaxed text-danger-300"
            >
              {error.message}
            </p>
          ))}
        </div>
      )}

      <button
        type="submit"
        disabled={pending || awaitingOperations}
        className="w-full rounded-control bg-mint-400 px-5 py-3 text-base font-semibold text-ink-950 transition duration-200 hover:bg-mint-300 focus-visible:ring-2 focus-visible:ring-mint-300 focus-visible:ring-offset-2 focus-visible:ring-offset-ink-950 focus-visible:outline-none active:scale-[0.98] disabled:cursor-not-allowed disabled:bg-ink-800 disabled:text-ink-600"
      >
        {pending ? "Calculating..." : "Calculate"}
      </button>
    </form>
  );
}
