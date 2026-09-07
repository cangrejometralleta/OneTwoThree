# Scripts

> Provisional. Written from Practice, not yet Weathered.

- Every Program Answers the same two Scripts.  
  `build.sh` Makes the Artefact. `run.sh` Starts the Service.
- Learn them once and every Repository Opens the same Way.  
  A Reader who Knows one Project can Start the next.
- The Script Enters its own Directory first.  
  It Runs from anywhere, or it Runs from one Place only.

```sh
cd -- "$(dirname -- "$0")"
```

- One Function per Step, Named by what it Checks.  
  `verify_secret`, `install_packages`, `build_binary`.
- The Calls Live at the Bottom, one per Line.  
  Read the last four Lines and you Know the Script.
- Build Refuses to Build what does not Pass.  
  Format, then Types, then Tests, then the Artefact.
- Run Refuses to Start what will Fail at Startup.  
  A missing Secret Costs one Line here and a Stack Trace there.
- Run Reads `.env` when Present and Says nothing when Absent.  
  An exported Variable Wins; see [Constants](constants.md).
- A Refusal Names what it Wanted and Shows the Line that Works.  
  ❌ Says what Broke; the next Line Says what to Type.
- The Gate Belongs to the Script, never to the Reader's Memory.

The Check Costs one Command.
`ls build.sh run.sh` in every Program Directory.
A Missing one is a Step Living in someone's Head.
