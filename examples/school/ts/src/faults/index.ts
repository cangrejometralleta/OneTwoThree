export type FaultKind =
  | 'NameEmpty' | 'CodeEmpty' | 'RutInvalid' | 'AgeTooLow' | 'PageInvalid'
  | 'BodyBroken' | 'PathBroken' | 'TokenInvalid' | 'StudentUnknown'
  | 'CourseUnknown' | 'RouteUnknown' | 'RutTaken' | 'RutUnregistered'
  | 'RegistryUnavailable';

const reasons: Record<FaultKind, string> = {
  NameEmpty: 'name is Empty', CodeEmpty: 'code is Empty',
  RutInvalid: 'rut Fails its Check Digit', AgeTooLow: 'age Must be 18 or more',
  PageInvalid: 'page or size is not valid', BodyBroken: 'body is not valid JSON',
  PathBroken: 'path Holds no valid Identity', TokenInvalid: 'token is Invalid or Expired',
  StudentUnknown: 'student not Found', CourseUnknown: 'course not Found',
  RouteUnknown: 'route not Found', RutTaken: 'rut is already Registered',
  RutUnregistered: 'rut is not in the Registry', RegistryUnavailable: 'registry cannot Answer now'
};
export class Fault extends Error {
  readonly kind: FaultKind;
  constructor(kind: FaultKind, reason = reasons[kind]) { super(reason); this.name = 'Fault'; this.kind = kind; }
}
export function reasonOf(kind: FaultKind): string { return reasons[kind]; }
