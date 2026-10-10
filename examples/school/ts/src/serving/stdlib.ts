import { createServer, type IncomingMessage, type Server, type ServerResponse } from 'node:http';
import { Handler, encodeResponse } from '../handler/handler.ts';
import type { Request } from '../transport/index.ts';
export function serve(handler: Handler, port: number): Promise<Server> {
  return new Promise((resolve, reject) => {
    const server = createServer(async (incoming: IncomingMessage, outgoing: ServerResponse) => {
      try {
        const chunks: Buffer[] = []; let length = 0;
        for await (const chunk of incoming) {
          const data = Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk);
          if (length <= 1024 * 1024) { const remaining = 1024 * 1024 + 1 - length; chunks.push(data.subarray(0,remaining)); }
          length += data.length;
        }
        const req: Request = { method: incoming.method ?? '', path: incoming.url ?? '/', headers: incoming.headers,
          body: Buffer.concat(chunks) };
        const response = encodeResponse(await handler.handle(req));
        outgoing.statusCode = response.status; for (const [key,value] of Object.entries(response.headers)) outgoing.setHeader(key,value);
        outgoing.end(response.body);
      } catch { outgoing.statusCode = 500; outgoing.setHeader('content-type','application/json'); outgoing.end('{"error":"internal error"}'); }
    });
    server.once('error',reject); server.listen(port,'0.0.0.0',() => { server.removeListener('error',reject); resolve(server); });
  });
}
