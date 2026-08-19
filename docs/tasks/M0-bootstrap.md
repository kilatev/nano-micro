# M0 — Bootstrap

## Goal

Establish a clean public Micro fork that builds and passes its existing tests
without functional changes.

## Work

1. Inspect upstream build and test documentation.
2. Build the unmodified editor.
3. Run the applicable upstream tests.
4. Record verified commands in `AGENTS.md`.
5. Add the agreed project documents.
6. Confirm the resulting binary can open, edit, save, and reopen a file.

## Non-goals

- No feature implementation.
- No dependency upgrades.
- No refactoring.
- No formatting of unrelated files.
- No M1 sidebar work.

## Done when

- The repository builds from a clean checkout.
- Existing tests pass, or pre-existing failures are documented.
- Manual edit/save/reopen verification succeeds.
- Git diff contains only project-context files.

## Initial Codex CLI prompt

```text
Read AGENTS.md and every document it requires.

We are working only on docs/tasks/M0-bootstrap.md.

Inspect the repository and determine the exact upstream build and test
commands. Do not modify files yet. Produce a small execution plan, identify any
uncertainty, and stop for my review.
```

## Implementation prompt after plan approval

```text
Implement the approved M0 plan.

Keep changes limited to M0. Run every relevant verification and stop before
M1. Do not commit. Return the diff summary, commands run, manual verification
checklist, and remaining risks.
```

