# SDK Documentation Publishing

Status: active

Owner: rtk_cloud_frontend

Last reviewed: 2026-07-21

The canonical SDK user manual lives in `content/manual/sdk/`. Do not copy its
prose into package READMEs or another website directory. The website renders
the Markdown under `/manual/sdk`, and the documentation generator produces
offline HTML, API symbol indexes, and PDF downloads from the same source.

## Build

From the frontend repository with an adjacent `rtk_cloud_client` checkout:

```sh
python3 tools/build_sdk_docs.py check \
  --source content/manual/sdk \
  --sdk-repo ../rtk_cloud_client

python3 tools/build_sdk_docs.py build \
  --source content/manual/sdk \
  --sdk-repo ../rtk_cloud_client \
  --version 0.1.0-rc1 \
  --output dist/sdk-docs
```

The build requires Python 3 with PyYAML, Pandoc, and WeasyPrint. It generates a
source-aware public symbol inventory for each SDK package, records the exact
SDK commit, produces stable and versioned PDF names, writes checksums, and
updates `dist/sdk-docs/current`.

Generated output is intentionally ignored by Git. Release packaging includes
the generated tree and fails when the HTML, complete PDF, or documentation
manifest is missing.

## Authoring

- Add English pages below `content/manual/sdk/` with YAML frontmatter.
- Register every page in `index.en.yaml`; its order controls navigation and PDF
  chapter order.
- Use repository-relative links for manual-owned assets.
- Keep credentials, keys, presigned URLs, and customer data out of examples.
- Regenerate after any exported SDK API change and review the symbol counts.

## Verification

Run the Python unit tests, Go tests, documentation build, `pdfinfo`, PDF text
extraction, and page rendering with `pdftoppm`. Inspect the rendered pages for
clipping, broken code blocks, missing glyphs, bad page breaks, and unreadable
tables before publishing.
