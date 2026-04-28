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
			Summary:      "Upload firmware, create rollout campaigns, monitor jobs, and protect devices with version validation.",
			Description:  "OTA is presented as a full firmware lifecycle service: binaries are registered, campaigns are targeted to devices or groups, jobs can be monitored or cancelled, and firmware-side checks help prevent incorrect images from being applied.",
			Highlights:   []string{"Immediate, scheduled, and user-controlled rollout modes", "Campaign status tracking and archive flow", "Project and version validation on device"},
			Capabilities: []string{"Firmware binary registration", "Targeting by group, model, or firmware version", "Job progress, cancellation, and history"},
			Outcomes:     []string{"Lower maintenance cost", "Safer staged rollouts", "Recover faster from firmware issues"},
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
