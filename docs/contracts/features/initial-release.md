# Initial release

Status: approved for implementation (2026-10-28).

## Goal and user-visible behavior

Deliver a small Firefox-based site wrapper for Linux, Windows, and Android. It opens a configured HTTP(S) site in a clean application window or a full-screen kiosk. The same settings screen provides URL, page zoom, display mode, language, password change, and a read-only application version.

## Scenarios and edge cases

- First launch requires a valid URL and a new settings password before opening the site.
- Later settings access requires that password. A wrong password leaves settings unchanged.
- Desktop accepts one URL argument for a single launch. It takes precedence over the saved URL and is not stored.
- Android accepts an external URL through an Intent for a single launch; it is not stored.
- In window mode, the site appears without browser tabs or address bar. In kiosk mode, it fills the screen.
- Invalid URLs, invalid zoom values, absent Firefox, loading failures, and unavailable settings storage produce understandable errors.
- English is selected by default. English, German, Kazakh, and Italian are the only language choices. The choice localizes the settings UI and sets the preferred language requested from the site.
- Navigation initiated by the site stays in the wrapper where feasible; navigation to unsupported schemes is delegated to the OS or refused with an error.

## Non-goals

- Building Firefox or Gecko from source.
- Device or OS lockdown. Kiosk mode hides browser chrome but does not stop a user with OS access.
- Protecting settings from a user who can edit application files directly.
- Automatic publishing or code signing without release credentials.

## Affected surfaces

No durable surface contracts exist yet. At completion, record observed behavior under `docs/contracts/surfaces/` and remove this temporary feature contract.

## Interfaces and constraints

- Saved settings include URL, zoom, mode, language, and a salted password verifier; never store plaintext passwords.
- Zoom is an integer percentage from 30 to 500, inclusive. Modes are `window` and `kiosk`.
- Language codes are `en`, `de`, `kk`, and `it`. The app version comes from build metadata and is not editable.
- Accept only absolute `http://` and `https://` URLs without embedded credentials.
- Desktop uses ready-made Firefox in a dedicated profile. Android uses GeckoView.
- Desktop settings open through a keyboard shortcut; Android kiosk settings open through a deliberate press gesture. Initial setup opens automatically. The settings UI shows URL, zoom control with numeric value, mode selection, language selection, password change, read-only version, and save/cancel actions.
- Build and test locally inside Docker where practical. CI produces Linux, Windows, and Android artifacts; platform-specific behavior still requires platform validation.

## Acceptance criteria

- On each platform, first setup, settings unlock, save, restart, window mode, kiosk mode, and zoom work.
- Desktop URL argument and Android Intent URL override the saved URL only for that launch.
- Selected language is reflected in the settings UI and in the language preference visible to websites, subject to sites' own language selection rules.
- Invalid input is rejected before modifying persisted settings.
- Automated tests cover configuration validation, password verification, URL priority, native messaging or equivalent protocol, and settings UI behavior where practical.
- CI runs formatting, static checks, tests, and artifact builds. Release instructions list remaining signing and device checks.

## Validation

- Run all project unit and integration tests, formatters, and static checks in Docker or CI.
- Smoke-test Linux visually; test Windows and Android artifacts on their actual platforms before declaring a release ready.
- Review the implementation against this contract and update durable surface contracts before deleting this file.
