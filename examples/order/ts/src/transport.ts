// Request is what a Handler Receives.
// No Framework Type Appears here, and that is the whole Point.
export type Request = {
  path: Record<string, string>;
  query: Record<string, string>;
  body: string;
};

// Response is what a Handler Returns.
export type Response = {
  status: number;
  body: unknown;
};

// Handler is the only Signature the Business Layer Knows.
export type Handler = (request: Request) => Response;

// Route Binds one Method and one Pattern to one Handler.
// Patterns Use {name}, and the Adapter Translates from there.
export type Route = {
  method: "GET" | "POST";
  pattern: string;
  handle: Handler;
};

// listPatternParams Finds every {name} inside a Route Pattern.
export function listPatternParams(pattern: string): string[] {
  return pattern
    .split("/")
    .filter((part) => part.startsWith("{") && part.endsWith("}"))
    .map((part) => part.slice(1, -1));
}
