package features

type Feature struct {
	Slug         string
	Title        string
	Icon         string
	Kicker       string
	Summary      string
	Description  string
	ImagePath    string
	ImageAlt     string
	SourceLabel  string
	SourceURL    string
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
			Icon:         "provision",
			Kicker:       "Frame onboarding around the contract-backed foundation.",
			Summary:      "Cloud registry and activation foundations are contract-backed; local Wi-Fi/BLE onboarding, claim UX, transfer/reset policy, and product readiness remain integration or roadmap scope.",
			Description:  "Provision presents Realtek Connect+ onboarding as a layered path from registry entry to activated cloud device. The current foundation is the account registry, cross-service provisioning command flow, video-cloud activation boundary, scoped tokens, and transport readiness contract; local Wi-Fi/BLE setup, QR/SoftAP UX, ownership transfer, factory reset policy, and aggregate product readiness are not presented as generally available implementation until those owner repositories land them.",
			ImagePath:    "/static/assets/feature-provision-flow.jpg",
			ImageAlt:     "Provisioning dashboard concept with mobile pairing steps, QR onboarding, and device activation status cards.",
			SourceLabel:  "Product onboarding interface contract",
			SourceURL:    "https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/PRODUCT_ONBOARDING.md",
			Highlights:   []string{"Contract-backed cloud registry, activation, token, and transport readiness boundaries", "Integration-ready claim material concepts for QR, serial, activation code, and future factory identity", "Roadmap treatment for local Wi-Fi/BLE setup, SoftAP UX, transfer/reset policy, and aggregate readiness"},
			Capabilities: []string{"Use account-side registry APIs plus cross-service DeviceProvisionRequested and DeviceProvisionSucceeded or DeviceProvisionFailed events for the cloud activation foundation", "Keep SDK claim parsing separate from account-side ownership and binding decisions", "Describe local onboarding and product readiness as planned owner-repository work instead of broadly available website functionality"},
			Outcomes:     []string{"Preserve the product onboarding vision without overclaiming implementation status", "Give firmware, app, SDK, account, and video-cloud teams one shared availability vocabulary", "Make evaluation conversations explicit about what is available now, integration-ready, or roadmap"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Available foundation",
					Title:   "Cloud-side provisioning is the implemented contract boundary",
					Intro:   "The stable story starts with account registry records and the cross-service activation flow rather than a single all-in-one provisioning endpoint.",
					Items: []string{
						"Account-side device registration, cross-service provisioning requests, video activation results, scoped token issuance, and owner transport readiness are the public cloud-side behaviors to discuss today.",
						"Provisioning remains multi-service orchestration across account manager, the cross-service channel, video cloud, device credentials, and transport state.",
						"Video activation alone is not full product readiness; the product state must distinguish registry, claim, local setup, activation, online, failure, and deactivation stages.",
					},
				},
				{
					Eyebrow: "Integration-ready",
					Title:   "Claim material has a defined interface, not final ownership policy",
					Intro:   "The product onboarding contract gives SDK and app teams common vocabulary for claim input while leaving authorization decisions with account-side policy.",
					Items: []string{
						"Claim material may represent QR payloads, serial numbers, activation codes, MAC addresses, or future factory identity inputs.",
						"SDK parsers should normalize supported claim material and return stable errors for malformed or unsupported inputs.",
						"Account-side follow-up work still owns reuse rules, already-claimed rejection, transfer behavior, factory reset semantics, and delete-versus-deactivate policy.",
					},
				},
				{
					Eyebrow: "Roadmap",
					Title:   "Local onboarding remains owner-repository implementation work",
					Intro:   "The public page keeps the product vision visible while avoiding a general-availability claim for local setup UX.",
					Items: []string{
						"BLE provisioning, SoftAP provisioning, local Wi-Fi credential transport, QR onboarding UX, ECDH or challenge-response handshakes, and manufacturing CA policy are not yet stable website-available implementation claims.",
						"Android and iOS are the primary targets for real local onboarding implementations; native and JavaScript/TypeScript packages should report explicit unsupported capability where needed.",
						"Full product readiness should wait for local setup results, claim/bind policy, cloud activation, and transport online state to be joined by an owner repository or integration service.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Availability",
				Title:   "Separate what is available, integration-ready, and roadmap",
				Intro:   "Realtek Connect+ provisioning is presented as a layered product path. The cloud-side foundation is contract-backed now; local onboarding UX and final ownership policy stay clearly marked until implementation lands.",
				Columns: []string{"Layer", "Public status", "Customer-facing boundary"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Cloud registry and activation foundation", "Available foundation", "Account registry, cross-service provisioning commands and events, video activation, scoped credentials, and transport readiness define the current cloud-side path."}},
					{Cells: []string{"Claim material parsing", "Integration-ready", "QR, serial, activation code, MAC, and future factory identity inputs have shared interface vocabulary, but ownership policy is account-side follow-up work."}},
					{Cells: []string{"Local Wi-Fi/BLE and SoftAP onboarding", "Roadmap", "Local discovery, credential handoff, QR onboarding UX, and mobile setup sessions are not described as generally available implementation in this website."}},
					{Cells: []string{"Transfer, reset, and product readiness", "Roadmap", "Already-claimed handling, ownership transfer, factory reset, delete/deactivate separation, and aggregate readiness projection require follow-up policy and service work."}},
				},
			},
		},
		{
			Slug:         "ota",
			Title:        "OTA",
			Icon:         "ota",
			Kicker:       "Ship firmware updates with rollout control.",
			Summary:      "Upload firmware, extract release metadata, target staged rollouts, and manage dynamic OTA jobs with force, normal, and user-controlled policies.",
			Description:  "OTA is positioned as a production firmware operations surface rather than a simple update button. Teams register firmware packages, review extracted version and model metadata, define rollout policies, and watch job progress from pilot cohorts through archive-ready release history.",
			ImagePath:    "/static/assets/feature-ota-control-center.jpg",
			ImageAlt:     "Firmware rollout control center with staged release timeline, device cohorts, and OTA job analytics.",
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
			Icon:         "fleet",
			Kicker:       "Operate connected products after launch.",
			Summary:      "Node registration, certificate provisioning, device registry, OTA job coordination, batch operations, and operator widgets for commercial fleets.",
			Description:  "Fleet Management expands the public operations story beyond the website's sales-lead admin page. It describes how product teams register nodes, issue device identity material, organize the registry, coordinate firmware operations, and review fleet-wide operator widgets in a future IoT platform console.",
			ImagePath:    "/static/assets/feature-fleet-management.png",
			ImageAlt:     "Fleet management dashboard with connected device groups, health status tiles, tags, and batch operation queue.",
			Highlights:   []string{"Node registration, certificate bootstrapping, and registry lifecycle state", "Groups, tags, metadata, sharing, and batch operator actions", "OTA job visibility, firmware image tracking, and fleet health widgets"},
			Capabilities: []string{"Register nodes, bind manufacturing records, and rotate or revoke device certificates as products move through activation and support workflows", "Search the device registry by model, region, firmware, customer, ownership, or rollout cohort while applying bulk tags and service actions", "Review activation counts, firmware distribution, alert queues, and queued OTA or support work without implying this marketing site is the production operator console"},
			Outcomes:     []string{"Give product and operations teams a credible public admin-platform narrative", "Keep the website lead admin boundary separate from future IoT console scope", "Connect provisioning, release, and support workflows under one fleet-operations story"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Registration",
					Title:   "Register nodes with durable device identity boundaries",
					Intro:   "The fleet story now starts at the point where hardware leaves the factory and enters a managed registry.",
					Items: []string{
						"Record serial number, model, MAC, factory lot, and claim state when a node is first registered into the platform catalog.",
						"Issue bootstrap certificates or device credentials, then support rotation or revocation workflows when products are repaired, replaced, or reworked.",
						"Keep claim tokens, activation checkpoints, and ownership-transfer state tied to the device record so support teams can trace onboarding history.",
					},
				},
				{
					Eyebrow: "Operations",
					Title:   "Use one registry to drive firmware and support workflows",
					Intro:   "Registry views are positioned as the operator surface that connects provisioning, OTA, and customer support work.",
					Items: []string{
						"Search the device registry by region, firmware, product family, installer, or customer account and save groups for repeat operations.",
						"Coordinate firmware images and OTA jobs from the same operations surface so release managers can move from device search to rollout action without spreadsheet handoffs.",
						"Apply batch tags, metadata edits, ownership updates, reboot requests, or service-state changes when support teams need to act on a cohort instead of one device at a time.",
					},
				},
				{
					Eyebrow: "Operator Widgets",
					Title:   "Show the metrics an IoT admin console would surface first",
					Intro:   "Dashboard widgets keep the page concrete while staying honest about what this repository actually ships today.",
					Items: []string{
						"Summarize activation counts, firmware mix, online-versus-offline ratios, alert backlogs, and support escalations in operator-facing statistics widgets.",
						"Use these widgets to highlight where OTA jobs, registration failures, or field alerts need attention before they become customer-visible incidents.",
						"The existing /admin/leads page only covers website sales leads; the future IoT platform admin console described here remains a public product narrative rather than a shipped control plane in this repo.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Operations Surface",
				Title:   "Map each admin workflow to the right platform boundary",
				Intro:   "Realtek Connect+ now describes the public operations story in concrete operator terms while keeping the website's own admin runtime clearly scoped to lead review.",
				Columns: []string{"Workflow", "What operators manage", "Why it matters", "Website boundary"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Node registration", "Serials, models, bootstrap certificates, claim state, and factory readiness for newly manufactured devices.", "Gives support and onboarding teams a traceable source of truth before devices reach customers.", "Described as future platform-console scope; this repo does not expose a live node-registration UI."}},
					{Cells: []string{"Device registry", "Searchable inventory, groups, tags, timezone or ownership metadata, and node-sharing visibility.", "Keeps fleets searchable and lets teams target the right cohort for support or release actions.", "Public feature content only; the current Go app ships marketing pages plus protected lead review."}},
					{Cells: []string{"Release operations", "Firmware images, OTA jobs, rollout cohorts, and batch actions tied to registry segments.", "Lets release managers coordinate campaigns without breaking the device inventory context.", "OTA and fleet operations are described credibly, but the website runtime is not a production release console."}},
					{Cells: []string{"Statistics widgets", "Activation totals, firmware distribution, alert backlog, and operator task queues.", "Surfaces the metrics product, support, and operations teams watch first after launch.", "These widgets describe expected operator dashboards; the shipped admin endpoint remains /admin/leads for website sales workflow only."}},
				},
			},
		},
		{
			Slug:         "smart-home",
			Title:        "Smart Home Experience",
			Icon:         "home",
			Kicker:       "Give end users clear control after onboarding.",
			Summary:      "Remote control, local control fallback, schedules, scenes, grouping, device sharing, push notifications, and alerts for connected home products.",
			Description:  "Smart Home Experience describes the consumer-facing product surface that sits on top of Realtek Connect+ cloud and app foundations. It covers how households control devices, automate routines, share access, and stay informed without implying that this marketing website is the shipped mobile app runtime.",
			ImagePath:    "/static/assets/feature-smart-home-experience.png",
			ImageAlt:     "Smart home app control surface with connected home devices, scenes, schedules, grouping, and notification cards.",
			Highlights:   []string{"Remote and local control paths for everyday device actions", "Schedules, scenes, groups, and shared-home coordination", "Push notifications and alerts that bring users back into the branded app at the right moment"},
			Capabilities: []string{"Device control surfaces for power, mode, status, and household context across mobile and local-network touchpoints", "Automation building blocks for recurring schedules, multi-device scenes, room or home grouping, and temporary or permanent node sharing", "Alert delivery for onboarding completion, offline state, abnormal events, maintenance reminders, and actionable support flows"},
			Outcomes:     []string{"Make connected products feel useful after first setup", "Reduce support friction when households share devices or automate routines", "Show Realtek Connect+ as both an operator platform and an end-user product experience"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Control Modes",
					Title:   "Support both cloud reach and at-home responsiveness",
					Intro:   "The page frames control as a product experience with multiple paths instead of a single mobile button.",
					Items: []string{
						"Use remote control for away-from-home power, mode, and status changes when devices stay connected through the Realtek Connect+ cloud path.",
						"Keep local control available on the home network so core actions can stay responsive during WAN degradation or when products intentionally prioritize nearby control.",
						"Align these control surfaces with the branded app rather than implying that this Go website is the live smart-home client.",
					},
				},
				{
					Eyebrow: "Automation",
					Title:   "Turn single devices into routines households can depend on",
					Intro:   "Automation content focuses on the user workflows buyers expect once a product ships beyond basic provisioning.",
					Items: []string{
						"Create recurring schedules around daily routines, quiet hours, occupancy assumptions, or energy-saving windows.",
						"Bundle scenes so users can trigger coordinated actions across lights, climate, appliances, or custom device categories from one tap.",
						"Group devices by room, home, or product set so the app can present household-level control instead of one-node-at-a-time management.",
					},
				},
				{
					Eyebrow: "Household Sharing",
					Title:   "Make multi-user homes manageable without losing accountability",
					Intro:   "Sharing and notification flows are treated as part of the consumer product model, not as back-office admin tools.",
					Items: []string{
						"Support node sharing so primary owners can invite family members, installers, or temporary guests with bounded access expectations.",
						"Use push notifications for onboarding completion, automation results, offline alerts, abnormal events, and OTA prompts that need the user back in the app.",
						"Surface alerts with enough device and household context to help users act quickly without turning every product event into noise.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "End-user Workflows",
				Title:   "Map the home experience to the right control pattern",
				Intro:   "Realtek Connect+ uses this page to describe how households move between direct control, automation, and shared-home coordination while the website itself remains a product narrative.",
				Columns: []string{"Workflow", "What the user does", "Why it matters", "Platform boundary"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Remote control", "Open the branded app away from home to adjust power, mode, or device state through the cloud path.", "Keeps products useful when the user is not on the local network.", "Described as app and cloud capability scope; this repo does not ship the native control client."}},
					{Cells: []string{"Local control", "Use home-network control paths for fast response and resilient fallback when internet conditions are poor.", "Improves perceived reliability for products that users expect to react immediately.", "Positioned as a product capability story layered on top of Realtek Connect+ device and app integration work."}},
					{Cells: []string{"Schedules and scenes", "Set timed routines and multi-device actions that match household habits instead of manually repeating the same steps.", "Turns connected devices into repeatable home workflows instead of one-off remote commands.", "Automation behavior is described here without claiming a finished rules engine in this website runtime."}},
					{Cells: []string{"Grouping and sharing", "Organize devices by room or home and invite additional household members into the right node set.", "Makes multi-device homes and multi-user access understandable at consumer scale.", "The marketing site explains the end-user model while leaving production identity and permissions to future app/platform implementations."}},
					{Cells: []string{"Push notifications and alerts", "Receive actionable updates about onboarding, offline status, abnormal events, or pending actions that need attention.", "Brings users back into the app only when context matters and supports support-readiness planning.", "Notification delivery is presented as product capability scope, not as a promise that this website emits live device alerts today."}},
				},
			},
		},
		{
			Slug:         "user-management",
			Title:        "User Management",
			Icon:         "user-shield",
			Kicker:       "Handle the account lifecycle around connected products.",
			Summary:      "Platform content for sign up, sign in, OTP verification, social login, password recovery, account changes, and account deletion.",
			Description:  "User Management describes the account lifecycle capabilities product teams usually need around a Realtek-based connected product. It covers identity onboarding, recovery, and privacy operations for future product apps and services. This website does not expose end-user sign-in or account management flows today.",
			ImagePath:    "/static/assets/feature-user-management.png",
			ImageAlt:     "User management console with profile cards, security verification, sharing permissions, and account lifecycle controls.",
			Highlights:   []string{"Self-service sign up and sign in journeys for branded mobile apps", "One-time password verification for account activation, recovery, and high-risk actions", "Third-party login and account-linking paths for partner or consumer ecosystems"},
			Capabilities: []string{"Forgot-password, change-password, and session-management controls", "Account deletion and retention workflows that hand off cleanly to support and compliance teams", "User profile, consent, and device-ownership state that stays separate from this marketing website"},
			Outcomes:     []string{"Shorten time to a production-ready account system", "Keep user lifecycle scope explicit during architecture reviews", "Avoid confusing product platform capabilities with the website's own lead-capture flows"},
		},
		{
			Slug:         "app-sdk",
			Title:        "App SDK",
			Icon:         "phone-code",
			Kicker:       "Build branded mobile experiences faster.",
			Summary:      "iOS and Android SDK modules, sample app baselines, push notifications, rebrand guidance, and app publishing paths for connected products.",
			Description:  "App SDK now frames the mobile experience as a launch surface product teams can brand, extend, and publish without rebuilding every connected-app primitive from scratch. It covers iOS and Android SDK layers, sample app structure, push workflows, and release planning while staying explicit that this repo is a server-rendered website, not a shipped mobile framework.",
			ImagePath:    "/static/assets/feature-app-sdk.png",
			ImageAlt:     "Mobile app SDK workspace with app screens, code modules, push notification blocks, and publishing checklist.",
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
			Icon:         "chart",
			Kicker:       "See the health of products in the field.",
			Summary:      "Activation statistics, firmware distribution, crash reports, logs, reboot reasons, RSSI, and memory signals.",
			Description:  "Insights gives engineering and support teams a view into fleet quality. It highlights operational statistics and device health signals that help teams prioritize fixes and understand real deployment behavior.",
			ImagePath:    "/static/assets/feature-insights-dashboard.jpg",
			ImageAlt:     "Operations insights dashboard with fleet health charts, alert cards, and device telemetry summaries.",
			Highlights:   []string{"Activation and association statistics", "Crash, reboot, and log visibility", "Firmware distribution and device health metrics"},
			Capabilities: []string{"RSSI, memory, uptime, and reboot reason signals", "Version adoption and rollout health", "Support-oriented device history"},
			Outcomes:     []string{"Find field issues earlier", "Support customers with evidence", "Measure firmware quality"},
		},
		{
			Slug:         "private-cloud",
			Title:        "Private Cloud",
			Icon:         "cloud-lock",
			Kicker:       "Deploy with enterprise ownership and control.",
			Summary:      "Compare public evaluation with private commercial deployment, regional hosting, custom domains, and enterprise upgrade planning.",
			Description:  "Private Cloud explains how Realtek Connect+ moves from a shared evaluation story into dedicated commercial deployment. It positions data ownership, regional placement, custom domain control, and upgrade planning as enterprise buying criteria rather than unsupported promises about this website's own runtime.",
			ImagePath:    "/static/assets/feature-private-cloud-architecture.jpg",
			ImageAlt:     "Private cloud architecture showing dedicated regions, branded domain entry points, and enterprise control boundaries.",
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
			Icon:         "nodes",
			Kicker:       "Connect products to the wider IoT ecosystem.",
			Summary:      "Matter Fabric positioning, voice assistant paths, MQTT over TLS, REST APIs, and webhooks for product and platform integrations.",
			Description:  "Integrations explains how Realtek Connect+ fits into smart home ecosystems and enterprise backends. It frames Matter Fabric participation, voice assistant connections, secure protocol access, and webhook delivery as supported integration patterns without claiming this website already operates every downstream service.",
			ImagePath:    "/static/assets/feature-integrations.png",
			ImageAlt:     "Integration hub connecting generic Matter, voice assistant, REST API, MQTT over TLS, webhook, app, and enterprise system endpoints.",
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
						"Use the page to show Realtek Connect+ interoperability intent for smart home buyers without implying this website alone ships live assistant certification.",
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
