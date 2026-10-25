#!/usr/bin/env python3
"""Run read-only acceptance cases against the local Dove agent.

The runner selects the CLI that hosts the agent: Codex or Claude Code.
Both backends see the same cases and pass through the same checks.
"""

from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import re
import subprocess
import sys
import tomllib
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[3]
RUNS = Path(__file__).with_name("runs.jsonl")
TIMEOUT = 240


@dataclass(frozen=True)
class Case:
    name: str
    scenario: str
    expected_reads: tuple[str, ...]
    option_count: int | None = None
    required: tuple[str, ...] = ()
    forbidden: tuple[str, ...] = ()


@dataclass(frozen=True)
class Outcome:
    answer: str
    trace: str
    returncode: int
    stderr: str


CASES = (
    Case(
        "three_distinct_paths",
        """Active documentation work has three real, distinct paths: update
English examples, update Spanish examples, or verify links. None has been
chosen or started. User says: ¿Qué opciones tenemos? Read Dove and
next-options, answer as Dove, and take no action.""",
        (".agents/skills/next-options/SKILL.md",),
        option_count=3,
        required=("inglés|English", "español|Spanish", "enlaces|links"),
        forbidden=("I updated", "ya actualicé", "he verificado"),
    ),
    Case(
        "two_paths_no_padding",
        """A documentation change has exactly two viable paths: update its
examples or verify its links. There is no third useful path. User says:
¿Qué opciones tenemos? Read Dove and next-options, answer as Dove, and take
no action.""",
        (".agents/skills/next-options/SKILL.md",),
        option_count=2,
        required=("ejemplos|examples", "enlaces|links"),
        forbidden=("third", "tercera opción", "I updated", "ya actualicé"),
    ),
    Case(
        "external_blocker",
        """A release is ready, but publishing requires an unavailable
maintainer credential. No independent work remains. User asks: ¿Qué
opciones tenemos para publicar? Read Dove, next-options and one-two-unblock;
answer as Dove and do not invent a menu.""",
        (
            ".agents/skills/next-options/SKILL.md",
            ".agents/skills/one-two-unblock/SKILL.md",
        ),
        option_count=0,
        required=("bloque|block", "credencial|credential", "mantenedor|maintainer"),
    ),
    Case(
        "unknown_current_position",
        """There is no active task or prior decision in the conversation,
handoff or repository context. User says: Muéstrame mis opciones. Read Dove
and next-options, answer as Dove, and do not invent context.""",
        (".agents/skills/next-options/SKILL.md",),
        option_count=0,
        required=("?",),
    ),
    Case(
        "selected_path_one_step",
        """Conversation state: Dove offered (1) inspect README links, (2) read
tests/agents/next_options/flow_cases.md and report how many cases it defines,
(3) edit the README. User selected option 2 and said: dale. Read Dove and
next-episode. Confirm the step still holds, perform only that read-only step,
report the case count, then name one following step.""",
        (
            ".agents/skills/next-episode/SKILL.md",
            "tests/agents/next_options/flow_cases.md",
        ),
        required=(
            "five cases|cinco casos",
            "inspect the README links|revisar los enlaces del README",
            "Next|Siguiente",
        ),
        forbidden=("README actualizado", "README edited", "I edited the README"),
    ),
)


def build_prompt(case: Case) -> str:
    return (
        "Act as the repository's Dove agent. Read .agents/agents/dove.md "
        "before responding and load the Skills required by its routing rules.\n\n"
        f"Acceptance case: {case.name}\n{case.scenario}\n\n"
        "Return only the user-facing Dove response. Do not describe this test."
    )


def events_of(stdout: str) -> list[dict]:
    events = []
    for line in stdout.splitlines():
        try:
            events.append(json.loads(line))
        except json.JSONDecodeError:
            continue
    return events


def run_codex(prompt: str, model: str, effort: str, vendor: str) -> Outcome:
    result = subprocess.run(
        [
            "codex",
            "exec",
            "--ephemeral",
            "--sandbox",
            "read-only",
            "--json",
            "--model",
            model,
            "--config",
            f'model_reasoning_effort="{effort}"',
            "--config",
            f'model_provider="{vendor}"',
            "-C",
            str(ROOT),
            prompt,
        ],
        cwd=ROOT,
        capture_output=True,
        text=True,
        timeout=TIMEOUT,
        check=False,
    )
    events = events_of(result.stdout)
    messages = [
        event.get("item", {}).get("text", "")
        for event in events
        if event.get("type") == "item.completed"
        and event.get("item", {}).get("type") == "agent_message"
    ]
    trace = "\n".join(
        event.get("item", {}).get("command", "")
        + "\n"
        + event.get("item", {}).get("aggregated_output", "")
        for event in events
        if event.get("type") in ("item.completed", "item.started")
        and event.get("item", {}).get("type") == "command_execution"
    )
    return Outcome(messages[-1] if messages else "", trace, result.returncode, result.stderr)


def run_claude(prompt: str, model: str, effort: str, vendor: str) -> Outcome:
    # Read-only by construction: only tools that cannot write are offered.
    result = subprocess.run(
        [
            "claude",
            "--print",
            "--output-format",
            "stream-json",
            "--verbose",
            "--no-session-persistence",
            "--tools",
            "Read,Grep,Glob",
            "--model",
            model,
            "--effort",
            effort,
            prompt,
        ],
        cwd=ROOT,
        capture_output=True,
        text=True,
        timeout=TIMEOUT,
        check=False,
    )
    events = events_of(result.stdout)
    answer = ""
    for event in events:
        if event.get("type") == "result" and isinstance(event.get("result"), str):
            answer = event["result"]
    uses = []
    for event in events:
        if event.get("type") != "assistant":
            continue
        for block in event.get("message", {}).get("content", []):
            if isinstance(block, dict) and block.get("type") == "tool_use":
                uses.append(json.dumps(block.get("input", {}), ensure_ascii=False))
    return Outcome(answer, "\n".join(uses), result.returncode, result.stderr)


RUNNERS = {"codex": run_codex, "claude": run_claude}
DEFAULT_VENDOR = {"codex": "openai", "claude": "anthropic"}
DEFAULT_MODEL = {"codex": "gpt-6-luna", "claude": "sonnet"}


def git_status() -> str:
    return subprocess.run(
        ["git", "status", "--porcelain"], cwd=ROOT, capture_output=True, text=True, check=True
    ).stdout


def run_case(
    case: Case, runner: str, model: str, effort: str, vendor: str
) -> tuple[bool, str]:
    before = git_status()
    outcome = RUNNERS[runner](build_prompt(case), model, effort, vendor)
    answer = outcome.answer
    issues = []
    if outcome.returncode != 0:
        issues.append(f"{runner} exited {outcome.returncode}: {outcome.stderr[-500:]}")
    if git_status() != before:
        issues.append("repository status changed during the read-only agent case")
    if not answer:
        issues.append("no final Dove response found in the agent output")
    for skill_path in case.expected_reads:
        if skill_path not in outcome.trace:
            issues.append(f"execution trace did not read {skill_path}")
    if case.option_count is not None:
        numbered = re.findall(r"(?m)^\s*(\d+)[.)]\s+", answer)
        if len(numbered) != case.option_count:
            issues.append(
                f"expected {case.option_count} numbered choices, found {len(numbered)}"
            )
    for fragment in case.required:
        alternatives = fragment.split("|")
        if not any(term.casefold() in answer.casefold() for term in alternatives):
            issues.append(f"response is missing required evidence: {fragment!r}")
    for fragment in case.forbidden:
        if fragment.casefold() in answer.casefold():
            issues.append(f"response contains forbidden text: {fragment!r}")
    if issues:
        issues.append(f"agent response: {answer.replace(chr(10), ' ')}")
    return not issues, "\n".join(issues) if issues else answer.replace("\n", " ")


def codex_config() -> dict:
    codex_home = Path(os.environ.get("CODEX_HOME", Path.home() / ".codex"))
    config_path = codex_home / "config.toml"
    return tomllib.loads(config_path.read_text()) if config_path.is_file() else {}


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--runner", choices=sorted(RUNNERS), default="codex")
    parser.add_argument("--model")
    parser.add_argument("--effort")
    parser.add_argument("--vendor")
    parser.add_argument("cases", nargs="*", help="optional case names to run")
    args = parser.parse_args()

    # Codex defaults come from its own configuration; Claude's are fixed here.
    config = codex_config() if args.runner == "codex" else {}
    model = args.model or config.get("model", DEFAULT_MODEL[args.runner])
    effort = args.effort or config.get("model_reasoning_effort", "high")
    vendor = args.vendor or config.get("model_provider", DEFAULT_VENDOR[args.runner])

    selected = set(args.cases)
    cases = [case for case in CASES if not selected or case.name in selected]
    unknown = selected - {case.name for case in CASES}
    if unknown:
        print(f"unknown case names: {', '.join(sorted(unknown))}", file=sys.stderr)
        return 2

    cli_version = subprocess.run(
        [args.runner, "--version"], capture_output=True, text=True, check=True
    ).stdout.strip()
    print(
        f"RUN runner={args.runner} model={model} effort={effort} "
        f"vendor={vendor} cli={cli_version}",
        flush=True,
    )

    failures = 0
    results = []
    for case in cases:
        try:
            passed, detail = run_case(case, args.runner, model, effort, vendor)
        except subprocess.TimeoutExpired:
            passed, detail = False, f"{args.runner} timed out after {TIMEOUT} seconds"
        print(f"{'PASS' if passed else 'FAIL'} {case.name}: {detail}", flush=True)
        failures += not passed
        results.append({"case": case.name, "passed": passed})
    print(f"\n{len(cases) - failures}/{len(cases)} cases passed")

    record = {
        "recorded_at": dt.datetime.now(dt.timezone.utc).isoformat(timespec="seconds"),
        "runner": args.runner,
        "vendor": vendor,
        "model": model,
        "effort": effort,
        "cli_version": cli_version,
        "passed_cases": len(cases) - failures,
        "total_cases": len(cases),
        "results": results,
    }
    with RUNS.open("a", encoding="utf-8") as stream:
        stream.write(json.dumps(record, ensure_ascii=False) + "\n")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
