import { loadOrderConfig } from "./config.js";
import { serveOrderRoutes } from "./adapters.js";
import { OrderAPI } from "./handlers.js";
import { SqliteOrders } from "./store-sqlite.js";

// main Casts the Players, then Steps off the Stage.
async function main(): Promise<void> {
  const config = loadOrderConfig({ ...process.env });

  const store = SqliteOrders.shapeOrderTable(config.databasePath);
  const api = new OrderAPI(store);

  await serveOrderRoutes(api.declareOrderRoutes(), config.port);
  console.log(`✅ Orders Listening on :${config.port}`);
}

main().catch((failure) => {
  console.error("❌ Orders Refused to Start:", failure);
  process.exit(1);
});
