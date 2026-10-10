import test from 'node:test';
import assert from 'node:assert/strict';
import { Handler } from '../src/handler/handler.ts';
import { SchoolService } from '../src/app/service.ts';
import { HmacTokenIssuer } from '../src/tokens/token.ts';
import { Fault } from '../src/faults/index.ts';
import { studentFromBody, validRut } from '../src/school/index.ts';
import { parseJson, stringifyJson } from '../src/codec/json.ts';
import { statusForFault } from '../src/app/service.ts';
import { FileStore } from '../src/store/store.ts';
import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';
import { loadSettings } from '../src/settings/loader.ts';
import type { Request } from '../src/transport/index.ts';

const valid = { rut:'12345678-5', name:'  Ana ', age:18n, courseId:1n };
const course = { id:1n, code:'CS', name:'Systems' };
function fixture(options: {registered?: boolean; registryError?: boolean; noticeError?: boolean} = {}) {
  const events: string[] = []; const rows = new Map<bigint,any>(); let next = 1n;
  const students = {
    async selectStudent(id:bigint) { events.push('student.select'); return rows.get(id); },
    async insertStudent(s:any) { events.push('student.insert'); if ([...rows.values()].some(x=>x.rut===s.rut)) throw new Fault('RutTaken'); const v={...s,id:next++}; rows.set(v.id,v); return v; },
    async updateStudent(id:bigint,s:any) { events.push('student.update'); if(!rows.has(id)) return undefined; const v={...s,id}; rows.set(id,v); return v; },
    async deleteStudent(id:bigint) { events.push('student.delete'); return rows.delete(id); },
    async listStudents() { events.push('student.list'); return [...rows.values()]; }
  };
  const courses = {
    async selectCourse(id:bigint) { events.push('course.select'); return id===1n?course:undefined; },
    async insertCourse(c:any) { events.push('course.insert'); return {...c,id:1n}; },
    async listCourses() { events.push('course.list'); return [course]; }
  };
  const registry = { async confirm(rut:string) { events.push('registry.confirm'); if(options.registryError) throw new Fault('RegistryUnavailable'); return options.registered !== false && rut==='12345678-5'; } };
  const notifier = { async announce(id:bigint) { events.push('notifier.announce'); if(options.noticeError) throw new Error('private error'); } };
  const logs:string[]=[];
  const service = new SchoolService(students,courses,registry,notifier,{minimumAgeYears:18n},line=>logs.push(line));
  const tokens = new HmacTokenIssuer('spec-secret',3600n); const handler = new Handler(service,tokens);
  const auth = `Bearer ${tokens.issue()}`;
  const request = (method:string,path:string,body?:unknown,authorization:string|null=auth):Request => ({ method,path,headers:authorization?{authorization}:{},body:body===undefined?new Uint8Array():new TextEncoder().encode(typeof body==='string'?body:stringifyJson(body)) });
  return { events,logs,rows,service,tokens,handler,request };
}

test('RUT check digit and canonical lowercase are enforced', async () => {
  assert.equal(validRut('12345678-5'),true); assert.equal(validRut('12345678-K'),false); assert.equal(validRut('6-K'),true);
  assert.equal(studentFromBody({...valid,rut:'6-K'},{minimumAgeYears:18n}).rut,'6-k');
  const f=fixture(); const kept=await f.service.enroll({...valid,rut:'12345678-5'}); assert.equal(kept.rut,'12345678-5');
});
test('student and course form rules use the configured age floor and declared order', async () => {
  assert.throws(()=>studentFromBody({...valid,age:20n},{minimumAgeYears:21n}),(e:any)=>e.kind==='AgeTooLow' && e.message==='age Must be 21 or more');
  const f=fixture();
  await assert.rejects(f.service.createCourse({code:' ',name:' '}),(e:any)=>e.kind==='CodeEmpty');
  await assert.rejects(f.service.createCourse({code:'CS',name:'  '}),(e:any)=>e.kind==='NameEmpty');
  assert.equal(f.events.includes('course.insert'),false);
});
test('admission follows form, course, registry, store and announcement order', async () => {
  const f=fixture(); const result=await f.service.enroll(valid);
  assert.equal(result.id,1n); assert.deepEqual(f.events,['course.select','registry.confirm','student.insert','notifier.announce']);
});
test('invalid form precedes course and registry; course precedes registry', async () => {
  const f=fixture(); await assert.rejects(f.service.enroll({...valid,rut:'bad'}), (e:any)=>e.kind==='RutInvalid');
  assert.deepEqual(f.events,[]);
  await assert.rejects(f.service.enroll({...valid,courseId:2n}), (e:any)=>e.kind==='CourseUnknown');
  assert.deepEqual(f.events,['course.select']);
});
test('registry denial and outage are controlled faults and never write', async () => {
  const denied=fixture({registered:false}); await assert.rejects(denied.service.enroll(valid),(e:any)=>e.kind==='RutUnregistered');
  assert.equal(denied.events.includes('student.insert'),false);
  const outage=fixture({registryError:true}); await assert.rejects(outage.service.enroll(valid),(e:any)=>e.kind==='RegistryUnavailable');
});
test('failed announcement is logged and enrollment still succeeds', async () => {
  const f=fixture({noticeError:true}); const result=await f.service.enroll(valid);
  assert.equal(result.id,1n); assert.deepEqual(f.logs,['student 1 is enrolled']);
});
test('rewrite checks registry before store and never announces', async () => {
  const f=fixture(); await f.service.enroll(valid); f.events.length=0;
  await f.service.save(1n,{...valid,name:'New'});
  assert.deepEqual(f.events,['course.select','registry.confirm','student.update']);
});
test('token format matches the cross-runtime vector and expiration edge', () => {
  const tokens=new HmacTokenIssuer('spec-secret',4102444800n-100n);
  const pinned='eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjo0MTAyNDQ0ODAwfQ.fgl9OtJeU8f90EIkPjELn34yWKjIbn3QcYn5UFoLAyY';
  assert.equal(tokens.verify(pinned,4102444800n),'student-registry');
  assert.equal(tokens.verify(pinned,4102444801n),undefined);
  assert.equal(tokens.verify('eyJzdWIiOiJzdHVkZW50LXJlZ2lzdHJ5IiwiZXhwIjoxfQ.V8VQ8Ie_fo88z4EBJ0QbtIc5ONTOqQ_-l_WLuhwz5t4',2n),undefined);
});
test('route unknown wins before auth, guard wins before path and body', async () => {
  const f=fixture();
  assert.deepEqual((await f.handler.handle(f.request('PATCH','/students'))).body,{error:'route not Found'});
  assert.deepEqual((await f.handler.handle(f.request('GET','/students/abc',undefined,null))).body,{error:'token is Invalid or Expired'});
  assert.equal((await f.handler.handle(f.request('GET','/students/abc'))).status,400);
  assert.equal((await f.handler.handle(f.request('POST','/students','{'))).body && ((await f.handler.handle(f.request('POST','/students','{'))).body as any).error,'body is not valid JSON');
});
test('token route and course route use the declared successes', async () => {
  const f=fixture();
  assert.equal((await f.handler.handle(f.request('POST','/token',undefined,null))).status,201);
  const made=await f.handler.handle(f.request('POST','/courses',{code:'CS',name:'Systems'})); assert.equal(made.status,201);
  assert.equal((await f.handler.handle(f.request('GET','/courses?page=0&size=1'))).status,200);
});
test('page syntax is strict, first repeated value wins, and out-of-range windows are empty', async () => {
  const f=fixture();
  assert.equal((await f.handler.handle(f.request('GET','/students?page=abc'))).body && ((await f.handler.handle(f.request('GET','/students?page=abc'))).body as any).error,'page or size is not valid');
  assert.deepEqual((await f.handler.handle(f.request('GET','/students?page=999999999999999999999999&size=1'))).body,[]);
  assert.deepEqual((await f.handler.handle(f.request('GET','/students?page=0&page=bad'))).body,[]);
});
test('JSON codec preserves signed 64-bit integer values and rejects a second value', () => {
  const n=parseJson('{"n":9223372036854775807}') as any;
  assert.equal(n.n,9223372036854775807n); assert.equal(stringifyJson(n),' {"n":9223372036854775807}'.trim());
  assert.throws(()=>parseJson('{} {}'));
});
test('each Fault kind has its declared edge status', () => {
  const expected:Record<string,number>={NameEmpty:400,CodeEmpty:400,RutInvalid:400,AgeTooLow:400,PageInvalid:400,BodyBroken:400,PathBroken:400,TokenInvalid:401,StudentUnknown:404,CourseUnknown:404,RouteUnknown:404,RutTaken:409,RutUnregistered:422,RegistryUnavailable:503};
  for(const [kind,status] of Object.entries(expected)) assert.equal(statusForFault(kind as any),status,kind);
});
test('student routes preserve path, body and page check order and status', async () => {
  const f=fixture();
  const created=await f.handler.handle(f.request('POST','/students',valid)); assert.equal(created.status,201);
  const read=await f.handler.handle(f.request('GET','/students/01')); assert.equal(read.status,200); assert.equal((read.body as any).id,1n);
  assert.equal((await f.handler.handle(f.request('PUT','/students/nope','{}'))).status,400);
  assert.equal((await f.handler.handle(f.request('PUT','/students/90',{...valid,rut:'invalid'}))).status,400);
  assert.equal((await f.handler.handle(f.request('PUT','/students/90',valid))).status,404);
  assert.equal((await f.handler.handle(f.request('DELETE','/students/1'))).status,204);
  assert.equal((await f.handler.handle(f.request('DELETE','/students/1'))).status,404);
});
test('body schema, UTF-8, integer bounds and size are strict while Content-Type is ignored', async () => {
  const f=fixture(); const req=f.request('POST','/courses','{"code":"CS","name":"Systems"}'); req.headers['content-type']='text/plain';
  assert.equal((await f.handler.handle(req)).status,201);
  for(const body of ['null','[]','{"code":"CS","name":"Systems","extra":true}','{"code":"CS"}','{"code":"CS","name":null}','{} {}','{"code":"CS","name":"X","id":1}'])
    assert.deepEqual((await f.handler.handle(f.request('POST','/courses',body))).body,{error:'body is not valid JSON'});
  const tooLarge=f.request('POST','/courses',' '.repeat(1024*1024+1)); assert.equal((await f.handler.handle(tooLarge)).status,400);
  const badNumber=f.request('POST','/students','{"rut":"12345678-5","name":"A","age":9223372036854775808,"courseId":1}'); assert.equal((await f.handler.handle(badNumber)).status,400);
});
test('FileStore persists records, orders lists and proves RUT uniqueness', async () => {
  const dir=await mkdtemp(join(tmpdir(),'school-ts-')); const path=join(dir,'school.json');
  try {
    const store=await FileStore.open(path); const keptCourse=await store.insertCourse({code:'C',name:'Course'});
    const a=await store.insertStudent({id:0n,...valid,name:'Zed'});
    await store.insertStudent({id:0n,...valid,rut:'12345670-7',name:'Ana'});
    assert.equal(keptCourse.id,1n); assert.deepEqual((await store.listStudents()).map(s=>s.name),['Ana','Zed']);
    await assert.rejects(store.insertStudent({id:0n,...valid,name:'Other'}),(e:any)=>e.kind==='RutTaken');
    const reopened=await FileStore.open(path); assert.equal((await reopened.selectStudent(a.id))?.name,'Zed');
    const disk=await readFile(path,'utf8'); assert.match(disk,/"nextStudentId":"3"/);
  } finally { await rm(dir,{recursive:true,force:true}); }
});
test('startup loads strict layered configuration and digit-only environment overrides', async () => {
  const dir=await mkdtemp(join(tmpdir(),'school-config-'));
  try {
    const program=join(dir,'program'); const data=join(dir,'data'); await import('node:fs/promises').then(fs=>fs.mkdir(program,{recursive:true}));
    await import('node:fs/promises').then(fs=>fs.mkdir(join(data,'config'),{recursive:true}));
    await import('node:fs/promises').then(fs=>fs.mkdir(join(data,'constants'),{recursive:true}));
    await writeFile(join(data,'config','school.defaults.json'),'{"port":8000,"databasePath":"db","registryPath":"registry.json","tokenLifeSeconds":600,"serverAdapter":"stdlib"}');
    await writeFile(join(data,'config','school.development.json'),'{"port":8001}');
    await writeFile(join(data,'constants','school.json'),'{"minimumAgeYears":18}');
    await writeFile(join(program,'.env'),'SCHOOL_PORT=8002\nTOKEN_SECRET=from-file\n');
    const config=await loadSettings(pathToFileURL(join(program,'main.ts')).href,{APP_ENV:'development',DATA_ROOT:data});
    assert.equal(config.port,8002); assert.equal(config.tokenSecret,'from-file'); assert.equal(config.databasePath,join(process.cwd(),'db'));
    await assert.rejects(loadSettings(pathToFileURL(join(program,'main.ts')).href,{APP_ENV:'development',DATA_ROOT:data,SCHOOL_UNKNOWN:'x'}),/unknown environment variable SCHOOL_UNKNOWN/);
    await assert.rejects(loadSettings(pathToFileURL(join(program,'main.ts')).href,{APP_ENV:'development',DATA_ROOT:data,SCHOOL_PORT:''}),/digits only/);
  } finally { await rm(dir,{recursive:true,force:true}); }
});
test('module imports keep vendors at their package boundaries', async () => {
  const read=async(path:string)=>readFile(new URL(path,import.meta.url),'utf8');
  const [core,door,transport,codec,app,server]=await Promise.all([
    read('../src/school/index.ts'),read('../src/handler/handler.ts'),read('../src/transport/index.ts'),
    read('../src/codec/json.ts'),read('../src/app/service.ts'),read('../src/serving/stdlib.ts')
  ]);
  assert.doesNotMatch(core,/from ['"]node:|from ['"].*\/(store|serving|tokens|campus)\//);
  assert.doesNotMatch(door,/from ['"].*\/(store|campus|serving)\//);
  assert.doesNotMatch(transport,/\bimport\b/);
  assert.match(codec,/JSON\.parse/); assert.match(codec,/JSON\.stringify/);
  assert.doesNotMatch(app,/from ['"]node:/);
  assert.match(server,/from ['"]node:http['"]/);
  await import('../src/serving/stdlib.ts');
});
