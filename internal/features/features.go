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
			Summary:      "Upload firmware, extract release metadata, target staged rollouts, and manage dynamic OTA jobs with force, normal, and user-controlled policies.",
			Description:  "OTA is positioned as a production firmware operations surface rather than a simple update button. Teams register firmware packages, review extracted version and model metadata, define rollout policies, and watch job progress from pilot cohorts through archive-ready release history.",
			Highlights:   []string{"Firmware upload pipeline with extracted version, model, and release metadata", "Version, model, region, and cohort targeting for staged campaigns", "Force, normal, scheduled, user-controlled, and time-window rollout modes"},
			Capabilities: []string{"Dynamic OTA policies for always-on or intermittently connected fleets", "Per-job status, device outcomes, cancellation, and archive history", "Compatibility validation and operator approvals before rollout"},
			Outcomes:     []string{"Lower firmware support cost", "Reduce fleet-wide regression risk", "Coordinate releases across consumer and commercial deployments"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Workflow",
					Title:   "From signed image upload to job execution",
					Intro:   "The OTA page now describes the full release workflow product teams expect before a binary reaches devices.",
					Items: []string{
						"Upload signed firmware images and extract embedded project, version, model, checksum, and release-note metadata.",
						"Attach rollout notes, force, normal, or user-controlled install policy, and maintenance-window guidance before approval.",
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
					{Cells: []string{"Force", "Start immediately and require installation once the device reaches the approved state.", "Urgent hotfixes across a narrow pilot or the whole eligible fleet.", "Critical security patches or rapid rollback replacements."}},
					{Cells: []string{"Normal", "Make the update available with standard retry and install behavior while preserving operator oversight.", "Broad staged rollouts where devices should update promptly without hard-forcing the install moment.", "Routine firmware releases that still need rollout telemetry and cancellation controls."}},
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
			Slug:         "user-management",
			Title:        "User Management",
			Kicker:       "Handle the account lifecycle around connected products.",
			Summary:      "Platform content for sign up, sign in, OTP verification, social login, password recovery, account changes, and account deletion.",
			Description:  "User Management describes the account lifecycle capabilities product teams usually need around a Realtek-based connected product. It covers identity onboarding, recovery, and privacy operations for future product apps and services. This website does not expose end-user sign-in or account management flows today.",
			Highlights:   []string{"Self-service sign up and sign in journeys for branded mobile apps", "One-time password verification for account activation, recovery, and high-risk actions", "Third-party login and account-linking paths for partner or consumer ecosystems"},
			Capabilities: []string{"Forgot-password, change-password, and session-management controls", "Account deletion and retention workflows that hand off cleanly to support and compliance teams", "User profile, consent, and device-ownership state that stays separate from this marketing website"},
			Outcomes:     []string{"Shorten time to a production-ready account system", "Keep user lifecycle scope explicit during architecture reviews", "Avoid confusing product platform capabilities with the website's own lead-capture flows"},
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
			Summary:      "Compare public evaluation with private commercial deployment, regional hosting, custom domains, and enterprise upgrade planning.",
			Description:  "Private Cloud explains how Realtek Connect+ moves from a shared evaluation story into dedicated commercial deployment. It positions data ownership, regional placement, custom domain control, and upgrade planning as enterprise buying criteria rather than unsupported promises about this website's own runtime.",
			Highlights:   []string{"Public evaluation versus dedicated private commercial deployment", "Data ownership, regional hosting boundaries, and custom domain control", "Commercial onboarding, upgrade path, and deployment support expectations"},
			Capabilities: []string{"Dedicated environment planning for customer-operated or managed private regions", "Reverse-proxy TLS termination, network policy alignment, and branded service endpoints", "Release promotion and maintenance-window planning across evaluation and production environments"},
			Outcomes:     []string{"Match enterprise procurement requirements", "Keep ownership boundaries explicit", "Create a credible path from pilot to production"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Commercial Models",
					Title:   "Start with evaluation, then move into owned deployment boundaries",
					Intro:   "The page frames public evaluation as a fast proof-of-concept path and private deployment as the commercial operating model for products with stricter ownership and compliance requirements.",
					Items: []string{
						"Use a shared evaluation environment to validate device flows, dashboards, and integration assumptions before commercial rollout.",
						"Transition to a dedicated deployment once product teams need tenant isolation, formal support processes, and customer-specific change windows.",
						"Keep the website explicit that these are platform deployment models, not evidence that this repo already ships a full private cloud control plane.",
					},
				},
				{
					Eyebrow: "Ownership",
					Title:   "Define where data lives and how the service is branded",
					Intro:   "Private deployment content is grounded in enterprise concerns about who operates the stack and where customer traffic terminates.",
					Items: []string{
						"Document customer-owned data boundaries for device metadata, operator access, and retained support exports.",
						"Offer custom domains and branded entry points so the deployment can align with the customer's DNS, certificate, and support model.",
						"Choose regional placement around residency, latency, and operational coverage requirements instead of forcing every product through one public region.",
					},
					Accent: true,
				},
				{
					Eyebrow: "Upgrade Path",
					Title:   "Promote proven configurations into commercial production",
					Intro:   "The upgrade path is described as an engineering and operations workflow rather than a one-click migration promise.",
					Items: []string{
						"Carry validated device models, app configuration, and integration settings from evaluation into a dedicated deployment plan.",
						"Use release promotion, maintenance windows, and rollback checkpoints to move from pilot tenants into production operations safely.",
						"Align the commercial cutover with customer security review, support readiness, and staged onboarding of real fleets.",
					},
				},
				{
					Eyebrow: "Deployment FAQ",
					Title:   "Answer the questions enterprise buyers raise first",
					Intro:   "FAQ-style guidance keeps the page practical without implying unsupported hosting features inside the website itself.",
					Items: []string{
						"Production TLS still terminates at a reverse proxy, ingress, or deployment platform in front of the Go website runtime.",
						"Private regions can follow customer-approved network boundaries and upgrade calendars instead of a shared public release schedule.",
						"Commercial support covers deployment planning, environment hardening expectations, and the path for future platform customization requests.",
					},
				},
			},
			Table: FeatureTable{
				Eyebrow: "Deployment Paths",
				Title:   "Compare evaluation and private commercial operating models",
				Intro:   "Realtek Connect+ positions deployment choice as a commercial decision: shared evaluation speeds discovery, while private environments add ownership, branding, and regional control for production programs.",
				Columns: []string{"Model", "Operations boundary", "What teams get", "Best fit"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Public evaluation", "Shared environment for workshops, pilot demos, and early integration discovery.", "Fast access to core platform flows without committing to a customer-specific operating boundary.", "Early evaluations, internal validation, and short proof-of-concept cycles."}},
					{Cells: []string{"Managed private deployment", "Dedicated commercial environment operated with agreed support windows and customer-specific policies.", "Tenant isolation, custom domain options, regional placement choices, and a clearer production support model.", "Teams that want private deployment outcomes without owning the day-to-day platform operations stack."}},
					{Cells: []string{"Customer-operated private region", "Customer-selected infrastructure and network boundary with coordinated release and upgrade planning.", "Maximum control over residency, access policies, and environment-level change management.", "Products with strict enterprise governance, regulated data boundaries, or regional hosting mandates."}},
				},
			},
		},
		{
			Slug:         "integrations",
			Title:        "Integrations",
			Kicker:       "Connect products to the wider IoT ecosystem.",
			Summary:      "Matter Fabric positioning, voice assistant paths, MQTT over TLS, REST APIs, and webhooks for product and platform integrations.",
			Description:  "Integrations explains how Realtek Connect+ fits into smart home ecosystems and enterprise backends. It frames Matter Fabric participation, voice assistant connections, secure protocol access, and webhook delivery as supported integration patterns without claiming this website already operates every downstream service.",
			Highlights:   []string{"Matter ecosystem positioning with Fabric-aware deployment planning", "Voice assistant, REST API, MQTT over TLS, and webhook integration paths", "Explicit ownership boundaries between product apps, cloud services, and customer systems"},
			Capabilities: []string{"Matter bridge/controller planning, commissioning touchpoints, and ecosystem mapping", "Secure REST and MQTT interfaces for product, support, and operations systems", "Webhook-driven event handoff into CRM, ticketing, and analytics workflows"},
			Outcomes:     []string{"Meet ecosystem interoperability expectations", "Connect business systems without custom one-off glue", "Keep integration scope credible for platform evaluations"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Matter Fabric",
					Title:   "Position Realtek products inside the customer's chosen ecosystem",
					Intro:   "The integrations page now explains Matter as an ecosystem contract with clear boundaries around who owns commissioning, controller roles, and long-term lifecycle UX.",
					Items: []string{
						"Describe how devices can participate in a Matter Fabric while still keeping Realtek app, cloud, and support flows in scope where products need them.",
						"Frame bridge and controller roles as deployment decisions so teams can map product categories to the right ecosystem entry point without overclaiming current implementation depth.",
						"Keep commissioning, credential ownership, and ecosystem-specific UX tied to the selected Matter platform instead of implying this marketing site is a live Matter control plane.",
					},
				},
				{
					Eyebrow: "Voice Assistants",
					Title:   "Connect branded products to familiar consumer control surfaces",
					Intro:   "Assistant integrations are positioned as cloud-to-cloud or skill/action patterns that extend a product's reach without replacing its own app and identity model.",
					Items: []string{
						"Cover Alexa and Google Assistant paths for products that need voice control, routine support, and ecosystem discovery alongside the branded app experience.",
						"Explain that assistant platforms own the voice UX while product teams keep device traits, account linking, and support escalation flows aligned with their own roadmap.",
						"Use the page to show interoperability intent for smart home buyers without copying ESP RainMaker wording or promising live assistant certification from this repo alone.",
					},
				},
				{
					Eyebrow: "Protocols",
					Title:   "Offer direct data and control paths for external systems",
					Intro:   "Business integrations are described as secure interfaces teams can expose around the platform rather than as ad hoc export buttons.",
					Items: []string{
						"Document REST APIs for authenticated product, support, and operations workflows that need request-response access to platform state.",
						"Position MQTT over TLS for policy-scoped telemetry, near-real-time command paths, and event fan-out into downstream infrastructure.",
						"Use webhooks for signed lifecycle, alert, and workflow events so CRM, ticketing, analytics, and fulfillment systems can react without polling.",
					},
				},
				{
					Eyebrow: "Ownership Boundaries",
					Title:   "Keep trust, credentials, and responsibilities explicit",
					Intro:   "The page stays careful about what Realtek Connect+ is describing versus what this repository actually implements today.",
					Items: []string{
						"Make authentication, topic policy, and receiving-system ownership part of the integration story so buyers can evaluate security posture early.",
						"Separate product-facing APIs and events from the website's own contact/admin runtime to avoid implying that marketing pages are the production integration surface.",
						"Preserve room for customer-specific deployment choices, event contracts, and certification work instead of claiming universal out-of-the-box availability.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Integration Paths",
				Title:   "Choose the ecosystem contract that fits the product",
				Intro:   "Each path represents a different trust boundary. Realtek Connect+ uses this page to show design intent and deployment options without claiming every protocol surface is already live in this repository.",
				Columns: []string{"Path", "Interaction model", "Ownership boundary", "Best fit"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Matter Fabric", "Commission devices into a customer-selected Matter fabric while preserving Realtek app and cloud workflows where needed.", "Fabric credentials, controller role, and household UX stay aligned with the chosen Matter ecosystem.", "Products that need standards-based smart home interoperability across ecosystems."}},
					{Cells: []string{"Voice assistants", "Use cloud-to-cloud or skill/action integrations to expose device traits to Alexa or Google Assistant.", "Assistant platforms own the voice UX while product teams keep device identity, support flows, and roadmap control.", "Consumer products that need branded apps plus familiar voice control."}},
					{Cells: []string{"REST APIs", "Expose authenticated HTTPS endpoints for device, user, and operations workflows in external systems.", "API contracts, auth policy, and lifecycle governance stay under the product team's platform boundary.", "Partner portals, support tooling, ERP/CRM integration, and enterprise admin workflows."}},
					{Cells: []string{"MQTT over TLS", "Publish telemetry and command streams over policy-scoped topics with TLS-protected client authentication.", "Broker policy, topic ownership, and retention rules stay inside the selected platform deployment model.", "Operational telemetry, event streaming, and near-real-time orchestration."}},
					{Cells: []string{"Webhooks", "Push signed lifecycle and alert events into downstream SaaS or internal automation.", "Receiving systems own the follow-up workflow while Connect+ owns the event contract and delivery path.", "Ticketing, notifications, analytics ingestion, and fulfillment triggers."}},
				},
			},
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
