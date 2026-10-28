# CI test packages

Status: active.

## Boundary and consumers

The `Checks` workflow produces downloadable packages for people testing Site Wrapper on Linux, Windows, and Android. These packages are test builds, not releases.

## Outputs and guarantees

- A successful run exposes `site-wrapper-linux-test` and `site-wrapper-windows-test`. Each contains the platform executable, the same signed desktop XPI beside it, and `TESTING.md`.
- Android has separate `site-wrapper-android-<ABI>-test` artifacts for `arm64-v8a`, `armeabi-v7a`, and `x86_64`. Each contains one installable debug APK and `TESTING.md`.
- Artifacts are retained for 30 days. The workflow fails if a signed XPI no longer matches extension sources or if an expected APK is missing.
- `TESTING.md` identifies setup steps, the checks to run, and what to include in a failure report.

## Limits

- CI cross-builds do not prove that window mode, kiosk mode, navigation, or zoom work on a physical device.
- Windows executables and Android debug APKs have no production signing identity. Android debug builds from different CI runs may require uninstalling the previous build.
- Packages contain no saved site URL, settings password, or release credentials.

## Related contracts

- [Initial release](../features/initial-release.md)
- [Code quality](code-quality.md)
