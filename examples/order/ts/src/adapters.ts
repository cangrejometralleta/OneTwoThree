import { createServer, type IncomingMessage, type ServerResponse } from "node:http";

import { listPatternParams, type Request, type Response, type Route } from "./transport.js";

// serveOrderRoutes Mounts every Route on node:http alone.
// One Service, one Framework: the Portability Point already Lives in School.
// Reference: https://nodejs.org/api/http.html#httpcreateserveroptions-requestlistener
export async function serveOrderRoutes(routes: Route[], port: number): Promise<void> {
  const server = createServer((req, res) => dispatchOrderRequest(routes, req, res));

  await new Promise<void>((resolve) => server.listen(port, resolve));
}

// dispatchOrderRequest Matches by hand, because node:http Has no Router.
async function dispatchOrderRequest(routes: Route[], req: IncomingMessage, res: ServerResponse) {
  const url = new URL(req.url ?? "/", "http://localhost");
  const found = matchRoutePattern(routes, req.method ?? "GET", url.pathname);

  if (!found) return writeReplyAsJSON(res, { status: 404, body: { error: "no such Route" } });

  const body = await readRequestStream(req);
  writeReplyAsJSON(res, found.route.handle({
    path: found.path,
    query: Object.fromEntries(url.searchParams),
    body,
  }));
}

// matchRoutePattern is the Router node:http Does not Ship.
function matchRoutePattern(routes: Route[], method: string, pathname: string) {
  for (const route of routes) {
    const names = listPatternParams(route.pattern);
    const shape = new RegExp(`^${route.pattern.replace(/\{[^}]+\}/g, "([^/]+)")}$`);
    const found = route.method === method ? shape.exec(pathname) : null;

    if (found) {
      const path = Object.fromEntries(names.map((name, index) => [name, found[index + 1]!]));
      return { route, path };
    }
  }

  return null;
}

// readRequestStream Collects the Body node:http Delivers in Pieces.
function readRequestStream(req: IncomingMessage): Promise<string> {
  return new Promise((resolve) => {
    let text = "";
    req.on("data", (chunk) => { text += chunk; });
    req.on("end", () => resolve(text));
  });
}

// writeReplyAsJSON is the one Place that Touches a ServerResponse.
function writeReplyAsJSON(res: ServerResponse, reply: Response) {
  res.writeHead(reply.status, { "Content-Type": "application/json" });
  res.end(reply.body === null ? "" : JSON.stringify(reply.body));
}
