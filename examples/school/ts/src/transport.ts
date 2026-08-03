// Request is what a Handler Receives.
// No Framework Type Appears here, and that is the whole Point.
export type Request = {
  path: Record<string, string>;
  query: Record<string, string>;
  token: string;
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
// Patterns Use {name}, and each Adapter Translates from there.
export type Route = {
  method: "GET" | "POST" | "PUT" | "DELETE";
  pattern: string;
  handle: Handler;
};

// Server is the Contract every HTTP Adapter Fulfils.
// Swapping Frameworks Means Swapping the Value behind this Interface.
export interface Server {
  serveRoutes(routes: Route[], port: number): Promise<void>;
}

// listPatternParams Finds every {name} inside a Route Pattern.
export function listPatternParams(pattern: string): string[] {
  return pattern
    .split("/")
    .filter((part) => part.startsWith("{") && part.endsWith("}"))
    .map((part) => part.slice(1, -1));
}

// readBearerToken Pulls the Token out of the Authorization Header.
// Reference: https://datatracker.ietf.org/doc/html/rfc6750#section-2.1
export function readBearerToken(header: string | undefined): string {
  return (header ?? "").replace(/^Bearer /, "");
}
