import { answerWith, requireProvenCaller, type Telling } from '../app/answers.ts';
import type { TokenIssuer } from '../app/ports.ts';
import type { SchoolService } from '../app/service.ts';
import { readJsonBody, readPageRequest, readPathNumber, ROUTE_UNKNOWN } from '../app/validation.ts';
import { UNNUMBERED_STUDENT, type CourseId, type StudentId } from '../school/school.ts';
import type { Handler, Libretto, Request } from '../transport/transport.ts';
import type { CourseView, StudentView, TokenView } from '../wire/wire.ts';
import {
  buildCourseRecord, buildStudentRecord, COURSE_FIELDS, renderCourseView, renderStudentView, STUDENT_FIELDS,
} from './crossings.ts';

/** CALLER_SUBJECT is the one Subject a Minted Token Names. */
export const CALLER_SUBJECT = 'student-registry';

const HTTP_OK = 200;
const HTTP_CREATED = 201;
const HTTP_NO_CONTENT = 204;

/**
 * SchoolHandler is the Door: it Orchestrates. It Decodes, Calls one Service method
 * and Answers with a Value or Fails. It Reaches no Provider and no Vendor.
 */
export class SchoolHandler {
  readonly #school: SchoolService;
  readonly #tokens: TokenIssuer;

  constructor(school: SchoolService, tokens: TokenIssuer) {
    this.#school = school;
    this.#tokens = tokens;
  }

  /** declareSchoolRoutes is the Libretto: every Door, its happy Status and its Guard. */
  declareSchoolRoutes(): Libretto {
    return {
      routes: [
        { method: 'POST', pattern: '/token', handle: answerWith(HTTP_CREATED, (request) => this.mintAccessToken(request)) },

        { method: 'GET', pattern: '/students', handle: this.guarded(HTTP_OK, (request) => this.listStudentRecords(request)) },
        { method: 'POST', pattern: '/students', handle: this.guarded(HTTP_CREATED, (request) => this.addStudentRecord(request)) },
        { method: 'GET', pattern: '/students/{id}', handle: this.guarded(HTTP_OK, (request) => this.showStudentRecord(request)) },
        { method: 'PUT', pattern: '/students/{id}', handle: this.guarded(HTTP_OK, (request) => this.saveStudentRecord(request)) },
        { method: 'DELETE', pattern: '/students/{id}', handle: this.guarded(HTTP_NO_CONTENT, (request) => this.dropStudentRecord(request)) },

        { method: 'GET', pattern: '/courses', handle: this.guarded(HTTP_OK, (request) => this.listCourseRecords(request)) },
        { method: 'POST', pattern: '/courses', handle: this.guarded(HTTP_CREATED, (request) => this.addCourseRecord(request)) },
        { method: 'GET', pattern: '/courses/{id}', handle: this.guarded(HTTP_OK, (request) => this.showCourseRecord(request)) },
      ],
      unknown: answerWith(HTTP_OK, async () => this.refuseUnknownRoute()),
    };
  }

  /** guarded Names the Caller before the Script runs, then Crosses to a Handler. */
  guarded(status: number, tell: Telling): Handler {
    return answerWith(status, requireProvenCaller(this.#tokens, tell));
  }

  /** mintAccessToken Hands out a Token. */
  async mintAccessToken(_request: Request): Promise<TokenView> {
    return { token: this.#tokens.issueAccessToken(CALLER_SUBJECT) };
  }

  /** listStudentRecords Tells one Page of the Roll. */
  async listStudentRecords(request: Request): Promise<StudentView[]> {
    const students = await this.#school.listStudents(readPageRequest(request));

    return students.map(renderStudentView);
  }

  /** addStudentRecord Enrols a new Student. */
  async addStudentRecord(request: Request): Promise<StudentView> {
    const body = readJsonBody(request, STUDENT_FIELDS);

    return renderStudentView(await this.#school.enrollStudent(buildStudentRecord(body, UNNUMBERED_STUDENT)));
  }

  /** showStudentRecord Tells the Story of one Student. */
  async showStudentRecord(request: Request): Promise<StudentView> {
    const id = readPathNumber(request) as StudentId;

    return renderStudentView(await this.#school.readStudent(id));
  }

  /** saveStudentRecord Rewrites the Student the Path Names. */
  async saveStudentRecord(request: Request): Promise<StudentView> {
    const id = readPathNumber(request) as StudentId;
    const body = readJsonBody(request, STUDENT_FIELDS);

    return renderStudentView(await this.#school.saveStudent(buildStudentRecord(body, id)));
  }

  /** dropStudentRecord Ends an Enrolment. It Answers no Body. */
  async dropStudentRecord(request: Request): Promise<undefined> {
    await this.#school.deleteStudent(readPathNumber(request) as StudentId);

    return undefined;
  }

  /** listCourseRecords Tells one Page of the Catalogue. */
  async listCourseRecords(request: Request): Promise<CourseView[]> {
    const courses = await this.#school.listCourses(readPageRequest(request));

    return courses.map(renderCourseView);
  }

  /** addCourseRecord Opens a new Course. */
  async addCourseRecord(request: Request): Promise<CourseView> {
    return renderCourseView(await this.#school.createCourse(buildCourseRecord(readJsonBody(request, COURSE_FIELDS))));
  }

  /** showCourseRecord Tells the Story of one Course. */
  async showCourseRecord(request: Request): Promise<CourseView> {
    return renderCourseView(await this.#school.readCourse(readPathNumber(request) as CourseId));
  }

  /** refuseUnknownRoute Answers a Door the Libretto does not Hold, before any Guard. */
  refuseUnknownRoute(): never {
    throw ROUTE_UNKNOWN;
  }
}
