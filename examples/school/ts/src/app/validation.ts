import { readJsonText, type JsonValue } from '../codec/json.ts';
import { refuseInvalidInput, reportMissingRecord } from '../faults/faults.ts';
import { buildPage, DEFAULT_PAGE_SIZE, LARGEST_WHOLE, PAGE_IS_INVALID, WHOLE_SET, type Page } from '../school/school.ts';
import { MAX_BODY_BYTES, type Request } from '../transport/transport.ts';

/** The Form of a Request, Checked before its Meaning: the Path, the Body, the Page. */
export const BODY_IS_BROKEN = refuseInvalidInput('body is not valid JSON');
export const PATH_IS_BROKEN = refuseInvalidInput('path Holds no valid Identity');
export const ROUTE_UNKNOWN = reportMissingRecord('route not Found');

const DIGITS = /^[0-9]+$/;
const BEARER_SCHEME = 'bearer ';

/** readClampedNumber Reads Digits as far as 64 bits can Hold them. */
function readClampedNumber(digits: string): bigint {
  const number = BigInt(digits);

  return number > LARGEST_WHOLE ? LARGEST_WHOLE : number;
}

/**
 * readPathNumber Reads the positive whole Number in the Path's {id}; 012 is 12.
 * A Number beyond 64 bits Names nobody, so it Reads as the largest, which no one Holds.
 */
export function readPathNumber(request: Request): bigint {
  const text = request.path['id'] ?? '';
  if (!DIGITS.test(text)) {
    throw PATH_IS_BROKEN;
  }

  const number = readClampedNumber(text);
  if (number === 0n) {
    throw PATH_IS_BROKEN;
  }

  return number;
}

/** FieldKinds Names each Field of a Body and whether it is Text or a whole Number. */
export type FieldKinds = Readonly<Record<string, 'text' | 'whole'>>;

/** BodyOf Turns FieldKinds into the Shape it Reads. */
export type BodyOf<K extends FieldKinds> = { readonly [F in keyof K]: K[F] extends 'text' ? string : bigint };

/**
 * readJsonBody Reads one JSON Object, in UTF-8, whose Fields are exactly those Named, none null,
 * each a String or a whole Number of 64 bits. Anything Less, or More, is Body is Broken.
 */
export function readJsonBody<K extends FieldKinds>(request: Request, kinds: K): BodyOf<K> {
  if (request.body.length > MAX_BODY_BYTES) {
    throw BODY_IS_BROKEN;
  }

  const sent = readJsonText(request.body);
  const names = Object.keys(kinds);
  if (!(sent instanceof Map) || sent.size !== names.length) {
    throw BODY_IS_BROKEN;
  }

  const body: Record<string, string | bigint> = {};
  for (const name of names) {
    body[name] = readField(sent.get(name), kinds[name]);
  }

  return body as BodyOf<K>;
}

function readField(value: JsonValue | undefined, kind: 'text' | 'whole' | undefined): string | bigint {
  if (kind === 'text' && typeof value === 'string') {
    return value;
  }
  if (kind === 'whole' && typeof value === 'bigint' && value >= -LARGEST_WHOLE - 1n && value <= LARGEST_WHOLE) {
    return value;
  }

  throw BODY_IS_BROKEN;
}

/** readPageRequest Reads the Window a list Route Asks for. Neither Named Asks for the whole Set; one alone Takes the other's Default. */
export function readPageRequest(request: Request): Page {
  const page = request.query['page'];
  const size = request.query['size'];
  if (page === undefined && size === undefined) {
    return WHOLE_SET;
  }

  return buildPage(page === undefined ? 0n : readPageNumber(page), size === undefined ? DEFAULT_PAGE_SIZE : readPageNumber(size));
}

function readPageNumber(text: string): bigint {
  if (!DIGITS.test(text)) {
    throw PAGE_IS_INVALID;
  }

  return readClampedNumber(text);
}

/** readBearerToken Reads the Token after a Bearer Scheme of any Case, or Nothing. */
export function readBearerToken(request: Request): string {
  const header = request.headers['authorization'] ?? '';
  if (header.slice(0, BEARER_SCHEME.length).toLowerCase() !== BEARER_SCHEME) {
    return '';
  }

  return header.slice(BEARER_SCHEME.length);
}
