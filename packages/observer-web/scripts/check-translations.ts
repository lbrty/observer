/**
 * check-translations.ts
 *
 * Compares all locale JSON files against the English reference and reports:
 *   - Keys missing from a locale (need translating)
 *   - Keys present in a locale but not in English (orphaned — safe to delete)
 *
 * Usage:
 *   bun run scripts/check-translations.ts
 *   bun run scripts/check-translations.ts --orphans   # also show orphaned keys
 */

import { readFileSync, readdirSync } from "fs";
import { basename, resolve } from "path";

const LOCALES_DIR = resolve(import.meta.dir, "../src/locales");
const REFERENCE = "en.json";
const showOrphans = process.argv.includes("--orphans");

// Flatten nested JSON to dot-notation keys

function flatten(obj: unknown, prefix = ""): Record<string, string> {
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(obj as Record<string, unknown>)) {
    const key = prefix ? `${prefix}.${k}` : k;
    if (v !== null && typeof v === "object" && !Array.isArray(v)) {
      Object.assign(out, flatten(v, key));
    } else {
      out[key] = String(v);
    }
  }
  return out;
}

function load(file: string): Record<string, string> {
  const raw = readFileSync(resolve(LOCALES_DIR, file), "utf-8");
  return flatten(JSON.parse(raw));
}

const files = readdirSync(LOCALES_DIR).filter((f) => f.endsWith(".json"));
const reference = load(REFERENCE);
const refKeys = new Set(Object.keys(reference));

const locales = files.filter((f) => f !== REFERENCE);

let totalMissing = 0;
let totalOrphaned = 0;
let hasIssues = false;

for (const file of locales) {
  const locale = load(file);
  const localeKeys = new Set(Object.keys(locale));

  const missing = [...refKeys].filter((k) => !localeKeys.has(k));
  const orphaned = showOrphans ? [...localeKeys].filter((k) => !refKeys.has(k)) : [];

  if (missing.length === 0 && orphaned.length === 0) {
    console.log(`✓  ${basename(file).padEnd(10)}  ${localeKeys.size} keys — complete`);
    continue;
  }

  hasIssues = true;
  console.log(`\n✗  ${basename(file)}`);

  if (missing.length > 0) {
    totalMissing += missing.length;
    console.log(`   missing (${missing.length}):`);
    for (const k of missing) {
      console.log(`     - ${k}`);
      console.log(`       en: ${reference[k]}`);
    }
  }

  if (orphaned.length > 0) {
    totalOrphaned += orphaned.length;
    console.log(`   orphaned (${orphaned.length}):`);
    for (const k of orphaned) {
      console.log(`     - ${k}`);
    }
  }
}

console.log("\n" + "─".repeat(52));
console.log(`reference: ${REFERENCE}  (${refKeys.size} keys)`);
console.log(`locales checked: ${locales.length}`);

if (totalMissing > 0) console.log(`missing total: ${totalMissing}`);
if (totalOrphaned > 0) console.log(`orphaned total: ${totalOrphaned}`);
if (!hasIssues) console.log("all locales complete");

process.exit(hasIssues ? 1 : 0);
