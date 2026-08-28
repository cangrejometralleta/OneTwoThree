#!/usr/bin/env python3
"""Color shaped Output by Co-occurrence. The Rules live in color.md.

A Token that Repeats is a Group of its Occurrences. The first three
such Groups, by Order of Appearance, take the three Hues. The fourth
Stays bare, because the Ceiling makes Color a Limiter.

Reads stdin, Writes stdout. Passes Through untouched when stdout is
not a Terminal or NO_COLOR is Set, so a Pipe never Eats an Escape.
--force Paints into a Pipe; --legend Names the Key on stderr. Bare
Arguments Choose the Tokens by Hand, when the Human Sees a Group the
Count missed. NO_COLOR Beats --force: the Reader's standing Preference
Outranks the Writer's Flag.
"""

import os
import re
import sys

# Blue, Magenta, Cyan. Green, Red and Yellow stay Reserved
# for the Emoji Trio — see color.md, Rule five.
HUES = ("\033[38;5;39m", "\033[38;5;170m", "\033[38;5;80m")
RESET = "\033[0m"

HUE_CEILING = 3
GROUP_FLOOR = 2
# A Slash Separates Tokens; a Dot, Dash, Colon or Underscore Binds them,
# so shape_test.go:41 and v0.1.0-rc.3 Survive whole while a Path Splits.
TOKEN_PATTERN = re.compile(r"[A-Za-z_][A-Za-z0-9_.:-]*[A-Za-z0-9_]|[A-Za-z_]")
PATH_MARKS = "/."


def name_referent(token, is_path_bound):
    """A Referent Carries a Mark, Sits in a Path, or Runs long and lowercase.

    Connectors never Qualify, so the Color Lands on what the Reader
    would Look up, never on the Grammar between.
    """
    if any(mark in token for mark in "._:-") or any(c.isdigit() for c in token):
        return True
    if is_path_bound:
        return True
    return len(token) >= 8 and token.islower()


def bind_to_path(text, match):
    """A Token Sits in a Path when a Slash Leads it, or a Dot Opens it."""
    before = text[match.start() - 1] if match.start() else ""
    after = text[match.end()] if match.end() < len(text) else ""
    return before in PATH_MARKS or after == "/"


def find_repeated_tokens(text):
    """Collect the Referents that Appear twice or more, in first-seen Order."""
    counts = {}
    order = []
    for match in TOKEN_PATTERN.finditer(text):
        token = match.group(0)
        if not name_referent(token, bind_to_path(text, match)):
            continue
        if token not in counts:
            order.append(token)
        counts[token] = counts.get(token, 0) + 1

    return [token for token in order if counts[token] >= GROUP_FLOOR]


def assign_hues(tokens):
    """Hand the first three Tokens a Hue each. The rest Stay bare."""
    return dict(zip(tokens[:HUE_CEILING], HUES))


def paint_text(text, hues):
    """Wrap every whole-Token Occurrence in its Hue."""
    if not hues:
        return text

    pattern = re.compile(
        r"(?<![A-Za-z0-9_])(" + "|".join(re.escape(t) for t in hues) + r")(?![A-Za-z0-9_])"
    )
    return pattern.sub(lambda m: hues[m.group(1)] + m.group(1) + RESET, text)


def write_legend(hues, stream):
    """Name which Token took which Hue, so the Key is never Guessed."""
    for token, hue in hues.items():
        stream.write(f"{hue}{token}{RESET}\n")


def main():
    wants_legend = "--legend" in sys.argv[1:]
    forces_color = "--force" in sys.argv[1:]
    chosen = [a for a in sys.argv[1:] if not a.startswith("-")]
    text = sys.stdin.read()

    if os.environ.get("NO_COLOR") or not (forces_color or sys.stdout.isatty()):
        sys.stdout.write(text)
        return 0

    tokens = chosen or find_repeated_tokens(text)
    hues = assign_hues(tokens)
    sys.stdout.write(paint_text(text, hues))

    if wants_legend:
        write_legend(hues, sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
