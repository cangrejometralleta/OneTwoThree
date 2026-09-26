# What the Muchi refactor left open (UNDISTILLED)

One API was refactored into packages by responsibility.
The Canon took three Rules from it: Layers, Providers, Translations.
This is what did not Fit yet.

- The catalog package still Returns concrete vendor types.  
  Composition may know implementations; is that a Leak or the design?
- A lease contract suite exists, with one adapter behind it.  
  A contract with one fulfiller Describes behaviour; it proves no swap.  
  Wait for the second adapter before calling it a Pattern.
- `constants.go` survived the move with a generic name.  
  The Constants rule already asks for a descriptive one.

The Translations check found four broken links there, and nothing else.  
They are fixed. Muchi is where the rule came from, so it is still one Practice.

Distill when a second project repeats one of these.
