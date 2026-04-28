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
			Summary:      "iOS and Android SDK modules, sample app baselines, push notifications, rebrand guidance, and app publishing paths for connected products.",
			Description:  "App SDK now frames the mobile experience as a launch surface product teams can brand, extend, and publish without rebuilding every connected-app primitive from scratch. It covers iOS and Android SDK layers, sample app structure, push workflows, and release planning while staying explicit that this repo is a server-rendered website, not a shipped mobile framework.",
			Highlights:   []string{"iOS and Android SDK coverage for onboarding, control, and account flows", "Sample app and rebrand paths for faster branded launches", "Push notification, release-readiness, and app publishing guidance"},
			Capabilities: []string{"Shared mobile primitives for sign-in, provisioning, device control, and sharing", "Reference app structure that can be configured, re-skinned, or extended by product teams", "Release planning for App Store and Google Play submission, rollout, and support operations"},
			Outcomes:     []string{"Shorten mobile launch timelines", "Keep app developers and product teams aligned on ownership", "Avoid one-off app/cloud integration work for every product line"},
			Sections: []FeatureSection{
				{
					Eyebrow: "SDK Foundations",
					Title:   "Deliver branded mobile apps without rebuilding the connected product stack",
					Intro:   "The page now positions mobile SDK work as a reusable platform layer instead of a vague app-ready claim.",
					Items: []string{
						"Cover shared onboarding, authentication, device control, and account-linking primitives through iOS and Android SDK layers instead of promising a full client framework in this repo.",
						"Describe how mobile teams can map common device models, provisioning states, and control surfaces into product-specific app experiences with less custom cloud glue.",
						"Keep the website explicit that the public Go application is describing app enablement scope rather than serving as the runtime for a native mobile client.",
					},
				},
				{
					Eyebrow: "Sample App",
					Title:   "Start from a reference experience, then adapt it to the product",
					Intro:   "A reference app story helps app developers and product teams understand what is reused versus what is branded.",
					Items: []string{
						"Use a sample app to accelerate white-label or branded launches while preserving room for custom navigation, design systems, and product-specific device flows.",
						"Separate shared account, onboarding, and device-management scaffolding from customer-owned copy, visual identity, and feature-priority decisions.",
						"Support phased delivery where teams start from a proven baseline, then replace screens, modules, and integrations as their roadmap matures.",
					},
				},
				{
					Eyebrow: "Notifications",
					Title:   "Treat alerts and lifecycle messaging as part of the app product surface",
					Intro:   "Push and in-app notification flows are described as a coordinated mobile, cloud, and support capability.",
					Items: []string{
						"Plan push notifications around onboarding completion, shared-device events, OTA prompts, alerts, and support workflows that need deep links back into the branded app.",
						"Connect notification payloads to user permissions, device ownership state, and customer support escalation paths so mobile UX stays consistent.",
						"Show how product teams can align notification tone, branding, and preference controls with their own market and compliance requirements.",
					},
				},
				{
					Eyebrow: "Publishing",
					Title:   "Coordinate store launch work across engineering and product teams",
					Intro:   "Publishing guidance keeps the App SDK page tied to real launch execution instead of stopping at SDK selection.",
					Items: []string{
						"Coordinate bundle identifiers, signing assets, store metadata, review checklists, and staged rollout plans for both the App Store and Google Play.",
						"Use the page's contact path as a CTA for app developers and product teams who need to align branding, release readiness, and backend capability scope.",
						"Keep store ownership, privacy disclosures, crash monitoring, and release approvals assigned to the product team instead of implying the website manages mobile operations.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Delivery Paths",
				Title:   "Choose the mobile delivery path that fits launch speed and brand control",
				Intro:   "Each path uses the same Realtek Connect+ app foundations, but ownership shifts between starter assets, product-specific UX, and release operations.",
				Columns: []string{"Delivery path", "What ships first", "Ownership boundary", "Best fit"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Reference sample app", "Adopt the shared sample flow with minimal configuration to validate onboarding, control, and cloud connectivity quickly.", "Realtek starter patterns stay prominent while the product team owns signing, store presence, and customer support messaging.", "Workshops, pilot programs, and fast proof-of-concept mobile validation."}},
					{Cells: []string{"Rebranded starter app", "Replace visual identity, copy, device catalog, and support flows on top of the proven starter structure.", "Core app architecture remains shared, while the product team owns brand expression, release cadence, and app-market positioning.", "Teams that need a branded launch quickly without rewriting common connected-app workflows."}},
					{Cells: []string{"Custom app on shared SDK", "Build a product-specific mobile experience while reusing shared SDK modules for identity, provisioning, telemetry, and device control.", "The product team owns UX architecture and roadmap differentiation, while shared SDK contracts define the cloud and device touchpoints.", "Long-lived product lines that need deeper customization, ecosystem hooks, or multi-brand app strategies."}},
				},
			},
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
