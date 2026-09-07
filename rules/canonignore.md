# Canonignore

> Provisional. Written from Practice, not yet Weathered.

- A Repository Carries more than it Teaches.  
  `.canonignore` Says which is which.
- `.gitignore` Says what the Repository does not Carry.  
  `.canonignore` Says what it Carries and does not Teach.
- Read every Path in it. Copy none.  
  Ignoring is not Hiding; it is Withholding Authority.

```text
.gitignore      absent from the Repository
.canonignore    present, and not to be Imitated
```

- Same Syntax, so no one Learns a second one.  
  Globs, `#` for a Comment, `!` to Take one back.
- Three Kinds Belong here, and a fourth Belongs nowhere.  
  The Material Read but never Copied, the undistilled Note,
  the generated Artefact.  
  Canon that Embarrasses you does not go here; it goes away.
- Name the Reason above the Pattern.  
  A Path with no Reason Rots into a Path no one Dares Delete.
- Take back the File that Explains the Rest.  
  `jokes/**` Ignores every Joke and Keeps the README that Frames them.
- The Entry Earns its Place by being Read wrongly once.  
  Add a Path the Day an Agent Copies from it.

```sh
git ls-files --cached --ignored --exclude-from=.canonignore
git ls-files --others --ignored --exclude-from=.canonignore
```

git Reads the File without being Told it is new.
The first Command Lists the Committed; the second, the merely Present.
A Path in neither Answer is Canon, and Governs.
