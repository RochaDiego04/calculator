import { useCalculatorForm } from "../hooks/useCalculatorForm";
import type { CalculatorFormProps } from "../types";

export function CalculatorForm({
  operations,
  onSubmit,
  errors = [],
  pending = false,
}: CalculatorFormProps) {
  const form = useCalculatorForm(operations, onSubmit);
  const visibleErrors =
    form.validationErrors.length > 0 ? form.validationErrors : errors;

  return (
    <form onSubmit={form.handleSubmit} noValidate>
      <div>
        <label htmlFor="operation">Operation</label>
        <select
          id="operation"
          value={form.selectedOperation?.name ?? ""}
          onChange={(event) => form.handleOperationChange(event.target.value)}
        >
          {operations.map((item) => (
            <option key={item.name} value={item.name}>
              {item.label}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label htmlFor="operand-a">A</label>
        <input id="operand-a" inputMode="decimal" name="a" />
      </div>

      {form.isBinary && (
        <div>
          <label htmlFor="operand-b">B</label>
          <input id="operand-b" inputMode="decimal" name="b" />
        </div>
      )}

      {visibleErrors.length > 0 && (
        <div role="alert">
          <ul>
            {visibleErrors.map((error, index) => (
              <li key={`${error.code}-${error.field ?? "general"}-${index}`}>
                {error.message}
              </li>
            ))}
          </ul>
        </div>
      )}

      <button type="submit" disabled={pending}>
        {pending ? "Calculating..." : "Calculate"}
      </button>
    </form>
  );
}
