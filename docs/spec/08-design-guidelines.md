# DOS design guidelines

This guide defines the design inputs for the original DOS FreightFlow Control
implementation. It preserves the supplied DOS artwork while leaving the new
interface composition, component code, copy, and visual tokens to the
implementation team.

## Approved DOS assets

The asset package supplied as `DOS-Logo_2026-07.zip` contained three PNG files.
They are preserved in `design-assets/` with descriptive names:

| File | Intended use |
| --- | --- |
| `delegateops-lockup.png` | Full DelegateOps Business Support Services lockup |
| `dos-mark.png` | Compact DOS mark for navigation or compact identity surfaces |
| `dos-workflow-mark.png` | Workflow-oriented mark for approved product contexts |

Use these files as supplied. Do not redraw, trace, recolor, crop, or combine
them without confirming brand ownership and approval. Do not treat generated
artwork as proof of trademark, copyright, or usage rights.

## Suggested enterprise design references

Use public design-system guidance as a reference for interaction quality, not as
DOS branding or a source of copied interface code.

- **IBM Carbon** is the suggested primary reference for dense operational
  tables, structured forms, status communication, and enterprise workflows.
- **Microsoft Fluent 2** is a useful secondary reference for semantic design
  tokens, adaptive layout, spacing, and component relationships.
- **WCAG 2.2** is the accessibility baseline. Target the applicable Level AA
  success criteria unless a product decision records a stricter requirement.

Reference links:

- [IBM Carbon data-table accessibility](https://carbondesignsystem.com/components/data-table/accessibility/)
- [Fluent 2 design tokens and foundations](https://fluent2.microsoft.design/get-started/design)
- [Fluent 2 layout guidance](https://fluent2.microsoft.design/layout)
- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)

Do not use IBM or Microsoft logos, names, proprietary artwork, or copied
component source in the DOS product. Build an original DOS token layer and
original component composition.

## DOS visual direction

The supplied artwork establishes a blue-to-cyan-to-violet direction on a dark
navy foundation. The implementation should derive semantic tokens from the
approved assets and confirm contrast for each foreground/background pair.

Suggested token families:

- `brand` — DOS identity and selected navigation;
- `action` — primary actions and links;
- `focus` — keyboard focus indication;
- `status` — success, warning, error, and informational states;
- `surface` — page, panel, table, and elevated layers;
- `content` — primary, secondary, muted, and inverse text;
- `border` — structural and interactive boundaries.

Do not use gradient color alone to communicate meaning. Every important state
must also have a text label, accessible name, icon, position, or other
non-color cue.

## Interaction and responsive rules

- Use sentence case and professional US English throughout the interface.
- Give each screen one visually clear primary action for its current context.
- Keep financial actions attributable, reviewable, and reversible where the
  workflow allows it.
- Make keyboard focus visible and preserve a logical heading and landmark
  structure.
- Give controls accessible names and expose validation errors next to the
  relevant field.
- Provide a usable mobile alternative for wide tables; never rely on hidden
  columns alone.
- Test at desktop `1440x900`, mobile `390x844`, and 200% zoom.
- Confirm no horizontal overflow, adequate target size, readable line length,
  and usable focus order at each responsive breakpoint.
- Pair status color with text and, when useful, an icon or shape.

## LLM implementation handoff

The implementation model must read this guide after the product, domain, and
security requirements. It must record any changed token, layout, component,
asset, or accessibility decision in its decision log and update the acceptance
matrix when behavior changes.
