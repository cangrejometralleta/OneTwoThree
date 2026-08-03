import type { Order, OrderID } from "./domain.js";

// A Provider is an Interface the Core Declares
// and something outside Fulfils. See RULES.md, Providers.

// Page Carries Pagination without Naming a Database.
export type Page = { number: number; size: number };

// OrderStore Keeps Orders wherever Orders Live.
// Fulfilled by SqliteOrders.
// Reference: https://nodejs.org/api/sqlite.html
export interface OrderStore {
  insertOrderRow(order: Order): Order;
  selectOrderRow(id: OrderID): Order;
  selectOrderPage(page: Page): Order[];
}
