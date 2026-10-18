# The Second Fulfiller Tests the Swap

A Contract with one Fulfiller Describes one working implementation.
It says nothing about another taking its Place.

Run the same Contract Suite against each Adapter.
When two independent Adapters pass it, the tested behavior travels across
the Boundary.

A Contract cannot prove what it does not Observe.
A passing Pair supports a Swap only for the cases the Suite Covers.

- Keep the Suite independent of each Adapter.
- Run the same Suite against every Fulfiller.
- Treat failures and behavior outside the Contract separately.

In this repository, the shared
[adapter suite](../examples/school/go/serving/adaptertest/adaptertest.go)
runs against [Stdlib](../examples/school/go/serving/stdlib_test.go),
[Chi](../examples/school/go/serving/chiserver/chi_test.go) and
[Gin](../examples/school/go/serving/ginserver/gin_test.go).

The contract belongs to the Core; each Adapter fulfils it.
See [Providers](../rules/providers.md) and [Layers](../rules/layers.md).
