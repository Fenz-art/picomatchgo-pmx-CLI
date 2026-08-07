import { defineConfig } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals.js";

export default defineConfig({
  ...nextVitals,
  ignorePatterns: [
    ".next/**",
    "out/**",
    "build/**",
    "node_modules/**",
    "next-env.d.ts",
  ],
});
