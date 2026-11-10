# Code Generation Needs a Size and Correctness Trial (MOVING)

The summary benchmark gave a promising but narrow signal: Dove used slightly
fewer tokens than Caveman and fewer than vanilla across six three-paragraph
story summaries. Code may expose a different cost. A short patch that skips an
acceptance case is not an improvement, so size and correctness need to be read
together.

## Planned comparison

- **Conditions:** vanilla, Caveman and Dove, in fresh worktrees.
- **Model:** `gpt-6-luna`.
- **Effort:** high.
- **Runs:** three independent samples per condition and scenario.
- **Fixed inputs:** task, starter commit, tools, time limit and acceptance checks.
- **Review:** hide the condition label until correctness and scope are scored.

These are the settings for the planned code trial. The earlier story-summary
benchmark ran in a different session; its serving model and effort were not
exposed, so its results do not belong to Luna High.

## School as the scale reference

School gives three useful scales. A small case resembles one rule in `school/`,
such as RUT validation or token expiry. A medium case resembles one application
use case, such as enrollment, or one request parser. A large case is a complete
runtime generated from the contract, specification and shape: nine routes,
fourteen faults, settings, composition and tests.

The full three-language School rewrite adds language and regeneration effects.
Keep it outside a single style run; use one Go runtime as the large anchor and
compare results within each size before comparing across sizes.

## Scenario cards

Each task would start from a fresh, versioned fixture. The model could change
production code, but not the task text, acceptance checks or measurements.

### Generate

| ID | Size | Task | Acceptance |
| --- | --- | --- | --- |
| G-S | Small | Implement a pure Chilean RUT check-digit validator using `examples/specs/school/SPEC.md`. | Accept valid lower- and upper-case check digits; reject malformed input, wrong check digits and empty input; preserve leading zeroes. |
| G-M | Medium | Implement School enrollment using the ordering and failure rules in `examples/specs/school/SPEC.md`: validate form, check course, ask registry, store, then announce. | Fake ports prove call order; denial differs from registry failure; announce once after storage; a failed notice does not fail enrollment. |
| G-L | Large | Generate the Go School runtime from `examples/specs/school/openapi.yaml`, `SPEC.md` and `SHAPE.md`. | All nine routes and fourteen faults match the contract; specified boundary cases pass; the Core imports no wire, framework, store or socket package. |

### Repair

| ID | Size | Task | Acceptance |
| --- | --- | --- | --- |
| R-S | Small | Repair the token-expiry comparison described in `examples/school/BEFORE.md`. | A token is accepted before its deadline and rejected after it; check the exact boundary with a fixed clock. |
| R-M | Medium | Repair a seeded School page parser that treats malformed values as zero and loses repeated-key behavior. Use the Page rules in `examples/specs/school/SPEC.md`. | Reject `+1`, `1e1`, whitespace, empty values and hex; use the first repeated value; ignore unknown keys; preserve one-sided defaults and large-number behavior. |
| R-L | Large | Repair a seeded contract drift across route, handler, application and storage boundaries in one School runtime. | Restore parity with recorded responses for supplied cases; preserve fault reasons, status, content type and body; all acceptance checks pass. |

### Refactor

| ID | Size | Task | Acceptance |
| --- | --- | --- | --- |
| F-S | Small | Extract RUT check-digit arithmetic from a seeded School student validator into one named pure function. | Existing valid and invalid results remain identical; add no dependency or public API. |
| F-M | Medium | Separate page-query parsing from other request-form checks in School `app/validation.go`. | Existing request cases and fault answers remain identical; handlers and business rules gain no parsing logic. |
| F-L | Large | Regenerate one seeded School runtime from the Contract, Specification and Shape while preserving its API behavior. | All supplied contract cases pass; vendor boundaries remain; no unrelated public surface or dependency appears. |

## What is still missing

The repair and refactor cards need seeded starter fixtures with deliberate
defects or structure. Those fixtures must live outside the School reference so
the benchmark never damages its working implementation. Acceptance checks also
need to be packaged before running any condition.

For each run, record passing and failing checks first, then response tokens,
changed-code tokens, added and removed production lines, files, dependencies,
omitted behaviors, unnecessary public surface and completion time. Report
medians and ranges. Keep correctness separate from size: a smaller patch that
fails a check is a failure, not an efficiency win.

Distill after the fixtures and a first controlled run show whether the size
question deserves a durable rule or pattern.

## Preliminary G-L Run — 2026-11-10

One large School generation was run for each condition from independent copies
of commit `117e11abba17ae2df155008e8ae2b2218739cd91`. Each copy began with the
same contract, Specification, Shape, project Rules, Go acceptance tests and
build script; production Go files were removed. Tests and other fixture inputs
were held fixed. The model was `gpt-6-luna`, effort `high`, through Codex CLI
0.166.0. This is one sample per condition (`n=1`), not the planned three, and
only covers G-L.

| Condition | Input tokens | Cached input | Output tokens | Reasoning output | Production files | Lines | Bytes | Acceptance |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Vanilla | 5,444,400 | 5,323,392 | 23,937 | 3,208 | 15 | 1,510 | 42,052 | `./build.sh` passed |
| Caveman | 3,138,010 | 3,049,216 | 18,934 | 3,990 | 15 | 1,514 | 42,162 | `GOFLAGS=-buildvcs=false ./build.sh` passed |
| Dove | 3,216,060 | 3,114,240 | 24,430 | 6,243 | 17 | 1,534 | 43,667 | Format and vet passed; tests failed |

Token totals come from each Codex CLI turn's usage record. Input includes the
prompt and repeated context loaded during tool use; cached input is part of
input, not additional tokens. Output includes reasoning, tool calls and final
text, so it is not a code-only token count. The generated files' line and byte
counts provide a separate size measure. Since all production source began
absent, these are generated totals, not net diff lines.

Vanilla and Caveman ultimately passed format, vet, tests and binary build.
Both needed code corrections after earlier gates failed. Dove's production
binary built and format and vet passed, but the three HTTP adapter acceptance
tests failed: each returned status `201` with body `{"id":"007"}` where the
shared adapter contract expected another answer. Other package tests passed.
An escalated rerun confirmed this is a code mismatch, not the sandbox's local
socket restriction.

The official [Caveman skill](https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md)
instructs that persisted code, comments and documents remain in normal prose.
So this condition tested the actual Caveman voice during the run, but it did
not ask the model to write Caveman-style source. The generated code's small
size differences are not evidence of a code-style advantage.

Treat this as preliminary only. Complete Dove's failing adapter behavior, run
G-S and G-M, and collect two more independent samples per condition before
comparing medians or drawing a style conclusion. The fixtures and JSONL usage
records for this run are in `/tmp/onetwothree-school-bench-Oq9Lf4/` on the
benchmark machine; the result table above preserves the measurements here.
