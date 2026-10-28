# UI rules

Status: target for the initial release; implementation is in progress.

## Boundary, consumers, and ownership

Applies to every application-owned screen on desktop and Android. Website content remains owned by the site.

## Observable rules

- Use a modern, restrained visual language: one primary accent color, system-friendly typography, consistent spacing, and little decoration.
- Each screen has one visually clear primary action. Secondary actions remain easy to find without competing with it.
- Layouts adapt to narrow phones and desktop windows. Long German and Kazakh labels wrap without clipping controls or actions.
- Interactive elements have visible keyboard focus, accessible names, adequate contrast, and touch sizes appropriate to the platform.
- Honor reduced-motion preferences. Motion may clarify state changes but is never required to understand an error or complete a task.
- The app-owned UI is localized in English, German, Kazakh, and Italian. English is the default and fallback. Translation-key completeness is checked automatically.

## Failure behavior

- Error presentation follows the same typography and spacing as normal states; it is visually distinct without relying on color alone.
- Missing translation keys fall back to English and are reported by automated checks.

## Non-goals

- Restyling or translating the loaded website.
- A decorative dashboard or multi-step visual tour.

## Related contracts

- [UX rules](ux.md)
- [Settings screen](settings-screen.md)
- [Initial release](../features/initial-release.md)
