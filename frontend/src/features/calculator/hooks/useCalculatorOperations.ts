import { useEffect, useState } from "react";
import { getOperations } from "../../../lib/api";
import type { OperationDescriptor } from "../../../lib/schemas";

export function useCalculatorOperations() {
  const [operations, setOperations] = useState<OperationDescriptor[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;

    void getOperations()
      .then((nextOperations) => {
        if (active) setOperations(nextOperations);
      })
      .catch(() => {
        if (active) setOperations([]);
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, []);

  return { operations, loading };
}