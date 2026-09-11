import { describe, expect, it } from "vitest";
import { formatResult } from "./format";

describe("formatResult", () => {
  it.each([
    [0.30000000000000004, "0.3"],
    [12, "12"],
    [1.23456789012345, "1.23456789012"],
  ])("formats %s as %s", (value, expected) => {
    expect(formatResult(value)).toBe(expected);
  });
});
