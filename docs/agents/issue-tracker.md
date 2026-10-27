# Issue tracker: Local Markdown

Issues and specs live in `.scratch/`. Keep these files eligible for
version control; the user will commit them later.

## Conventions

- One feature per directory: `.scratch/<feature-slug>/`.
- Specs: `.scratch/<feature-slug>/spec.md`.
- Tickets: `.scratch/<feature-slug>/issues/<NN>-<slug>.md`,
  numbered from `01`, one file per ticket.
- Record triage state in a `Status:` line near the top.
  Read `triage-labels.md` for the supported values.
- Append comments under `## Comments`.

When a skill says to publish, write the appropriate local file.
When fetching a ticket, read its referenced path. If a number matches
multiple features, clarify which feature is intended.

## Wayfinding

- Map: `.scratch/<effort>/map.md`, with Notes, Decisions-so-far,
  and Fog sections.
- Child tickets use the same numbered issue paths.
- Record `Type: research`, `prototype`, `grilling`, or `task`.
- Wayfinding lifecycle uses `Status: open`, `claimed`, or `resolved`.
- Record dependencies as `Blocked by: NN, NN`.
- A ticket is unblocked when all listed dependencies are resolved.
- Select the lowest-numbered open, unblocked ticket.
- Set `Status: claimed` and save before starting work.
- To resolve, append an `## Answer`, set `Status: resolved`, and add
  a summary with a ticket link to the map's Decisions-so-far section.
