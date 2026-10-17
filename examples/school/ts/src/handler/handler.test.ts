import assert from 'node:assert/strict';
import { test } from 'node:test';

import { buildRequest } from '../transport/transport.ts';
import { TestingSchool } from './testing-school.ts';

const ENROL = '{"rut":"12345678-5","name":"Ada","age":20,"courseId":1}';
const REWRITE = '{"rut":"12345678-5","name":"Ada Lovelace","age":21,"courseId":1}';
const ONE = { path: { id: '1' } };

test('eachRouteAnswersWithItsDeclaredStatus', async () => {
  const school = new TestingSchool();
  const stories = [
    ['a Token is Minted', 'POST', '/token', {}, '', 201],
    ['someone Enrols', 'POST', '/students', {}, ENROL, 201],
    ['the Roll is Read', 'GET', '/students', {}, '', 200],
    ['one Student is Read', 'GET', '/students/{id}', ONE, '', 200],
    ['an Enrolment is Rewritten', 'PUT', '/students/{id}', ONE, REWRITE, 200],
    ['the Catalogue is Read', 'GET', '/courses', {}, '', 200],
    ['a Course Opens', 'POST', '/courses', {}, '{"code":"FIS-201","name":"Physics"}', 201],
    ['one Course is Read', 'GET', '/courses/{id}', ONE, '', 200],
    ['an Enrolment Ends', 'DELETE', '/students/{id}', ONE, '', 204],
  ] as const;

  for (const [told, method, pattern, parts, body, status] of stories) {
    const reply = await school.callRoute(method, pattern, parts, body);
    assert.equal(reply.status, status, `${told}: ${JSON.stringify(reply.body, (_, v) => typeof v === 'bigint' ? `${v}` : v)}`);
  }
});

test('theLibrettoDeclaresExactlyTheNineRoutesOfTheSpec', () => {
  assert.equal(new TestingSchool().door.declareSchoolRoutes().routes.length, 9);
});

test('aRewriteChangesTheStudentAndADeleteLeavesNoBody', async () => {
  const school = new TestingSchool();
  await school.callRoute('POST', '/students', {}, ENROL);

  assert.deepEqual((await school.callRoute('PUT', '/students/{id}', ONE, REWRITE)).body,
    { id: 1n, rut: '12345678-5', name: 'Ada Lovelace', age: 21n, courseId: 1n });
  assert.deepEqual(await school.callRoute('DELETE', '/students/{id}', ONE), { status: 204, body: undefined });
});

test('theViewEchoesTheRutInLowerCase', async () => {
  const school = new TestingSchool();
  await school.callRoute('POST', '/students', {}, '{"rut":"10000013-K","name":"Ada","age":20,"courseId":1}');

  assert.deepEqual((await school.callRoute('GET', '/students/{id}', ONE)).body,
    { id: 1n, rut: '10000013-k', name: 'Ada', age: 20n, courseId: 1n });
  assert.equal(school.front.asked[0], '10000013-k');
});

test('anEmptyListAnswersAnEmptyArray', async () => {
  assert.deepEqual((await new TestingSchool().callRoute('GET', '/students')).body, []);
});

test('enrolmentSurvivesAFailedAnnouncement', async () => {
  const school = new TestingSchool();
  school.front.noticeFailure = new Error('mail is down');

  assert.equal((await school.callRoute('POST', '/students', {}, ENROL)).status, 201);
});

test('everyRouteButTheTokenRefusesAnUnnamedCaller', async () => {
  for (const route of new TestingSchool().door.declareSchoolRoutes().routes) {
    if (route.pattern === '/token') {
      continue;
    }
    assert.deepEqual(await route.handle(buildRequest(ONE)), { status: 401, body: { error: 'token is Invalid or Expired' } },
      `${route.method} ${route.pattern}`);
  }
});

test('anUnknownRouteAnswersBeforeTheGuard', async () => {
  assert.deepEqual(await new TestingSchool().door.declareSchoolRoutes().unknown(buildRequest()),
    { status: 404, body: { error: 'route not Found' } });
});

test('aGuardedHandlerReceivesTheNamedCaller', async () => {
  const school = new TestingSchool();
  let seen = '';
  const reply = await school.door.guarded(200, async (request) => {
    seen = request.caller;

    return 'told';
  })(buildRequest({ headers: school.bearer() }));

  assert.equal(reply.status, 200);
  assert.equal(seen, 'student-registry');
});

test('anUncontrolledFailureAnswersFiveHundredAndSaysNothingOfTheDriver', async () => {
  const school = new TestingSchool();
  const reply = await school.door.guarded(200, async () => {
    throw new Error('the driver Broke at /var/db');
  })(buildRequest({ headers: school.bearer() }));

  assert.deepEqual(reply, { status: 500, body: { error: 'internal error' } });
});
