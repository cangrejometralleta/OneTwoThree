export type Request = { method: string; path: string; headers: Record<string, string | string[] | undefined>; body: Uint8Array };
export type Response = { status: number; headers: Record<string, string>; body?: unknown };
