package content

import "testing"

func TestLocalizedCatalogsKeepFeatureAndDocSlugs(t *testing.T) {
	en := CatalogFor(DefaultLocale())
	for _, locale := range SupportedLocales()[1:] {
		catalog := CatalogFor(locale)
		if len(catalog.Features) != len(en.Features) {
			t.Fatalf("%s features = %d, want %d", locale.Code, len(catalog.Features), len(en.Features))
		}
		if len(catalog.Docs) != len(en.Docs) {
			t.Fatalf("%s docs = %d, want %d", locale.Code, len(catalog.Docs), len(en.Docs))
		}
		for index, feature := range catalog.Features {
			if feature.Slug != en.Features[index].Slug {
				t.Fatalf("%s feature slug %d = %q, want %q", locale.Code, index, feature.Slug, en.Features[index].Slug)
			}
			if feature.Title == "" || feature.Summary == "" || feature.ImageAlt == "" {
				t.Fatalf("%s feature %q missing localized content", locale.Code, feature.Slug)
			}
		}
		for index, section := range catalog.Docs {
			if section.Slug != en.Docs[index].Slug {
				t.Fatalf("%s doc slug %d = %q, want %q", locale.Code, index, section.Slug, en.Docs[index].Slug)
			}
			if section.Title == "" || section.Summary == "" || section.Description == "" {
				t.Fatalf("%s doc %q missing localized content", locale.Code, section.Slug)
			}
		}
	}
}

func TestLocaleFromPath(t *testing.T) {
	tests := []struct {
		path       string
		code       string
		publicPath string
		ok         bool
	}{
		{path: "/", code: "en", publicPath: "/", ok: true},
		{path: "/features/ota", code: "en", publicPath: "/features/ota", ok: true},
		{path: "/zh-tw/features/ota", code: "zh-TW", publicPath: "/features/ota", ok: true},
		{path: "/zh-cn/contact", code: "zh-CN", publicPath: "/contact", ok: true},
		{path: "/fr/features", ok: false},
	}
	for _, tc := range tests {
		locale, publicPath, ok := LocaleFromPath(tc.path)
		if ok != tc.ok {
			t.Fatalf("%s ok = %v, want %v", tc.path, ok, tc.ok)
		}
		if !ok {
			continue
		}
		if locale.Code != tc.code || publicPath != tc.publicPath {
			t.Fatalf("%s = (%s, %s), want (%s, %s)", tc.path, locale.Code, publicPath, tc.code, tc.publicPath)
		}
	}
}
