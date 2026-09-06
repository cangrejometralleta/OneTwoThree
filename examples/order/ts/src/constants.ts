import { checkIntegerRange, readGlobalValues } from "./settings.js";

// OrderConstants Holds shared Meaning, independent of Deployment.
export type OrderConstants = Readonly<{
  itemNameWidth: number;
  itemTotalWidth: number;
}>;

export const OrderValues = readGlobalValues<OrderConstants>(
  new URL("../../constants/order.json", import.meta.url), {
  itemNameWidth: checkIntegerRange(1, 200),
  itemTotalWidth: checkIntegerRange(1, 200),
  },
);
