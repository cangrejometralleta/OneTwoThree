import assert from "node:assert/strict";
import test from "node:test";

import {
  ErrOrderUnknown,
  applyMemberRate, buildOrderReceipt, checkOrderRecord, formatItemLine, reportOrderState, sumItemPrices,
  type Order, type OrderID,
} from "./domain.js";
import { OrderAPI } from "./handlers.js";
import type { OrderStore, Page } from "./providers.js";
import type { Request } from "./transport.js";

// FakeOrders Stands in for SQLite.
// No Database, no Framework, no Port: the Providers Allow it.
class FakeOrders implements OrderStore {
  private orders = new Map<OrderID, Order>();

  insertOrderRow(order: Order): Order {
    this.orders.set(order.id, order);
    return order;
  }

  selectOrderRow(id: OrderID): Order {
    const order = this.orders.get(id);
    if (!order) throw ErrOrderUnknown;
    return order;
  }

  selectOrderPage(_: Page): Order[] {
    return [...this.orders.values()];
  }
}

// buildTestingOrders Hands the API one Fake and nothing else.
function buildTestingOrders(): OrderAPI {
  return new OrderAPI(new FakeOrders());
}

function askOrders(body: unknown, path: Record<string, string> = {}): Request {
  return { path, query: {}, body: JSON.stringify(body) };
}

test("sumItemPrices adds every Line", () => {
  const items = [{ name: "Widget", price: 100, qty: 2 }, { name: "Gadget", price: 50, qty: 1 }];
  assert.equal(sumItemPrices(items), 250);
});

test("applyMemberRate lowers the Total", () => {
  assert.equal(applyMemberRate(200, 10), 180);
});

test("checkOrderRecord guards each Condition", () => {
  const items = [{ name: "Widget", price: 1, qty: 1 }];
  const good = { id: "ab12", memberId: "m1", percent: 0, items } as unknown as Order;

  assert.equal(checkOrderRecord(good), true);
  assert.equal(checkOrderRecord({ ...good, id: "" as OrderID }), false);
  assert.equal(checkOrderRecord({ ...good, memberId: "" }), false);
  assert.equal(checkOrderRecord({ ...good, items: [] }), false);
});

// The Text that Ties this Domain to the Snippet RULES.md already Cites.
test("buildOrderReceipt renders three Sections", () => {
  const items = [{ name: "Widget", price: 100, qty: 2 }];
  const receipt = buildOrderReceipt("ab12", items, 10);
  const lines = receipt.split("\n");

  assert.equal(lines.length, 2);
  assert.equal(lines[0], formatItemLine(items[0]!));
  assert.equal(lines[1], reportOrderState("ab12", 180, null));
});

test("addOrderRecord opens a Purchase", () => {
  const reply = buildTestingOrders().addOrderRecord(
    askOrders({ memberId: "m1", percent: 10, items: [{ name: "Widget", price: 100, qty: 2 }] }),
  );
  assert.equal(reply.status, 201);
});

test("addOrderRecord refuses an empty Basket", () => {
  const reply = buildTestingOrders().addOrderRecord(askOrders({ memberId: "m1", percent: 10, items: [] }));
  assert.equal(reply.status, 400);
});

test("addOrderRecord refuses a bad Percent", () => {
  const reply = buildTestingOrders().addOrderRecord(
    askOrders({ memberId: "m1", percent: 150, items: [{ name: "Widget", price: 100, qty: 1 }] }),
  );
  assert.equal(reply.status, 400);
});

test("showOrderRecord reports an Absence", () => {
  const reply = buildTestingOrders().showOrderRecord(askOrders(null, { id: "unknown" }));
  assert.equal(reply.status, 404);
});

test("listOrderRecords refuses a negative Page", () => {
  const reply = buildTestingOrders().listOrderRecords({ path: {}, query: { page: "-1" }, body: "" });
  assert.equal(reply.status, 400);
});
