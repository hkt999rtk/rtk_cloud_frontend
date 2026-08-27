package content

import (
	"strings"
	"testing"

	"realtek-connect/internal/features"
)

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

func TestCloudPlansLeadWithRealtekManagedServiceInEveryLocale(t *testing.T) {
	for _, locale := range SupportedLocales() {
		feature, ok := CatalogFor(locale).FeatureBySlug("private-cloud")
		if !ok {
			t.Fatalf("%s catalog missing cloud plans feature", locale.Code)
		}
		if len(feature.Sections) < 2 {
			t.Fatalf("%s cloud plans sections = %d, want at least 2", locale.Code, len(feature.Sections))
		}
		managed := strings.ToLower(feature.Sections[0].Title + " " + feature.Sections[0].Intro)
		private := strings.ToLower(feature.Sections[1].Title + " " + feature.Sections[1].Intro)
		if !strings.Contains(managed, "realtek") || !strings.Contains(managed, "managed") && !strings.Contains(managed, "託管") && !strings.Contains(managed, "托管") {
			t.Fatalf("%s first cloud plan is not the Realtek-managed service: %s", locale.Code, managed)
		}
		if !strings.Contains(private, "private cloud") {
			t.Fatalf("%s second cloud plan is not Private Cloud: %s", locale.Code, private)
		}
		copy := strings.ToLower(feature.Description + " " + strings.Join(feature.Highlights, " "))
		for _, forbidden := range []string{"one-time license", "annual maintenance", "unit price"} {
			if strings.Contains(copy, forbidden) {
				t.Fatalf("%s cloud plans copy contains unapproved pricing detail %q: %s", locale.Code, forbidden, copy)
			}
		}
	}
}

func TestManagedCloudSimplifiedChineseCopy(t *testing.T) {
	feature, ok := CatalogFor(localeByCode(t, "zh-CN")).FeatureBySlug("private-cloud")
	if !ok {
		t.Fatal("zh-CN catalog missing cloud plans feature")
	}
	copy := strings.Join([]string{feature.Summary, feature.Description}, " ")
	for _, want := range []string{"Realtek 托管", "客户", "实际使用量付费"} {
		if !strings.Contains(copy, want) {
			t.Fatalf("zh-CN managed cloud copy %q does not contain %q", copy, want)
		}
	}
	for _, unwanted := range []string{"託管", "客戶", "實際", "付費"} {
		if strings.Contains(copy, unwanted) {
			t.Fatalf("zh-CN managed cloud copy %q still contains %q", copy, unwanted)
		}
	}
}

func TestOTACopyPromotesImplementedCampaignPolicies(t *testing.T) {
	catalog := CatalogFor(DefaultLocale())
	ota, ok := catalog.FeatureBySlug("ota")
	if !ok {
		t.Fatal("default catalog missing ota feature")
	}

	assertOTATableStatus(t, ota, "Scheduled policy", "Available foundation")
	assertOTATableStatus(t, ota, "Time-window policy", "Available foundation")
	assertOTATableStatus(t, ota, "User-consent policy", "Available foundation")
	assertOTATableStatus(t, ota, "Archive", "Available foundation")

	copy := strings.Join([]string{
		ota.Summary,
		ota.Description,
		strings.Join(ota.Highlights, " "),
		strings.Join(ota.Capabilities, " "),
	}, " ")
	for _, want := range []string{
		"available foundations",
		"Available campaign policy surfaces",
		"implemented",
		"roadmap scope",
	} {
		if !strings.Contains(copy, want) {
			t.Fatalf("OTA copy does not contain availability boundary %q: %s", want, copy)
		}
	}
}

func TestLocalizedOTACopyPromotesImplementedCampaignPolicies(t *testing.T) {
	tests := []struct {
		localeCode string
		wants      []string
	}{
		{
			localeCode: "zh-TW",
			wants:      []string{"現有基礎", "已可用的 campaign surface", "roadmap 範圍"},
		},
		{
			localeCode: "zh-CN",
			wants:      []string{"现有基础", "已可用的 campaign surface", "roadmap 范圍"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.localeCode, func(t *testing.T) {
			catalog := CatalogFor(localeByCode(t, tc.localeCode))
			ota, ok := catalog.FeatureBySlug("ota")
			if !ok {
				t.Fatal("localized catalog missing ota feature")
			}
			copy := strings.Join([]string{
				ota.Summary,
				ota.Description,
				strings.Join(ota.Highlights, " "),
				strings.Join(ota.Capabilities, " "),
			}, " ")
			for _, want := range tc.wants {
				if !strings.Contains(copy, want) {
					t.Fatalf("%s OTA copy does not contain availability boundary %q: %s", tc.localeCode, want, copy)
				}
			}
		})
	}
}

func assertOTATableStatus(t *testing.T, feature features.Feature, concept, status string) {
	t.Helper()

	for _, row := range feature.Table.Rows {
		if len(row.Cells) < 2 || row.Cells[0] != concept {
			continue
		}
		if row.Cells[1] != status {
			t.Fatalf("OTA table status for %q = %q, want %q", concept, row.Cells[1], status)
		}
		return
	}
	t.Fatalf("OTA table missing concept %q", concept)
}

func localeByCode(t *testing.T, code string) Locale {
	t.Helper()

	for _, locale := range SupportedLocales() {
		if locale.Code == code {
			return locale
		}
	}
	t.Fatalf("missing supported locale %q", code)
	return Locale{}
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

func TestToSimplifiedCoversProvisioningContractTerms(t *testing.T) {
	got := ToSimplified("合約支撐的基礎")
	want := "合约支撑的基础"
	if got != want {
		t.Fatalf("ToSimplified() = %q, want %q", got, want)
	}
}
