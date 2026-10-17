import assert from 'node:assert/strict';
import { test } from 'node:test';

import { buildRequest, type Response } from '../transport/transport.ts';
import { TestingSchool } from './testing-school.ts';

const VALID = '{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}';

type Arrival = (school: TestingSchool) => Promise<Response>;

const enrolWith = (body: string): Arrival => (school) => school.callRoute('POST', '/students', {}, body);
const showing = (pattern: string, id: string): Arrival => (school) => school.callRoute('GET', pattern, { path: { id } });

/**
 * Every Fault the Spec Names, Reached the Way a Caller Reaches it.
 * The Status and the Reason are Spelled here, never Read from the Fault under Test,
 * so a Declaration that Drifts Fails here instead of Agreeing with itself.
 */
test('eachFaultReachesTheEdgeWhole', async () => {
  const cases: [string, Arrival, number, string][] = [
    ['someone Enrols with no Name', enrolWith('{"rut":"12345678-5","name":"","age":20,"courseId":1}'), 400, 'name is Empty'],
    ['someone Opens a Course with no Code', (school) => school.callRoute('POST', '/courses', {}, '{"code":" ","name":"Physics"}'), 400, 'code is Empty'],
    ['someone Enrols with a RUT that Fails its Digit', enrolWith('{"rut":"12345678-9","name":"Ada","age":20,"courseId":1}'), 400, 'rut Fails its Check Digit'],
    ['someone Enrols below the Enrolment Age', enrolWith('{"rut":"12345678-5","name":"Ada","age":9,"courseId":1}'), 400, 'age Must be 18 or more'],
    ['someone Asks for a Page that Counts backwards', (school) => school.callRoute('GET', '/students', { query: { page: '-1' } }), 400, 'page or size is not valid'],
    ['someone Sends a Body that is not JSON', enrolWith('{"rut":'), 400, 'body is not valid JSON'],
    ['someone Sends a Body with an Identity of their own', enrolWith('{"id":5,"rut":"12345678-5","name":"Ada","age":20,"courseId":1}'), 400, 'body is not valid JSON'],
    ['someone Asks for a Student by a Path that Holds no Number', showing('/students/{id}', 'abc'), 400, 'path Holds no valid Identity'],
    ['someone Arrives with no Token at all', (school) => school.door.declareSchoolRoutes().routes[1]!.handle(buildRequest()), 401, 'token is Invalid or Expired'],
    ['someone Asks for a Student who never Enrolled', showing('/students/{id}', '42'), 404, 'student not Found'],
    ['someone Asks for a Student beyond 64 bits', showing('/students/{id}', '99999999999999999999'), 404, 'student not Found'],
    ['someone Asks for a Course that never Opened', showing('/courses/{id}', '99'), 404, 'course not Found'],
    ['someone Enrols into a Course that never Opened', enrolWith('{"rut":"12345678-5","name":"Ada","age":20,"courseId":-1}'), 404, 'course not Found'],
    ['someone Knocks on a Door the Libretto does not Hold', (school) => school.door.declareSchoolRoutes().unknown(buildRequest()), 404, 'route not Found'],
    ['someone Enrols with a RUT already Registered', async (school) => {
      await enrolWith(VALID)(school);

      return enrolWith('{"rut":"12345678-5","name":"Grace","age":22,"courseId":1}')(school);
    }, 409, 'rut is already Registered'],
    ['someone Enrols with a RUT the Registry Denies', (school) => {
      school.front.registered = false;

      return enrolWith(VALID)(school);
    }, 422, 'rut is not in the Registry'],
    ['someone Rewrites a Student nobody Holds with a RUT the Registry Denies', (school) => {
      school.front.registered = false;

      return school.callRoute('PUT', '/students/{id}', { path: { id: '7' } }, VALID);
    }, 422, 'rut is not in the Registry'],
    ['someone Enrols while the Registry cannot Answer', (school) => {
      school.front.down = new Error('registry socket closed');

      return enrolWith(VALID)(school);
    }, 503, 'registry cannot Answer now'],
    ['someone Drops a Student who already Left', (school) => school.callRoute('DELETE', '/students/{id}', { path: { id: '7' } }), 404, 'student not Found'],
    ['someone Rewrites an Enrolment that never Existed', (school) => school.callRoute('PUT', '/students/{id}', { path: { id: '7' } }, VALID), 404, 'student not Found'],
  ];

  for (const [story, arrive, status, reason] of cases) {
    assert.deepEqual(await arrive(new TestingSchool()), { status, body: { error: reason } }, story);
  }
});
