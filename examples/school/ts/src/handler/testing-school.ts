import { FrontDesk, Shelf } from '../app/fakes.ts';
import { SchoolService } from '../app/service.ts';
import { AccessTokens } from '../tokens/access-tokens.ts';
import { buildRequest, type Request, type Response } from '../transport/transport.ts';
import { CALLER_SUBJECT, SchoolHandler } from './handler.ts';

/** TestingSchool Hands the Door its Fakes: no Database, no Framework, no Port. */
export class TestingSchool {
  readonly books = new Shelf();
  readonly front = new FrontDesk();
  readonly tokens = new AccessTokens('test', 60n);
  readonly door = new SchoolHandler(new SchoolService({
    rules: { minimumAgeYears: 18n }, students: this.books, courses: this.books,
    registry: this.front, notices: this.front, log: () => {},
  }), this.tokens);

  /** bearer Arrives with a Token the Door Minted. */
  bearer(): Readonly<Record<string, string>> {
    return { authorization: `Bearer ${this.tokens.issueAccessToken(CALLER_SUBJECT)}` };
  }

  /** callRoute Arrives the Way a Client Arrives: through the Libretto, the Answer and the Guard. */
  callRoute(method: string, pattern: string, parts: Partial<Request> = {}, body = ''): Promise<Response> {
    const request = buildRequest({ ...parts, headers: this.bearer(), body: new TextEncoder().encode(body) });
    const route = this.door.declareSchoolRoutes().routes.find((each) => each.method === method && each.pattern === pattern);
    if (route === undefined) {
      throw new Error(`no Route Answers ${method} ${pattern}`);
    }

    return route.handle(request);
  }
}
