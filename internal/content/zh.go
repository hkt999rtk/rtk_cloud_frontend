package content

import (
	"strings"

	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
)

type localizedFeature struct {
	Title        string
	Kicker       string
	Summary      string
	Description  string
	ImageAlt     string
	SourceLabel  string
	SourceURL    string
	Highlights   []string
	Capabilities []string
	Outcomes     []string
	Flows        []features.FeatureFlow
	Sections     []features.FeatureSection
	Table        features.FeatureTable
	Tables       []features.FeatureTable
	RelatedLinks []features.FeatureRelatedLink
}

type localizedDoc struct {
	Title        string
	Kicker       string
	Summary      string
	Description  string
	Highlights   []string
	Deliverables []string
	Audience     []string
	Table        docs.Table
}

func zhTWCatalog() Catalog {
	return Catalog{
		Locale:   supportedLocales[1],
		Text:     portalApprovedText("zh-TW"),
		Pages:    portalApprovedPages("zh-TW"),
		Features: localizedFeatures(customerZHFeatures()),
		Docs:     localizedDocs(zhTWDocs()),
	}
}

func zhCNCatalog() Catalog {
	catalog := zhTWCatalog()
	catalog.Locale = supportedLocales[2]
	catalog.Text = portalApprovedText("zh-CN")
	catalog.Pages = portalApprovedPages("zh-CN")
	catalog.Features = simplifiedFeatures(catalog.Features)
	catalog.Docs = simplifiedDocs(catalog.Docs)
	return catalog
}

func zhTWFeatures() map[string]localizedFeature {
	return map[string]localizedFeature{
		"provision": {
			Title:        "Provision 配網",
			Kicker:       "以合約支撐的基礎來描述裝置導入。",
			Summary:      "雲端 registry 與啟用基礎已有合約邊界；本地 Wi-Fi/BLE 配網、claim UX、轉移/重設政策與完整產品就緒狀態仍屬整合或 roadmap 範圍。",
			Description:  "Provision 將 Realtek Connect+ 導入描述為從 registry 到雲端啟用裝置的分層流程。現階段可對外說明的是帳號 registry、跨服務 provision command、WebRTC Video over TURN 啟用邊界、service-scoped provisioning credential 與 transport readiness 合約；本地 Wi-Fi/BLE 設定、QR/SoftAP UX、所有權轉移、factory reset 政策與彙整產品就緒狀態，在負責 repo 完成前不描述為全面可用實作。",
			ImageAlt:     "顯示裝置 registry、導入狀態、儀表板面板與行動端輔助畫面的企業平台介面。",
			Highlights:   []string{"雲端 registry、啟用、service credential 與 transport readiness 的合約邊界", "QR、序號、activation code 與未來 factory identity 的 claim material 概念", "本地 Wi-Fi/BLE、SoftAP UX、轉移/重設政策與彙整 readiness 維持 roadmap 標示"},
			Capabilities: []string{"以帳號 registry API 與跨服務 DeviceProvisionRequested、DeviceProvisionSucceeded 或 DeviceProvisionFailed 事件描述雲端啟用基礎", "將 SDK claim parsing 與帳號端所有權及綁定政策分開", "把本地配網與產品就緒狀態描述為規劃中的 owner-repo 工作，而非網站已全面提供的功能"},
			Outcomes:     []string{"保留產品導入願景且不誇大實作狀態", "讓 firmware、App、SDK、account 與 WebRTC Video over TURN 團隊共用 availability 詞彙", "讓評估討論清楚區分目前可用、整合就緒與 roadmap"},
		},
		"ota": {
			Title:        "OTA 韌體更新",
			Kicker:       "區分韌體生命週期基礎與 campaign roadmap。",
			Summary:      "韌體上傳、catalog、target enablement、rollout 狀態、report、cancel 與 download 屬於現有基礎；scheduled、time-window、user-consent 與 archive campaign policy 已可用；approval workflow、dashboard、analytics 與分階段百分比 rollout 仍屬 roadmap 範圍。",
			Description:  "OTA 以 interface-first 韌體 campaign 路徑呈現。現階段公開文案可描述 upload、enablement、rollout query/report、cancel 與 download route 代表的韌體生命週期基礎；scheduled、time-window、user-consent 與 archive policy 已可視為可用的 campaign surface；approval workflow、dashboard、analytics 與分階段百分比 rollout 仍需標示為 roadmap。",
			ImageAlt:     "顯示韌體發布、裝置 registry、裝置群健康圖表與警示面板的企業營運主控台。",
			SourceLabel:  "韌體 campaign 介面合約",
			Highlights:   []string{"韌體上傳、catalog、target enablement、rollout 狀態、report、cancel 與 download 的現有基礎", "scheduled、time-window、user-consent 與 archive 屬於已可用的 campaign surface", "審核流程、dashboard、analytics、分階段百分比 rollout 與自動 cohort ramping 維持 roadmap 標示"},
			Capabilities: []string{"以既有 firmware route 作為可用實作邊界，而非暗示完整 campaign engine 已完成", "將 scheduled、time-window、user-consent 與 archive 行為視為已可用的 campaign policy vocabulary", "cancel 保持為生命週期基礎；approval workflow、dashboard、analytics 與分階段百分比 rollout 則標示為 roadmap 範圍"},
			Outcomes:     []string{"保留 OTA campaign 願景且不誇大第一階段實作", "讓 firmware、SDK、backend 與產品團隊共用 availability 詞彙", "讓評估討論清楚區分目前可用、已實作與 roadmap"},
		},
		"fleet-management": {
			Title:        "裝置群管理",
			Kicker:       "在產品上市後營運連網裝置。",
			Summary:      "節點註冊、憑證配置、裝置 registry、OTA 工作協調、批次操作與營運指標卡。",
			Description:  "裝置群管理描述產品團隊如何註冊節點、發放裝置身分、組織 registry、協調韌體作業並檢視整體營運狀態。",
			ImageAlt:     "顯示裝置 registry、裝置群健康圖表、發布狀態與營運警示的企業營運主控台。",
			Highlights:   []string{"節點註冊、憑證啟動與生命週期狀態", "群組、標籤、中繼資料、分享與批次操作", "OTA 工作、韌體映像與裝置健康指標"},
			Capabilities: []string{"註冊節點並綁定製造資料", "依型號、區域、韌體、客戶或 cohort 搜尋裝置", "檢視啟用數、韌體分布、警示佇列與支援工作"},
			Outcomes:     []string{"建立可信的營運平台敘事", "清楚區分網站 lead admin 與未來 IoT console", "串接配網、發布與支援流程"},
		},
		"smart-home": {
			Title:        "智慧家庭體驗",
			Kicker:       "讓使用者在導入後清楚控制產品。",
			Summary:      "遠端控制、本地控制備援、排程、情境、群組、裝置分享、推播通知與警示。",
			Description:  "智慧家庭體驗描述位於 Connect+ 雲端與 App 基礎上的消費端產品介面，涵蓋控制、自動化、分享與通知。",
			ImageAlt:     "顯示行動 App 情境、連網裝置面板與營運儀表板的企業平台介面。",
			Highlights:   []string{"日常操作的遠端與本地控制路徑", "排程、情境、群組與家庭分享", "推播通知與可行動警示"},
			Capabilities: []string{"裝置電源、模式、狀態與家庭情境控制", "週期排程、多裝置情境、房間或家庭群組", "導入完成、離線、異常與維護提醒"},
			Outcomes:     []string{"讓連網產品在首次設定後持續有用", "降低多使用者與自動化支援摩擦", "同時呈現營運平台與終端產品體驗"},
		},
		"user-management": {
			Title:        "使用者管理",
			Kicker:       "管理連網產品周邊的帳號生命週期。",
			Summary:      "註冊、登入、OTP、第三方登入、密碼復原、帳號變更與帳號刪除等平台能力展示。",
			Description:  "使用者管理描述連網產品常見的帳號生命週期能力。此網站目前不提供終端使用者登入或帳號管理實作；角色指派與權限政策屬於產品 authorization contract 範圍，不是此網站已提供的 live ACL 系統。",
			ImageAlt:     "顯示帳號、裝置、安全與裝置群管理面板的企業營運主控台。",
			Highlights:   []string{"品牌 App 的自助註冊與登入", "帳號啟用、復原與高風險操作的 OTP 驗證", "第三方登入與帳號連結路徑"},
			Capabilities: []string{"忘記密碼、變更密碼與 session 管理", "帳號刪除與保留流程", "使用者 profile、同意與裝置所有權狀態", "產品角色與權限維持在 account-side authorization contract，不由本網站宣稱已完成 ACL"},
			Outcomes:     []string{"縮短正式帳號系統規劃時間", "在架構審查中明確帳號範圍", "避免混淆平台能力與網站 lead capture"},
		},
		"app-sdk": {
			Title:        "App SDK",
			Kicker:       "更快打造品牌化行動體驗。",
			Summary:      "iOS/Android SDK、App 端與裝置端參考範例、推播規劃、rebrand 指引與 App 上架路徑。",
			Description:  "App SDK 將行動體驗定位為可品牌化、可擴充與可發布的產品介面，並透過 Android/iOS/WebApp 家庭 App 範例、Linux 裝置模擬器與 PRO2 裝置範例協助團隊先驗證 SDK 流程。",
			ImageAlt:     "顯示通用 App client、模擬器、參考裝置與中央平台節點的 Realtek Connect+ 範例生態系圖。",
			Highlights:   []string{"涵蓋導入、控制與帳號流程的 iOS/Android SDK", "Android、iOS、WebApp、Linux 模擬器與 PRO2 裝置範例", "推播、發布準備與上架指引"},
			Capabilities: []string{"登入、配網、裝置控制與分享的共用行動元件", "用 App 與裝置參考範例驗證跨端 SDK 用法", "App Store 與 Google Play 發布規劃"},
			Outcomes:     []string{"縮短行動與裝置整合驗證時程", "對齊 App、韌體與產品團隊的範例邊界", "減少每條產品線的一次性 App/雲端/裝置整合工作"},
			Sections: []features.FeatureSection{
				{Eyebrow: "SDK 基礎", Title: "不用重建連網產品堆疊，也能交付品牌化行動 App", Intro: "此頁將行動 SDK 定位為可重用的平台層，而不是籠統宣稱已經有完整 App runtime。", Items: []string{"透過 iOS 與 Android SDK 層涵蓋導入、認證、裝置控制與帳號連結 primitive。", "協助行動團隊把常見裝置模型、配網狀態與控制介面映射到產品專屬體驗。", "保持清楚邊界：此 Go 網站描述 App enablement 範圍，不是 native mobile client runtime。"}},
				{Eyebrow: "參考範例", Title: "參考範例應用先證明 SDK 用法，再進入產品整合", Intro: "範例生態系讓 App、韌體與產品團隊用具體的 App 端與裝置端 reference 驗證流程。", Items: []string{"以 rtk_cloud_client repository 作為 sample code、specification 與 sample README 的 source of truth；此網站只摘要客戶評估路徑，不 hosting SDK source。", "使用 Android 智慧家庭範例、iOS 智慧家庭範例與 WebApp Ops Lab 範例驗證配網、裝置列表/細節、燈具與空調控制、相機監看與 debug report。", "Linux 模擬器只驗證 command、state、report 與 snapshot metadata，不模擬 camera frame 或 WebRTC signaling。PRO2 host smoke 只驗證 adapter 與 signaling lifecycle；physical media 需使用實體 PRO2 device。", "這些是 SDK usage references，不是正式 app-store app 或 white-label release package。"}},
				{Eyebrow: "通知", Title: "把提醒與生命週期訊息當成 App 產品介面的一部分", Intro: "推播與 in-app notification 流程需要行動、雲端與支援能力協同規劃。", Items: []string{"圍繞導入完成、分享事件、OTA 提醒、警報與支援流程規劃推播通知。", "將通知 payload 連接到產品 authorization decision、裝置所有權狀態與支援升級路徑。", "讓產品團隊可依市場與合規需求調整通知語氣、品牌與偏好設定。"}},
				{Eyebrow: "發布", Title: "跨工程與產品團隊協調上架工作", Intro: "上架指引讓 App SDK 頁面與實際 release execution 連結，而不是停在 SDK 選型。", Items: []string{"協調 bundle identifier、簽章資產、商店 metadata、審核 checklist 與 App Store / Google Play staged rollout。", "用 contact path 讓 App 開發與產品團隊對齊品牌、release readiness 與後端能力範圍。", "產品團隊仍擁有 store ownership、privacy disclosure、crash monitoring 與 release approval。"}, Accent: true},
			},
			Table: features.FeatureTable{
				Eyebrow: "範例生態系",
				Title:   "橫跨 App 與裝置介面的參考範例應用",
				Intro:   "每個範例都協助客戶驗證特定 Realtek Connect+ 整合路徑，同時把正式 App 所有權與正式雲端 wire contract 分開。",
				Columns: []string{"範例", "端點", "驗證重點", "邊界"},
				Rows: []features.FeatureTableRow{
					{Cells: []string{"Android 智慧家庭範例", "App", "配網 adapter 狀態、裝置列表/細節、燈具與空調控制、相機監看與 debug report。", "Native Kotlin reference；不是 app-store deliverable。"}},
					{Cells: []string{"iOS 智慧家庭範例", "App", "Swift SDK 的設定 profile、裝置控制、相機邊界與 debug evidence。", "Native Swift reference；客戶產品團隊擁有 release UX 與簽章。"}},
					{Cells: []string{"WebApp Ops Lab 範例", "App", "雲端側導入、MQTT payload inspection、模擬控制、相機輔助流程與 debug report。", "瀏覽器 reference，不包含 BLE 或 SoftAP onboarding。"}},
					{Cells: []string{"Linux 模擬器", "裝置", "不需硬體即可驗證燈具、空調、meter、state、log、report 與 snapshot metadata。", "不模擬 camera frame 或 WebRTC signaling，也不是 media-capable peer。"}},
					{Cells: []string{"PRO2 裝置範例", "裝置", "Host smoke 驗證 adapter wiring 與 signaling lifecycle；實體硬體才能驗證 snapshot、camera/audio 與 answerer。", "Host smoke 不等於 physical media validation；vendor media call 仍由 firmware 擁有。"}},
				},
			},
			RelatedLinks: []features.FeatureRelatedLink{
				{Title: "Video Cloud", Summary: "將 RTK signaling 與 stored-video SDK workflow 接到產品觀看體驗。", Href: "/features/video-cloud"},
				{Title: "SDK 能力工作流", Summary: "在 SDK Manual 查看 package support 與整合邊界。", Href: "/manual/sdk/capability-workflows"},
			},
		},
		"video-cloud": {
			Title:        "Video Cloud",
			Kicker:       "Live WebRTC 與加密 stored video 是兩條獨立產品路徑。",
			Summary:      "清楚呈現雲端 signaling、ICE/TURN、session lifecycle、加密 clip、snapshot、playback URL，以及 App、SDK、Cloud 與 Device 權責。",
			Description:  "Realtek Connect+ 提供 Android 與 iOS SDK，用於雲端 signaling 與 stored-video workflow。Product app 將這些 SDK 與平台 WebRTC 及 media component 整合，交付最終觀看體驗。",
			ImageAlt:     "相機分別連接 Live signaling 與加密 stored-video 路徑，最後進入行動裝置觀看介面。",
			Highlights:   []string{"Live WebRTC 使用 HTTPS control API 與目前 device owner 的 MQTT 或 WebSocket transport", "Recording、stored clip 與 snapshot 使用獨立於 live viewing 的加密上傳及播放 lifecycle", "Live media frame 保持在 peer-to-peer 或 TURN relay 路徑，不會由雲端自動保存"},
			Capabilities: []string{"HTTPS ICE、session create、answer retrieval、close、expiry 與 failure-state helper", "透過 current-owner MQTT 或 WebSocket 傳送相同 webrtc_offer signaling payload；不 fan-out，也不自動 fallback 到非 owner transport", "Authorize、presigned upload、complete、list、thumbnail、短效播放、URL refresh 與 delete workflow"},
			Outcomes:     []string{"整合雲端 signaling，不需另建 signaling backend", "由 App 團隊掌握 platform WebRTC rendering 與產品 UX", "讓 live viewing、recording、clip 與 snapshot 在營運和技術上保持清楚區隔"},
			Flows: []features.FeatureFlow{
				{
					Eyebrow: "Live WebRTC",
					Title:   "單一 signaling lifecycle，media 由兩端負責",
					Intro:   "Realtek Connect+ 協調 signaling 與 ICE/TURN；App 與 Device 連接並呈現實際 media。",
					Steps: []features.FeatureFlowStep{
						{Title: "取得 ICE", Body: "App 透過 RTK SDK 呼叫 HTTPS ICE API。"},
						{Title: "建立 offer", Body: "App 的 platform WebRTC component 建立 SDP offer。"},
						{Title: "開啟 session", Body: "RTK SDK 透過 HTTPS 建立 signaling session。"},
						{Title: "傳送 offer", Body: "Cloud 經目前 device owner 的 MQTT 或 WebSocket transport 送出 webrtc_offer。"},
						{Title: "建立 answer", Body: "Device SDK 或 firmware 接收 offer、連接 camera/audio track 並建立 answer。"},
						{Title: "提交 answer", Body: "Device 經 HTTPS answer API 回覆 SDP answer。"},
						{Title: "協商 media", Body: "App 經 RTK SDK 取得 answer，完成 platform media negotiation。"},
						{Title: "關閉", Body: "App 或 Device 關閉 session；expiry 與 timeout 也會終止 stale session。"},
					},
				},
				{
					Eyebrow: "Stored video",
					Title:   "完整 media object 加密後直接上傳",
					Intro:   "Stored clip 由獨立 recording 與 upload workflow 建立；live session 絕不會自動變成 clip。",
					Steps: []features.FeatureFlowStep{
						{Title: "錄製", Body: "產品錄製完整 MP4 clip 或擷取 JPEG snapshot。"},
						{Title: "加密", Body: "Client 加密 media object，wrapped-key material 不寫入 log。"},
						{Title: "授權", Body: "RTK SDK 取得 upload lifecycle 與短效 presigned PUT URL。"},
						{Title: "上傳", Body: "Client 在 upload lifecycle 內執行一次 object-storage PUT。"},
						{Title: "完成", Body: "SDK 標記 upload complete，service 驗證後公開 ready object。"},
						{Title: "瀏覽與播放", Body: "App 進行 list、filter、pagination、thumbnail、短效 range URL refresh、playback 與 delete。"},
					},
				},
			},
			Sections: []features.FeatureSection{
				{Eyebrow: "Snapshot 與 Clip", Title: "依使用者動作選擇正確 media object", Intro: "兩者共享授權和安全交付概念，但不能互換。", Items: []string{"Snapshot 是用於預覽、alert 或單一時間點檢視的 JPEG 靜態圖；目前 technical default 上限為 5 MiB。", "Clip 是用於 recorded playback 的完整 MP4 media object；目前 technical default 上限為 256 MiB。", "除非產品明確錄製並執行 stored-video workflow，否則 Live WebRTC 不會產生任何 stored object。"}},
				{Eyebrow: "Retention 與 URL 預設值", Title: "將目前數值視為部署設定，不是價格或 SLA", Intro: "部署營運者選擇適用 storage policy；商務、region、backup 與 recovery 需求需另外確認。", Items: []string{"Retention profile 可依 deployment 設定為 1、7 或 30 天。", "Signed upload 與 playback URL 目前預設有效 10 分鐘；client 應更新 playback URL，不應永久保存。", "Authorized upload lifecycle 目前預設 30 分鐘，未完成時進入 failed 或 expired。", "本頁所有 limits 都是目前 technical defaults，不是價格、quota 承諾、backup 承諾、region 承諾或 SLA。"}},
				{Eyebrow: "Current release boundary", Title: "清楚知道產品整合仍負責什麼", Intro: "SDK 減少 cloud workflow boilerplate，同時保留明確的 media 與 UX ownership。", Items: []string{"目前版本不包含 server-side transcoding、S3 multipart upload、simulcast negotiation、renegotiation，也不提供完整的 SDK 內建 WebRTC media renderer。", "Product app 負責 platform WebRTC component、media renderer、audio policy、foreground/background lifecycle 與最終觀看 UX。", "Device SDK 或 firmware 負責 offer handling、answer generation、camera/audio tracks、codec 與裝置 resource limits。"}, Accent: true},
			},
			Table: features.FeatureTable{
				Eyebrow: "權責矩陣",
				Title:   "連接每一層，同時保留清楚邊界",
				Intro:   "同一組權責分工適用於產品架構、SDK 評估與支援除錯。",
				Columns: []string{"層級", "提供能力", "整合邊界"},
				Rows: []features.FeatureTableRow{
					{Cells: []string{"Realtek Connect+ Cloud", "HTTPS signaling、current-owner MQTT/WebSocket delivery、TURN credential、session state 與 stored-video lifecycle API。", "協調 control 與 storage workflow；不接收或保存 Live media frame。"}},
					{Cells: []string{"RTK Android/iOS SDK", "Authentication、ICE、offer/answer/close helper、stable error 與 clip workflow helper。", "將 session 與 media-object data 交給 product app；不包含完整 WebRTC renderer。"}},
					{Cells: []string{"Product app", "產品專屬 live 與 recorded viewing experience。", "整合 platform WebRTC、renderer、audio policy、lifecycle、controls、empty/error states 與 playback URL refresh。"}},
					{Cells: []string{"Device SDK / firmware", "Offer handling、answer generation、camera/audio tracks、recording、encryption 與 resource enforcement。", "負責 physical media validation、codec、track attachment 與 constrained-device behavior。"}},
				},
			},
			Tables: []features.FeatureTable{{
				Eyebrow: "Sample truth matrix",
				Title:   "依各 sample 實際能力進行驗證",
				Intro:   "Fixture 與 host-smoke sample 是有用的整合證據，但不等同 physical media validation。",
				Columns: []string{"Sample", "可見狀態", "驗證內容", "不代表"},
				Rows: []features.FeatureTableRow{
					{Cells: []string{"Android playback", "Real SDK + Media3", "Stored clip list 與 playback integration。", "Live WebRTC rendering。"}},
					{Cells: []string{"iOS playback", "Real SDK + AVPlayer", "Stored clip list 與 playback integration。", "Live WebRTC rendering。"}},
					{Cells: []string{"Android/iOS Live", "RTK signaling integration / fixture UI", "Session data、offer/answer、close 與 UI lifecycle。", "完整 media rendering 或 physical camera validation。"}},
					{Cells: []string{"WebApp Ops Lab", "Fixture-backed", "Signaling helper demonstration 與 operations workflow。", "Production WebRTC peer 或 native onboarding。"}},
					{Cells: []string{"Linux simulator", "Device workflow simulator", "Command、state、report 與 validation evidence。", "Camera frame 或 WebRTC signaling。"}},
					{Cells: []string{"PRO2 host smoke", "Adapter + lifecycle smoke", "Adapter wiring 與 signaling lifecycle。", "Physical camera、audio、codec 或 rendered-media validation。"}},
				},
			}},
			RelatedLinks: []features.FeatureRelatedLink{
				{Title: "App SDK", Summary: "查看 mobile integration surface 與 reference sample ecosystem。", Href: "/features/app-sdk"},
				{Title: "Capability workflows", Summary: "比較 SDK package capability 與權責邊界。", Href: "/manual/sdk/capability-workflows"},
				{Title: "Video workflows", Summary: "依照 live signaling 與 stored-video implementation guidance 整合。", Href: "/manual/sdk/video-workflows"},
			},
		},
		"insights": {
			Title:        "營運洞察",
			Kicker:       "看見 field 中產品的健康狀態。",
			Summary:      "啟用統計、韌體分布、crash report、log、重啟原因、RSSI 與記憶體訊號。",
			Description:  "Insights 讓工程與支援團隊掌握裝置群品質，透過營運統計與裝置健康訊號優先處理 field 問題。",
			ImageAlt:     "顯示裝置群健康圖表、韌體分布、遙測摘要與警示面板的企業營運主控台。",
			Highlights:   []string{"啟用與關聯統計", "Crash、重啟與 log 可視性", "韌體分布與裝置健康指標"},
			Capabilities: []string{"RSSI、記憶體、uptime 與重啟原因", "版本採用與發布健康", "支援導向的裝置歷史"},
			Outcomes:     []string{"更早發現 field 問題", "用證據支援客戶", "衡量韌體品質"},
		},
		"private-cloud": {
			Title:       "雲端方案與用量",
			Kicker:      "優先選擇 Realtek 託管服務，需要時再導入私有雲。",
			Summary:     "由 Realtek 建置、託管與維運，客戶按實際使用量付費；若有專屬治理需求，也可選擇客戶自有 Private Cloud。",
			Description: "Realtek 託管服務是建議的導入方式，產品團隊不必自行建置及維運雲端平台，即可開始使用 Realtek Connect+。Realtek 負責服務的託管、維護與營運，客戶依實際使用量付費。需要自有基礎架構、資料位置或治理邊界的組織，仍可選擇 Private Cloud。",
			ImageAlt:    "比較 Realtek 託管服務與客戶自有 Private Cloud 的 Realtek Connect+ 雲端方案。",
			Highlights:  []string{"推薦由 Realtek 建置並維運的託管服務", "依實際使用量彈性付費", "保留客戶自有 Private Cloud 作為第二選項"},
			Capabilities: []string{
				"由 Realtek 負責平台託管、維護與營運生命週期",
				"依實際使用量彈性付費的託管服務",
				"支援客戶雲端或地端環境的 Private Cloud 規劃",
			},
			Outcomes: []string{"不需先建立雲端維運團隊即可開始使用", "讓雲端成本與實際平台使用情況對齊", "需要企業治理時仍可轉向客戶自有環境"},
			Sections: []features.FeatureSection{
				{
					Eyebrow: "推薦方案",
					Title:   "Realtek 託管服務——使用多少，支付多少",
					Intro:   "由 Realtek 負責雲端平台營運，產品團隊可以專注在裝置、App 與客戶體驗。",
					Items: []string{
						"Realtek 負責雲端建置、託管、維護與日常平台營運。",
						"客戶依實際服務使用量彈性付費；詳細服務方案請洽詢 Realtek 團隊。",
						"適合希望快速開始、不想自行承擔底層雲端維運工作的產品團隊。",
					},
					Accent: true,
				},
				{
					Eyebrow: "Private Cloud",
					Title:   "部署在由客戶掌控的基礎架構",
					Intro:   "當客戶有專屬基礎架構、資料位置或治理要求時，可選擇 Private Cloud。",
					Items: []string{
						"使用標準 Container 或 VM workload 部署至客戶選擇的雲端或地端環境。",
						"基礎架構、網路政策、資料位置與營運邊界由客戶掌控。",
						"部署範圍、支援內容與商務條款透過個別 Private Cloud 導入案確認。",
					},
				},
			},
			Table: features.FeatureTable{
				Eyebrow: "雲端選項",
				Title:   "依團隊需求選擇營運模式",
				Intro:   "Realtek 託管服務是建議起點；當客戶控制權是首要考量時，仍可選擇 Private Cloud。",
				Columns: []string{"方案", "營運責任", "服務方式", "適合對象"},
				Rows: []features.FeatureTableRow{
					{Cells: []string{"Realtek 託管服務（推薦）", "由 Realtek 託管、維護與營運。", "依實際使用量彈性付費；詳細方案請洽 Realtek 團隊。", "希望快速導入且不自建雲端維運能力的團隊。"}},
					{Cells: []string{"Private Cloud", "客戶掌控基礎架構，並與 Realtek 約定支援邊界。", "依客戶需求制定部署與支援方案。", "具有專屬基礎架構、資料位置或治理需求的組織。"}},
				},
			},
		},
		"security": {
			Title:        "安全與 PKI",
			Kicker:       "以 X.509 憑證與受管理的 PKI 階層維繫裝置身分與雲端信任。",
			Summary:      "裝置憑證、兩層 CA 階層、mTLS 驗證、憑證生命週期作業，以及 OCSP / CRL 撤銷基礎，構成 Realtek Connect+ 的部署安全模型。",
			Description:  "安全與 PKI 說明 Realtek Connect+ 如何以建立在 X.509 憑證上的公開金鑰基礎架構來驗證裝置身分、保護雲端通訊，並在不自訂通訊協定的前提下支援大規模撤銷。每台裝置在佈建時都會取得由平台 CA 階層簽發的唯一憑證。雲端連線一律透過該憑證進行 mutual TLS 驗證，讓平台可以確認硬體身分、套用裝置端點政策，並以標準工具在大規模裝置群中輪替或撤銷憑證；裝置憑證不定義 human user roles，也不宣告產品 ACL 已可用。",
			ImageAlt:     "顯示裝置身分、PKI 信任邊界、雲端服務與整合端點的企業架構圖。",
			SourceLabel:  "平台安全合約",
			Highlights:   []string{"兩層 X.509 CA 階層：離線 root CA 與線上 issuing CA，負責裝置憑證簽發", "每台裝置在製造或首次啟用時取得並綁定硬體身分的專屬憑證", "所有裝置對雲端連線都使用 mutual TLS，雙向驗證 X.509 憑證", "憑證生命週期作業包含簽發、更新、輪替與撤銷，並支援 OCSP 與 CRL"},
			Capabilities: []string{"維持兩層 CA 階層，讓 root CA 保持離線，並由 issuing CA 在佈建流程中按需簽發裝置憑證", "將每張裝置憑證綁定序號、MAC 位址與型號，讓雲端無需共享密鑰即可驗證硬體身分", "在 MQTT 與 HTTPS 端點強制 mutual TLS，讓未通過驗證的裝置無法存取平台 API 或訊息代理", "將硬體身分政策與 human role assignment、產品 ACL 可用性清楚分開", "依憑證到期排程或安全事件支援更新與輪替，而不必完整重新佈建裝置", "發布 CRL 端點並執行 OCSP responder，讓依賴方可即時查驗憑證狀態", "可透過裝置群管理主控台撤銷單一裝置憑證，或批次撤銷遭入侵的製造批次"},
			Outcomes:     []string{"以個別 X.509 憑證作為每台裝置的身分基礎，降低共享密鑰風險", "用標準 PKI 產物滿足企業與營運團隊的安全審查", "讓撤銷保持快速且範圍明確，單一受損裝置不會暴露整個裝置群", "保留可追蹤的簽發與撤銷記錄，方便審計與支援追溯"},
		},
		"integrations": {
			Title:        "生態系整合",
			Kicker:       "把產品接入更廣泛的物聯網生態系。",
			Summary:      "Matter Fabric 定位、語音助理、MQTT over TLS、REST API 與 webhook 整合路徑。",
			Description:  "整合頁說明 Realtek Connect+ 如何與智慧家庭生態系與企業後端銜接，包含 Matter、語音助理、安全協定與 webhook 事件交付。Bearer token、MQTT topic scope 與 webhook signature 在此被描述為 service/integration credential，不是 human product role。",
			ImageAlt:     "顯示安全裝置連線、雲端服務、API、MQTT、webhook 與生態系端點的企業架構圖。",
			Highlights:   []string{"Matter 生態系與 Fabric 部署規劃", "語音助理、REST API、MQTT over TLS 與 webhook", "產品 App、雲端服務與客戶系統的權責邊界"},
			Capabilities: []string{"Matter bridge/controller 規劃與 commissioning touchpoint", "供產品、支援與營運系統使用的安全 REST/MQTT 介面", "把事件交付到 CRM、ticketing 與 analytics 的 webhook", "把 bearer token scope、MQTT topic scope 與 webhook signature 視為整合 credential，不當作使用者角色"},
			Outcomes:     []string{"符合互通性期待", "連接業務系統而不需一次性 glue code", "為平台評估保留可信整合範圍"},
		},
	}
}

func zhTWDocs() map[string]localizedDoc {
	return map[string]localizedDoc{
		"product-overview": {
			Title:        "產品總覽",
			Kicker:       "在深入實作前先定位平台範圍。",
			Summary:      "Realtek Connect+ 評估所需的平台架構、能力地圖與商業封裝指引。",
			Description:  "產品總覽協助團隊比較平台範圍、架構邊界與推出優先順序，說明韌體、雲端、App、營運與企業部署如何銜接。",
			Highlights:   []string{"平台架構敘事", "跨配網、OTA、App、洞察與雲端的能力地圖", "面向產品團隊的商業評估框架"},
			Deliverables: []string{"架構圖與生命週期摘要", "能力比較表", "硬體、行動與雲端利害關係人的評估指引"},
			Audience:     []string{"產品經理", "解決方案架構師", "技術業務團隊"},
		},
		"development": docZH("開發", "以一個交付計畫組織韌體、雲端與 App 工作流。"),
		"apis": {
			Title:        "API",
			Kicker:       "透過結構化整合介面開放雲端能力。",
			Summary:      "REST、MQTT over TLS、webhook 與服務合約文件入口。",
			Description:  "API 章節描述 Realtek Connect+ 周邊的合約層，讓 dashboard、支援工具、商業系統與裝置事件工作流能理解整合介面。Human product role 屬於 account-side authorization contract；bearer token scope 維持 service credential。",
			Highlights:   []string{"裝置、使用者、OTA 與 analytics 的 REST resource 類別", "MQTT over TLS 裝置訊息與狀態更新", "營運事件交付的 webhook pattern"},
			Deliverables: []string{"產品 authorization 邊界總覽", "Endpoint family 與 payload 預期", "支援與 CRM 工作流整合範例", "清楚區分 human role assignment、device credential 與 service bearer scope"},
			Audience:     []string{"整合雲端服務的 backend 團隊", "建立商業系統 hook 的 partner 工程師", "回答 API 範圍問題的技術帳戶團隊"},
		},
		"sdks": {
			Title:        "SDK",
			Kicker:       "記錄打造連網產品體驗所需的開發者介面。",
			Summary:      "行動 SDK、韌體 SDK、可重用 client 元件，以及用於產品驗證的參考範例。",
			Description:  "SDK 章節串接 App、韌體、裝置模擬與雲端整合工作，協助產品團隊在投入正式 App 或裝置整合前先驗證 Realtek Connect+ 流程。",
			Highlights:   []string{"iOS 與 Android 行動整合路徑", "韌體端服務與身分建構區塊", "Android、iOS、WebApp、Linux 模擬器與 PRO2 裝置參考範例"},
			Deliverables: []string{"依產品類型選擇 SDK 的指引", "品牌化體驗的客製化邊界", "展示配網、裝置列表/細節、燈具與空調控制、相機流程、debug report 與 MQTT payload 檢視的參考範例"},
			Audience:     []string{"整合品牌 App 的行動工程師", "對齊裝置端依賴的嵌入式工程師", "規劃重用與客製化範圍的專案負責人"},
			Table: docs.Table{
				Eyebrow: "參考範例",
				Title:   "SDK 範例矩陣",
				Intro:   "範例生態系用來證明 App 端與裝置端 SDK 用法，同時把正式 App 所有權、上架與正式雲端合約分開。",
				Columns: []string{"範例", "介面", "驗證內容"},
				Rows: []docs.TableRow{
					{Cells: []string{"Android 智慧家庭範例", "原生行動 App", "配網 adapter 狀態、裝置列表/細節、燈具與空調控制、相機監看、debug report 與 redacted evidence 收集。"}},
					{Cells: []string{"iOS 智慧家庭範例", "原生行動 App", "使用 Swift SDK 驗證相同智慧家庭流程，包含設定 profile、裝置控制、相機邊界與 debug evidence。"}},
					{Cells: []string{"WebApp Ops Lab 範例", "瀏覽器 App", "雲端側導入、裝置 registry 檢視、MQTT payload inspection、模擬控制、相機輔助流程與 debug report，不包含 BLE 或 SoftAP。"}},
					{Cells: []string{"Linux 裝置模擬器", "裝置參考", "驗證燈具、空調、meter、state、log、report 與 snapshot metadata；不模擬 camera frame 或 WebRTC signaling。"}},
					{Cells: []string{"PRO2 相機裝置範例", "裝置韌體參考", "Host smoke 驗證 adapter 與 signaling lifecycle；camera/audio media validation 需要實體硬體。"}},
				},
			},
		},
		"firmware": docZH("韌體", "釐清裝置軟體堆疊必須提供的能力。"),
		"cli":      docZH("CLI", "用可重複的命令列流程支援開發者與營運者。"),
		"deployment": {
			Title:        "部署",
			Kicker:       "記錄評估與商用部署背後的生產執行設定。",
			Summary:      "生產部署設定、持久化 SQLite 儲存、反向代理 TLS、健康檢查、備份/還原與回滾注意事項。",
			Description:  "部署說明如何在不把 TLS 終止搬進應用程式本身的前提下，以生產設定執行 Go 站台。內容涵蓋容器或 VM 執行模式、用來保存聯絡與分析資料的持久化 SQLite volume、健康檢查端點、邊緣快取預期，以及生產營運所需的備份/還原流程。",
			Highlights:   []string{"具備模板與靜態資產的容器或 VM 部署設定", "用於聯絡與分析事件的持久化 SQLite volume", "TLS 終止、快取標頭與探針由反向代理或託管邊緣負責", "生產發布的備份、還原與回滾檢查點"},
			Deliverables: []string{"生產部署檢查表", "SQLite 備份與還原作業手冊", "供營運人員使用的回滾與驗證步驟"},
			Audience:     []string{"規劃託管環境的平台團隊", "處理生產復原的營運團隊", "評估資料與復原邊界的企業買家"},
		},
		"release-notes": {
			Title:        "版本資訊",
			Kicker:       "追蹤韌體、雲端、App 與營運介面的產品演進。",
			Summary:      "版本化變更紀錄、升級說明、相容性聲明與發布溝通模式。",
			Description:  "版本資訊定義平台定期更新後團隊期待的文件結構，涵蓋產品版本變更、升級影響與相容性說明。",
			Highlights:   []string{"逐版本產品變更摘要", "升級注意事項與相容性提醒", "面向客戶的發布溝通結構"},
			Deliverables: []string{"雲端、App 與韌體版本資訊模板", "依受眾區分的升級影響", "歷史版本歸檔策略"},
			Audience:     []string{"客戶成功團隊", "工程團隊", "產品主管"},
		},
	}
}

func docZH(title, kicker string) localizedDoc {
	return localizedDoc{
		Title:        title,
		Kicker:       kicker,
		Summary:      "整理此工作流的目標、交付內容與團隊權責，協助 Realtek Connect+ 評估與導入。",
		Description:  kicker + " 本章節讓相關團隊能在評估早期對齊範圍、責任與後續實作入口。",
		Highlights:   []string{"工作流範圍", "跨團隊責任", "導入檢查點"},
		Deliverables: []string{"實作指引", "評估 checklist", "交付里程碑"},
		Audience:     []string{"工程團隊", "產品團隊", "解決方案與支援團隊"},
	}
}

func localizedFeatures(overrides map[string]localizedFeature) []features.Feature {
	base := features.All()
	out := make([]features.Feature, 0, len(base))
	for _, feature := range base {
		if item, ok := overrides[feature.Slug]; ok {
			feature.Title = item.Title
			feature.Kicker = item.Kicker
			feature.Summary = item.Summary
			feature.Description = item.Description
			feature.ImageAlt = item.ImageAlt
			if item.SourceLabel != "" {
				feature.SourceLabel = item.SourceLabel
			}
			if item.SourceURL != "" {
				feature.SourceURL = item.SourceURL
			}
			feature.Highlights = item.Highlights
			feature.Capabilities = item.Capabilities
			feature.Outcomes = item.Outcomes
			feature.Flows = item.Flows
			feature.Sections = item.Sections
			feature.Table = item.Table
			feature.Tables = item.Tables
			if len(item.RelatedLinks) > 0 {
				feature.RelatedLinks = item.RelatedLinks
			}
		}
		out = append(out, feature)
	}
	return out
}

func localizedDocs(overrides map[string]localizedDoc) []docs.Section {
	base := docs.All()
	out := make([]docs.Section, 0, len(base))
	for _, section := range base {
		if item, ok := overrides[section.Slug]; ok {
			section.Title = item.Title
			section.Kicker = item.Kicker
			section.Summary = item.Summary
			section.Description = item.Description
			section.Highlights = item.Highlights
			section.Deliverables = item.Deliverables
			section.Audience = item.Audience
			if len(item.Table.Rows) > 0 {
				section.Table = item.Table
			}
		}
		out = append(out, section)
	}
	return out
}

func simplifiedFeatures(input []features.Feature) []features.Feature {
	out := make([]features.Feature, len(input))
	for index, feature := range input {
		feature.Title = toSimplified(feature.Title)
		feature.Kicker = toSimplified(feature.Kicker)
		feature.Summary = toSimplified(feature.Summary)
		feature.Description = toSimplified(feature.Description)
		feature.ImageAlt = toSimplified(feature.ImageAlt)
		feature.SourceLabel = toSimplified(feature.SourceLabel)
		feature.Highlights = simplifiedSlice(feature.Highlights)
		feature.Capabilities = simplifiedSlice(feature.Capabilities)
		feature.Outcomes = simplifiedSlice(feature.Outcomes)
		feature.Flows = simplifiedFeatureFlows(feature.Flows)
		feature.Sections = simplifiedFeatureSections(feature.Sections)
		feature.Table = simplifiedFeatureTable(feature.Table)
		for tableIndex, table := range feature.Tables {
			feature.Tables[tableIndex] = simplifiedFeatureTable(table)
		}
		for linkIndex, link := range feature.RelatedLinks {
			link.Title = toSimplified(link.Title)
			link.Summary = toSimplified(link.Summary)
			feature.RelatedLinks[linkIndex] = link
		}
		out[index] = feature
	}
	return out
}

func simplifiedDocs(input []docs.Section) []docs.Section {
	out := make([]docs.Section, len(input))
	for index, section := range input {
		section.Title = toSimplified(section.Title)
		section.Kicker = toSimplified(section.Kicker)
		section.Summary = toSimplified(section.Summary)
		section.Description = toSimplified(section.Description)
		section.Highlights = simplifiedSlice(section.Highlights)
		section.Deliverables = simplifiedSlice(section.Deliverables)
		section.Audience = simplifiedSlice(section.Audience)
		section.Table = simplifiedDocsTable(section.Table)
		out[index] = section
	}
	return out
}

func simplifiedFeatureSections(input []features.FeatureSection) []features.FeatureSection {
	out := make([]features.FeatureSection, len(input))
	for index, section := range input {
		section.Eyebrow = toSimplified(section.Eyebrow)
		section.Title = toSimplified(section.Title)
		section.Intro = toSimplified(section.Intro)
		section.Items = simplifiedSlice(section.Items)
		out[index] = section
	}
	return out
}

func simplifiedFeatureFlows(input []features.FeatureFlow) []features.FeatureFlow {
	out := make([]features.FeatureFlow, len(input))
	for index, flow := range input {
		flow.Eyebrow = toSimplified(flow.Eyebrow)
		flow.Title = toSimplified(flow.Title)
		flow.Intro = toSimplified(flow.Intro)
		for stepIndex, step := range flow.Steps {
			step.Title = toSimplified(step.Title)
			step.Body = toSimplified(step.Body)
			flow.Steps[stepIndex] = step
		}
		out[index] = flow
	}
	return out
}

func simplifiedFeatureTable(input features.FeatureTable) features.FeatureTable {
	input.Eyebrow = toSimplified(input.Eyebrow)
	input.Title = toSimplified(input.Title)
	input.Intro = toSimplified(input.Intro)
	input.Columns = simplifiedSlice(input.Columns)
	for index, row := range input.Rows {
		row.Cells = simplifiedSlice(row.Cells)
		input.Rows[index] = row
	}
	return input
}

func simplifiedDocsTable(input docs.Table) docs.Table {
	input.Eyebrow = toSimplified(input.Eyebrow)
	input.Title = toSimplified(input.Title)
	input.Intro = toSimplified(input.Intro)
	input.Columns = simplifiedSlice(input.Columns)
	for index, row := range input.Rows {
		row.Cells = simplifiedSlice(row.Cells)
		input.Rows[index] = row
	}
	return input
}

func simplifiedSlice(input []string) []string {
	out := make([]string, len(input))
	for index, value := range input {
		out[index] = toSimplified(value)
	}
	return out
}

func toSimplified(value string) string {
	replacer := strings.NewReplacer(
		"環境", "环境", "規劃", "规划", "適合", "适合", "範圍", "范围",
		"總覽", "总览", "資訊", "信息", "系統", "系统", "架構", "架构",
		"彙整", "汇总", "彙", "汇", "劃", "划", "環", "环", "適", "适",
		"圍", "围", "總", "总", "資", "资", "統", "统", "構", "构",
		"持續", "持续", "連網", "联网", "彈性", "弹性", "詳細", "详细",
		"群組", "群组", "視覺", "视觉", "體驗", "体验", "可視性", "可视性",
		"並", "并", "組", "组", "覺", "觉", "細", "细", "彈", "弹", "連", "连", "當", "当",
		"客戶", "客户", "託管", "托管", "維護", "维护", "維運", "维运",
		"實際", "实际", "付費", "付费", "專屬", "专属", "建議", "建议",
		"優先", "优先", "計價單位", "计价单位", "正式費率", "正式费率",
		"另行確認", "另行确认", "Realtek 負責", "Realtek 负责", "生命週期", "生命周期",
		"治理邊界", "治理边界", "當客戶", "当客户", "額度", "额度",
		"雲", "云", "網", "网", "聯", "联", "體", "体", "繁體", "繁体", "簡體", "简体",
		"裝", "装", "置", "置", "與", "与", "導", "导", "入", "入", "韌", "韧",
		"檔", "档", "訊", "讯", "號", "号", "覽", "览", "現", "现", "實", "实",
		"區", "区", "註", "注", "冊", "册", "憑", "凭", "證", "证", "營", "营",
		"運", "运", "儀", "仪", "錶", "表", "錶", "表", "態", "态", "態", "态",
		"啟", "启", "關", "关", "係", "系", "擴", "扩", "階", "阶", "發", "发",
		"佈", "布", "發", "发", "佇", "队", "列", "列", "稱", "称", "權", "权",
		"隱", "隐", "穩", "稳", "務", "务", "產", "产", "團", "团", "隊", "队",
		"評", "评", "估", "估", "將", "将", "資料", "资料", "對", "对", "齊", "齐",
		"種", "种", "語", "语", "頁", "页", "點", "点", "選", "选", "擇", "择",
		"開", "开", "標", "标", "籤", "签", "檢", "检", "視", "视", "異", "异",
		"常", "常", "員", "员", "應", "应", "該", "该", "產", "产", "品", "品",
		"場", "场", "內", "内", "後", "后", "裡", "里", "說", "说", "讀", "读",
		"寫", "写", "為", "为", "這", "这", "個", "个", "從", "从", "過", "过",
		"還", "还", "讓", "让", "會", "会", "處", "处", "儲", "储", "請", "请",
		"謝", "谢", "錄", "录", "檢", "检", "欄", "栏", "選", "选", "擇", "择",
		"輸", "输", "討", "讨", "論", "论", "類", "类", "別", "别", "礎", "础",
		"絡", "络", "們", "们", "麼", "么", "蓋", "盖", "綁", "绑", "廠", "厂",
		"擁", "拥", "錄", "录", "狀", "状", "遙", "遥", "測", "测", "時", "时",
		"縮", "缩", "備", "备", "進", "进", "線", "线", "協", "协",
		"業", "业", "術", "术", "載", "载", "詢", "询", "問", "问", "聲", "声",
		"產品行銷改善", "产品营销改善", "服務品質觀察", "服务质量观察", "產品行銷", "产品营销", "服務品質", "服务质量", "單筆", "单笔",
		"單", "单", "聲明", "声明", "處理", "处理", "聯絡", "联络", "請求", "请求",
		"行銷", "营销", "品質", "质量", "銷", "销", "觀", "观", "質", "质", "筆", "笔", "無", "无", "於", "于", "沒", "没", "記", "记", "屬", "属",
		"預", "预", "期", "期", "個月", "个月", "訊息", "讯息", "嵌", "嵌",
		"瀏", "浏", "約", "约", "撐", "撑",
		"範", "范", "陣", "阵", "參", "参", "號", "号", "燈", "灯", "鏡", "镜",
		"智慧", "智能", "相機", "相机", "模擬", "模拟", "驗", "验",
	)
	return replacer.Replace(value)
}

func ToSimplified(value string) string {
	return toSimplified(value)
}
