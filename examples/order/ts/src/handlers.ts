import {
  BusinessError, ErrBodyIsBroken, ErrOrderIsInvalid, ErrPageIsInvalid,
  checkOrderRecord, checkPercentRange, generateOrderId,
  type Order, type OrderID,
} from "./domain.js";
import { buildOrderRecord, renderOrderView, renderOrderViews, type OrderBody } from "./dto.js";
import type { OrderStore, Page } from "./providers.js";
import type { Request, Response, Route } from "./transport.js";

// OrderAPI Holds the Cast the Story Needs.
// One Provider, no Library: a Test can Hand it one Fake.
export class OrderAPI {
  constructor(private readonly orders: OrderStore) {}

  // declareOrderRoutes is the Libretto.
  declareOrderRoutes(): Route[] {
    return [
      { method: "POST", pattern: "/orders", handle: (r) => this.addOrderRecord(r) },
      { method: "GET", pattern: "/orders", handle: (r) => this.listOrderRecords(r) },
      { method: "GET", pattern: "/orders/{id}", handle: (r) => this.showOrderRecord(r) },
    ];
  }

  // addOrderRecord Opens a new Order and Prices it on the Spot.
  addOrderRecord(request: Request): Response {
    return attempt(() => ({
      status: 201,
      body: renderOrderView(this.orders.insertOrderRow(this.readOrderBody(request))),
    }));
  }

  // showOrderRecord Tells the Story of one Purchase.
  showOrderRecord(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: renderOrderView(this.orders.selectOrderRow(request.path["id"] as OrderID)),
    }));
  }

  // listOrderRecords Tells one Page of the Ledger.
  listOrderRecords(request: Request): Response {
    return attempt(() => ({
      status: 200,
      body: renderOrderViews(this.orders.selectOrderPage(readPageRequest(request))),
    }));
  }

  // readOrderBody Decodes the Wire, Mints an Identity, then Checks the Shape.
  private readOrderBody(request: Request): Order {
    const body = readJsonBody<OrderBody>(request);
    const order = buildOrderRecord(body, generateOrderId());

    if (!checkOrderRecord(order)) throw ErrOrderIsInvalid;

    const broken = checkPercentRange(order.percent);
    if (broken) throw broken;

    return order;
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
    return {
      status: failure.kind === "absent" ? 404 : 400,
      body: { error: failure.message },
    };
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

// readPageRequest Reads Pagination, Defaulting to the whole Set.
function readPageRequest(request: Request): Page {
  const page = { number: Number(request.query["page"] ?? 0), size: Number(request.query["size"] ?? 0) };

  if (page.number < 0 || page.size < 0) throw ErrPageIsInvalid;

  return page;
}
