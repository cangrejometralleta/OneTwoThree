# What the Muchi refactor left open (UNDISTILLED)

One API was refactored into packages by responsibility.
The Canon took three Rules from it: Layers, Providers, Translations.
This is what did not Fit yet.

- The catalog package still Returns concrete vendor types.  
  Composition may know implementations; is that a Leak or the design?
- A lease contract suite exists, with one adapter behind it.  
  A contract with one fulfiller Describes behaviour; it proves no swap.  
  Wait for the second adapter before calling it a Pattern.
- ~~`constants.go` survived the move with a generic name.~~  
  Resolved in `muchi-api` (now its own repo): the file carries a doc
  comment and one descriptive constant, `UserAgent`.

The Translations check found four broken links there, and nothing else.  
They are fixed. Muchi is where the rule came from, so it is still one Practice.

The Muchi API and front are now split into two public repos
(`cangrejometralleta/muchi-api`, `metaliaw/muchi`). Checked both:
the front is Python/Vite, with no catalog package, no lease adapter,
no Go `constants.go` — it can't repeat either open point. The split
is Muchi becoming two repos, not a second project; the two remaining
points below still wait for one.

Distill when a second project repeats one of these.
