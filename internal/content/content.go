package content

import (
	"strings"

	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
)

type Locale struct {
	Code   string
	Lang   string
	Prefix string
	Label  string
}

type AlternateLink struct {
	HrefLang string
	Label    string
	Href     string
	Current  bool
}

type PageMeta struct {
	Title       string
	Description string
}

type Catalog struct {
	Locale   Locale
	Text     map[string]string
	Pages    map[string]PageMeta
	Features []features.Feature
	Docs     []docs.Section
}

var supportedLocales = []Locale{
	{Code: "en", Lang: "en", Prefix: "", Label: "English"},
	{Code: "zh-TW", Lang: "zh-Hant", Prefix: "/zh-tw", Label: "繁體中文"},
	{Code: "zh-CN", Lang: "zh-Hans", Prefix: "/zh-cn", Label: "简体中文"},
}

func SupportedLocales() []Locale {
	locales := make([]Locale, len(supportedLocales))
	copy(locales, supportedLocales)
	return locales
}

func DefaultLocale() Locale {
	return supportedLocales[0]
}

func LocaleFromPath(path string) (Locale, string, bool) {
	clean := "/" + strings.Trim(strings.TrimSpace(path), "/")
	if clean == "/" {
		return DefaultLocale(), "/", true
	}
	for _, locale := range supportedLocales[1:] {
		if clean == locale.Prefix {
			return locale, "/", true
		}
		if strings.HasPrefix(clean, locale.Prefix+"/") {
			trimmed := strings.TrimPrefix(clean, locale.Prefix)
			if trimmed == "" {
				trimmed = "/"
			}
			return locale, trimmed, true
		}
	}
	firstSegment := strings.Trim(strings.Split(strings.TrimPrefix(clean, "/"), "/")[0], " ")
	if strings.HasPrefix(clean, "/zh-") || strings.HasPrefix(clean, "/en/") || clean == "/en" || len(firstSegment) == 2 {
		return Locale{}, "", false
	}
	return DefaultLocale(), clean, true
}

func PathForLocale(locale Locale, publicPath string) string {
	path := "/" + strings.Trim(strings.TrimSpace(publicPath), "/")
	if path == "/" {
		return locale.Prefix + "/"
	}
	return locale.Prefix + path
}

func CatalogFor(locale Locale) Catalog {
	switch locale.Code {
	case "zh-TW":
		return zhTWCatalog()
	case "zh-CN":
		return zhCNCatalog()
	default:
		return enCatalog()
	}
}

func (c Catalog) T(key string) string {
	if value, ok := c.Text[key]; ok {
		return value
	}
	return enText()[key]
}

func (c Catalog) Page(key string) PageMeta {
	if page, ok := c.Pages[key]; ok {
		return page
	}
	return enPages()[key]
}

func (c Catalog) FeatureBySlug(slug string) (features.Feature, bool) {
	for _, feature := range c.Features {
		if feature.Slug == slug {
			return feature, true
		}
	}
	return features.Feature{}, false
}

func (c Catalog) DocBySlug(slug string) (docs.Section, bool) {
	for _, section := range c.Docs {
		if section.Slug == slug {
			return section, true
		}
	}
	return docs.Section{}, false
}

func enCatalog() Catalog {
	return Catalog{
		Locale:   supportedLocales[0],
		Text:     enText(),
		Pages:    enPages(),
		Features: features.All(),
		Docs:     docs.All(),
	}
}

func enPages() map[string]PageMeta {
	return map[string]PageMeta{
		"home": {
			Title:       "Realtek Connect+ | IoT Cloud Platform",
			Description: "Realtek Connect+ is an IoT cloud platform for provisioning, OTA, fleet management, app SDKs, insights, private cloud, and integrations.",
		},
		"features": {
			Title:       "Features | Realtek Connect+",
			Description: "Explore provisioning, OTA, fleet management, app SDK, insights, private cloud, and ecosystem integrations for Realtek-based IoT products.",
		},
		"docs": {
			Title:       "Developer Docs | Realtek Connect+",
			Description: "Browse Realtek Connect+ documentation entry points for product overview, development, APIs, SDKs, firmware, CLI, deployment, and release notes.",
		},
		"contact": {
			Title:       "Contact | Realtek Connect+",
			Description: "Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
		},
	}
}

func enText() map[string]string {
	return map[string]string{
		"skip.main":                  "Skip to main content",
		"brand.home":                 "Realtek Connect+ home",
		"nav.docs":                   "Docs",
		"nav.features":               "Features",
		"nav.architecture":           "Architecture",
		"nav.contact":                "Contact",
		"footer.tagline":             "IoT cloud platform concept for Realtek-based connected products.",
		"footer.docs":                "Developer Docs",
		"footer.features":            "Features",
		"footer.contact":             "Contact Us",
		"home.eyebrow":               "IoT cloud platform for product teams",
		"home.lede":                  "Bring Realtek-based devices online with provisioning, OTA, fleet operations, app SDKs, insights, private cloud options, and ecosystem integrations.",
		"home.cta.primary":           "Contact Us",
		"home.cta.secondary":         "Explore Services",
		"home.chip.silicon":          "Silicon",
		"home.chip.sdk":              "Device SDK",
		"home.chip.cloud":            "Cloud",
		"home.chip.ops":              "Ops",
		"home.overview.eyebrow":      "Platform overview",
		"home.overview.title":        "From chipset to connected product lifecycle.",
		"home.overview.lede":         "Connect+ frames Realtek connectivity silicon, firmware enablement, cloud services, mobile app workflows, and operations tooling as one commercial IoT product path.",
		"home.surfaces.eyebrow":      "Platform surfaces",
		"home.surfaces.title":        "Show the product system, not only the feature list.",
		"home.surfaces.card.title":   "Provisioning, OTA, and fleet health share one visual language.",
		"home.surfaces.card.body":    "The public site now anchors the cloud story with product surfaces that connect devices, secure cloud operations, and dashboard workflows.",
		"home.surface.onboarding":    "Onboarding",
		"home.surface.rollouts":      "Rollouts",
		"home.surface.insights":      "Insights",
		"home.surface.security":      "Security",
		"home.principles.eyebrow":    "Design principles",
		"home.principles.title":      "Built for enterprise connected-product programs.",
		"home.principle.active":      "Manageability",
		"home.principle.scale":       "Scalability",
		"home.principle.security":    "Security",
		"home.principle.privacy":     "Privacy",
		"home.principle.cost":        "Cost control",
		"home.principle.custom":      "Customizability",
		"home.principle.panel":       "Operate firmware, users, fleets, and support workflows from one platform story.",
		"home.principle.body":        "Realtek Connect+ is presented as a lifecycle system: onboarding, cloud identity, OTA, app SDKs, metrics, and enterprise deployment work together instead of living as disconnected product pages.",
		"home.services.eyebrow":      "Core services",
		"home.services.title":        "Realtek cloud capabilities, packaged for connected-product teams.",
		"home.feature.details":       "View details",
		"home.arch.eyebrow":          "Architecture",
		"home.arch.title":            "A direct flow from device onboarding to fleet operations.",
		"home.arch.device.title":     "Realtek Device SDK",
		"home.arch.device.body":      "Identity, provisioning, firmware services, and device signals.",
		"home.arch.cloud.title":      "Secure Cloud",
		"home.arch.cloud.body":       "Registry, OTA campaigns, user-device association, APIs.",
		"home.arch.app.title":        "App SDK + Dashboard",
		"home.arch.app.body":         "Mobile onboarding, product control, insights, support workflows.",
		"home.deploy.eyebrow":        "Public vs private cloud",
		"home.deploy.title":          "Start with evaluation, mature into controlled deployment.",
		"home.deploy.public":         "Public evaluation",
		"home.deploy.public.title":   "Validate product fit before committing to a private footprint.",
		"home.deploy.public.body":    "Use the public-facing website and documentation structure to align firmware, mobile, cloud, and product stakeholders before commercial deployment planning.",
		"home.deploy.docs":           "Deployment Docs",
		"home.deploy.private":        "Private commercial cloud",
		"home.deploy.private.title":  "Plan data ownership, custom domains, regional placement, and support boundaries.",
		"home.deploy.private.body":   "The private cloud narrative gives enterprise buyers a clear path from concept validation to branded, controlled, commercially supported operation.",
		"home.deploy.discuss":        "Discuss Private Cloud",
		"home.use.eyebrow":           "Use cases",
		"home.use.title":             "Built for commercial connected-device teams.",
		"home.use.smart.body":        "App onboarding, sharing, push notifications, voice assistant paths, and OTA maintenance.",
		"home.use.industrial.body":   "Fleet grouping, metadata, secure updates, private cloud, and operations visibility.",
		"home.use.appliance.body":    "Long lifecycle firmware maintenance, activation data, branded apps, and support diagnostics.",
		"home.docs.eyebrow":          "Developer portal",
		"home.docs.title":            "Give product, firmware, app, and cloud teams one documentation spine.",
		"home.docs.open":             "Open section",
		"home.cta.eyebrow":           "Early access",
		"home.cta.title":             "Plan a Realtek Connect+ product path.",
		"home.cta.body":              "Register your interest in provisioning, OTA, private cloud, app SDK, or fleet operations.",
		"features.eyebrow":           "Features",
		"features.title":             "Cloud services for the full IoT product lifecycle.",
		"features.body":              "Realtek Connect+ presents the core capabilities product teams expect from a commercial connected-device platform.",
		"features.open":              "Open feature",
		"feature.discuss.prefix":     "Discuss",
		"feature.all":                "All Features",
		"feature.highlights":         "Highlights",
		"feature.highlights.title":   "What this service covers",
		"feature.capabilities":       "Capabilities",
		"feature.capabilities.title": "Platform building blocks",
		"feature.outcomes":           "Outcomes",
		"feature.outcomes.title":     "Why product teams use it",
		"feature.next":               "Next step",
		"feature.cta.prefix":         "Evaluate",
		"feature.cta.suffix":         "for your product roadmap.",
		"feature.cta.body":           "Share your product category, target deployment, and cloud requirements with the Realtek Connect+ team.",
		"docs.eyebrow":               "Developer docs",
		"docs.title":                 "Documentation entry points for cloud, firmware, app, and deployment teams.",
		"docs.body":                  "Realtek Connect+ now includes a server-rendered documentation portal structure covering platform overview, development, APIs, SDKs, firmware, CLI workflows, deployment, and release notes.",
		"docs.cta.primary":           "Talk to the platform team",
		"docs.cta.secondary":         "See app platform context",
		"docs.portal.eyebrow":        "Portal structure",
		"docs.portal.title":          "Choose the track that matches your workstream.",
		"docs.why.eyebrow":           "Why this matters",
		"docs.why.title":             "Mirror the documentation surfaces teams expect before product launch.",
		"docs.shared.title":          "Shared platform narrative",
		"docs.shared.body":           "Give product, sales, and engineering stakeholders one place to understand lifecycle scope before implementation deep-dives.",
		"docs.depth.title":           "Workstream-specific depth",
		"docs.depth.body":            "Separate firmware, API, mobile SDK, deployment, and release concerns so each team can navigate directly to its implementation surface.",
		"docs.static.title":          "Static first version",
		"docs.static.body":           "Keep the docs portal compatible with the current Go templates and server-rendered architecture while leaving space for deeper follow-on content.",
		"doc.back":                   "Back to docs",
		"doc.discuss":                "Discuss implementation",
		"doc.coverage":               "Coverage",
		"doc.coverage.title":         "What this section explains",
		"doc.outputs":                "Expected outputs",
		"doc.outputs.title":          "What teams should be able to leave with",
		"doc.audience":               "Primary audience",
		"doc.audience.title":         "Who should start here",
		"doc.next":                   "Next sections",
		"doc.next.title":             "Continue through the developer portal.",
		"doc.view":                   "View section",
		"contact.eyebrow":            "Contact",
		"contact.title":              "Register interest in Realtek Connect+.",
		"contact.body":               "Tell us which service matters most for your product team. Requests are stored locally in SQLite for this first version.",
		"contact.context.eyebrow":    "Early access",
		"contact.context.title":      "For IoT product planning, firmware maintenance, and private cloud evaluation.",
		"contact.context.body":       "Use this form for provisioning, OTA, fleet management, app SDK, insights, private cloud, or integration discussions.",
		"contact.thanks":             "Thanks",
		"contact.recorded":           "Your Realtek Connect+ request has been recorded.",
		"contact.review":             "Review features",
		"contact.error.summary":      "Please review the details below before submitting.",
		"contact.website":            "Website",
		"contact.name":               "Name",
		"contact.company":            "Company",
		"contact.email":              "Email",
		"contact.interest":           "Interest",
		"contact.select":             "Select a service",
		"contact.message":            "Message",
		"contact.submit":             "Submit Request",
	}
}
