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
			Kicker:       "Separate firmware lifecycle from campaign roadmap.",
			Summary:      "Firmware upload, catalog, target enablement, rollout status, report, cancel, and download are available foundations; scheduled, time-window, user-consent, and archive campaign policy surfaces are now available too, while approval workflow, dashboards, analytics, and staged percentage rollout remain roadmap scope.",
			Description:  "OTA is presented as an interface-first firmware campaign path. Current public copy can describe the firmware lifecycle foundation already represented by upload, enablement, rollout query/report, cancel, and download routes, while the implemented scheduled, time-window, user-consent, and archive policy surfaces are promoted as available campaign behavior. Approval workflow, dashboards, analytics, and staged percentage rollout stay clearly labeled as roadmap work.",
			ImagePath:    "/static/assets/feature-ota-control-center.jpg",
			ImageAlt:     "Firmware rollout control center with staged release timeline, device cohorts, and OTA job status cards.",
			SourceLabel:  "Firmware campaign interface contract",
			SourceURL:    "https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/FIRMWARE_CAMPAIGN.md",
			Highlights:   []string{"Available firmware lifecycle foundation for upload, catalog, target enablement, rollout status, report, cancel, and download", "Available campaign policy surfaces for schedule, time-window, user-consent, and archive operations", "Roadmap framing for approval workflow, dashboards, analytics, and staged percentage rollout"},
			Capabilities: []string{"Use existing firmware routes as the available implementation boundary instead of implying a complete campaign engine", "Promote scheduled, time-window, user-consent, and archive behavior as available campaign policy vocabulary now that backend and SDK support has landed", "Keep cancel available as lifecycle behavior while treating approval workflow, dashboards, analytics, and staged percentage rollout as roadmap scope"},
			Outcomes:     []string{"Keep OTA campaign vision visible without overclaiming phase-one implementation", "Give firmware, SDK, backend, and product teams one shared availability vocabulary", "Make buyer conversations explicit about what is available now, implemented, or roadmap"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Available Foundation",
					Title:   "Use the current firmware lifecycle as the implementation boundary",
					Intro:   "The OTA story starts with the firmware surfaces that exist today, not with a claim that the full commercial campaign engine is complete.",
					Items: []string{
						"Describe publish, enablement, whitelist, rollout query/report, cancel, and download behavior as the available firmware lifecycle foundation.",
						"Keep release metadata, model and version checks, target enablement, and device-reported rollout status tied to the existing backend route inventory.",
						"Treat force and normal policy language as basic delivery vocabulary unless a deeper campaign policy engine is explicitly implemented.",
					},
				},
				{
					Eyebrow: "Implemented Policy Surfaces",
					Title:   "Promote the campaign policy surfaces that are now implemented",
					Intro:   "The campaign contract defines policy names so teams can align implementation, and the now-implemented scheduled, time-window, user-consent, and archive surfaces can be described as available rather than roadmap.",
					Items: []string{
						"Scheduled and time-window OTA are available campaign-policy surfaces with backend enforcement and SDK handling in place.",
						"User-consent-required OTA is now a supported policy surface rather than a placeholder mobile UX.",
						"Archive is available for closing or hiding completed campaigns without deleting audit history.",
					},
				},
				{
					Eyebrow: "Roadmap Guardrails",
					Title:   "Keep phase-two operations out of available-now claims",
					Intro:   "The page preserves the commercial OTA direction while clearly marking unsupported campaign operations as follow-up scope.",
					Items: []string{
						"Approval workflow, operator dashboards, analytics, and success-rate reporting are roadmap capabilities, not phase-one availability claims.",
						"Staged percentage rollout and automatic cohort ramping stay out of the available feature list until a campaign engine implements them.",
						"Device firmware install protocol, firmware signing policy, and mobile consent UX remain owned by their implementation layers rather than this public website.",
					},
					Accent: true,
				},
			},
			Table: FeatureTable{
				Eyebrow: "Availability Labels",
				Title:   "Map each OTA concept to the right implementation status",
				Intro:   "Realtek Connect+ uses the firmware campaign contract as shared vocabulary while separating available lifecycle behavior and implemented policy surfaces from later campaign operations.",
				Columns: []string{"Concept", "Status", "Public copy stance", "Follow-up boundary"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Firmware lifecycle foundation", "Available foundation", "Firmware upload/catalog, target enablement, device rollout status, report, cancel, and download can be discussed as current baseline behavior.", "A richer campaign engine still needs owner-repo implementation before stronger public claims."}},
					{Cells: []string{"Scheduled policy", "Available foundation", "Describe as an available campaign vocabulary for a supported start time or maintenance schedule.", "Richer operator UX can extend the supported surface later."}},
					{Cells: []string{"Time-window policy", "Available foundation", "Explain as an available campaign time-window constraint for eligible installs.", "Time-zone handling and device/app behavior can extend the supported surface later."}},
					{Cells: []string{"User-consent policy", "Available foundation", "Name it as a supported policy surface rather than a placeholder consent experience.", "Mobile UX refinements can extend the supported surface later."}},
					{Cells: []string{"Cancel", "Available foundation", "Keep cancel tied to stopping eligible pending firmware rollouts in the current lifecycle foundation.", "Do not imply broader campaign pause, approval, or analytics workflows."}},
					{Cells: []string{"Archive", "Available foundation", "Describe archive as an available campaign-management surface for closing or hiding completed campaigns without deleting audit history.", "Active campaign views can extend the supported surface later."}},
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
			Summary:      "iOS and Android SDK modules, app-side and device-side reference samples, push notification planning, rebrand guidance, and publishing paths for connected products.",
			Description:  "App SDK now frames the mobile experience as a launch surface product teams can brand, extend, and publish without rebuilding every connected-app primitive from scratch. It covers iOS and Android SDK layers, Android/iOS/WebApp home app samples, Linux device simulation, a PRO2 camera device demo, push workflows, and release planning while staying explicit that this repo is a server-rendered website, not a shipped mobile framework.",
			ImagePath:    "/static/assets/connectplus-sample-ecosystem.png",
			ImageAlt:     "Realtek Connect+ sample ecosystem diagram with Android, iOS, WebApp, Linux simulator, PRO2 device demo, and cloud hub nodes.",
			Highlights:   []string{"iOS and Android SDK coverage for onboarding, control, and account flows", "Android, iOS, WebApp, Linux Simulator, and PRO2 Device Demo reference samples", "Push notification, release-readiness, and app publishing guidance"},
			Capabilities: []string{"Shared mobile primitives for sign-in, provisioning, device control, and sharing", "Reference app and device samples that validate SDK usage across app-side and device-side flows", "Release planning for App Store and Google Play submission, rollout, and support operations"},
			Outcomes:     []string{"Shorten mobile and device integration validation", "Keep app developers, firmware teams, and product owners aligned on sample boundaries", "Avoid one-off app/cloud/device integration work for every product line"},
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
					Eyebrow: "Reference Samples",
					Title:   "Reference sample applications prove SDK usage before product integration",
					Intro:   "The sample ecosystem helps app developers, firmware teams, and product owners validate flows with concrete app-side and device-side references.",
					Items: []string{
						"Use the rtk_cloud_client repository as the source of truth for sample code, specifications, and sample-specific README files; this website summarizes the customer evaluation path rather than hosting SDK source.",
						"Use Android Home Automation, iOS Home Automation, and WebApp Ops Lab samples to validate provisioning, device list/detail, light and AC control, camera monitor, and debug report flows.",
						"Use the Linux Simulator and PRO2 Device Demo to validate device-side command handling, sample MQTT payloads, snapshot upload, status/log/event reporting, and the WebRTC answerer boundary.",
						"Treat these as SDK usage references, not production app-store apps or white-label release packages.",
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
				Eyebrow: "Sample Ecosystem",
				Title:   "Reference sample applications across app and device surfaces",
				Intro:   "Each sample helps customers validate a specific Realtek Connect+ integration path while keeping production app ownership and formal cloud wire contracts separate.",
				Columns: []string{"Sample", "Side", "Validation focus", "Boundary"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Android Home Automation sample", "App", "Provisioning adapter states, device list/detail, light and AC controls, camera monitor, and debug report.", "Native Kotlin reference; not an app-store deliverable."}},
					{Cells: []string{"iOS Home Automation sample", "App", "Swift SDK usage for setup profiles, device control, camera boundaries, and debug evidence.", "Native Swift reference; customer product teams own release UX and signing."}},
					{Cells: []string{"WebApp Ops Lab sample", "App", "Cloud-side onboarding, MQTT payload inspection, simulated controls, camera helpers, and debug report.", "Browser reference without BLE or SoftAP onboarding."}},
					{Cells: []string{"Linux Simulator", "Device", "Light, AC, and camera command handling without hardware, including local state and validation output.", "Simulator for evidence and development; not a production device firmware package."}},
					{Cells: []string{"PRO2 Device Demo", "Device", "Device-bound token, owner transport, snapshot upload, camera logs/events, and WebRTC answerer boundary.", "Firmware reference; concrete Realtek SDK calls stay firmware-owned."}},
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
			Kicker:       "Any cloud or on-premises. No serverless lock-in.",
			Summary:      "Realtek Connect+ private deployment runs as a container or VM workload on GCP, Azure, AWS, or your own data center — no cloud vendor dependency, no serverless runtime required.",
			Description:  "Private Cloud explains the two deployment tiers and Realtek Connect+'s infrastructure model. The platform runs as a standard container or VM workload, giving customers full choice of cloud provider or on-premises infrastructure. This is explicitly different from serverless-native IoT platforms that tie private deployment to a single cloud account.",
			ImagePath:    "/static/assets/feature-private-cloud-architecture.jpg",
			ImageAlt:     "Private cloud architecture showing container and VM workloads running across multiple cloud providers and on-premises data centers.",
			Highlights:   []string{"VM/container deployment on GCP, Azure, AWS, or on-premises — no cloud lock-in", "Free evaluation tier with up to 200 devices on request and no expiry", "Commercial tier with one-time license plus annual maintenance and no minimum scale"},
			Capabilities: []string{"Container or VM workload deployment on any major cloud provider or on-premises data center", "Reverse-proxy TLS termination, network policy alignment, and branded service endpoints", "Release promotion and maintenance-window planning across evaluation and production environments"},
			Outcomes:     []string{"Avoid cloud vendor lock-in for IoT infrastructure", "Keep ownership and residency boundaries explicit", "Create a credible path from pilot to production on your own terms"},
			Sections: []FeatureSection{
				{
					Eyebrow: "Infrastructure Model",
					Title:   "Standard container and VM workloads — no serverless dependency",
					Intro:   "Realtek Connect+ private deployments run as conventional container or VM processes. There is no serverless runtime requirement and no dependency on a specific cloud provider's managed services.",
					Items: []string{
						"Deploy on GCP, Azure, AWS, or your own on-premises data center using standard container orchestration (Kubernetes, Docker Compose) or VM images.",
						"Unlike serverless-native IoT platforms that require the customer to own a specific cloud account (e.g., AWS Lambda + DynamoDB), Realtek Connect+ has no infrastructure prerequisites beyond a host that can run containers or VMs.",
						"This means regulated industries, enterprises with existing data center contracts, and teams with multi-cloud policies can deploy without restructuring their infrastructure strategy.",
					},
				},
				{
					Eyebrow: "Commercial Models",
					Title:   "Start with evaluation, then move into owned deployment boundaries",
					Intro:   "Public evaluation gives fast PoC access on shared infrastructure. Private commercial deployment is a dedicated environment on infrastructure you choose.",
					Items: []string{
						"Use the public evaluation environment to validate device flows, dashboards, and integration assumptions before committing to a customer-specific operating boundary.",
						"Transition to a dedicated private deployment once product teams need tenant isolation, formal support processes, and customer-specific change windows.",
						"The evaluation tier is suitable for development and pilot work; commercial products with real device fleets require a private commercial agreement.",
					},
				},
				{
					Eyebrow: "Plans & Limits",
					Title:   "Evaluation tier limits and the path to commercial scale",
					Intro:   "Concrete limits so developer teams can plan a pilot without surprises, and a clear handoff into the commercial conversation when the pilot grows.",
					Items: []string{
						"Evaluation accounts start with a 5-device default quota and can be raised up to 200 devices on request.",
						"Evaluation access does not expire — request a quota raise or a commercial conversation when your fleet grows; we do not auto-cancel evaluation accounts.",
						"Evaluation use is limited to development, proof-of-concept, and internal validation; commercial product shipments and customer-facing fleets require a private commercial agreement.",
						"Self-service signup with email verification is on the roadmap; pre-launch evaluation accounts are issued by the Realtek Connect+ team via the contact form.",
						"There is no minimum scale for the commercial tier; small fleets can move out of evaluation as soon as commercial use begins, even before they reach the 200-device evaluation ceiling.",
					},
				},
				{
					Eyebrow: "Pricing Factors",
					Title:   "How commercial pricing is shaped",
					Intro:   "We do not publish a price list — every commercial deployment is sized to the customer's actual scope. The factors below are the inputs the sales team uses when preparing a quote, so buyers can frame the conversation before getting on the phone.",
					Items: []string{
						"Fleet size — total addressable device count for the deployment, including planned expansion within the contract term.",
						"Deployment topology — single-region managed deployment, multi-region, or fully customer-operated infrastructure across one or more clouds or on-premises sites.",
						"Support coverage — the response-time, escalation, and on-call expectations the customer needs in their support agreement.",
						"Customization scope — branding/white-label, custom domain handling, and any product-specific platform extensions beyond the standard release.",
						"Term length — typical contract structure is a one-time platform license fee plus annual maintenance; multi-year terms are quoted separately.",
					},
				},
				{
					Eyebrow: "SDK Licensing",
					Title:   "What you can build with",
					Intro:   "SDK distribution today and the planned posture at general availability.",
					Items: []string{
						"Realtek Connect+ device SDK packages (Native C, Android, iOS, JavaScript/TypeScript, Go) are currently distributed under evaluation terms.",
						"An open-source SDK release is planned at general availability, so commercial customers and the wider community can integrate without bespoke license negotiation.",
						"The platform backend stays a proprietary commercial product; private deployments install signed builds rather than building from source.",
					},
				},
				{
					Eyebrow: "Support",
					Title:   "What support looks like at each tier",
					Intro:   "Support coverage is intentionally tiered: evaluation gets a self-serve community lane, commercial gets contracted accountability.",
					Items: []string{
						"Evaluation support is community-tier: documentation, integration guides, and the public issue tracker for the SDK once it is open. There is no response-time commitment on the evaluation tier.",
						"Commercial support is contract-defined: response-time, uptime, and escalation paths live in the customer agreement rather than a published tier matrix.",
						"Customers needing a specific SLA structure should raise it during the commercial conversation so it can be priced and committed inside the agreement.",
					},
				},
				{
					Eyebrow: "Ownership",
					Title:   "Define where data lives and how the service is branded",
					Intro:   "Private deployment gives customers control over data residency, access boundaries, and service identity.",
					Items: []string{
						"Device metadata, operator access logs, and retained support exports stay inside the customer-owned environment — no data crosses to a shared Realtek-operated region.",
						"Custom domains and branded entry points let the deployment align with the customer's DNS, certificate, and support model.",
						"Choose regional placement around residency, latency, and operational coverage requirements rather than accepting a fixed shared region.",
					},
					Accent: true,
				},
				{
					Eyebrow: "Upgrade Path",
					Title:   "Promote proven configurations into commercial production",
					Intro:   "The upgrade path is an engineering and operations workflow, not a one-click migration promise.",
					Items: []string{
						"Carry validated device models, app configuration, and integration settings from evaluation into a dedicated deployment plan.",
						"Use release promotion, maintenance windows, and rollback checkpoints to move from pilot tenants into production operations safely.",
						"Align the commercial cutover with customer security review, support readiness, and staged onboarding of real fleets.",
					},
				},
				{
					Eyebrow: "Deployment FAQ",
					Title:   "Questions enterprise buyers raise first",
					Intro:   "Direct answers to infrastructure and procurement questions.",
					Items: []string{
						"What infrastructure does private deployment run on? Standard containers or VMs — bring your own GCP, Azure, AWS, or on-premises host.",
						"Is there a cloud vendor requirement? No. Realtek Connect+ has no dependency on a specific cloud provider's managed services or serverless platform.",
						"Commercial support covers deployment planning, environment hardening expectations, and the path for future platform customization requests.",
					},
				},
			},
			Table: FeatureTable{
				Eyebrow: "Deployment Paths",
				Title:   "Compare evaluation and private commercial operating models",
				Intro:   "Realtek Connect+ positions deployment choice as a commercial decision. The infrastructure model stays the same — containers and VMs — regardless of which tier or cloud the customer chooses.",
				Columns: []string{"Model", "Infrastructure", "Device quota & cost", "Best fit"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Public evaluation", "Shared environment hosted by Realtek.", "5 devices by default, up to 200 on request. Free, non-commercial use only.", "Early evaluations, internal validation, and short proof-of-concept cycles."}},
					{Cells: []string{"Managed private deployment", "Container/VM on customer-selected cloud (GCP, Azure, AWS) or on-premises, operated with agreed support windows.", "No device floor. Commercial agreement: one-time license fee plus annual maintenance — contact sales for quote.", "Teams that want private deployment outcomes without owning the day-to-day platform operations stack."}},
					{Cells: []string{"Customer-operated private region", "Customer-owned container/VM infrastructure — any cloud or data center — with coordinated release and upgrade planning.", "No device floor. Commercial license + maintenance plus the customer's own infrastructure costs.", "Products with strict enterprise governance, regulated data boundaries, or multi-cloud / on-prem mandates."}},
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
		{
			Slug:        "security",
			Title:       "Security & PKI",
			Icon:        "certificate",
			Kicker:      "Ground device identity and cloud trust in X.509 certificates and a managed PKI hierarchy.",
			Summary:     "Device certificates, a two-tier CA hierarchy, mutual TLS authentication, certificate lifecycle operations, and OCSP/CRL revocation infrastructure for Realtek Connect+ deployments.",
			Description: "Security & PKI describes how Realtek Connect+ uses a Public Key Infrastructure built on X.509 certificates to establish verifiable device identity, protect cloud communications, and support fleet-wide revocation without custom protocol work. Each device receives a unique certificate signed by the platform CA hierarchy at provisioning time. Mutual TLS authenticates every cloud connection using that certificate so the platform can verify hardware identity, enforce policy, and rotate or revoke credentials across large fleets with standard tooling.",
			ImagePath:   "/static/assets/feature-private-cloud-architecture.jpg",
			ImageAlt:    "PKI hierarchy diagram showing root CA, intermediate CA, device certificates, and mutual TLS cloud connections.",
			SourceLabel: "Platform security contract",
			Highlights: []string{
				"Two-tier X.509 CA hierarchy: offline root CA and online issuing CA for device certificate issuance",
				"Per-device certificates provisioned at manufacture or first activation and bound to hardware identity",
				"Mutual TLS (mTLS) for every device-to-cloud connection — both sides present and verify X.509 certificates",
				"Certificate lifecycle operations: issuance, renewal, rotation, and revocation with OCSP and CRL support",
			},
			Capabilities: []string{
				"Operate a two-tier CA hierarchy where the root CA stays offline and the issuing CA signs device certificates on demand during provisioning workflows",
				"Bind each device certificate to its serial number, MAC address, and model so the cloud can verify hardware identity without shared secrets",
				"Enforce mutual TLS on MQTT and HTTPS endpoints so unauthenticated devices cannot reach platform APIs or message brokers",
				"Support certificate renewal and rotation workflows triggered by expiry schedules or security events without requiring full device re-provisioning",
				"Publish CRL endpoints and run an OCSP responder so relying parties can check certificate validity in real time",
				"Revoke individual device certificates or batch-revoke a compromised manufacturing lot through the fleet management console",
			},
			Outcomes: []string{
				"Eliminate shared-secret risks by grounding every device identity in an individual X.509 certificate",
				"Satisfy enterprise and operator security reviews with standard PKI artefacts — certificate chains, CRL distribution points, OCSP endpoints",
				"Keep revocation fast and scoped: a single compromised device certificate does not expose the rest of the fleet",
				"Give audit teams a traceable issuance and revocation log tied to hardware manufacturing records",
			},
			Sections: []FeatureSection{
				{
					Eyebrow: "CA Hierarchy",
					Title:   "Two-tier PKI anchored in an offline root CA",
					Intro:   "The platform CA design separates long-term trust anchors from online issuance operations so the root key is never exposed to network threats.",
					Items: []string{
						"The root CA is kept offline and used only to sign the intermediate issuing CA certificate and any cross-certification material. Its private key never touches a network-connected host.",
						"The issuing CA operates online within the platform deployment boundary. It receives provisioning requests, validates manufacturing metadata, and signs device certificates during activation.",
						"Certificate chains are short — root → issuing CA → device — so relying parties and firmware TLS stacks can validate with minimal chain processing overhead.",
						"CA certificate and key material for private cloud deployments can be generated inside a customer-controlled HSM or key management service, keeping root key custody with the operator.",
					},
				},
				{
					Eyebrow: "Device Identity",
					Title:   "Per-device X.509 certificates bound to hardware identity",
					Intro:   "Each manufactured device receives a unique certificate that encodes its hardware identity in the Subject and SubjectAltName fields.",
					Items: []string{
						"The certificate Subject carries serial number, model, and manufacturing lot. SubjectAltName extensions carry MAC address and device type URIs for policy matching.",
						"Certificates are generated at one of two injection points: factory provisioning before shipment, or cloud-side issuance during first-activation onboarding for products that carry a bootstrap credential.",
						"The device private key is generated on-device and never leaves the hardware. Only the certificate signing request (CSR) is transmitted to the issuing CA.",
						"Certificate validity periods are set by product lifecycle expectations. Consumer devices typically receive 5–10 year certificates; industrial deployments can use shorter periods with automated renewal.",
					},
				},
				{
					Eyebrow: "Mutual TLS",
					Title:   "mTLS on every device-to-cloud connection",
					Intro:   "The platform enforces certificate-based mutual authentication on all MQTT broker and HTTPS API endpoints used by devices.",
					Items: []string{
						"MQTT broker connections require the device to present its X.509 certificate at TLS handshake. The broker validates the certificate chain against the platform trust store and rejects connections from unrecognized or revoked certificates.",
						"HTTPS device APIs apply the same client certificate requirement so every REST interaction carries a verifiable hardware identity.",
						"Policy enforcement at the broker and API gateway uses the verified certificate Subject fields — serial number, model, lot — to scope which topics, endpoints, and operations the device is permitted to use.",
						"Server-side certificates on broker and API endpoints are signed by a separate service CA so devices can pin the expected server CA without conflating their own issuing CA with server identity.",
					},
					Accent: true,
				},
				{
					Eyebrow: "Lifecycle Operations",
					Title:   "Issuance, renewal, rotation, and revocation at fleet scale",
					Intro:   "Certificate lifecycle management is designed to work at fleet scale through scheduled operations and event-driven triggers rather than per-device manual workflows.",
					Items: []string{
						"Expiry-based renewal is triggered by the device or the fleet management console before the certificate validity window closes. The device generates a fresh CSR and the issuing CA produces a replacement certificate without requiring re-provisioning.",
						"Forced rotation can be triggered by a security event — compromised lot, key exposure, CA algorithm migration — and dispatched as an OTA-style campaign to affected cohorts.",
						"Certificate revocation is recorded in the platform CRL and reflected at the OCSP responder within the configured propagation window. Revoked devices are disconnected at the next TLS handshake.",
						"The manufacturing record and activation log are retained alongside the certificate issuance history so support teams can reconstruct the credential chain for any device in the fleet.",
					},
				},
			},
			Table: FeatureTable{
				Eyebrow: "PKI components",
				Title:   "Platform security building blocks",
				Intro:   "Each component is independently scoped so private cloud operators can substitute their own CA, HSM, or revocation infrastructure where required.",
				Columns: []string{"Component", "Role", "Operator flexibility"},
				Rows: []FeatureTableRow{
					{Cells: []string{"Root CA", "Offline trust anchor; signs the issuing CA certificate only.", "Can be customer-operated and held in an HSM outside the platform deployment boundary."}},
					{Cells: []string{"Issuing CA", "Online CA that signs device certificates during provisioning and renewal.", "Can be hosted inside the platform or replaced by a customer-operated intermediate CA that chains to the root."}},
					{Cells: []string{"Device certificate", "Per-device X.509 credential binding hardware identity to the platform PKI.", "Certificate profile (validity, key usage, SANs) is configurable per product line or deployment."}},
					{Cells: []string{"Mutual TLS enforcement", "Client certificate requirement on MQTT broker and HTTPS device API endpoints.", "Certificate validation policy and topic/endpoint ACLs are configurable per fleet or deployment tier."}},
					{Cells: []string{"OCSP responder", "Real-time certificate status endpoint used by relying parties to check revocation.", "Can be operated by the platform or delegated to a customer PKI service."}},
					{Cells: []string{"CRL distribution", "Periodic revocation list published for clients that cache status offline.", "CRL publication interval and distribution points are configurable per deployment."}},
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
