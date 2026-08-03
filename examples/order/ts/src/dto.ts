import { applyMemberRate, buildOrderReceipt, sumItemPrices, type Item, type Order, type OrderID } from "./domain.js";

// OrderBody is what a Client Sends.
// Items Carry no Identity of their own, so they Cross unchanged.
export type OrderBody = { memberId: string; percent: number; items: Item[] };

// OrderView is what a Client Gets back.
// Total and Receipt are Derived, never Stored twice.
export type OrderView = {
  id: string;
  memberId: string;
  percent: number;
  items: Item[];
  total: number;
  receipt: string;
};

// buildOrderRecord Promotes untrusted Input into a Business Order.
export function buildOrderRecord(body: OrderBody, id: OrderID): Order {
  return { id, memberId: body.memberId, percent: body.percent, items: body.items };
}

// renderOrderView Demotes a Business Order back to the Wire.
export function renderOrderView(order: Order): OrderView {
  const total = applyMemberRate(sumItemPrices(order.items), order.percent);
  const receipt = buildOrderReceipt(order.id, order.items, order.percent);

  return { id: order.id, memberId: order.memberId, percent: order.percent, items: order.items, total, receipt };
}

// renderOrderViews Repeats the Move for a whole Page.
export function renderOrderViews(orders: Order[]): OrderView[] {
  return orders.map(renderOrderView);
}
