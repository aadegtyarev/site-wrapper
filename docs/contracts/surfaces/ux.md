# UX rules

Status: target for the initial release; implementation is in progress.

## Boundary, consumers, and ownership

Governs the user's path through setup, launch, settings, recovery, and errors across all three platforms.

## Observable flows and guarantees

- First launch presents one setup screen. A single save action validates the inputs and opens the site.
- Later launches go directly to the saved site. Settings require one deliberate entry action, password unlock, then one screen for all controls.
- The current task remains visible after a recoverable error. Valid non-secret input is preserved, and retry does not require repeating unrelated steps.
- Every recoverable error is shown in the selected app language with a concrete next step. No caught error is silently discarded.
- Missing Firefox offers a link to Mozilla's official download page. Load failures offer retry and access to settings. Failed persistence states that changes were not saved.
- Unexpected failures show an error ID and write local diagnostics that omit passwords and URL credentials. Raw stack traces stay out of user-facing messages.
- The preferred language is a request to the site, not a promise that the site will switch languages.

## Interfaces and compatibility

- Window mode has a visible settings entry. Kiosk mode has a documented deliberate gesture or keyboard shortcut that does not obscure normal website controls.
- Keyboard, touch, and screen-reader flows must complete the same tasks.

## Non-goals

- Locking the user out of the operating system.
- Hiding a failure to preserve an uncluttered screen.

## Related contracts

- [UI rules](ui.md)
- [Settings screen](settings-screen.md)
- [Initial release](../features/initial-release.md)
