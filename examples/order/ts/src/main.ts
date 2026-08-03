import { serveOrderRoutes } from "./adapters.js";
import { OrderAPI } from "./handlers.js";
import { SqliteOrders } from "./store-sqlite.js";

// main Casts the Players, then Steps off the Stage.
async function main(): Promise<void> {
  const store = SqliteOrders.shapeOrderTable("order.db");
  const api = new OrderAPI(store);

  await serveOrderRoutes(api.declareOrderRoutes(), 8081);
  console.log("✅ Orders Listening on :8081");
}

main().catch((failure) => {
  console.error("❌ Orders Refused to Start:", failure);
  process.exit(1);
});
