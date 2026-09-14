# Entrypoints

> Provisional. Written from Practice, not yet Weathered.

- A Program Declares its Entry Points, or it has none a Reader can Trust.  
  What you Enter through must be Evident before you Ask.
- Evident Means Named in the Root, at the Top, in one List.  
  A Newcomer Reads that List and Runs the Program.
- An Entry Point nobody Listed is a Door someone Remembers.  
  Memory is not a Door.
- See [Scripts](scripts.md) for what the two Doors Do.  
  This Rule Says how the Doors are Declared and Multiplied.

## One Logic, one Shim per Platform

- A Shim is the thin File a Platform Knows how to Open,  
  and it does nothing but Call the Logic that Lives elsewhere.
- It Exists because the Platform Demands a `.cmd`, not a `.sh`.  
  The Demand is about the Extension, never about the Work.
- Name, Call, Exit Code. A fourth Line is the Shim starting to Think,  
  and a Shim that Thinks is a Wrapper, which is a second Program.
- The Logic Lives once, in the Language it Thinks best in.  
  Every other Platform Calls it. None of them Translates it.
- A Script Translated is a second Script, and the second one Lies.  
  Both Start Identical and Part at the first Fix only one Received.
- A Shim cannot Diverge, because a Shim Decides nothing.
- One Helper Finds the Interpreter, for every Shim.  
  Five Copies of that Search go stale on the first new Path.
- The Pattern Asks for an Interpreter already Installed over there.  
  Where you cannot Assume one, Rewrite, and Test for the Divergence.

```cmd
@echo off
rem Deploys the Project. The Logic lives in the .sh.
call "%~dp0run-bash.cmd" deploy.sh %*
exit /b %errorlevel%
```

## The Order of the Search

- What you Choose is not an Interpreter. It is a Toolbox.  
  Prefer the Interpreter that Sees the Tools the Script Names.
- First the one that Shares the System PATH.
- Then one from the PATH, Launchers Discarded.
- Then the foreign Environment, Path Translated and Said out loud.
- Then Nothing: an Address to Install from, and a non-zero Code.
- An Interpreter that Runs but cannot Find its Tools is worse than none.  
  It Fails further in.

## The Invariants

- The Shim Propagates the Exit Code.  
  Without it, a Windows CI Passes always.
- The Shim Resolves against itself with `%~dp0`, never the current Directory.
- The Line Endings Live in `.gitattributes`, not in each Machine's Editor.

  ```
  *.cmd text eol=crlf
  *.sh  text eol=lf
  ```

- The `.cmd` stays ASCII. Another Codepage Dirties the Accents.
- A Shim is Read from here and Proven only over there.  
  Until someone Runs it on the Platform, Say it is Unverified.

The Check Costs two Commands.
`rg -n 'build\.sh|run\.sh' README.md` Names the Doors in the Root.
`wc -l *.cmd` Says no Shim Passed three Lines plus `@echo off`.
Every `.sh` with a Shim Beside it Holds the Logic alone.
