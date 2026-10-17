import assert from 'node:assert/strict';
import { readdirSync, readFileSync } from 'node:fs';
import { dirname, join, relative } from 'node:path';
import { test } from 'node:test';

/**
 * Node cannot Let the Compiler Hold every Boundary, so this Test Reads the Imports.
 * It Says which Vendor may Live where, and the package File is Part of the Check.
 */
const SOURCES = 'src';
const IMPORT = /^import\s+(?:type\s+)?(?:[^'"]+\s+from\s+)?['"]([^'"]+)['"]/gm;

const VENDOR_HOMES: Readonly<Record<string, string>> = {
  'drizzle-orm': 'store',
  'better-sqlite3': 'store',
  express: 'serving/express',
  'node:http': 'serving',
  'node:crypto': 'tokens',
};

function readImports(): { home: string; file: string; imported: string }[] {
  const found: { home: string; file: string; imported: string }[] = [];
  for (const entry of readdirSync(SOURCES, { recursive: true, encoding: 'utf8' })) {
    if (!entry.endsWith('.ts') || entry.endsWith('.test.ts')) {
      continue;
    }
    const text = readFileSync(join(SOURCES, entry), 'utf8');
    for (const match of text.matchAll(IMPORT)) {
      found.push({ home: relative(SOURCES, join(SOURCES, dirname(entry))), file: entry, imported: match[1] ?? '' });
    }
  }

  return found;
}

test('everyVendorLivesInOneModule', () => {
  for (const { home, file, imported } of readImports()) {
    for (const [vendor, allowed] of Object.entries(VENDOR_HOMES)) {
      if (imported === vendor || imported.startsWith(`${vendor}/`)) {
        assert.ok(home === allowed || home.startsWith(`${allowed}/`), `${file} Imports ${imported}, which Lives in ${allowed}`);
      }
    }
  }
});

test('theCoreNamesNoVendorAndNoSocket', () => {
  for (const { home, file, imported } of readImports()) {
    if (home === 'school') {
      assert.match(imported, /^\.\.\/faults\/|^\.\//, `${file} Imports ${imported}`);
    }
  }
});

test('theShapesTheLayersSpeakInImportNothing', () => {
  for (const { home, file, imported } of readImports()) {
    assert.ok(!['transport', 'wire', 'faults'].includes(home), `${file} Imports ${imported}`);
  }
});

test('everyVendorIsDeclaredInThePackageFile', () => {
  const manifest = JSON.parse(readFileSync('package.json', 'utf8')) as { dependencies: Record<string, string> };

  for (const vendor of ['drizzle-orm', 'better-sqlite3', 'express']) {
    assert.ok(vendor in manifest.dependencies, `${vendor} Must be a Declared Dependency`);
  }
});
