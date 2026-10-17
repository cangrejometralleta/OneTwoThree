/** MAX_BODY_BYTES Bounds what a Client may Send: one mebibyte. */
export const MAX_BODY_BYTES = 1 << 20;

/**
 * Request is what the Door Hands the Script: no Framework Type Crosses it.
 * Header Names Arrive in lower case; the Body is at most one Byte past the Bound.
 */
export interface Request {
  readonly path: Readonly<Record<string, string>>;
  readonly query: Readonly<Record<string, string>>;
  readonly headers: Readonly<Record<string, string>>;
  readonly body: Uint8Array;
  readonly caller: string;
}

/** Response is what the Script Hands back: a Status, and a Body to Render or undefined. */
export interface Response {
  readonly status: number;
  readonly body: unknown;
}

/** Handler Answers one Request. */
export type Handler = (request: Request) => Promise<Response>;

/** Route Declares one Door: Method, Pattern and the Handler behind it. */
export interface Route {
  readonly method: string;
  readonly pattern: string;
  readonly handle: Handler;
}

/**
 * Libretto is every Route a Service Declares, and the one Answer for a Door it does not Hold.
 * The Unknown Answers before any Guard, because there is no Route to Guard.
 */
export interface Libretto {
  readonly routes: readonly Route[];
  readonly unknown: Handler;
}

/** Running is a Server that Listens until it is Closed. */
export interface Running {
  readonly port: number;
  close(): Promise<void>;
}

/** Server Serves a Libretto on a Port; each Framework Fills it in its own Module. */
export interface Server {
  serveRoutes(libretto: Libretto, port: number): Promise<Running>;
}

/** buildRequest Fills the Fields a Test or an Adapter Leaves out. */
export function buildRequest(parts: Partial<Request> = {}): Request {
  return {
    path: parts.path ?? {},
    query: parts.query ?? {},
    headers: parts.headers ?? {},
    body: parts.body ?? new Uint8Array(),
    caller: parts.caller ?? '',
  };
}
