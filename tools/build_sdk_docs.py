#!/usr/bin/env python3
"""Build and validate the single-source RTK Cloud SDK user manual."""

from __future__ import annotations

import argparse
import dataclasses
import datetime as dt
import hashlib
import html
import json
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Iterable

import yaml


LOCALES = ("en",)
PACKAGES = ("native", "android", "ios", "javascript", "go", "freertos-pro2")
REQUIRED_TOPICS = {
    "overview",
    "getting-started",
    "authentication-security",
    "lifecycle-errors",
    "capability-workflows",
    "video-workflows",
    "sample-applications",
    "troubleshooting",
    *(f"packages/{name}" for name in PACKAGES),
}
FRONTMATTER_RE = re.compile(r"\A---\s*\n(.*?)\n---\s*\n(.*)\Z", re.DOTALL)
LOCAL_LINK_RE = re.compile(r"!?\[[^]]*]\(([^)]+)\)")
HTML_LINK_RE = re.compile(r"(?:href|src)=[\"']([^\"']+)[\"']", re.IGNORECASE)


class DocsError(RuntimeError):
    pass


@dataclasses.dataclass(frozen=True)
class Page:
    slug: str
    title: str
    description: str
    package: str | None
    source_path: Path
    markdown: str


@dataclasses.dataclass(frozen=True)
class PackageReference:
    package: str
    source_paths: tuple[str, ...]
    symbols: tuple[str, ...]


def run(command: list[str], *, cwd: Path | None = None, stdout_path: Path | None = None) -> str:
    try:
        completed = subprocess.run(
            command,
            cwd=cwd,
            check=True,
            text=True,
            stdout=subprocess.PIPE if stdout_path is None else stdout_path.open("w", encoding="utf-8"),
            stderr=subprocess.PIPE,
        )
    except FileNotFoundError as exc:
        raise DocsError(f"required tool not found: {command[0]}") from exc
    except subprocess.CalledProcessError as exc:
        detail = (exc.stderr or "").strip()
        raise DocsError(f"command failed ({' '.join(command)}): {detail}") from exc
    return completed.stdout or ""


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def parse_frontmatter(path: Path) -> tuple[dict[str, object], str]:
    match = FRONTMATTER_RE.match(path.read_text(encoding="utf-8"))
    if not match:
        raise DocsError(f"{path}: missing YAML frontmatter")
    metadata = yaml.safe_load(match.group(1)) or {}
    if not isinstance(metadata, dict):
        raise DocsError(f"{path}: frontmatter must be a mapping")
    for field in ("title", "description"):
        if not str(metadata.get(field, "")).strip():
            raise DocsError(f"{path}: {field} is required")
    return metadata, match.group(2).strip() + "\n"


def load_pages(source: Path) -> tuple[dict[str, object], list[Page]]:
    index_path = source / "index.en.yaml"
    if not index_path.is_file():
        raise DocsError(f"manual index not found: {index_path}")
    index = yaml.safe_load(index_path.read_text(encoding="utf-8")) or {}
    if not isinstance(index, dict) or not isinstance(index.get("sections"), list):
        raise DocsError(f"{index_path}: sections list is required")

    pages: list[Page] = []
    seen: set[str] = set()
    for entry in index["sections"]:
        if not isinstance(entry, dict):
            raise DocsError(f"{index_path}: every section must be a mapping")
        slug = str(entry.get("slug", "")).strip().strip("/")
        if not slug or slug in seen or ".." in Path(slug).parts:
            raise DocsError(f"{index_path}: invalid or duplicate slug {slug!r}")
        seen.add(slug)
        page_path = source / f"{slug}.en.md"
        if not page_path.is_file():
            raise DocsError(f"indexed page not found: {page_path}")
        metadata, markdown = parse_frontmatter(page_path)
        package = str(entry.get("package", "")).strip() or None
        if package is not None and package not in PACKAGES:
            raise DocsError(f"{index_path}: unsupported package {package!r}")
        pages.append(Page(
            slug=slug,
            title=str(metadata["title"]).strip(),
            description=str(metadata["description"]).strip(),
            package=package,
            source_path=page_path,
            markdown=markdown,
        ))

    missing = REQUIRED_TOPICS.difference(seen)
    if missing:
        raise DocsError(f"manual index is missing required topics: {', '.join(sorted(missing))}")
    unindexed = {
        path.relative_to(source).as_posix().removesuffix(".en.md")
        for path in source.rglob("*.en.md")
    }.difference(seen)
    if unindexed:
        raise DocsError(f"manual pages are not indexed: {', '.join(sorted(unindexed))}")
    return index, pages


def validate_links(source: Path, pages: Iterable[Page]) -> None:
    known = {page.slug for page in pages}
    errors: list[str] = []
    for page in pages:
        for raw_target in LOCAL_LINK_RE.findall(page.markdown):
            target = raw_target.split()[0].strip("<>").split("#", 1)[0]
            if target.startswith("/content-assets/manual/sdk/"):
                asset = source / "assets" / target.removeprefix("/content-assets/manual/sdk/")
                if not asset.is_file():
                    errors.append(f"{page.source_path}: broken manual asset {raw_target}")
                continue
            if not target or target.startswith(("http://", "https://", "mailto:", "/")):
                continue
            resolved = (page.source_path.parent / target).resolve()
            if target.endswith(".md"):
                try:
                    relative = resolved.relative_to(source.resolve()).as_posix()
                except ValueError:
                    relative = ""
                slug = re.sub(r"\.(en\.)?md$", "", relative)
                if slug in known:
                    continue
            if not resolved.exists():
                errors.append(f"{page.source_path}: broken local link {raw_target}")
    if errors:
        raise DocsError("\n".join(errors))


def validate_generated_html(root: Path) -> None:
    errors: list[str] = []
    for page in root.rglob("*.html"):
        for target in HTML_LINK_RE.findall(page.read_text(encoding="utf-8")):
            clean = target.split("#", 1)[0].split("?", 1)[0]
            if not clean or clean.startswith(("http://", "https://", "mailto:", "data:", "/")):
                continue
            resolved = (page.parent / clean).resolve()
            if clean.endswith("/"):
                resolved = resolved / "index.html"
            if not resolved.exists():
                errors.append(f"{page}: broken generated link {target}")
    if errors:
        raise DocsError("\n".join(errors))


def git_value(repo: Path, *args: str) -> str:
    return run(["git", "-C", str(repo), *args]).strip()


def package_version(sdk_repo: Path, requested: str | None) -> str:
    if requested:
        if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._-]*", requested):
            raise DocsError("version may contain only letters, digits, dots, underscores, and dashes")
        return requested
    package_json = sdk_repo / "packages/javascript/package.json"
    if package_json.is_file():
        value = json.loads(package_json.read_text(encoding="utf-8")).get("version")
        if value and value != "0.0.0":
            return str(value)
    return "development"


def public_symbols(text: str, patterns: Iterable[str]) -> set[str]:
    result: set[str] = set()
    for pattern in patterns:
        for match in re.finditer(pattern, text, re.MULTILINE):
            value = " ".join(match.group(0).strip().split())
            if value:
                result.add(value)
    return result


def build_references(sdk_repo: Path) -> list[PackageReference]:
    specs = {
        "native": (
            ("packages/native/include/rtk_cloud_client/rtkc.h", "packages/native/include/rtk_cloud_client/rtkc.hpp"),
            (r"^(?:rtkc_status_t|void|int|size_t|const\s+char\*)\s+rtkc_[A-Za-z0-9_]+\s*\([^;]*;", r"^typedef\s+(?:struct|enum)\s+[A-Za-z0-9_]+", r"^class\s+[A-Za-z0-9_]+", r"^\s+Status\s+[a-z_][A-Za-z0-9_]*\s*\([^;{]*"),
        ),
        "android": (
            ("packages/android/rtk-cloud-client/src/main/kotlin/com/rtk/cloud/client/RtkCloudClientPackage.kt", "packages/android/rtk-cloud-client/src/main/kotlin/com/rtk/cloud/client/AndroidPkiDeviceAuth.kt"),
            (r"^(?:public\s+)?(?:data\s+|enum\s+|sealed\s+)?class\s+[A-Za-z0-9_]+[^\n]*", r"^\s*(?:public\s+)?(?:suspend\s+)?fun\s+[A-Za-z0-9_]+[^\n]*", r"^fun\s+interface\s+[A-Za-z0-9_]+[^\n]*", r"^(?:public\s+)?(?:interface|object)\s+[A-Za-z0-9_]+[^\n]*"),
        ),
        "ios": (
            ("packages/ios/Sources/RTKCloudClient/RTKCloudClientPackage.swift", "packages/ios/Sources/RTKCloudClient/PKI.swift"),
            (r"^public\s+(?:struct|class|enum|protocol|actor)\s+[A-Za-z0-9_]+[^\n]*", r"^\s*public\s+(?:static\s+)?func\s+[A-Za-z0-9_]+[^\n]*", r"^\s*public\s+(?:convenience\s+)?init\s*\([^\n]*", r"^\s*public\s+(?:static\s+)?(?:let|var)\s+[A-Za-z0-9_]+[^\n]*"),
        ),
        "javascript": (
            ("packages/javascript/src/index.ts",),
            (r"^export\s+(?:declare\s+)?(?:class|interface|type|enum|function|const)\s+[A-Za-z0-9_]+[^\n]*",),
        ),
        "go": (
            tuple(
                path.relative_to(sdk_repo).as_posix()
                for path in sorted((sdk_repo / "packages/golang/rtkc").rglob("*.go"))
                if not path.name.endswith("_test.go") and "/internal/" not in path.as_posix()
            ),
            (r"^type\s+[A-Z][A-Za-z0-9_]*[^\n]*", r"^func\s+(?:\([^)]*\)\s*)?[A-Z][A-Za-z0-9_]*[^\n]*"),
        ),
        "freertos-pro2": (
            tuple(
                path.relative_to(sdk_repo).as_posix()
                for path in sorted((sdk_repo / "packages/freertos/pro2_demo/include").rglob("*.h"))
            ),
            (r"^(?:rtkc_status_t|void|int|size_t)\s+[A-Za-z0-9_]+\s*\([^;]*;", r"^typedef\s+(?:struct|enum)\s+[A-Za-z0-9_]+"),
        ),
    }
    references: list[PackageReference] = []
    for package, (relative_paths, patterns) in specs.items():
        symbols: set[str] = set()
        for relative in relative_paths:
            path = sdk_repo / relative
            if not path.is_file():
                raise DocsError(f"SDK public source is missing: {path}")
            symbols.update(public_symbols(path.read_text(encoding="utf-8"), patterns))
        if not symbols:
            raise DocsError(f"no public symbols discovered for {package}")
        references.append(PackageReference(package, tuple(relative_paths), tuple(sorted(symbols))))
    return references


def reference_markdown(reference: PackageReference, sdk_commit: str) -> str:
    lines = [
        f"# {reference.package.replace('-', ' ').title()} API symbol index",
        "",
        f"Generated from SDK commit `{sdk_commit}`. This inventory is generated from public source files and must not be edited.",
        "",
        "## Public source files",
        "",
        *(f"- `{path}`" for path in reference.source_paths),
        "",
        f"## Exported symbols ({len(reference.symbols)})",
        "",
        *(f"- `{symbol}`" for symbol in reference.symbols),
        "",
    ]
    return "\n".join(lines)


def write_css(path: Path) -> None:
    path.write_text("""
:root { color-scheme: light; --ink:#172033; --muted:#5b6475; --brand:#0068b5; --line:#d8dee8; }
* { box-sizing: border-box; }
body { color:var(--ink); font:16px/1.62 -apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif; margin:0 auto; max-width:1120px; padding:40px 56px; }
nav { background:#f4f7fb; border:1px solid var(--line); border-radius:12px; padding:18px 24px; }
h1,h2,h3 { line-height:1.2; margin-top:1.65em; page-break-after:avoid; }
h1 { color:var(--brand); font-size:2.25rem; border-bottom:3px solid var(--brand); margin-bottom:.85em; padding-bottom:.35em; }
h2 { font-size:1.55rem; }
a { color:var(--brand); }
pre { background:#111827; color:#f8fafc; border-radius:8px; overflow:auto; padding:16px; white-space:pre-wrap; }
code { font-family:"SFMono-Regular",Consolas,monospace; font-size:.9em; }
table { border-collapse:collapse; width:100%; }
th,td { border:1px solid var(--line); padding:8px 10px; text-align:left; vertical-align:top; }
th { background:#eef4fa; }
img { height:auto; max-width:100%; }
.doc-meta { color:var(--muted); font-size:.9rem; }
@page { size:A4; margin:18mm 17mm 20mm; @bottom-center { content:"Realtek Connect+ SDK Manual  •  " counter(page); color:#6b7280; font-size:9pt; } }
@media print { body { font-size:10.5pt; max-width:none; padding:0; } nav { display:none; } h1 { break-before:page; } h1:first-child { break-before:auto; } pre,table,img { break-inside:avoid; } a { color:inherit; text-decoration:none; } }
""".strip() + "\n", encoding="utf-8")


def pandoc(markdown_path: Path, output: Path, *, title: str, css: Path, toc: bool = True, resource_path: Path | None = None, embed_resources: bool = False) -> None:
    command = [
        "pandoc", str(markdown_path), "--from=gfm", "--to=html5", "--standalone",
        "--metadata", f"title={title}", "--css", os.path.relpath(css, output.parent), "--output", str(output),
    ]
    if resource_path is not None:
        command.extend(["--resource-path", str(resource_path)])
    if embed_resources:
        command.append("--embed-resources")
    if toc:
        command.extend(["--toc", "--toc-depth=3"])
    run(command)


def markdown_for_page(page: Page, pages: list[Page]) -> str:
    navigation = ["<nav><strong>SDK manual</strong><ul>"]
    root_depth = len(Path(page.slug).parts)
    prefix = "../" * root_depth
    for candidate in pages:
        target = prefix + candidate.slug + "/index.html"
        navigation.append(f'<li><a href="{html.escape(target)}">{html.escape(candidate.title)}</a></li>')
    navigation.append("</ul></nav>")
    body = page.markdown.replace("](/content-assets/manual/sdk/", f"]({prefix}assets/manual/")
    return "\n".join(navigation) + "\n\n" + f"# {page.title}\n\n{page.description}\n\n" + body


def build(source: Path, sdk_repo: Path, output_root: Path, version: str) -> Path:
    index, pages = load_pages(source)
    validate_links(source, pages)
    sdk_repo = sdk_repo.resolve()
    if not (sdk_repo / ".git").exists():
        raise DocsError(f"SDK repository is not a Git checkout: {sdk_repo}")
    sdk_commit = git_value(sdk_repo, "rev-parse", "HEAD")
    references = build_references(sdk_repo)

    root = output_root / version
    if root.exists():
        shutil.rmtree(root)
    html_root = root / "html"
    pdf_root = root / "pdf"
    assets_root = html_root / "assets"
    work_root = root / ".work"
    for path in (html_root, pdf_root, assets_root, work_root):
        path.mkdir(parents=True, exist_ok=True)
    css = assets_root / "sdk-docs.css"
    write_css(css)
    source_assets = source / "assets"
    if source_assets.is_dir():
        shutil.copytree(source_assets, assets_root / "manual", dirs_exist_ok=True)

    combined_parts: list[str] = []
    package_parts: dict[str, list[str]] = {name: [] for name in PACKAGES}
    for page in pages:
        target = html_root / page.slug / "index.html"
        target.parent.mkdir(parents=True, exist_ok=True)
        md_path = work_root / f"{page.slug.replace('/', '-')}.md"
        rendered_markdown = markdown_for_page(page, pages)
        md_path.write_text(rendered_markdown, encoding="utf-8")
        pandoc(md_path, target, title=f"{page.title} | RTK Cloud SDK", css=css)
        printable_markdown = page.markdown.replace("](/content-assets/manual/sdk/", "](assets/")
        combined_parts.append(f"# {page.title}\n\n{page.description}\n\n{printable_markdown}")
        if page.package:
            package_parts[page.package].append(combined_parts[-1])

    reference_pages: dict[str, str] = {}
    for reference in references:
        reference_md = reference_markdown(reference, sdk_commit)
        reference_pages[reference.package] = reference_md
        target = html_root / "reference" / reference.package / "index.html"
        target.parent.mkdir(parents=True, exist_ok=True)
        md_path = work_root / f"reference-{reference.package}.md"
        md_path.write_text(reference_md, encoding="utf-8")
        pandoc(md_path, target, title=f"{reference.package.title()} API reference", css=css)

    reference_index = work_root / "reference-index.md"
    reference_index.write_text("\n".join([
        "# SDK API references",
        "",
        f"Generated from SDK commit `{sdk_commit}`.",
        "",
        *(f"- [{name.replace('-', ' ').title()}]({name}/)" for name in PACKAGES),
        "",
    ]), encoding="utf-8")
    pandoc(reference_index, html_root / "reference" / "index.html", title="SDK API references", css=css)

    landing = work_root / "index.md"
    landing_lines = [
        "# RTK Cloud SDK User Manual",
        "",
        str(index.get("description", "Integration and API guidance for RTK Cloud SDK users.")),
        "",
        "## Manual chapters",
        "",
        *(f"- [{page.title}]({page.slug}/index.html) — {page.description}" for page in pages),
        "",
        "## API references",
        "",
        *(f"- [{name.replace('-', ' ').title()}](reference/{name}/index.html)" for name in PACKAGES),
        "",
    ]
    landing.write_text("\n".join(landing_lines), encoding="utf-8")
    pandoc(landing, html_root / "index.html", title="RTK Cloud SDK User Manual", css=css)

    combined_md = work_root / "combined.md"
    combined_md.write_text("\n\n".join(combined_parts + [reference_pages[name] for name in PACKAGES]), encoding="utf-8")
    combined_html = work_root / "combined.html"
    pandoc(combined_md, combined_html, title=f"RTK Cloud SDK User Manual {version}", css=css, resource_path=source, embed_resources=True)
    combined_pdf = pdf_root / f"rtk-cloud-sdk-user-manual-{version}.pdf"
    run(["weasyprint", str(combined_html), str(combined_pdf)])

    for name in PACKAGES:
        package_md = work_root / f"package-{name}.md"
        shared = [combined_parts[i] for i, page in enumerate(pages) if page.slug in {"overview", "getting-started", "authentication-security", "lifecycle-errors", "troubleshooting"}]
        package_md.write_text("\n\n".join(shared + package_parts[name] + [reference_pages[name]]), encoding="utf-8")
        package_html = work_root / f"package-{name}.html"
        pandoc(package_md, package_html, title=f"RTK Cloud {name.title()} SDK {version}", css=css, resource_path=source, embed_resources=True)
        run(["weasyprint", str(package_html), str(pdf_root / f"rtk-cloud-{name}-sdk-{version}.pdf")])

    shutil.copy2(combined_pdf, pdf_root / "rtk-cloud-sdk-user-manual.pdf")
    for name in PACKAGES:
        shutil.copy2(
            pdf_root / f"rtk-cloud-{name}-sdk-{version}.pdf",
            pdf_root / f"rtk-cloud-{name}-sdk.pdf",
        )
    download_index = pdf_root / "index.html"
    download_index.write_text("<!doctype html><html><head><meta charset=\"utf-8\"><title>SDK PDF downloads</title></head><body><h1>SDK PDF downloads</h1><ul>" + "".join([
        '<li><a href="rtk-cloud-sdk-user-manual.pdf">Complete SDK user manual</a></li>',
        *(f'<li><a href="rtk-cloud-{name}-sdk.pdf">{html.escape(name.replace("-", " ").title())} SDK manual</a></li>' for name in PACKAGES),
    ]) + "</ul></body></html>\n", encoding="utf-8")

    validate_generated_html(html_root)

    shutil.rmtree(work_root)
    files = sorted(path for path in root.rglob("*") if path.is_file())
    generated_at = dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    manifest = {
        "schema_version": 1,
        "sdk_version": version,
        "sdk_commit": sdk_commit,
        "website_commit": git_value(Path(__file__).resolve().parents[1], "rev-parse", "HEAD"),
        "generated_at": generated_at,
        "locales": list(LOCALES),
        "packages": list(PACKAGES),
        "tools": {
            "pandoc": run(["pandoc", "--version"]).splitlines()[0],
            "weasyprint": run(["weasyprint", "--version"]).strip(),
            "api_inventory": "source-aware built-in extractor v1",
        },
        "files": [
            {"path": path.relative_to(root).as_posix(), "bytes": path.stat().st_size, "sha256": sha256(path)}
            for path in files
        ],
    }
    manifest_path = root / "manifest.json"
    manifest_path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    files.append(manifest_path)
    checksum_lines = [f"{sha256(path)}  {path.relative_to(root).as_posix()}" for path in sorted(files)]
    (root / "SHA256SUMS").write_text("\n".join(checksum_lines) + "\n", encoding="utf-8")
    current = output_root / "current"
    if current.exists() or current.is_symlink():
        current.unlink() if current.is_symlink() or current.is_file() else shutil.rmtree(current)
    try:
        current.symlink_to(version)
    except OSError:
        shutil.copytree(root, current)
    return root


def check(source: Path, sdk_repo: Path) -> None:
    _, pages = load_pages(source)
    validate_links(source, pages)
    git_value(sdk_repo.resolve(), "rev-parse", "HEAD")
    references = build_references(sdk_repo.resolve())
    missing = set(PACKAGES).difference(reference.package for reference in references)
    if missing:
        raise DocsError(f"missing API references: {', '.join(sorted(missing))}")
    for tool in ("pandoc", "weasyprint"):
        if shutil.which(tool) is None:
            raise DocsError(f"required tool not found: {tool}")


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=__doc__)
    subparsers = result.add_subparsers(dest="command", required=True)
    for name in ("build", "check"):
        command = subparsers.add_parser(name)
        command.add_argument("--source", type=Path, required=True)
        command.add_argument("--sdk-repo", type=Path, required=True)
        if name == "build":
            command.add_argument("--version")
            command.add_argument("--output", type=Path, required=True)
    return result


def main() -> int:
    args = parser().parse_args()
    try:
        if args.command == "check":
            check(args.source.resolve(), args.sdk_repo.resolve())
            print("SDK documentation source check passed")
        else:
            version = package_version(args.sdk_repo.resolve(), args.version)
            output = build(args.source.resolve(), args.sdk_repo.resolve(), args.output.resolve(), version)
            print(output)
    except (DocsError, OSError, ValueError, yaml.YAMLError) as exc:
        print(f"SDK documentation error: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
