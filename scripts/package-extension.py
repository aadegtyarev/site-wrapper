from pathlib import Path
from sys import argv
from zipfile import ZIP_DEFLATED, ZipFile


FILES = (
    "manifest.json",
    "background.js",
    "content.js",
    "error.html",
    "error.js",
    "settings.html",
    "settings.css",
    "settings.js",
    "i18n.mjs",
)


def main() -> None:
    if len(argv) != 2:
        raise SystemExit("usage: package-extension.py OUTPUT.xpi")
    root = Path(__file__).resolve().parents[1]
    output = Path(argv[1]).resolve()
    output.parent.mkdir(parents=True, exist_ok=True)
    with ZipFile(output, "w", ZIP_DEFLATED) as archive:
        for name in FILES:
            archive.write(root / "extension" / name, name)


if __name__ == "__main__":
    main()
