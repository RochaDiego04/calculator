import { AmbientBackground } from "./components/AmbientBackground";
import { CalculatorForm } from "./features/calculator/components/CalculatorForm";
import { HistoryPanel } from "./features/calculator/components/HistoryPanel";
import { useCalculator } from "./features/calculator/hooks/useCalculator";
import { useCalculatorOperations } from "./features/calculator/hooks/useCalculatorOperations";
import { formatOperands } from "./features/calculator/utils/expression";
import { formatResult } from "./lib/format";

function App() {
  const calculator = useCalculator();
  const operationData = useCalculatorOperations();

  return (
    <main className="relative flex min-h-dvh flex-col overflow-hidden">
      <AmbientBackground />

      <div className="relative mx-auto flex w-full max-w-5xl flex-1 flex-col mt-28 px-4 py-12 sm:px-6 md:py-16">
        <header className="max-w-xl">
          <h1 className="text-4xl font-semibold tracking-tight text-ink-50 md:text-5xl">
            Calculator
          </h1>
          <p className="mt-3 leading-relaxed text-ink-400">
            Seven operations, evaluated by a Go service. Results keep full
            float64 precision and are rounded for display only.
          </p>
        </header>

        <div className="mt-8 grid items-start gap-6 lg:mt-12 lg:grid-cols-[minmax(0,1fr)_20rem]">
          <section className="glass rounded-card p-5 sm:p-7">
            <div className="rounded-control border border-white/5 bg-ink-950/45 px-5 py-6">
              <p className="text-sm text-ink-600">Result</p>
              {calculator.result ? (
                <>
                  <output className="numerals mt-1 block text-4xl font-medium break-all text-mint-300 sm:text-5xl">
                    {formatResult(calculator.result.result)}
                  </output>
                  <p className="numerals mt-2 text-sm text-ink-600">
                    {formatOperands(
                      calculator.result,
                      operationData.operations,
                    )}
                  </p>
                </>
              ) : (
                <p
                  aria-hidden="true"
                  className="numerals mt-1 text-4xl font-medium text-ink-800 sm:text-5xl"
                >
                  0
                </p>
              )}
            </div>

            <div className="mt-6">
              <CalculatorForm
                operations={operationData.operations}
                onSubmit={calculator.submit}
                errors={calculator.errors}
                pending={calculator.pending}
              />
            </div>

            {operationData.loading && (
              <p className="mt-4 text-sm text-ink-600">
                Loading operations from the API.
              </p>
            )}
          </section>

          <HistoryPanel
            entries={calculator.history}
            operations={operationData.operations}
            onClear={calculator.clearHistory}
          />
        </div>
      </div>
    </main>
  );
}

export default App;
