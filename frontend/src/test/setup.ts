import { afterEach } from "vitest";
import "@testing-library/jest-dom/vitest";

afterEach(() => {
  window.localStorage.clear();
});

// jsdom ships no matchMedia, and the ambient background queries it for
// reduced-motion and pointer capability.
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});
