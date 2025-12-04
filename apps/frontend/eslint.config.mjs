import nextCoreWebVitals from "eslint-config-next/core-web-vitals";
import nextTypescript from "eslint-config-next/typescript";
import { dirname } from "path";
import { fileURLToPath } from "url";
import pluginQuery from "@tanstack/eslint-plugin-query";
// import tseslint from "typescript-eslint";
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const eslintConfig = [...nextCoreWebVitals, ...nextTypescript, ...compat.config({
  rules: {
    "@typescript-eslint/no-explicit-any": "warn",
    "@typescript-eslint/no-unused-vars": "warn",
    "@typescript-eslint/no-unsafe-member-access": "warn",
    "@next/next/no-html-link-for-pages": "warn",
    "@typescript-eslint/no-unused-expressions": "warn",
  }
}), // ...tseslint.configs.recommendedTypeChecked,
...pluginQuery.configs["flat/recommended"], {
  languageOptions: {
    parserOptions: {
      projectService: true,
      tsconfigRootDir: import.meta.dirname,
    },
  },
}, {
  ignores: ["node_modules/**", ".next/**", "out/**", "build/**", "next-env.d.ts"]
}];

export default eslintConfig;
