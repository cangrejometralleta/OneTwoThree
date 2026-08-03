// Brand Gives TypeScript what Go Gets from a named Type:
// a String that Refuses to Stand in for another String.
declare const brand: unique symbol;
type Brand<T, Name> = T & { readonly [brand]: Name };

// OrderID is a String so an unassigned Order Compares Empty.
// checkOrderRecord Already Assumes this Shape.
export type OrderID = Brand<string, "OrderID">;

// Item is one Line of an Order.
export type Item = { name: string; price: number; qty: number };

// Order is the Business Truth about one Purchase.
export type Order = {
  id: OrderID;
  memberId: string;
  percent: number;
  items: Item[];
};

// BusinessError Carries a Name the Transport can Map.
export class BusinessError extends Error {
  constructor(readonly kind: "invalid" | "absent", message: string) {
    super(message);
  }
}

// The Business Fails in named Ways, never in Numbers.
export const ErrOrderIsInvalid = new BusinessError("invalid", "order Misses an Identity, a Member or its Items");
export const ErrPercentIsInvalid = new BusinessError("invalid", "percent Must sit between zero and a Hundred");
export const ErrOrderUnknown = new BusinessError("absent", "order not Found");
export const ErrPageIsInvalid = new BusinessError("invalid", "page Numbers Must not be negative");
export const ErrBodyIsBroken = new BusinessError("invalid", "body is not valid JSON");

// generateOrderId Mints a short Reference before the Row exists.
// A Purchase Needs its own Name the Moment it is Born.
export function generateOrderId(): OrderID {
  return Math.random().toString(16).slice(2, 10) as OrderID;
}

// checkOrderRecord Names each Condition, then Reads them together.
// The Break is not in the Expression, it is in the Vocabulary.
export function checkOrderRecord(order: Order): boolean {
  const identified = order.id !== "";
  const assigned = order.memberId !== "";
  const filled = order.items.length > 0;

  return identified && assigned && filled;
}

// checkPercentRange Refuses a Discount the Business cannot Grant.
export function checkPercentRange(percent: number): BusinessError | null {
  if (percent < 0 || percent > 100) return ErrPercentIsInvalid;

  return null;
}

// sumItemPrices Adds every Line into a Total.
export function sumItemPrices(items: Item[]): number {
  return items.reduce((total, it) => total + it.price * it.qty, 0);
}

// formatItemLine Renders one Item for the Receipt.
export function formatItemLine(it: Item): string {
  return `${it.name.padEnd(12)} x${it.qty} ${String(it.price * it.qty).padStart(6)}`;
}

// applyMemberRate Lowers a Total by a Percentage.
export function applyMemberRate(total: number, percent: number): number {
  return total - Math.floor((total * percent) / 100);
}

// reportOrderState Says how it Went, at a Glance.
export function reportOrderState(id: string, total: number, failure: unknown): string {
  if (failure) return `❌ Order ${id} Failed: ${failure}`;

  return `✅ Order ${id} Closed at ${total}`;
}

// buildOrderReceipt Reads as three Sections: Total, Lines, Result.
export function buildOrderReceipt(id: string, items: Item[], percent: number): string {
  const total = applyMemberRate(sumItemPrices(items), percent);

  const lines = items.map(formatItemLine);
  lines.push(reportOrderState(id, total, null));

  return lines.join("\n");
}
