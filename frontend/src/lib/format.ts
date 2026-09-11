export function formatResult(value: number): string {
  return String(Number(value.toPrecision(12)));
}