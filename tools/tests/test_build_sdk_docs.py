from __future__ import annotations

import importlib.util
import sys
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location("build_sdk_docs", ROOT / "tools/build_sdk_docs.py")
assert SPEC and SPEC.loader
docs = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = docs
SPEC.loader.exec_module(docs)


class BuildSDKDocsTest(unittest.TestCase):
    def test_canonical_manual_is_complete_and_links_are_valid(self) -> None:
        source = ROOT / "content/manual/sdk"
        index, pages = docs.load_pages(source)
        self.assertEqual(index["title"], "RTK Cloud SDK User Manual")
        self.assertEqual(set(docs.PACKAGES), {page.package for page in pages if page.package})
        docs.validate_links(source, pages)

    def test_frontmatter_requires_title_and_description(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "page.en.md"
            path.write_text("---\ntitle: Example\n---\nBody\n", encoding="utf-8")
            with self.assertRaisesRegex(docs.DocsError, "description is required"):
                docs.parse_frontmatter(path)

    def test_local_link_validation_rejects_missing_asset(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory)
            page_path = source / "page.en.md"
            page_path.write_text("unused", encoding="utf-8")
            page = docs.Page("page", "Page", "Description", None, page_path, "![missing](assets/no.png)\n")
            with self.assertRaisesRegex(docs.DocsError, "broken local link"):
                docs.validate_links(source, [page])

    def test_generated_html_validation_rejects_broken_link(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "index.html").write_text('<a href="missing/index.html">Missing</a>', encoding="utf-8")
            with self.assertRaisesRegex(docs.DocsError, "broken generated link"):
                docs.validate_generated_html(root)

    def test_public_source_inventory_covers_every_package(self) -> None:
        sdk_repo = ROOT.parent / "rtk_cloud_client"
        if not sdk_repo.is_dir():
            self.skipTest("rtk_cloud_client checkout is not available")
        references = docs.build_references(sdk_repo)
        self.assertEqual(set(docs.PACKAGES), {reference.package for reference in references})
        for reference in references:
            self.assertGreater(len(reference.symbols), 0, reference.package)


if __name__ == "__main__":
    unittest.main()
