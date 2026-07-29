import { createServer, type IncomingMessage, type ServerResponse } from "node:http";

import express from "express";
import Fastify from "fastify";

import { listPatternParams, readBearerToken, type Request, type Response, type Route, type Server } from "./transport.js";

// NodeServer Serves the Routes with node:http alone.
// Reference: https://nodejs.org/api/http.html#httpcreateserveroptions-requestlistener
export class NodeServer implements Server {
  async serveRoutes(routes: Route[], port: number): Promise<void> {
    const server = createServer((req, res) => this.dispatchNodeRequest(routes, req, res));

    await new Promise<void>((resolve) => server.listen(port, resolve));
  }

  // dispatchNodeRequest Matches by hand, because node:http Has no Router.
  private async dispatchNodeRequest(routes: Route[], req: IncomingMessage, res: ServerResponse) {
    const url = new URL(req.url ?? "/", "http://localhost");
    const found = matchRoutePattern(routes, req.method ?? "GET", url.pathname);

    if (!found) return writeReplyAsJSON(res, { status: 404, body: { error: "no such Route" } });

    const body = await readRequestStream(req);
    writeReplyAsJSON(res, found.route.handle({
      path: found.path,
      query: Object.fromEntries(url.searchParams),
      token: readBearerToken(req.headers.authorization),
      body,
    }));
  }
}

// ExpressServer Serves the same Routes with express.
// Reference: https://expressjs.com/en/5x/api.html
export class ExpressServer implements Server {
  async serveRoutes(routes: Route[], port: number): Promise<void> {
    const app = express();
    app.use(express.text({ type: "*/*" }));

    for (const route of routes) {
      app[route.method.toLowerCase() as "get"](translateRoutePattern(route.pattern), (req, res) => {
        const reply = route.handle({
          path: req.params as Record<string, string>,
          query: req.query as Record<string, string>,
          token: readBearerToken(req.headers.authorization),
          body: typeof req.body === "string" ? req.body : "",
        });

        res.status(reply.status).json(reply.body ?? undefined);
      });
    }

    await new Promise<void>((resolve) => { app.listen(port, () => resolve()); });
  }
}

// FastifyServer Serves the same Routes with fastify.
// Reference: https://fastify.dev/docs/latest/Reference/Routes/
export class FastifyServer implements Server {
  async serveRoutes(routes: Route[], port: number): Promise<void> {
    const app = Fastify();

    // Fastify Parses JSON before we See it, so its Parser Comes out first.
    // Reference: https://fastify.dev/docs/latest/Reference/ContentTypeParser/
    app.removeAllContentTypeParsers();
    app.addContentTypeParser("*", { parseAs: "string" }, (_, body, done) => done(null, body));

    for (const route of routes) {
      app.route({
        method: route.method,
        url: translateRoutePattern(route.pattern),
        handler: async (req, reply) => {
          const answer = route.handle({
            path: req.params as Record<string, string>,
            query: req.query as Record<string, string>,
            token: readBearerToken(req.headers.authorization),
            body: typeof req.body === "string" ? req.body : "",
          });

          return reply.code(answer.status).send(answer.body ?? undefined);
        },
      });
    }

    await app.listen({ port });
  }
}

// translateRoutePattern Turns {id} into :id, which express and fastify Want.
export function translateRoutePattern(pattern: string): string {
  return pattern
    .split("/")
    .map((part) => (part.startsWith("{") && part.endsWith("}") ? `:${part.slice(1, -1)}` : part))
    .join("/");
}

// matchRoutePattern is the Router node:http Does not Ship.
export function matchRoutePattern(routes: Route[], method: string, pathname: string) {
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

// selectServerAdapter Picks the Framework at Startup, never at Build Time.
export function selectServerAdapter(name: string): Server {
  if (name === "express") return new ExpressServer();
  if (name === "fastify") return new FastifyServer();

  return new NodeServer();
}
