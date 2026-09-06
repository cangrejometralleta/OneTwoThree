package main

// OrderConstants Holds shared Meaning, independent of Deployment.
// Keep the Snapshot private; Readers Receive a Value Copy.
type OrderConstants struct {
	ItemNameWidth  int `json:"itemNameWidth"`
	ItemTotalWidth int `json:"itemTotalWidth"`
}

var orderConstants = readGlobalValues[OrderConstants]("../constants/order.json", map[string]valueRule{
	"itemNameWidth":  checkIntegerRange(1, 200),
	"itemTotalWidth": checkIntegerRange(1, 200),
})

// ReadOrderConstants Returns the Startup Snapshot without a Mutation Path.
func ReadOrderConstants() OrderConstants {
	return orderConstants
}
