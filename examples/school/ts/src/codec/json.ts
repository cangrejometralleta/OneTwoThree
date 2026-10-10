// JSON is a vendor boundary. This module alone owns parsing and serialization.
export function parseJson(source: string): unknown {
  let i = 0;
  const ws = () => { while (/[\u0009\u000a\u000d\u0020]/.test(source[i] ?? '')) i++; };
  const string = (): string => { const start = i++; while (i < source.length) { if (source[i] === '\\') { i += 2; continue; } if (source[i++] === '"') return JSON.parse(source.slice(start, i)); } throw new Error('json'); };
  const value = (): unknown => {
    ws(); const ch = source[i];
    if (ch === '"') return string();
    if (ch === '{') { i++; ws(); const out: Record<string, unknown> = Object.create(null); if (source[i] === '}') { i++; return out; } while (true) { ws(); if (source[i] !== '"') throw new Error('json'); const key = string(); ws(); if (source[i++] !== ':') throw new Error('json'); out[key] = value(); ws(); const sep = source[i++]; if (sep === '}') return out; if (sep !== ',') throw new Error('json'); } }
    if (ch === '[') { i++; ws(); const out: unknown[] = []; if (source[i] === ']') { i++; return out; } while (true) { out.push(value()); ws(); const sep = source[i++]; if (sep === ']') return out; if (sep !== ',') throw new Error('json'); } }
    const rest = source.slice(i); const m = /^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?/.exec(rest); if (m) { i += m[0].length; return /[.eE]/.test(m[0]) ? Number(m[0]) : BigInt(m[0]); }
    for (const [word, v] of [['true', true], ['false', false], ['null', null]] as const) if (rest.startsWith(word)) { i += word.length; return v; }
    throw new Error('json');
  };
  const result = value(); ws(); if (i !== source.length) throw new Error('json'); return result;
}
export function stringifyJson(value: unknown): string {
  return JSON.stringify(value, (_key, item) => typeof item === 'bigint' ? { __jsonBigInt: item.toString() } : item)
    .replace(/\{"__jsonBigInt":"(-?\d+)"\}/g, '$1');
}
