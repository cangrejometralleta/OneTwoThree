#!/bin/sh
# Expose this repository's skills and agents to Claude Code for every project.
# Links only. Never copies, never overwrites, never writes through a link.
# Safe to run again. Target: $CLAUDE_HOME, or ~/.claude.
set -eu

ROOT=$(cd -- "$(dirname -- "$0")/../../../../.." && pwd -P)
TARGET=${CLAUDE_HOME:-$HOME/.claude}
MARK='## Dove Canon skills'
LEGACY_MARK='## OneTwoThree skills'

link_entry() {
    src=$1 dest=$2
    if [ -L "$dest" ] && [ "$(readlink "$dest")" = "$src" ]; then echo "= $dest"; return; fi
    if [ -e "$dest" ] || [ -L "$dest" ]; then echo "! $dest exists, skipped"; return; fi
    ln -s "$src" "$dest" && echo "+ $dest"
}

link_all() {
    from=$1 to=$2 pattern=$3
    mkdir -p "$to"
    for src in "$from"/$pattern; do
        [ -e "$src" ] || continue
        link_entry "$src" "$to/$(basename "$src")"
    done
}

print_instruction() {
    echo "$MARK"
    echo 'Skills and agents under ~/.claude are symbolic links into the Dove Canon.'
    echo 'A relative link inside a SKILL.md reads correctly: `..` resolves on the real directory.'
    echo 'Lexical navigation does not follow it; use `cd -P` or `readlink -f` to walk up out of a linked skill.'
    echo 'Commands and edits still target the current project, never the Canon repository.'
}

# The instruction belongs to this machine, never to a tracked file.
# A link here may point into the Canon itself; appending would rewrite it.
write_instruction() {
    file=$TARGET/CLAUDE.md
    home=$(cd -- "$TARGET" && pwd -P)
    if [ -L "$file" ] || { [ -e "$file" ] && [ "$(readlink -f "$file")" != "$home/CLAUDE.md" ]; }; then
        echo "! $file resolves to $(readlink -f "$file"), not written"
        echo "  Make it a real file that imports that path with @, then add:"
        print_instruction | sed 's/^/    /'
        return
    fi
    if grep -qsF "$MARK" "$file"; then echo "= $file"; return; fi
    if grep -qsF "$LEGACY_MARK" "$file"; then
        echo "! $file has an old OneTwoThree instruction block; left unchanged"
        echo "  Review it manually before running this installer again."
        return
    fi
    { echo ''; print_instruction; } >>"$file"
    echo "+ $file"
}

link_all "$ROOT/.agents/skills" "$TARGET/skills" '*'
link_all "$ROOT/.agents/agents" "$TARGET/agents" '*.md'
write_instruction
