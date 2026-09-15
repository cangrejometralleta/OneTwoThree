# Scripts

> Provisional. Written from practice, not yet Weathered.

- Every program Answers the same two scripts.  
  `build.sh` Makes the artefact. `run.sh` starts the service.
- Learn them Once and every repository opens the same way.  
  A reader who knows one project can Start the next.
- The script Enters its own directory first.  
  It runs from Anywhere, or it runs from one place only.

```sh
cd -- "$(dirname -- "$0")"
```

- One function per step, named by what it Checks.  
  `verify_secret`, `install_packages`, `build_binary`.
- The calls Live at the bottom, one per line.  
  Read the last four Lines and you know the script.
- Build Refuses to build what does not pass.  
  Format, then types, then tests, then the Artefact.
- Run Refuses to start what will fail at startup.  
  A missing secret Costs one line here and a stack trace there.
- Run Reads `.env` when present and says nothing when absent.  
  An exported variable Wins; see [Constants](constants.md).
- A refusal Names what it wanted and shows the line that works.  
  ❌ says what broke; the next Line says what to type.
- The gate Belongs to the script, never to the reader's memory.

The check Costs one command.  
`ls build.sh run.sh` in every program Directory.  
A missing one is a Step living in someone's head.
