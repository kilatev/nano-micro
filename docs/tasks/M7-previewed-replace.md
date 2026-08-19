# M7 — Previewed selective project replace

## Outcome

Preview project replacements, choose individual matches, and apply them safely
across open and unopened files.

## Depends on

- M6 workspace search complete.
- Session/recovery infrastructure available for transactional backups.

## Non-goals

- Unreviewed replace-all as the primary workflow.
- Applying stale matches silently.

## Acceptance

- Every selected match is revalidated before mutation.
- Open unsaved buffers are treated as authoritative editor state.
- Overlapping matches and regex capture substitutions are deterministic.
- All affected files have recoverable transactional backups.
- Partial failure is reported and cannot silently corrupt the remaining files.

## Detailed plan

Define only after M6 search behavior and recovery transactions are proven.

