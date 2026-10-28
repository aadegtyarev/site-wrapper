# Settings screen

Status: target for the initial release; implementation is in progress.

## Boundary, consumers, and ownership

The application-owned setup and settings screen. Desktop and Android use the same content and behavior with platform-appropriate controls.

## Inputs, outputs, and flows

- First setup shows URL, zoom, mode, language, and password creation on one screen. The primary action saves valid settings and opens the site.
- Later access asks for the settings password, then opens the same settings screen. Wrong passwords show an error and allow retry.
- URL accepts an absolute HTTP(S) address without embedded credentials. Zoom accepts 30–500%, inclusive, with 100% as default.
- Mode selects `window` or `kiosk`. Language selects English, German, Kazakh, or Italian; English is default. Changing language updates app-owned labels immediately.
- Version is shown read-only on the screen, sourced from build metadata. A link to the public source repository and Mozilla trademark attribution are visible on the same screen.
- Existing passwords can be changed after unlock. Password creation and change require confirmation and reject mismatches.
- Save validates all fields before writing. Cancel discards unsaved changes. A CLI or Intent URL override affects only the current launch and never rewrites the saved URL.

## Failure behavior and recovery

- Field errors appear beside the field. Failed saves preserve non-secret inputs, explain that changes were not saved, and offer retry.
- Passwords are never displayed by default or persisted in plaintext. The screen does not leak a password in logs or error messages.

## Security and privacy

- A salted verifier protects changes through the settings UI. It is not a boundary against someone who can modify application files.
- External URLs are validated before loading and do not bypass the settings password for persisted changes.

## Non-goals

- Browser history, cookies, device administration, or website account settings.

## Related contracts

- [UI rules](ui.md)
- [UX rules](ux.md)
- [Initial release](../features/initial-release.md)
