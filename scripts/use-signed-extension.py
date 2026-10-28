from hashlib import sha256
from json import loads
from pathlib import Path
from shutil import copyfile
from sys import argv
from zipfile import ZipFile

from package_extension import FILES


ROOT = Path(__file__).resolve().parents[1]
NAME = "site-wrapper@aadegtyarev.github.io.xpi"
SIGNED = ROOT / "release" / "signed" / NAME
SHA256 = "049bef4420be3ebab11db161c8d2474b7f07fb7e462c6b77a83c14345a4b69c8"
SIGNATURE_FILES = {"META-INF/mozilla.rsa", "META-INF/mozilla.sf", "META-INF/manifest.mf"}


def verify() -> None:
    if sha256(SIGNED.read_bytes()).hexdigest() != SHA256:
        raise ValueError("signed XPI checksum differs from the reviewed file")
    with ZipFile(SIGNED) as archive:
        if archive.testzip() is not None:
            raise ValueError("signed XPI has a damaged ZIP entry")
        names = set(archive.namelist())
        if not SIGNATURE_FILES <= names:
            raise ValueError("signed XPI is missing signature files")
        if set(FILES) - names:
            raise ValueError("signed XPI is missing extension files")
        unexpected = names - set(FILES) - {name for name in names if name.startswith("META-INF/")}
        if unexpected:
            raise ValueError(f"signed XPI has unexpected content: {sorted(unexpected)}")
        for name in FILES:
            actual = archive.read(name)
            expected = (ROOT / "extension" / name).read_bytes()
            if name == "manifest.json":
                if loads(actual) != loads(expected):
                    raise ValueError("signed XPI manifest differs from source")
            elif actual != expected:
                raise ValueError(f"signed XPI {name} differs from source")


def main() -> None:
    if len(argv) != 2:
        raise SystemExit("usage: use-signed-extension.py OUTPUT.xpi")
    try:
        verify()
    except (OSError, ValueError) as error:
        raise SystemExit(f"cannot package signed extension: {error}") from error
    output = Path(argv[1])
    output.parent.mkdir(parents=True, exist_ok=True)
    copyfile(SIGNED, output)


if __name__ == "__main__":
    main()
