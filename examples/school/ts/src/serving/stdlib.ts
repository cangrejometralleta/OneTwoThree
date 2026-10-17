import { createServer, type Server as NodeServer } from 'node:http';
import type { AddressInfo } from 'node:net';

import type { Libretto, Running, Server } from '../transport/transport.ts';
import { serveLibretto } from './reading.ts';

/**
 * StdlibServer Serves a Libretto with node:http alone.
 * Reference: https://nodejs.org/api/http.html#httpcreateserveroptions-requestlistener
 */
export class StdlibServer implements Server {
  serveRoutes(libretto: Libretto, port: number): Promise<Running> {
    return listen(createServer((incoming, outgoing) => void serveLibretto(libretto, incoming, outgoing)), port);
  }
}

/** listen Starts a Node Server and Hands back the Port it Bound. */
export function listen(server: NodeServer, port: number): Promise<Running> {
  return new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(port, () => {
      resolve({
        port: (server.address() as AddressInfo).port,
        close: () => new Promise((done) => {
          server.closeAllConnections();
          server.close(() => done());
        }),
      });
    });
  });
}
