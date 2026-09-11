import { CalculatorForm } from "./features/calculator/components/CalculatorForm";
import { HistoryPanel } from "./features/calculator/components/HistoryPanel";
import { formatResult } from "./lib/format";
import { useCalculator } from "./features/calculator/hooks/useCalculator";
import { useCalculatorOperations } from "./features/calculator/hooks/useCalculatorOperations";

function App() {
  const calculator = useCalculator();
  const operationData = useCalculatorOperations();

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-900 text-white">
      <section>
        <h1>Calculator</h1>
        <CalculatorForm
          operations={operationData.operations}
          onSubmit={calculator.submit}
          errors={calculator.errors}
          pending={calculator.pending}
        />
        {operationData.loading && <p>Loading operations...</p>}
        {calculator.result && (
          <output>{formatResult(calculator.result.result)}</output>
        )}
        <HistoryPanel
          entries={calculator.history}
          operations={operationData.operations}
          onClear={calculator.clearHistory}
        />
      </section>
    </main>
  );
}

export default App;
