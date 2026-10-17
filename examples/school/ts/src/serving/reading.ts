import type { IncomingMessage, ServerResponse } from 'node:http';

import { writeJsonText } from '../codec/json.ts';
import { buildRequest, MAX_BODY_BYTES, type Libretto, type Request, type Response } from '../transport/transport.ts';

/** MAX_SWALLOW_BYTES Bounds what is Read and Dropped past the Body, so a Client can Read the Answer. */
const MAX_SWALLOW_BYTES = 2 << 20;

const HTTP_NO_CONTENT = 204;

/**
 * serveLibretto Joins the three Beats every Adapter Shares: read, answer, write.
 * Every Framework Hands it the whole Node Request, so none of them Decides.
 */
export async function serveLibretto(libretto: Libretto, incoming: IncomingMessage, outgoing: ServerResponse): Promise<void> {
  const body = await readBoundedBody(incoming);
  const headers: Record<string, string> = {};
  for (let position = 0; position < incoming.rawHeaders.length; position += 2) {
    const name = (incoming.rawHeaders[position] ?? '').toLowerCase();
    headers[name] ??= incoming.rawHeaders[position + 1] ?? '';
  }

  const target = incoming.url ?? '/';
  const question = target.indexOf('?');
  const rawPath = question < 0 ? target : target.slice(0, question);
  const rawQuery = question < 0 ? '' : target.slice(question + 1);
  const request = buildRequest({ query: readQueryValues(rawQuery), headers, body });

  writeResponse(outgoing, await answerRequest(libretto, incoming.method ?? '', rawPath, request));
}

/** answerRequest Finds the Route a Method and a raw Path Name, or Answers the Unknown. */
export function answerRequest(libretto: Libretto, method: string, rawPath: string, request: Request): Promise<Response> {
  for (const route of libretto.routes) {
    const names = matchRoutePattern(route.pattern, rawPath);
    if (names !== undefined && route.method === method) {
      return route.handle({ ...request, path: names });
    }
  }

  return libretto.unknown(request);
}

/** matchRoutePattern Reads the Names a Pattern Declares: /students/{id} Holds id. /students/ is not /students. */
export function matchRoutePattern(pattern: string, rawPath: string): Record<string, string> | undefined {
  const wanted = pattern.split('/');
  const given = rawPath.split('/');
  if (wanted.length !== given.length) {
    return undefined;
  }

  const names: Record<string, string> = {};
  for (const [position, part] of wanted.entries()) {
    const segment = given[position] ?? '';
    if (part.startsWith('{') && part.endsWith('}')) {
      if (segment === '') {
        return undefined;
      }
      try {
        names[part.slice(1, -1)] = decodeURIComponent(segment);
      } catch {
        return undefined;
      }
    } else if (part !== segment) {
      return undefined;
    }
  }

  return names;
}

/** readQueryValues Takes the first Value of a repeated Key; a Pair that cannot be Decoded is Dropped. */
export function readQueryValues(rawQuery: string): Record<string, string> {
  const values: Record<string, string> = {};
  if (rawQuery === '') {
    return values;
  }

  for (const pair of rawQuery.split('&')) {
    const equals = pair.indexOf('=');
    const rawName = equals < 0 ? pair : pair.slice(0, equals);
    const rawValue = equals < 0 ? '' : pair.slice(equals + 1);
    try {
      const name = decodeURIComponent(rawName.replaceAll('+', ' '));
      values[name] ??= decodeURIComponent(rawValue.replaceAll('+', ' '));
    } catch {
      continue;
    }
  }

  return values;
}

/** readBoundedBody Reads at most one Byte past the Bound, and Drops a little more so the Socket stays Polite. */
async function readBoundedBody(incoming: IncomingMessage): Promise<Uint8Array> {
  const chunks: Buffer[] = [];
  let kept = 0;
  let seen = 0;
  for await (const chunk of incoming as AsyncIterable<Buffer>) {
    seen += chunk.length;
    if (kept <= MAX_BODY_BYTES) {
      const room = MAX_BODY_BYTES + 1 - kept;
      chunks.push(chunk.subarray(0, room));
      kept += Math.min(room, chunk.length);
    }
    if (seen > MAX_BODY_BYTES + MAX_SWALLOW_BYTES) {
      break;
    }
  }

  return new Uint8Array(Buffer.concat(chunks));
}

/** writeResponse Renders the Script's Response: JSON, and no Body for a 204. */
export function writeResponse(outgoing: ServerResponse, response: Response): void {
  if (response.status === HTTP_NO_CONTENT || response.body === undefined) {
    outgoing.writeHead(response.status).end();

    return;
  }

  const json = Buffer.from(writeJsonText(response.body), 'utf8');
  outgoing.writeHead(response.status, { 'Content-Type': 'application/json', 'Content-Length': json.length }).end(json);
}
