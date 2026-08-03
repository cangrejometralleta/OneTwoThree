import { DatabaseSync } from "node:sqlite";

import { ErrOrderUnknown, type Item, type Order, type OrderID } from "./domain.js";
import type { OrderStore, Page } from "./providers.js";

// OrderRow is the Storage Shape of an Order.
// Items Live as a JSON Column: one Table, no Join,
// the Shape a Receipt Needs and nothing more.
type OrderRow = { id: string; member_id: string; percent: number; items: string };

// SqliteOrders Fulfils OrderStore with one Connection.
// It is the only Class in this Program that Knows SQLite Exists.
// Reference: https://nodejs.org/api/sqlite.html
export class SqliteOrders implements OrderStore {
  constructor(private readonly db: DatabaseSync) {}

  // shapeOrderTable Creates the Table, Documenting the Schema in place.
  static shapeOrderTable(path: string): SqliteOrders {
    const db = new DatabaseSync(path);

    db.exec(`
      CREATE TABLE IF NOT EXISTS order_row (
        id        TEXT PRIMARY KEY,
        member_id TEXT NOT NULL,
        percent   INTEGER NOT NULL,
        items     TEXT NOT NULL
      );`);

    return new SqliteOrders(db);
  }

  insertOrderRow(order: Order): Order {
    const insert = this.db.prepare(
      "INSERT INTO order_row (id, member_id, percent, items) VALUES (?, ?, ?, ?)",
    );

    insert.run(order.id, order.memberId, order.percent, JSON.stringify(order.items));

    return order;
  }

  selectOrderRow(id: OrderID): Order {
    const row = this.db.prepare("SELECT * FROM order_row WHERE id = ?").get(id) as OrderRow | undefined;

    if (!row) throw ErrOrderUnknown;

    return decodeOrderRow(row);
  }

  selectOrderPage(page: Page): Order[] {
    const rows = this.db
      .prepare("SELECT * FROM order_row ORDER BY id DESC LIMIT ? OFFSET ?")
      .all(pageLimit(page), page.number * page.size) as OrderRow[];

    return rows.map(decodeOrderRow);
  }
}

// pageLimit Turns an absent Size into the whole Set.
function pageLimit(page: Page): number {
  return page.size > 0 ? page.size : -1;
}

// decodeOrderRow Turns Storage back into Business.
function decodeOrderRow(row: OrderRow): Order {
  return {
    id: row.id as OrderID,
    memberId: row.member_id,
    percent: row.percent,
    items: JSON.parse(row.items) as Item[],
  };
}
