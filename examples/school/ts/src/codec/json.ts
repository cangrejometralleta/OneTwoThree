/**
 * The strict JSON Reader and Writer, standard Library only. JSON.parse Rounds a 64 bit Number,
 * Keeps the last of a repeated Key and Hides whether 20 was Written 20.0, so this Module Reads
 * one Value Token by Token and Hands back Plain Values: a Map, an Array, a string, a bigint for a
 * whole Number, a number for any other, a boolean or null.
 * Reference: https://www.rfc-editor.org/rfc/rfc8259
 */

/** JsonValue is every Plain Value a Read can Hand back. */
export type JsonValue = null | boolean | bigint | number | string | JsonValue[] | Map<string, JsonValue>;

const STRICT_UTF8 = new TextDecoder('utf-8', { fatal: true, ignoreBOM: true });
const WHITESPACE = new Set([' ', '\t', '\n', '\r']);
const WHOLE_NUMBER = /^-?(0|[1-9][0-9]*)$/;
const ANY_NUMBER = /^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?/;
const ESCAPES: Readonly<Record<string, string>> = { '"': '"', '\\': '\\', '/': '/', b: '\b', f: '\f', n: '\n', r: '\r', t: '\t' };

class Broken extends Error {}

/** Cursor Walks one Text once. */
class Cursor {
  position = 0;
  readonly text: string;

  constructor(text: string) {
    this.text = text;
  }

  skipWhitespace(): void {
    while (this.position < this.text.length && WHITESPACE.has(this.text[this.position] ?? '')) {
      this.position++;
    }
  }

  peek(): string {
    return this.text[this.position] ?? '';
  }

  expect(word: string): void {
    if (!this.text.startsWith(word, this.position)) {
      throw new Broken();
    }
    this.position += word.length;
  }
}

/** readJsonText Reads one whole JSON Value from strict UTF-8, or undefined. */
export function readJsonText(data: Uint8Array): JsonValue | undefined {
  try {
    const cursor = new Cursor(STRICT_UTF8.decode(data));
    const value = readValue(cursor);
    cursor.skipWhitespace();

    return cursor.position === cursor.text.length ? value : undefined;
  } catch {
    return undefined;
  }
}

function readValue(cursor: Cursor): JsonValue {
  cursor.skipWhitespace();
  const next = cursor.peek();
  if (next === '{') {
    return readObject(cursor);
  }
  if (next === '[') {
    return readArray(cursor);
  }
  if (next === '"') {
    return readString(cursor);
  }
  for (const [word, value] of [['true', true], ['false', false], ['null', null]] as const) {
    if (cursor.text.startsWith(word, cursor.position)) {
      cursor.expect(word);

      return value;
    }
  }

  return readNumber(cursor);
}

function readObject(cursor: Cursor): Map<string, JsonValue> {
  const fields = new Map<string, JsonValue>();
  cursor.expect('{');
  cursor.skipWhitespace();
  if (cursor.peek() === '}') {
    cursor.expect('}');

    return fields;
  }

  for (;;) {
    cursor.skipWhitespace();
    const key = readString(cursor);
    if (fields.has(key)) {
      throw new Broken();
    }
    cursor.skipWhitespace();
    cursor.expect(':');
    fields.set(key, readValue(cursor));
    cursor.skipWhitespace();
    if (cursor.peek() === '}') {
      cursor.expect('}');

      return fields;
    }
    cursor.expect(',');
  }
}

function readArray(cursor: Cursor): JsonValue[] {
  const items: JsonValue[] = [];
  cursor.expect('[');
  cursor.skipWhitespace();
  if (cursor.peek() === ']') {
    cursor.expect(']');

    return items;
  }

  for (;;) {
    items.push(readValue(cursor));
    cursor.skipWhitespace();
    if (cursor.peek() === ']') {
      cursor.expect(']');

      return items;
    }
    cursor.expect(',');
  }
}

function readString(cursor: Cursor): string {
  cursor.expect('"');
  let text = '';
  for (;;) {
    const character = cursor.text[cursor.position++];
    if (character === undefined || character < ' ') {
      throw new Broken();
    }
    if (character === '"') {
      return text;
    }
    if (character !== '\\') {
      text += character;
      continue;
    }

    const escape = cursor.text[cursor.position++] ?? '';
    if (escape === 'u') {
      const hex = cursor.text.slice(cursor.position, cursor.position + 4);
      if (!/^[0-9a-fA-F]{4}$/.test(hex)) {
        throw new Broken();
      }
      text += String.fromCharCode(Number.parseInt(hex, 16));
      cursor.position += 4;
    } else if (escape in ESCAPES) {
      text += ESCAPES[escape];
    } else {
      throw new Broken();
    }
  }
}

function readNumber(cursor: Cursor): bigint | number {
  const found = ANY_NUMBER.exec(cursor.text.slice(cursor.position));
  if (found === null) {
    throw new Broken();
  }
  cursor.position += found[0].length;

  return WHOLE_NUMBER.test(found[0]) ? BigInt(found[0]) : Number(found[0]);
}

/** writeJsonText Renders a Plain Value or a Record as compact JSON; a bigint Writes its Digits. */
export function writeJsonText(value: unknown): string {
  if (typeof value === 'bigint') {
    return value.toString();
  }
  if (value instanceof Map) {
    return writeJsonText(Object.fromEntries(value));
  }
  if (Array.isArray(value)) {
    return `[${value.map(writeJsonText).join(',')}]`;
  }
  if (value !== null && typeof value === 'object') {
    const fields = Object.entries(value).filter(([, item]) => item !== undefined);

    return `{${fields.map(([key, item]) => `${JSON.stringify(key)}:${writeJsonText(item)}`).join(',')}}`;
  }

  return JSON.stringify(value);
}
