# Canonignore

> Provisional. Written from practice, not yet Weathered.

- `.gitignore` keeps paths out of the repository. `.canonignore` marks
  versioned paths the Canon carries but does not govern.
- Entries cover raw notes, contextual humour, lookup references and exact
  tracked support files. Consult them only when that material is relevant.
- Excluded material may inform the current work; it is not an instruction and
  must not be copied into governing Canon files.
- Add an entry only when a versioned path needs this boundary. Keep untracked
  build output in `.gitignore` or the example's own ignore file.
- To audit versioned exclusions, run:

  ```sh
  git ls-files --cached --ignored --exclude-from=.canonignore
  ```

Do not enumerate every ignored file on disk during routine work.
