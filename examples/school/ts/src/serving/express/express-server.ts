import { createServer } from 'node:http';

import express from 'express';

import type { Libretto, Running, Server } from '../../transport/transport.ts';
import { serveLibretto } from '../reading.ts';
import { listen } from '../stdlib.ts';

/**
 * ExpressServer Serves a Libretto through Express, and is the only Module that Imports it.
 * One catch-all Middleware Hands every Request to the Libretto, so Express Decides nothing.
 * Reference: https://expressjs.com/en/5x/api.html#app.use
 */
export class ExpressServer implements Server {
  serveRoutes(libretto: Libretto, port: number): Promise<Running> {
    const application = express();
    application.disable('x-powered-by');
    application.use((incoming, outgoing) => void serveLibretto(libretto, incoming, outgoing));

    return listen(createServer(application), port);
  }
}
