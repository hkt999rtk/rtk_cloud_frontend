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
	Flows        []FeatureFlow
	Sections     []FeatureSection
	Table        FeatureTable
	Tables       []FeatureTable
	RelatedLinks []FeatureRelatedLink
}

type FeatureFlow struct {
	Eyebrow, Title, Intro string
	Steps                 []FeatureFlowStep
}
type FeatureFlowStep struct{ Title, Body string }
type FeatureSection struct {
	Eyebrow, Title, Intro string
	Items                 []string
	Accent                bool
}
type FeatureTable struct {
	Eyebrow, Title, Intro string
	Columns               []string
	Rows                  []FeatureTableRow
}
type FeatureTableRow struct{ Cells []string }
type FeatureRelatedLink struct{ Title, Summary, Href string }

func All() []Feature {
	return []Feature{
		{
			Slug: "provision", Title: "Provision", Icon: "provision",
			Kicker:      "Bring devices online quickly and securely.",
			Summary:     "Create a clear onboarding path from first network setup to secure cloud registration.",
			Description: "From first-time setup to cloud registration, Realtek Connect+ Provision helps product teams create a consistent and secure device onboarding experience. Once registered, devices can connect to OTA updates, fleet management, and other cloud services. Realtek works with each team to plan the onboarding method and integration scope for its product.",
			ImagePath:   "/static/assets/portal-provisioning-desktop-v2.webp", ImageAlt: "Sanitized Realtek Connect+ device registration workspace",
			Highlights:   []string{"Guide users through a clear first-time setup", "Create a trusted device identity in the cloud", "Continue directly into updates and fleet operations"},
			Capabilities: []string{"Plan Wi-Fi, Bluetooth, QR code, or other suitable onboarding paths with the Realtek team", "Connect device registration with the services your product needs", "Give product and support teams a consistent view of onboarding progress"},
			Outcomes:     []string{"Reduce friction before a device is ready to use", "Build a repeatable onboarding experience across product lines", "Prepare every registered device for ongoing cloud operations"},
			Flows:        []FeatureFlow{{Eyebrow: "Customer journey", Title: "Three steps from setup to cloud services", Intro: "The exact experience is planned around your product while keeping the customer journey simple.", Steps: []FeatureFlowStep{{Title: "Complete first-time network setup", Body: "Help the device connect using the onboarding method that best fits the product."}, {Title: "Create device identity and register to the cloud", Body: "Establish a trusted device record so cloud services can recognize and manage it."}, {Title: "Connect OTA, fleet management, and other services", Body: "Move directly from registration into the services used throughout the product lifecycle."}}}},
		},
		marketingFeature("ota", "OTA firmware updates", "ota", "Update products safely after launch.", "Plan releases, roll out in stages, monitor progress, and respond when a device needs attention.", "Realtek Connect+ OTA helps product and operations teams deliver firmware with control. Choose the target devices, release in stages, and keep progress visible throughout the update.", "/static/assets/portal-ota-desktop-v2.webp", "Sanitized firmware and OTA operations workspace", []string{"Plan controlled firmware releases", "Roll out by device group or schedule", "See progress and handle exceptions"}),
		marketingFeature("fleet-management", "Fleet management", "fleet", "Keep every device fleet organized.", "Bring devices, groups, health signals, alerts, and coordinated actions into one operating view.", "Realtek Connect+ Fleet Management gives operations teams a practical way to organize connected products, identify devices that need attention, and coordinate routine actions across groups.", "/static/assets/portal-fleet-desktop-v2.webp", "Sanitized connected-device fleet workspace", []string{"Organize devices into useful groups", "Review health and alerts in one place", "Coordinate actions across the fleet"}),
		marketingFeature("smart-home", "Smart Home", "home", "Make connected-home experiences easier.", "Simplify first-time setup, everyday control, sharing, and automation for customers.", "Realtek Connect+ helps product teams shape a connected-home experience that feels consistent from setup to daily use, while leaving room for the workflows and integrations that make each brand distinct.", "/static/assets/feature-smart-home-experience.png", "Connected-home setup, control, sharing, and automation experience", []string{"Simplify device setup and control", "Support household sharing", "Connect common automation experiences"}),
		marketingFeature("user-management", "User management", "shield-user", "Manage people and device access clearly.", "Organize accounts, device relationships, sharing, and access rights as the product grows.", "Realtek Connect+ User Management helps teams keep customer accounts and device access understandable, from first ownership through sharing and support.", "/static/assets/portal-users-public-v1.webp", "Sanitized user and access management workspace", []string{"Connect customer accounts and devices", "Support sharing and access roles", "Give support teams a clearer operating view"}),
		marketingFeature("app-sdk", "App SDK", "phone-code", "Launch the branded app experience sooner.", "Use SDKs and reference paths to shorten integration time for onboarding, control, and account experiences.", "Realtek Connect+ App SDK helps mobile and web teams integrate the essential connected-product experiences without rebuilding every cloud interaction from the beginning.", "/static/assets/portal-sdk-public-v1.webp", "Sanitized SDK and application integration workspace", []string{"Accelerate onboarding and control experiences", "Reuse consistent account and device flows", "Move from evaluation to branded product faster"}),
		marketingFeature("video-cloud", "Video Cloud", "dashboard", "Connect live video to the product experience.", "Bring devices, secure cloud viewing, and application workflows together for connected-camera products.", "Realtek Connect+ Video Cloud supports product teams building live-view and device experiences across the camera, cloud, and customer application.", "/static/assets/portal-video-public-v1.webp", "Sanitized live video operations workspace", []string{"Connect device and cloud viewing paths", "Support live product experiences", "Keep video operations visible to the team"}),
		marketingFeature("insights", "Insights", "chart", "Turn device signals into operating insight.", "Understand fleet health, firmware distribution, and product trends in views teams can act on.", "Realtek Connect+ Insights helps product and operations teams see how connected devices are performing and where follow-up may be needed.", "/static/assets/portal-insights-public-v1.webp", "Sanitized fleet reports and insights workspace", []string{"Monitor device health trends", "Understand firmware distribution", "Focus teams on the issues that matter"}),
		marketingFeature("integrations", "Integrations", "nodes", "Connect the cloud to existing services.", "Use APIs and product events to connect Realtek Connect+ with the systems your teams already operate.", "Realtek Connect+ Integrations provide a practical path for connecting device operations with business systems, applications, and product services. The Realtek team plans the right integration approach with each customer.", "/static/assets/feature-integrations.png", "Cloud integration architecture connecting applications and product services", []string{"Connect existing product services", "Use APIs for the workflows teams need", "Plan integrations around the product architecture"}),
		marketingFeature("security", "Security", "cloud-lock", "Build on a trusted device-to-cloud foundation.", "Use device identity, PKI, and access control to protect connected-product operations.", "Realtek Connect+ Security establishes a practical foundation for recognizing devices, protecting cloud access, and managing who can perform sensitive actions.", "/static/assets/connectplus-architecture-diagram-corporate-v2.jpg", "Device identity and secure cloud access architecture", []string{"Establish trusted device identity", "Use certificates and PKI foundations", "Control access to product operations"}),
		cloudPlansFeature(),
	}
}

func marketingFeature(slug, title, icon, kicker, summary, description, imagePath, imageAlt string, highlights []string) Feature {
	return Feature{
		Slug: slug, Title: title, Icon: icon, Kicker: kicker, Summary: summary, Description: description, ImagePath: imagePath, ImageAlt: imageAlt,
		Highlights:   highlights,
		Capabilities: []string{"Plan the experience around the product and its customers", "Bring device and cloud operations into one consistent path", "Work with the Realtek team on the right integration scope"},
		Outcomes:     []string{"Shorten the path from integration to daily operations", "Give product and operations teams a shared view", "Keep the experience ready to evolve with the product"},
	}
}

func cloudPlansFeature() Feature {
	return Feature{
		Slug: "private-cloud", Title: "Cloud plans", Icon: "cloud-lock",
		Kicker:      "Choose the cloud operating model that fits your product.",
		Summary:     "Start with the recommended Realtek-managed service, with Private Cloud available for dedicated deployment needs.",
		Description: "The recommended Realtek-managed service is built, hosted, and operated by Realtek, with flexible payment based on actual usage. For products that require a dedicated environment or specific infrastructure boundaries, the Realtek team can also plan a Private Cloud deployment.",
		ImagePath:   "/static/assets/feature-private-cloud-architecture.jpg", ImageAlt: "Managed cloud and private cloud deployment options",
		Highlights:   []string{"Recommended Realtek-managed operating service", "Flexible payment based on actual usage", "Private Cloud available for dedicated requirements"},
		Capabilities: []string{"Realtek-managed setup and ongoing operations", "Usage, billing, invoices, and payments in the service portal", "Private deployment planning with the Realtek team"},
		Outcomes:     []string{"Start using the service without building the operating platform first", "Match the cloud model to product and business needs", "Keep a path to a dedicated environment when required"},
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
