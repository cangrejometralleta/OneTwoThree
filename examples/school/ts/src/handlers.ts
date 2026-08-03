import {
  BusinessError, ErrBodyIsBroken, ErrPageIsInvalid, ErrPathIsBroken,
  checkStudentRecord, type Student, type CourseID, type StudentID,
} from "./domain.js";
import {
  buildCourseRecord, buildStudentRecord, renderCourseView, renderStudentView,
} from "./dto.js";
import type { CourseStore, Page, StudentStore, TokenIssuer } from "./providers.js";
import type { Request, Response, Route } from "./transport.js";

// SchoolAPI Holds the Cast every Story Needs.
// Three Providers, no Libraries: a Test can Hand it three Fakes.
export class SchoolAPI {
  constructor(
    private readonly students: StudentStore,
    private readonly courses: CourseStore,
    private readonly tokens: TokenIssuer,
  ) {}

  // declareSchoolRoutes is the Libretto.
  // Read it once and you Know the whole Service.
  declareSchoolRoutes(): Route[] {
    return [
      { method: "POST", pattern: "/token", handle: (r) => this.mintAccessToken(r) },

      { method: "GET", pattern: "/students", handle: (r) => this.listStudentRecords(r) },
      { method: "POST", pattern: "/students", handle: (r) => this.addStudentRecord(r) },
      { method: "GET", pattern: "/students/{id}", handle: (r) => this.showStudentRecord(r) },
      { method: "PUT", pattern: "/students/{id}", handle: (r) => this.saveStudentRecord(r) },
      { method: "DELETE", pattern: "/students/{id}", handle: (r) => this.dropStudentRecord(r) },

      { method: "GET", pattern: "/courses", handle: (r) => this.listCourseRecords(r) },
      { method: "POST", pattern: "/courses", handle: (r) => this.addCourseRecord(r) },
      { method: "GET", pattern: "/courses/{id}", handle: (r) => this.showCourseRecord(r) },
    ];
  }

  // mintAccessToken Hands out a Token, and Guards nothing else.
  mintAccessToken(_: Request): Response {
    return attempt(() => ({
      status: 201,
      body: { token: this.tokens.issueAccessToken("student-registry") },
    }));
  }

  // listStudentRecords Tells one Page of the Roll.
  listStudentRecords(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: this.students.selectStudentPage(readPageRequest(request)).map(renderStudentView),
    }));
  }

  // showStudentRecord Tells the Story of one Student.
  showStudentRecord(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: renderStudentView(this.students.selectStudentRow(readPathNumber(request) as StudentID)),
    }));
  }

  // addStudentRecord Enrols someone new.
  addStudentRecord(request: Request): Response {
    return attempt(() => ({
      status: 201,
      body: renderStudentView(this.students.insertStudentRow(this.readStudentBody(request, 0))),
    }));
  }

  // saveStudentRecord Rewrites an Enrolment that already Exists.
  saveStudentRecord(request: Request): Response {
    return attempt(() => {
      const student = this.readStudentBody(request, readPathNumber(request));

      return { status: 200, body: renderStudentView(this.students.updateStudentRow(student)) };
    });
  }

  // dropStudentRecord Ends an Enrolment.
  dropStudentRecord(request: Request): Response {
    return attempt(() => {
      this.students.deleteStudentRow(readPathNumber(request) as StudentID);

      return { status: 204, body: null };
    });
  }

  // listCourseRecords Tells one Page of the Catalogue.
  listCourseRecords(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: this.courses.selectCoursePage(readPageRequest(request)).map(renderCourseView),
    }));
  }

  // showCourseRecord Tells the Story of one Course.
  showCourseRecord(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: renderCourseView(this.courses.selectCourseRow(readPathNumber(request) as CourseID)),
    }));
  }

  // addCourseRecord Opens a new Course.
  addCourseRecord(request: Request): Response {
    return attempt(() => ({
      status: 201,
      body: renderCourseView(this.courses.insertCourseRow(buildCourseRecord(readJsonBody(request), 0))),
    }));
  }

  // readStudentBody Decodes the Wire, Validates it, then Checks the Course.
  private readStudentBody(request: Request, id: number): Student {
    const student = buildStudentRecord(readJsonBody(request), id);

    const broken = checkStudentRecord(student);
    if (broken) throw broken;

    this.courses.selectCourseRow(student.course);

    return student;
  }
}

// attempt is the one Place that Turns a thrown Failure into a Reply.
// The Mapping Lives here once, never inside a Handler.
function attempt(story: () => Response): Response {
  try {
    return story();
  } catch (failure) {
    return buildFailureReply(failure);
  }
}

// buildFailureReply Maps a Business Failure onto an HTTP Code.
export function buildFailureReply(failure: unknown): Response {
  if (failure instanceof BusinessError) {
    const status = { invalid: 400, absent: 404, taken: 409 }[failure.kind];

    return { status, body: { error: failure.message } };
  }

  return { status: 500, body: { error: "unexpected Failure" } };
}

// readJsonBody Refuses anything the Wire Cannot Parse.
function readJsonBody<T>(request: Request): T {
  try {
    return JSON.parse(request.body || "{}") as T;
  } catch {
    throw ErrBodyIsBroken;
  }
}

// readPathNumber Pulls an Identity out of the Path.
function readPathNumber(request: Request): number {
  const raw = Number(request.path["id"]);

  if (!Number.isInteger(raw) || raw < 0) throw ErrPathIsBroken;

  return raw;
}

// readPageRequest Reads Pagination, Defaulting to the whole Set.
function readPageRequest(request: Request): Page {
  const page = { number: Number(request.query["page"] ?? 0), size: Number(request.query["size"] ?? 0) };

  if (page.number < 0 || page.size < 0) throw ErrPageIsInvalid;

  return page;
}
