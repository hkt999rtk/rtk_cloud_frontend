package features

type Feature struct {
	Slug         string
	Title        string
	Kicker       string
	Summary      string
	Description  string
	Highlights   []string
	Capabilities []string
	Outcomes     []string
	Sections     []FeatureSection
	Table        FeatureTable
}

type FeatureSection struct {
	Eyebrow string
	Title   string
	Intro   string
	Items   []string
	Accent  bool
}

type FeatureTable struct {
	Eyebrow string
	Title   string
	Intro   string
	Columns []string
	Rows    []FeatureTableRow
}

type FeatureTableRow struct {
	Cells []string
}

func All() []Feature {
	return []Feature{
		{
			Slug:         "provision",
			Title:        "Provision",
			Kicker:       "Onboard devices with less product friction.",
			Summary:      "Secure Wi-Fi/BLE onboarding, activation, and account binding for Realtek-based IoT products.",
			Description:  "Provision gives product teams a repeatable path from factory-ready hardware to a user-owned connected device. It covers first-time activation, local onboarding, cloud registration, and user-device association.",
			Highlights:   []string{"Wi-Fi and BLE onboarding flows", "Device binding and ownership transfer", "Activation state and first-run telemetry"},
			Capabilities: []string{"Claiming tokens and device identity handoff", "User-device association during app onboarding", "Timezone and metadata initialization"},
			Outcomes:     []string{"Reduce setup failures", "Shorten app onboarding", "Prepare devices for fleet operations"},
		},
		{
			Slug:         "ota",
			Title:        "OTA",
			Kicker:       "Ship firmware updates with rollout control.",
			Summary:      "Upload firmware, extract release metadata, target staged rollouts, and monitor dynamic OTA jobs with policy-level safeguards.",
			Description:  "OTA is positioned as a production firmware operations surface rather than a simple update button. Teams register firmware packages, review extracted version and model metadata, define rollout policies, and watch job progress from pilot cohorts through archive-ready release history.",
			Highlights:   []string{"Firmware upload pipeline with extracted version, model, and release metadata", "Version, model, region, and cohort targeting for staged campaigns", "Immediate, scheduled, user-controlled, and time-window rollout modes"},
			Capabilities: []string{"Dynamic OTA policies for always-on or intermittently connected fleets", "Per-job status, device outcomes, cancellation, and archive history", "Compatibility validation and operator approvals before rollout"},
			Outcomes:     []string{"Lower firmware support cost", "Reduce fleet-wide regression risk", "Coordinate releases across consumer and commercial deployments"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Workflow",
					Title:   "From signed image upload to job execution",
					Intro:   "The OTA page now describes the full release workflow product teams expect before a binary reaches devices.",
					Items: []string{
						"Upload signed firmware images and extract embedded project, version, model, checksum, and release-note metadata.",
						"Attach rollout notes, mandatory or optional install policy, and maintenance-window guidance before approval.",
						"Create job detail views that show pending, in-progress, succeeded, cancelled, and failed device outcomes by campaign wave.",
					},
				},
				{
					Eyebrow: "Targeting",
					Title:   "Match each rollout to the right fleet slice",
					Intro:   "Campaign definition is described as an operator workflow rather than a live cloud implementation promise.",
					Items: []string{
						"Target by product family, hardware model, current firmware version, customer tier, region, or support cohort.",
						"Blend pilot cohorts with staged expansion so operators can start narrow, inspect telemetry, and widen safely.",
						"Use dynamic OTA rules to keep offline devices eligible for the latest approved package when they reconnect.",
					},
				},
				{
					Eyebrow: "Controls",
					Title:   "Safety checks and release operations",
					Intro:   "Operational controls stay explicit about website scope while still showing a credible release-management story.",
					Items: []string{
						"Validate project and version compatibility before devices accept a package.",
						"Cancel active waves and archive completed campaigns without losing audit history.",
						"Retain operator-visible status, retry intent, and exception handling notes for support and QA teams.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Rollout Strategies",
				Title:   "Choose the delivery mode that fits the release",
				Intro:   "Realtek Connect+ describes rollout modes as policy templates. Dynamic OTA keeps device eligibility aligned with the latest approved campaign even when endpoints reconnect later.",
				Columns: []string{"Strategy", "Operator control", "Targeting pattern", "Best fit"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Immediate", "Launch now and monitor wave status in real time.", "Urgent hotfixes across a narrow pilot or the whole eligible fleet.", "Critical security patches or rapid rollback replacements."}},
					{Cells: []string{"Scheduled", "Approve once, then release during a defined maintenance window.", "Region- or customer-specific batches timed for low-support hours.", "Planned feature releases and coordinated commercial deployments."}},
					{Cells: []string{"User-controlled", "Notify users, keep the job pending, and let the app or device owner choose install timing.", "Consumer devices that must respect end-user availability and local context.", "Appliances and smart-home products where UX matters more than immediacy."}},
					{Cells: []string{"Time-window", "Allow installs only inside approved hours while preserving campaign eligibility outside the window.", "Always-on fleets, retail estates, or shared environments with operational blackout periods.", "Commercial fleets that need policy-driven upgrades without overnight surprises."}},
				},
			},
		},
		{
			Slug:         "fleet-management",
			Title:        "Fleet Management",
			Kicker:       "Operate connected products after launch.",
			Summary:      "Device registry, groups, tags, batch operations, sharing, and lifecycle metadata for commercial fleets.",
			Description:  "Fleet Management brings device and node operations into one product surface. Teams can organize devices, attach searchable metadata, execute batch operations, and manage ownership or sharing models.",
			Highlights:   []string{"Device registry and lifecycle state", "Groups, tags, metadata, and timezone handling", "Sharing and batch operations"},
			Capabilities: []string{"Group by model, region, firmware, or customer", "Bulk metadata and service changes", "User-node relationship visibility"},
			Outcomes:     []string{"Keep fleets searchable", "Support commercial support workflows", "Scale operations beyond launch"},
		},
		{
			Slug:         "app-sdk",
			Title:        "App SDK",
			Kicker:       "Build branded mobile experiences faster.",
			Summary:      "iOS and Android SDK building blocks, sample app patterns, push notifications, and rebrand-ready app flows.",
			Description:  "App SDK gives product teams a mobile path for onboarding, control, sharing, alerts, and account/device flows. The website positions it as a bridge between Realtek device firmware and a customer-owned mobile app.",
			Highlights:   []string{"iOS and Android integration path", "Sample app flows for onboarding and control", "Push notifications and branded app readiness"},
			Capabilities: []string{"Login and user-device association screens", "Device control and service state updates", "Rebrand/customize and app store publishing path"},
			Outcomes:     []string{"Accelerate app delivery", "Keep brand ownership", "Reduce custom cloud integration work"},
		},
		{
			Slug:         "insights",
			Title:        "Insights",
			Kicker:       "See the health of products in the field.",
			Summary:      "Activation statistics, firmware distribution, crash reports, logs, reboot reasons, RSSI, and memory signals.",
			Description:  "Insights gives engineering and support teams a view into fleet quality. It highlights operational statistics and device health signals that help teams prioritize fixes and understand real deployment behavior.",
			Highlights:   []string{"Activation and association statistics", "Crash, reboot, and log visibility", "Firmware distribution and device health metrics"},
			Capabilities: []string{"RSSI, memory, uptime, and reboot reason signals", "Version adoption and rollout health", "Support-oriented device history"},
			Outcomes:     []string{"Find field issues earlier", "Support customers with evidence", "Measure firmware quality"},
		},
		{
			Slug:         "private-cloud",
			Title:        "Private Cloud",
			Kicker:       "Deploy with enterprise ownership and control.",
			Summary:      "Customer-owned deployment options for data control, custom domains, regional constraints, and commercial support.",
			Description:  "Private Cloud positions Realtek Connect+ for commercial products that need ownership boundaries, customization, and deployment control beyond public evaluation environments.",
			Highlights:   []string{"Enterprise-owned deployment model", "Custom domain and cloud customization", "Data ownership and commercial support"},
			Capabilities: []string{"Dedicated deployment planning", "Environment-specific policies", "Integration with customer operations"},
			Outcomes:     []string{"Meet enterprise requirements", "Control data boundaries", "Prepare for commercial scale"},
		},
		{
			Slug:         "integrations",
			Title:        "Integrations",
			Kicker:       "Connect products to the wider IoT ecosystem.",
			Summary:      "Alexa, Google Assistant, Matter, REST APIs, MQTT over TLS, and webhooks for product and platform integrations.",
			Description:  "Integrations extend device data and control beyond the core cloud. The page presents smart home assistant connections, secure protocol access, and API/webhook paths for business systems.",
			Highlights:   []string{"Voice assistant and Matter-ready positioning", "REST API and MQTT over TLS", "Webhooks for event-driven workflows"},
			Capabilities: []string{"Smart home ecosystem touchpoints", "External operations and CRM hooks", "Secure cloud-to-cloud integration"},
			Outcomes:     []string{"Meet ecosystem expectations", "Connect operations tools", "Support differentiated product experiences"},
		},
	}
}

func BySlug(slug string) (Feature, bool) {
	for _, feature := range All() {
		if feature.Slug == slug {
			return feature, true
		}
	}
	return Feature{}, false
}
