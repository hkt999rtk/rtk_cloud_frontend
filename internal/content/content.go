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

// ContactInterestOption pairs a stable interest slug (stored in the leads
// table and validated against leads.AllowedInterests) with the locale-aware
// label shown in the contact form dropdown.
type ContactInterestOption struct {
	Slug  string
	Label string
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

// ContactInterestOptions returns the inquiry-type dropdown options for the
// contact form in the catalog's locale. Each option carries the stable slug
// stored in the leads table plus the locale-specific label shown to visitors.
// Order is stable so the dropdown stays predictable across renders.
func (c Catalog) ContactInterestOptions() []ContactInterestOption {
	keys := []string{
		"contact.interest.option.evaluation-access",
		"contact.interest.option.commercial-deployment",
		"contact.interest.option.partnership",
		"contact.interest.option.technical-question",
		"contact.interest.option.other",
	}
	slugs := []string{
		"evaluation-access",
		"commercial-deployment",
		"partnership",
		"technical-question",
		"other",
	}
	options := make([]ContactInterestOption, 0, len(slugs))
	for i, slug := range slugs {
		options = append(options, ContactInterestOption{
			Slug:  slug,
			Label: c.T(keys[i]),
		})
	}
	return options
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
		"privacy": {
			Title:       "Privacy Notice | Realtek Connect+",
			Description: "Review how Realtek Connect+ handles website inquiries, contact form data, retention, data requests, and local video behavior.",
		},
	}
}

func enText() map[string]string {
	return map[string]string{
		"skip.main":                                 "Skip to main content",
		"brand.home":                                "Realtek Connect+ home",
		"nav.docs":                                  "Docs",
		"nav.features":                              "Features",
		"nav.architecture":                          "Architecture",
		"nav.contact":                               "Contact",
		"footer.tagline":                            "IoT cloud platform concept for Realtek-based connected products.",
		"footer.docs":                               "Developer Docs",
		"footer.features":                           "Features",
		"footer.contact":                            "Contact Us",
		"footer.privacy":                            "Privacy",
		"home.eyebrow":                              "IoT cloud platform for product teams",
		"home.lede":                                 "Bring Realtek-based devices online with provisioning, OTA, fleet operations, app SDKs, insights, private cloud options, and ecosystem integrations.",
		"home.cta.primary":                          "Talk to sales",
		"home.cta.secondary":                        "See plans & limits",
		"home.chip.silicon":                         "Silicon",
		"home.chip.sdk":                             "Device SDK",
		"home.chip.cloud":                           "Cloud",
		"home.chip.ops":                             "Ops",
		"home.overview.eyebrow":                     "Platform overview",
		"home.overview.title":                       "From chipset to connected product lifecycle.",
		"home.overview.lede":                        "Connect+ frames Realtek connectivity silicon, firmware enablement, cloud services, mobile app workflows, and operations tooling as one commercial IoT product path.",
		"home.surfaces.eyebrow":                     "Platform surfaces",
		"home.surfaces.title":                       "Show the product system, not only the feature list.",
		"home.surfaces.card.title":                  "Provisioning, OTA, and fleet health share one visual language.",
		"home.surfaces.card.body":                   "The public site now anchors the cloud story with product surfaces that connect devices, secure cloud operations, and dashboard workflows.",
		"home.surface.onboarding":                   "Onboarding",
		"home.surface.rollouts":                     "Rollouts",
		"home.surface.insights":                     "Insights",
		"home.surface.security":                     "Security",
		"home.principles.eyebrow":                   "Design principles",
		"home.principles.title":                     "Built for enterprise connected-product programs.",
		"home.principle.active":                     "Manageability",
		"home.principle.scale":                      "Scalability",
		"home.principle.security":                   "Security",
		"home.principle.privacy":                    "Privacy",
		"home.principle.cost":                       "Cost control",
		"home.principle.custom":                     "Customizability",
		"home.principle.panel":                      "Operate firmware, users, fleets, and support workflows from one platform story.",
		"home.principle.body":                       "Realtek Connect+ is presented as a lifecycle system: onboarding, cloud identity, OTA, app SDKs, metrics, and enterprise deployment work together instead of living as disconnected product pages.",
		"home.services.eyebrow":                     "Core services",
		"home.services.title":                       "Realtek cloud capabilities, packaged for connected-product teams.",
		"home.feature.details":                      "View details",
		"home.arch.eyebrow":                         "Architecture",
		"home.arch.title":                           "A direct flow from device onboarding to fleet operations.",
		"home.arch.device.title":                    "Realtek Device SDK",
		"home.arch.device.body":                     "Identity, provisioning, firmware services, and device signals.",
		"home.arch.cloud.title":                     "Secure Cloud",
		"home.arch.cloud.body":                      "Registry, OTA campaigns, user-device association, APIs.",
		"home.arch.app.title":                       "App SDK + Dashboard",
		"home.arch.app.body":                        "Mobile onboarding, product control, insights, support workflows.",
		"home.film.eyebrow":                         "Brand foundation",
		"home.film.title":                           "Built on Realtek's connected intelligence.",
		"home.film.body":                            "Realtek Connect+ extends a semiconductor and connectivity foundation into a cloud platform story for product teams building connected devices at commercial scale.",
		"home.film.cta":                             "Watch brand film",
		"home.film.title.attr":                      "Realtek corporate brand film",
		"home.film.fallback":                        "Your browser does not support the video tag.",
		"home.film.point.silicon":                   "Semiconductor foundation",
		"home.film.point.ecosystem":                 "Connected-product ecosystem",
		"home.film.point.enterprise":                "Enterprise deployment trust",
		"home.deploy.eyebrow":                       "Public vs private cloud",
		"home.deploy.title":                         "Start with evaluation, deploy anywhere — any cloud or on-premises.",
		"home.deploy.public":                        "Public evaluation",
		"home.deploy.public.title":                  "Validate product fit before committing to a private footprint.",
		"home.deploy.public.body":                   "Free evaluation tier — start with 5 devices by default, raise up to 200 on request. No expiry. Use it to align firmware, mobile, cloud, and product stakeholders before commercial deployment planning.",
		"home.deploy.docs":                          "Deployment Docs",
		"home.deploy.private":                       "Private commercial cloud",
		"home.deploy.private.title":                 "Run on GCP, Azure, AWS, or your own data center. No cloud lock-in.",
		"home.deploy.private.body":                  "Realtek Connect+ private deployment runs as a standard container or VM workload. No serverless runtime required, no dependency on a specific cloud provider — bring your own infrastructure.",
		"home.deploy.discuss":                       "Discuss Private Cloud",
		"home.use.eyebrow":                          "Use cases",
		"home.use.title":                            "Built for commercial connected-device teams.",
		"home.use.smart.body":                       "App onboarding, sharing, push notifications, voice assistant paths, and OTA maintenance.",
		"home.use.industrial.body":                  "Fleet grouping, metadata, secure updates, private cloud, and operations visibility.",
		"home.use.appliance.body":                   "Long lifecycle firmware maintenance, activation data, branded apps, and support diagnostics.",
		"home.docs.eyebrow":                         "Developer portal",
		"home.docs.title":                           "Give product, firmware, app, and cloud teams one documentation spine.",
		"home.docs.open":                            "Open section",
		"home.cta.eyebrow":                          "Early access",
		"home.cta.title":                            "Plan a Realtek Connect+ product path.",
		"home.cta.body":                             "Register your interest in provisioning, OTA, private cloud, app SDK, or fleet operations.",
		"features.eyebrow":                          "Features",
		"features.title":                            "Cloud services for the full IoT product lifecycle.",
		"features.body":                             "Realtek Connect+ presents the core capabilities product teams expect from a commercial connected-device platform.",
		"features.open":                             "Open feature",
		"feature.discuss.prefix":                    "Discuss",
		"feature.all":                               "All Features",
		"feature.highlights":                        "Highlights",
		"feature.highlights.title":                  "What this service covers",
		"feature.capabilities":                      "Capabilities",
		"feature.capabilities.title":                "Platform building blocks",
		"feature.outcomes":                          "Outcomes",
		"feature.outcomes.title":                    "Why product teams use it",
		"feature.next":                              "Next step",
		"feature.cta.prefix":                        "Evaluate",
		"feature.cta.suffix":                        "for your product roadmap.",
		"feature.cta.body":                          "Share your product category, target deployment, and cloud requirements with the Realtek Connect+ team.",
		"docs.eyebrow":                              "Developer docs",
		"docs.title":                                "Documentation entry points for cloud, firmware, app, and deployment teams.",
		"docs.body":                                 "Realtek Connect+ now includes a server-rendered documentation portal structure covering platform overview, development, APIs, SDKs, firmware, CLI workflows, deployment, and release notes.",
		"docs.cta.primary":                          "Talk to the platform team",
		"docs.cta.secondary":                        "See app platform context",
		"docs.portal.eyebrow":                       "Portal structure",
		"docs.portal.title":                         "Choose the track that matches your workstream.",
		"docs.why.eyebrow":                          "Why this matters",
		"docs.why.title":                            "Mirror the documentation surfaces teams expect before product launch.",
		"docs.shared.title":                         "Shared platform narrative",
		"docs.shared.body":                          "Give product, sales, and engineering stakeholders one place to understand lifecycle scope before implementation deep-dives.",
		"docs.depth.title":                          "Workstream-specific depth",
		"docs.depth.body":                           "Separate firmware, API, mobile SDK, deployment, and release concerns so each team can navigate directly to its implementation surface.",
		"docs.static.title":                         "Static first version",
		"docs.static.body":                          "Keep the docs portal compatible with the current Go templates and server-rendered architecture while leaving space for deeper follow-on content.",
		"doc.back":                                  "Back to docs",
		"doc.discuss":                               "Discuss implementation",
		"doc.coverage":                              "Coverage",
		"doc.coverage.title":                        "What this section explains",
		"doc.outputs":                               "Expected outputs",
		"doc.outputs.title":                         "What teams should be able to leave with",
		"doc.audience":                              "Primary audience",
		"doc.audience.title":                        "Who should start here",
		"doc.next":                                  "Next sections",
		"doc.next.title":                            "Continue through the developer portal.",
		"doc.view":                                  "View section",
		"contact.eyebrow":                           "Contact",
		"contact.title":                             "Register interest in Realtek Connect+.",
		"contact.body":                              "Tell us which service matters most for your product team. Requests are stored locally in SQLite for this first version.",
		"contact.context.eyebrow":                   "Early access",
		"contact.context.title":                     "For IoT product planning, firmware maintenance, and private cloud evaluation.",
		"contact.context.body":                      "Use this form for provisioning, OTA, fleet management, app SDK, insights, private cloud, or integration discussions.",
		"contact.thanks":                            "Thanks",
		"contact.recorded":                          "Your Realtek Connect+ request has been recorded.",
		"contact.review":                            "Review features",
		"contact.error.summary":                     "Please review the details below before submitting.",
		"contact.website":                           "Website",
		"contact.name":                              "Name",
		"contact.company":                           "Company",
		"contact.email":                             "Email",
		"contact.interest":                          "Interest",
		"contact.select":                            "Select an inquiry type",
		"contact.interest.option.evaluation-access": "Evaluation access (free tier signup or quota raise)",
		"contact.interest.option.commercial-deployment": "Commercial deployment (private cloud, pricing, contract)",
		"contact.interest.option.partnership":           "Partnership or business development",
		"contact.interest.option.technical-question":    "Technical question (existing customer or evaluator)",
		"contact.interest.option.other":                 "Other",
		"contact.message":                               "Message",
		"contact.submit":                                "Submit Request",
		"contact.privacy":                               "By submitting this form, you understand that your inquiry will be handled according to the Realtek Connect+ privacy notice.",
		"contact.privacy.link":                          "Privacy notice",
		"privacy.eyebrow":                               "Privacy",
		"privacy.title":                                 "Privacy notice for Realtek Connect+ website inquiries.",
		"privacy.intro":                                 "This first-version website collects only the information needed to respond to Realtek Connect+ business inquiries and early access requests. When ANALYTICS_ENABLED is true, it also uses first-party SQLite analytics for basic site measurement.",
		"privacy.data.title":                            "Data we collect",
		"privacy.data.body":                             "The contact form can collect name, company, email, area of interest, and an optional message. Analytics event rows are limited to page_view, click_cta, scroll, and engaged events, canonical page or CTA keys, coarse scroll percent or duration when applicable, referrer origin only, and a short-lived ephemeral session id. Analytics does not collect names, email addresses, company names, form messages, raw IP addresses, full referrer URLs, or account identifiers.",
		"privacy.use.title":                             "How we use the data",
		"privacy.use.body":                              "We use submitted information to respond to inquiries, plan product discussions, understand interest in Realtek Connect+ services, and protect the website from spam or abuse. First-party analytics is used only to understand coarse traffic source, page engagement, CTA conversion, and content effectiveness for this website.",
		"privacy.retention.title":                       "Retention",
		"privacy.retention.body":                        "Website leads are intended to be retained for up to 24 months unless a longer period is needed for an active business discussion or required operational record. Raw analytics event rows are retained for 90 days by default through ANALYTICS_RETENTION_DAYS.",
		"privacy.rights.title":                          "Access, correction, or deletion requests",
		"privacy.rights.body":                           "To request access, correction, or deletion of submitted inquiry data, or to ask about first-party analytics, contact privacy@example.com. Replace this placeholder address with the official privacy contact before public launch.",
		"privacy.video.title":                           "Local brand video",
		"privacy.video.body":                            "The homepage brand film is hosted by this website as a local MP4 asset. The video player does not create a YouTube iframe or contact YouTube services. The site does not load third-party analytics services, advertising pixels, or fingerprinting scripts.",
		"privacy.admin.title":                           "Internal access",
		"privacy.admin.body":                            "Lead review is protected by an admin token. Admin pages are excluded from the sitemap and are marked noindex.",
		"privacy.legal.title":                           "Legal review",
		"privacy.legal.body":                            "This notice is a GDPR-lite implementation for the website prototype. It is not a complete legal compliance package and should be reviewed before public launch.",
	}
}
