# What the Muchi refactor left open (UNDISTILLED)

One API was refactored into packages by responsibility.
The Canon took three Rules from it: Layers, Providers, Translations.
This is what did not Fit yet.

- A lease contract suite exists, with one adapter behind it.  
  A contract with one fulfiller Describes behaviour; it proves no swap.  
  Wait for the second adapter before calling it a Pattern.

The catalog's vendor types and `constants.go` were muchi-api's own work,  
already Answered by Layers, Providers and Constants. They went back there.

Distill when a second adapter fulfils the same contract.
