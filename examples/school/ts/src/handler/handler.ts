import { Fault } from '../faults/index.ts';
import type { Course, Student } from '../school/index.ts';
import { parseJson, stringifyJson } from '../codec/json.ts';
import type { Request, Response } from '../transport/index.ts';
import type { TokenIssuer } from '../app/ports.ts';
import { SchoolService, statusForFault } from '../app/service.ts';

const I64_MIN = -(1n << 63n), I64_MAX = (1n << 63n) - 1n;
const json = (status: number, body: unknown): Response => ({ status, headers: { 'content-type':'application/json' }, body });
function routeFor(method: string, path: string): { route: string; id?: string | undefined } | undefined {
  const key = `${method.toUpperCase()} ${path}`;
  const exact = new Set(['POST /token','GET /students','POST /students','GET /courses','POST /courses']);
  if (exact.has(key)) return { route: key };
  const idPath = path.match(/^\/(students|courses)\/([^/]+)$/);
  if (idPath) {
    const allowed = idPath[1] === 'students' ? ['GET','PUT','DELETE'] : ['GET'];
    if (allowed.includes(method.toUpperCase())) return { route: `${method.toUpperCase()} /${idPath[1]}/{id}`, id: idPath[2]! };
  }
  return undefined;
}
function parseId(raw: string): bigint { if (!/^\d+$/.test(raw)) throw new Fault('PathBroken'); const id = BigInt(raw); if (id <= 0n) throw new Fault('PathBroken'); return id; }
function bodyObject(req: Request, allowed: string[]): Record<string, unknown> {
  if (req.body.byteLength > 1024 * 1024) throw new Fault('BodyBroken');
  let raw: string;
  try { raw = new TextDecoder('utf-8', { fatal:true }).decode(req.body); } catch { throw new Fault('BodyBroken'); }
  let value: unknown; try { value = parseJson(raw); } catch { throw new Fault('BodyBroken'); }
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Fault('BodyBroken');
  const obj = value as Record<string, unknown>;
  if (Object.keys(obj).length !== allowed.length || allowed.some(key => !Object.hasOwn(obj,key))) throw new Fault('BodyBroken');
  return obj;
}
function integer(value: unknown): value is bigint { return typeof value === 'bigint' && value >= I64_MIN && value <= I64_MAX; }
function studentBody(req: Request) {
  const b = bodyObject(req, ['rut','name','age','courseId']);
  if (typeof b.rut !== 'string' || typeof b.name !== 'string' || !integer(b.age) || !integer(b.courseId)) throw new Fault('BodyBroken');
  return { rut:b.rut, name:b.name, age:b.age, courseId:b.courseId };
}
function courseBody(req: Request) {
  const b = bodyObject(req, ['code','name']); if (typeof b.code !== 'string' || typeof b.name !== 'string') throw new Fault('BodyBroken');
  return { code:b.code, name:b.name };
}
function page(url: URL): { page: bigint; size: bigint } | undefined {
  const p = url.searchParams.get('page'), s = url.searchParams.get('size');
  if (p === null && s === null) return undefined;
  const parse = (v: string | null, fallback: bigint, min: bigint): bigint => {
    if (v === null) return fallback; if (!/^\d+$/.test(v)) throw new Fault('PageInvalid');
    const n = BigInt(v); if (n < min || n > I64_MAX) { if (n > I64_MAX) return n; throw new Fault('PageInvalid'); } return n;
  };
  return { page:parse(p,0n,0n), size:parse(s,10n,1n) };
}
function windowed<T>(rows: T[], p: {page: bigint;size: bigint} | undefined): T[] {
  if (!p) return rows;
  const start = p.page * p.size; if (start > BigInt(Number.MAX_SAFE_INTEGER)) return [];
  return rows.slice(Number(start), Number(start + p.size));
}
const studentView = (s: Student) => ({id:s.id,rut:s.rut,name:s.name,age:s.age,courseId:s.courseId});
const courseView = (c: Course) => ({id:c.id,code:c.code,name:c.name});
function failure(fault: Fault): Response { return json(statusForFault(fault.kind),{error:fault.message}); }
export class Handler {
  private service: SchoolService; private tokens: TokenIssuer;
  constructor(service: SchoolService, tokens: TokenIssuer) { this.service=service; this.tokens=tokens; }
  async handle(req: Request): Promise<Response> {
    const path = req.path.split('?',1)[0]!; const match = routeFor(req.method,path);
    if (!match) return failure(new Fault('RouteUnknown'));
    if (match.route !== 'POST /token') {
      const auth = req.headers.authorization;
      const m = typeof auth === 'string' ? /^Bearer\s+(.+)$/i.exec(auth) : null;
      if (!m || !this.tokens.verify(m[1]!)) return failure(new Fault('TokenInvalid'));
    }
    try {
      const url = new URL(req.path,'http://school.local');
      switch (match.route) {
        case 'POST /token': return json(201,{token:this.tokens.issue()});
        case 'GET /students': { const p = page(url); return json(200,windowed((await this.service.listStudents()).map(studentView),p)); }
        case 'POST /students': return json(201,studentView(await this.service.enroll(studentBody(req))));
        case 'GET /students/{id}': return json(200,studentView(await this.service.readStudent(parseId(match.id!))));
        case 'PUT /students/{id}': return json(200,studentView(await this.service.save(parseId(match.id!),studentBody(req))));
        case 'DELETE /students/{id}': await this.service.deleteStudent(parseId(match.id!)); return {status:204,headers:{}};
        case 'GET /courses': { const p = page(url); return json(200,windowed((await this.service.listCourses()).map(courseView),p)); }
        case 'POST /courses': return json(201,courseView(await this.service.createCourse(courseBody(req))));
        case 'GET /courses/{id}': return json(200,courseView(await this.service.readCourse(parseId(match.id!))));
      }
      return failure(new Fault('RouteUnknown'));
    } catch (e) { if (e instanceof Fault) return failure(e); return json(500,{error:'internal error'}); }
  }
}
export function encodeResponse(response: Response): { status: number; headers: Record<string,string>; body: Uint8Array } {
  return { status:response.status, headers:response.headers, body:response.status === 204 ? new Uint8Array() : new TextEncoder().encode(stringifyJson(response.body)) };
}
