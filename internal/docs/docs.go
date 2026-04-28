package docs

type Section struct {
	Slug         string
	Title        string
	Icon         string
	Kicker       string
	Summary      string
	Description  string
	Highlights   []string
	Deliverables []string
	Audience     []string
}

func All() []Section {
	return []Section{
		{
			Slug:         "product-overview",
			Title:        "Product Overview",
			Icon:         "document",
			Kicker:       "Position the platform before diving into implementation.",
			Summary:      "Platform architecture, capability map, and commercial packaging guidance for Realtek Connect+ evaluations.",
			Description:  "Product Overview is the landing space for teams comparing platform scope, architecture boundaries, and rollout priorities. It frames how device firmware, cloud services, apps, operations, and enterprise deployment fit together without implying the website itself runs those systems.",
			Highlights:   []string{"Platform architecture narrative", "Capability map across onboarding, OTA, app, insights, and cloud", "Commercial evaluation framing for product teams"},
			Deliverables: []string{"Architecture diagrams and lifecycle summaries", "Capability comparison tables for internal alignment", "Evaluation guidance for hardware, mobile, and cloud stakeholders"},
			Audience:     []string{"Product managers evaluating commercial fit", "Solution architects mapping Realtek platform scope", "Sales engineers preparing technical discovery"},
		},
		{
			Slug:         "development",
			Title:        "Development",
			Icon:         "grid",
			Kicker:       "Organize firmware, cloud, and app workstreams around one delivery plan.",
			Summary:      "Environment setup, team responsibilities, and implementation tracks for device, cloud, mobile, and operations teams.",
			Description:  "Development outlines how engineering teams move from proof-of-concept into execution. It explains the expected workstreams, the handoffs between firmware and cloud teams, and the checkpoints needed before release and deployment readiness reviews.",
			Highlights:   []string{"Suggested phase plan from prototype to production", "Shared responsibilities across firmware, mobile, cloud, and QA", "Integration checkpoints for provisioning, OTA, and operations"},
			Deliverables: []string{"Environment bootstrap checklists", "Cross-team implementation milestones", "Validation gates before field pilots"},
			Audience:     []string{"Engineering managers planning delivery", "Firmware and app leads coordinating dependencies", "QA teams defining release readiness"},
		},
		{
			Slug:         "apis",
			Title:        "APIs",
			Icon:         "api",
			Kicker:       "Expose cloud capabilities through structured integration surfaces.",
			Summary:      "REST, MQTT over TLS, webhook, and service contract documentation entry points for external systems.",
			Description:  "APIs describes the contract layer around Realtek Connect+. It positions cloud APIs as the integration surface for dashboards, support tooling, business systems, and device event workflows, while keeping the first implementation static and server-rendered.",
			Highlights:   []string{"REST resource categories for devices, users, OTA, and analytics", "MQTT over TLS for device messaging and state updates", "Webhook patterns for operational event delivery"},
			Deliverables: []string{"Authentication and authorization model overview", "Endpoint families and payload expectations", "Integration examples for support and CRM workflows"},
			Audience:     []string{"Backend teams integrating cloud services", "Partner engineers building business system hooks", "Technical account teams answering API scope questions"},
		},
		{
			Slug:         "sdks",
			Title:        "SDKs",
			Icon:         "package",
			Kicker:       "Document the developer surfaces used to build connected product experiences.",
			Summary:      "Mobile SDK, firmware SDK, and reusable client components for branded product delivery.",
			Description:  "SDKs serves as the catalog of implementation building blocks. It connects mobile, firmware, and service-side integration stories so product teams can understand which surfaces are reused, customized, or wrapped for their own connected product launch.",
			Highlights:   []string{"iOS and Android mobile integration path", "Firmware-side service and identity building blocks", "Sample client components for branded product apps"},
			Deliverables: []string{"SDK selection guidance by product type", "Customization boundaries for branded experiences", "Support expectations for onboarding and lifecycle features"},
			Audience:     []string{"Mobile engineers integrating branded apps", "Embedded developers mapping device-side dependencies", "Program leads planning reuse vs customization"},
		},
		{
			Slug:         "firmware",
			Title:        "Firmware",
			Icon:         "device",
			Kicker:       "Clarify what the device software stack must provide.",
			Summary:      "Provisioning, identity, telemetry, OTA agent, and diagnostics expectations for device firmware teams.",
			Description:  "Firmware documents the device-side responsibilities required to make cloud features credible. It maps the expected lifecycle hooks for onboarding, signal reporting, OTA safety, and device health so teams can scope implementation effort before integration begins.",
			Highlights:   []string{"Provisioning and claim flow responsibilities", "OTA validation, targeting, and rollback safety concepts", "Health signal and diagnostics reporting expectations"},
			Deliverables: []string{"Firmware capability checklist", "Integration points for cloud identity and telemetry", "Release-readiness topics for factory and field updates"},
			Audience:     []string{"Embedded teams implementing device services", "System architects aligning firmware and cloud contracts", "Operations teams reviewing field support hooks"},
		},
		{
			Slug:         "cli",
			Title:        "CLI",
			Icon:         "terminal",
			Kicker:       "Support operators and developers with repeatable command-line workflows.",
			Summary:      "Command-line entry points for local testing, release preparation, fleet actions, and support diagnostics.",
			Description:  "CLI frames the operational workflows that often accompany a commercial IoT platform. It covers the kinds of scripted tasks operators and developers expect, from local environment setup to release packaging, diagnostics collection, and bulk maintenance actions.",
			Highlights:   []string{"Local developer setup and smoke commands", "Bulk operational workflows for support teams", "Release and environment maintenance scripts"},
			Deliverables: []string{"Command catalog for setup, test, release, and support", "Operator guardrails for high-impact actions", "Examples for scripted environment workflows"},
			Audience:     []string{"Developers running local test flows", "Support engineers handling fleet maintenance", "DevOps teams standardizing release automation"},
		},
		{
			Slug:         "deployment",
			Title:        "Deployment",
			Icon:         "cloud-lock",
			Kicker:       "Show how evaluation environments mature into commercial cloud footprints.",
			Summary:      "Public evaluation, container packaging, persistent SQLite storage, reverse proxy TLS, and operations ownership guidance.",
			Description:  "Deployment explains how teams move from a public evaluation story into controlled commercial environments. It now includes the packaging and storage assumptions for the current Go site so infrastructure teams can run the app with a persistent SQLite volume while keeping TLS and ingress at the platform edge.",
			Highlights:   []string{"Public evaluation versus private commercial deployment", "Container packaging with runtime templates, static assets, and mounted SQLite storage", "TLS and ingress handled by deployment environment or reverse proxy"},
			Deliverables: []string{"Deployment topology comparison", "Infrastructure ownership matrix", "Operational readiness FAQ covering persistence, ingress, and native builds"},
			Audience:     []string{"Platform teams planning hosted environments", "Enterprise buyers evaluating data boundaries", "DevOps engineers preparing deployment standards"},
		},
		{
			Slug:         "release-notes",
			Title:        "Release Notes",
			Icon:         "refresh",
			Kicker:       "Track product evolution across firmware, cloud, app, and ops surfaces.",
			Summary:      "Versioned change logs, upgrade notes, compatibility statements, and rollout communication patterns.",
			Description:  "Release Notes defines the documentation structure teams expect once the platform is shipping updates regularly. It sets up a home for product version changes, upgrade implications, and compatibility notes across firmware, cloud, mobile app, and operational tooling.",
			Highlights:   []string{"Version-by-version product change summaries", "Upgrade notes and compatibility callouts", "Customer-facing communication structure for releases"},
			Deliverables: []string{"Release note template for cloud, app, and firmware updates", "Upgrade impact sections by audience", "Archive strategy for historical product changes"},
			Audience:     []string{"Customer success teams preparing release comms", "Engineering teams documenting compatibility changes", "Product leadership tracking roadmap delivery"},
		},
	}
}

func BySlug(slug string) (Section, bool) {
	for _, section := range All() {
		if section.Slug == slug {
			return section, true
		}
	}
	return Section{}, false
}
