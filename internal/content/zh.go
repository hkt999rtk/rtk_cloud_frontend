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
	Highlights   []string
	Capabilities []string
	Outcomes     []string
}

type localizedDoc struct {
	Title        string
	Kicker       string
	Summary      string
	Description  string
	Highlights   []string
	Deliverables []string
	Audience     []string
}

func zhTWCatalog() Catalog {
	return Catalog{
		Locale:   supportedLocales[1],
		Text:     zhTWText(),
		Pages:    zhTWPages(),
		Features: localizedFeatures(zhTWFeatures()),
		Docs:     localizedDocs(zhTWDocs()),
	}
}

func zhCNCatalog() Catalog {
	catalog := zhTWCatalog()
	catalog.Locale = supportedLocales[2]
	catalog.Text = simplifiedText(catalog.Text)
	catalog.Pages = simplifiedPages(catalog.Pages)
	catalog.Features = simplifiedFeatures(catalog.Features)
	catalog.Docs = simplifiedDocs(catalog.Docs)
	return catalog
}

func zhTWPages() map[string]PageMeta {
	return map[string]PageMeta{
		"home": {
			Title:       "Realtek Connect+ | 物聯網雲端平台",
			Description: "Realtek Connect+ 是面向 Realtek 裝置的物聯網雲端平台，涵蓋配網、OTA、裝置營運、App SDK、洞察、私有雲與整合能力。",
		},
		"features": {
			Title:       "功能服務 | Realtek Connect+",
			Description: "探索 Realtek Connect+ 為連網產品提供的配網、OTA、裝置管理、App SDK、洞察、私有雲與生態系整合能力。",
		},
		"docs": {
			Title:       "開發者文件 | Realtek Connect+",
			Description: "瀏覽 Realtek Connect+ 的產品總覽、開發、API、SDK、韌體、CLI、部署與版本資訊文件入口。",
		},
		"contact": {
			Title:       "聯絡我們 | Realtek Connect+",
			Description: "聯絡 Realtek Connect+ 團隊，討論配網、OTA、裝置營運、App SDK、洞察或私有雲評估。",
		},
		"privacy": {
			Title:       "隱私權聲明 | Realtek Connect+",
			Description: "了解 Realtek Connect+ 如何處理網站詢問、聯絡表單資料、保存期限、資料請求與本站託管影片。",
		},
	}
}

func zhTWText() map[string]string {
	return map[string]string{
		"skip.main":                  "跳到主要內容",
		"brand.home":                 "Realtek Connect+ 首頁",
		"nav.docs":                   "文件",
		"nav.features":               "功能",
		"nav.architecture":           "架構",
		"nav.contact":                "聯絡",
		"footer.tagline":             "面向 Realtek 連網產品的物聯網雲端平台概念網站。",
		"footer.docs":                "開發者文件",
		"footer.features":            "功能服務",
		"footer.contact":             "聯絡我們",
		"footer.privacy":             "隱私權",
		"home.eyebrow":               "為產品團隊打造的物聯網雲端平台",
		"home.lede":                  "透過配網、OTA、裝置營運、App SDK、洞察、私有雲與生態系整合，讓 Realtek 裝置更快進入可營運的連網產品生命週期。",
		"home.cta.primary":           "聯絡我們",
		"home.cta.secondary":         "探索服務",
		"home.chip.silicon":          "晶片",
		"home.chip.sdk":              "裝置 SDK",
		"home.chip.cloud":            "雲端",
		"home.chip.ops":              "營運",
		"home.overview.eyebrow":      "平台總覽",
		"home.overview.title":        "從晶片到連網產品生命週期。",
		"home.overview.lede":         "Connect+ 將 Realtek 連線晶片、韌體能力、雲端服務、行動 App 流程與營運工具整合成一條商用物聯網產品路徑。",
		"home.surfaces.eyebrow":      "平台介面",
		"home.surfaces.title":        "呈現完整產品系統，而不只是功能清單。",
		"home.surfaces.card.title":   "配網、OTA 與裝置健康使用一致的視覺語言。",
		"home.surfaces.card.body":    "網站以產品介面呈現裝置、雲端安全營運與儀表板工作流，讓平台價值更具體。",
		"home.surface.onboarding":    "導入配網",
		"home.surface.rollouts":      "韌體發布",
		"home.surface.insights":      "營運洞察",
		"home.surface.security":      "安全",
		"home.principles.eyebrow":    "設計原則",
		"home.principles.title":      "為企業級連網產品計畫而設計。",
		"home.principle.active":      "可管理性",
		"home.principle.scale":       "可擴展性",
		"home.principle.security":    "安全性",
		"home.principle.privacy":     "隱私",
		"home.principle.cost":        "成本控制",
		"home.principle.custom":      "可客製化",
		"home.principle.panel":       "用同一個平台敘事管理韌體、使用者、裝置群與支援流程。",
		"home.principle.body":        "Realtek Connect+ 以生命週期系統呈現：導入配網、雲端身分、OTA、App SDK、指標與企業部署相互銜接。",
		"home.services.eyebrow":      "核心服務",
		"home.services.title":        "為連網產品團隊封裝 Realtek 雲端能力。",
		"home.feature.details":       "查看細節",
		"home.arch.eyebrow":          "架構",
		"home.arch.title":            "從裝置導入到裝置群營運的直接流程。",
		"home.arch.device.title":     "Realtek 裝置 SDK",
		"home.arch.device.body":      "身分、配網、韌體服務與裝置訊號。",
		"home.arch.cloud.title":      "安全雲端",
		"home.arch.cloud.body":       "裝置登錄、OTA 活動、使用者裝置關聯與 API。",
		"home.arch.app.title":        "App SDK 與儀表板",
		"home.arch.app.body":         "行動配網、產品控制、洞察與支援流程。",
		"home.film.eyebrow":          "品牌基礎",
		"home.film.title":            "建立在 Realtek 的連網智慧之上。",
		"home.film.body":             "Realtek Connect+ 將半導體與連線技術基礎延伸為雲端平台敘事，協助產品團隊打造可商用規模化的連網裝置。",
		"home.film.cta":              "觀看品牌影片",
		"home.film.title.attr":       "Realtek 企業形象影片",
		"home.film.fallback":         "你的瀏覽器不支援 video 標籤。",
		"home.film.point.silicon":    "半導體技術基礎",
		"home.film.point.ecosystem":  "連網產品生態系",
		"home.film.point.enterprise": "企業部署信任",
		"home.deploy.eyebrow":        "公有評估與私有雲",
		"home.deploy.title":          "從評估開始，逐步走向受控部署。",
		"home.deploy.public":         "公有評估",
		"home.deploy.public.title":   "在投入私有部署前驗證產品適配性。",
		"home.deploy.public.body":    "透過公開網站與文件結構，讓韌體、行動、雲端與產品團隊在商用部署前對齊。",
		"home.deploy.docs":           "部署文件",
		"home.deploy.private":        "私有商用雲",
		"home.deploy.private.title":  "規劃資料所有權、自訂網域、區域部署與支援邊界。",
		"home.deploy.private.body":   "私有雲敘事讓企業買家看見從概念驗證到品牌化、受控且有商業支援營運的路徑。",
		"home.deploy.discuss":        "討論私有雲",
		"home.use.eyebrow":           "使用情境",
		"home.use.title":             "為商用連網裝置團隊打造。",
		"home.use.smart.body":        "App 配網、分享、推播通知、語音助理路徑與 OTA 維護。",
		"home.use.industrial.body":   "裝置群分組、中繼資料、安全更新、私有雲與營運可視性。",
		"home.use.appliance.body":    "長生命週期韌體維護、啟用資料、品牌 App 與支援診斷。",
		"home.docs.eyebrow":          "開發者入口",
		"home.docs.title":            "為產品、韌體、App 與雲端團隊建立共同文件主軸。",
		"home.docs.open":             "開啟章節",
		"home.cta.eyebrow":           "早期評估",
		"home.cta.title":             "規劃 Realtek Connect+ 產品路徑。",
		"home.cta.body":              "登記你對配網、OTA、私有雲、App SDK 或裝置群營運的需求。",
		"features.eyebrow":           "功能",
		"features.title":             "涵蓋完整物聯網產品生命週期的雲端服務。",
		"features.body":              "Realtek Connect+ 呈現商用連網裝置平台應具備的核心能力。",
		"features.open":              "開啟功能",
		"feature.discuss.prefix":     "討論",
		"feature.all":                "所有功能",
		"feature.highlights":         "重點",
		"feature.highlights.title":   "此服務涵蓋什麼",
		"feature.capabilities":       "能力",
		"feature.capabilities.title": "平台基礎能力",
		"feature.outcomes":           "成果",
		"feature.outcomes.title":     "產品團隊為何使用",
		"feature.next":               "下一步",
		"feature.cta.prefix":         "評估",
		"feature.cta.suffix":         "如何支援你的產品路線圖。",
		"feature.cta.body":           "分享你的產品類別、目標部署與雲端需求，讓 Realtek Connect+ 團隊協助評估。",
		"docs.eyebrow":               "開發者文件",
		"docs.title":                 "雲端、韌體、App 與部署團隊的文件入口。",
		"docs.body":                  "Realtek Connect+ 提供伺服器端渲染的文件入口，涵蓋平台總覽、開發、API、SDK、韌體、CLI、部署與版本資訊。",
		"docs.cta.primary":           "聯絡平台團隊",
		"docs.cta.secondary":         "查看 App 平台脈絡",
		"docs.portal.eyebrow":        "入口結構",
		"docs.portal.title":          "選擇符合你工作流的文件軌道。",
		"docs.why.eyebrow":           "為什麼重要",
		"docs.why.title":             "對齊產品上市前各團隊期待的文件介面。",
		"docs.shared.title":          "共同平台敘事",
		"docs.shared.body":           "讓產品、業務與工程利害關係人，在深入實作前先理解生命週期範圍。",
		"docs.depth.title":           "依工作流分層",
		"docs.depth.body":            "分開韌體、API、行動 SDK、部署與版本議題，讓各團隊直達自己的實作面。",
		"docs.static.title":          "靜態優先版本",
		"docs.static.body":           "維持 Go template 與伺服器端渲染架構相容，同時保留後續深入內容空間。",
		"doc.back":                   "返回文件",
		"doc.discuss":                "討論實作",
		"doc.coverage":               "涵蓋範圍",
		"doc.coverage.title":         "本章節說明內容",
		"doc.outputs":                "預期產出",
		"doc.outputs.title":          "團隊讀完後應能取得什麼",
		"doc.audience":               "主要讀者",
		"doc.audience.title":         "誰應該從這裡開始",
		"doc.next":                   "下一個章節",
		"doc.next.title":             "繼續瀏覽開發者入口。",
		"doc.view":                   "查看章節",
		"contact.eyebrow":            "聯絡",
		"contact.title":              "登記 Realtek Connect+ 評估需求。",
		"contact.body":               "告訴我們你的產品團隊最關注哪一項服務。第一版會將請求儲存在本機 SQLite。",
		"contact.context.eyebrow":    "早期評估",
		"contact.context.title":      "面向物聯網產品規劃、韌體維護與私有雲評估。",
		"contact.context.body":       "可用此表單討論配網、OTA、裝置群管理、App SDK、洞察、私有雲或整合需求。",
		"contact.thanks":             "謝謝",
		"contact.recorded":           "你的 Realtek Connect+ 請求已記錄。",
		"contact.review":             "查看功能",
		"contact.error.summary":      "送出前請檢查下列欄位。",
		"contact.website":            "網站",
		"contact.name":               "姓名",
		"contact.company":            "公司",
		"contact.email":              "Email",
		"contact.interest":           "關注服務",
		"contact.select":             "選擇服務",
		"contact.message":            "訊息",
		"contact.submit":             "送出需求",
		"contact.privacy":            "送出此表單即表示你理解我們會依 Realtek Connect+ 隱私權聲明處理你的詢問資料。",
		"contact.privacy.link":       "隱私權聲明",
		"privacy.eyebrow":            "隱私權",
		"privacy.title":              "Realtek Connect+ 網站詢問隱私權聲明。",
		"privacy.intro":              "第一版網站只收集回覆 Realtek Connect+ 商務詢問與早期評估請求所需的資訊。",
		"privacy.data.title":         "我們收集的資料",
		"privacy.data.body":          "聯絡表單可能收集姓名、公司、Email、關注服務與選填訊息。網站也會使用維運與疑難排解 HTTP 服務所需的基本伺服器紀錄。",
		"privacy.use.title":          "資料使用方式",
		"privacy.use.body":           "我們使用提交資料來回覆詢問、規劃產品討論、了解 Realtek Connect+ 服務需求，並保護網站免於 spam 或濫用。",
		"privacy.retention.title":    "保存期限",
		"privacy.retention.body":     "網站 leads 預期最多保存 24 個月，除非仍有進行中的商務討論或必要營運紀錄需要更長期間。",
		"privacy.rights.title":       "查詢、更正或刪除請求",
		"privacy.rights.body":        "如需查詢、更正或刪除已提交的詢問資料，請聯絡 privacy@example.com。正式公開前需將此 placeholder 信箱替換為正式隱私聯絡窗口。",
		"privacy.video.title":        "本站託管品牌影片",
		"privacy.video.body":         "首頁品牌影片以本站 local MP4 資產託管。影片播放器不會建立 YouTube iframe，也不會連線到 YouTube 服務。",
		"privacy.admin.title":        "內部存取",
		"privacy.admin.body":         "Lead review 以 admin token 保護。Admin 頁面不會放入 sitemap，並標示 noindex。",
		"privacy.legal.title":        "法務審閱",
		"privacy.legal.body":         "此聲明是網站 prototype 的 GDPR-lite 實作，不是完整法律合規套件，公開上線前應完成審閱。",
	}
}

func zhTWFeatures() map[string]localizedFeature {
	return map[string]localizedFeature{
		"provision": {
			Title:        "Provision 配網",
			Kicker:       "降低裝置導入與綁定摩擦。",
			Summary:      "為 Realtek 物聯網產品提供安全 Wi-Fi/BLE 配網、啟用與帳號綁定流程。",
			Description:  "Provision 讓產品團隊能從出廠硬體一路走到使用者擁有的連網裝置，涵蓋首次啟用、本地配網、雲端登錄與使用者裝置關聯。",
			ImageAlt:     "顯示手機配對、QR 導入與裝置啟用狀態卡片的配網儀表板。",
			Highlights:   []string{"Wi-Fi 與 BLE 配網流程", "裝置綁定與所有權轉移", "啟用狀態與首次遙測"},
			Capabilities: []string{"Claim token 與裝置身分交接", "App 導入流程中的使用者裝置關聯", "時區與中繼資料初始化"},
			Outcomes:     []string{"降低設定失敗率", "縮短 App 導入時間", "讓裝置準備進入營運管理"},
		},
		"ota": {
			Title:        "OTA 韌體更新",
			Kicker:       "以可控節奏發布韌體更新。",
			Summary:      "上傳韌體、擷取版本資料、指定分批發布目標，並管理強制、一般、排程與使用者控制的 OTA 工作。",
			Description:  "OTA 被定位為正式韌體營運介面，而不是單一更新按鈕。團隊可註冊韌體包、檢查版本與型號資料、設定發布政策並觀察活動狀態。",
			ImageAlt:     "顯示分階段發布時間軸、裝置群組與 OTA 分析的韌體發布控制中心。",
			Highlights:   []string{"韌體上傳與版本、型號、checksum、release note 擷取", "依版本、型號、區域與 cohort 進行分批發布", "支援強制、一般、排程、使用者控制與時間窗策略"},
			Capabilities: []string{"為常連線或間歇連線裝置設計動態 OTA 規則", "查看每個工作狀態、裝置結果、取消與封存紀錄", "發布前進行相容性驗證與操作核准"},
			Outcomes:     []string{"降低韌體支援成本", "降低全裝置群回歸風險", "協調消費與商用部署更新"},
		},
		"fleet-management": {
			Title:        "裝置群管理",
			Kicker:       "在產品上市後營運連網裝置。",
			Summary:      "節點註冊、憑證配置、裝置 registry、OTA 工作協調、批次操作與營運指標卡。",
			Description:  "裝置群管理描述產品團隊如何註冊節點、發放裝置身分、組織 registry、協調韌體作業並檢視整體營運狀態。",
			ImageAlt:     "顯示連線裝置群組、健康狀態、標籤與批次操作佇列的裝置群管理儀表板。",
			Highlights:   []string{"節點註冊、憑證啟動與生命週期狀態", "群組、標籤、中繼資料、分享與批次操作", "OTA 工作、韌體映像與裝置健康指標"},
			Capabilities: []string{"註冊節點並綁定製造資料", "依型號、區域、韌體、客戶或 cohort 搜尋裝置", "檢視啟用數、韌體分布、警示佇列與支援工作"},
			Outcomes:     []string{"建立可信的營運平台敘事", "清楚區分網站 lead admin 與未來 IoT console", "串接配網、發布與支援流程"},
		},
		"smart-home": {
			Title:        "智慧家庭體驗",
			Kicker:       "讓使用者在導入後清楚控制產品。",
			Summary:      "遠端控制、本地控制備援、排程、情境、群組、裝置分享、推播通知與警示。",
			Description:  "智慧家庭體驗描述位於 Connect+ 雲端與 App 基礎上的消費端產品介面，涵蓋控制、自動化、分享與通知。",
			ImageAlt:     "顯示智慧家庭 App 控制、裝置、情境、排程、群組與通知卡片的介面。",
			Highlights:   []string{"日常操作的遠端與本地控制路徑", "排程、情境、群組與家庭分享", "推播通知與可行動警示"},
			Capabilities: []string{"裝置電源、模式、狀態與家庭情境控制", "週期排程、多裝置情境、房間或家庭群組", "導入完成、離線、異常與維護提醒"},
			Outcomes:     []string{"讓連網產品在首次設定後持續有用", "降低多使用者與自動化支援摩擦", "同時呈現營運平台與終端產品體驗"},
		},
		"user-management": {
			Title:        "使用者管理",
			Kicker:       "管理連網產品周邊的帳號生命週期。",
			Summary:      "註冊、登入、OTP、第三方登入、密碼復原、帳號變更與帳號刪除等平台能力展示。",
			Description:  "使用者管理描述連網產品常見的帳號生命週期能力。此網站目前不提供終端使用者登入或帳號管理實作。",
			ImageAlt:     "顯示使用者資料、安全驗證、分享權限與帳號生命週期控制的管理介面。",
			Highlights:   []string{"品牌 App 的自助註冊與登入", "帳號啟用、復原與高風險操作的 OTP 驗證", "第三方登入與帳號連結路徑"},
			Capabilities: []string{"忘記密碼、變更密碼與 session 管理", "帳號刪除與保留流程", "使用者 profile、同意與裝置所有權狀態"},
			Outcomes:     []string{"縮短正式帳號系統規劃時間", "在架構審查中明確帳號範圍", "避免混淆平台能力與網站 lead capture"},
		},
		"app-sdk": {
			Title:        "App SDK",
			Kicker:       "更快打造品牌化行動體驗。",
			Summary:      "iOS/Android SDK、範例 App、推播通知、rebrand 指引與 App 上架路徑。",
			Description:  "App SDK 將行動體驗定位為可品牌化、可擴充與可發布的產品介面，協助團隊重用常見連網 App 能力。",
			ImageAlt:     "顯示 App 畫面、程式模組、推播區塊與上架 checklist 的行動 SDK 工作區。",
			Highlights:   []string{"涵蓋導入、控制與帳號流程的 iOS/Android SDK", "範例 App 與 rebrand 路徑", "推播、發布準備與上架指引"},
			Capabilities: []string{"登入、配網、裝置控制與分享的共用行動元件", "可設定、換皮或擴充的參考 App 結構", "App Store 與 Google Play 發布規劃"},
			Outcomes:     []string{"縮短行動 App 上市時間", "對齊 App 開發與產品權責", "減少每條產品線的一次性整合工作"},
		},
		"insights": {
			Title:        "營運洞察",
			Kicker:       "看見 field 中產品的健康狀態。",
			Summary:      "啟用統計、韌體分布、crash report、log、重啟原因、RSSI 與記憶體訊號。",
			Description:  "Insights 讓工程與支援團隊掌握裝置群品質，透過營運統計與裝置健康訊號優先處理 field 問題。",
			ImageAlt:     "顯示裝置群健康圖表、警示卡片與裝置遙測摘要的營運洞察儀表板。",
			Highlights:   []string{"啟用與關聯統計", "Crash、重啟與 log 可視性", "韌體分布與裝置健康指標"},
			Capabilities: []string{"RSSI、記憶體、uptime 與重啟原因", "版本採用與發布健康", "支援導向的裝置歷史"},
			Outcomes:     []string{"更早發現 field 問題", "用證據支援客戶", "衡量韌體品質"},
		},
		"private-cloud": {
			Title:        "私有雲",
			Kicker:       "以企業級所有權與控制部署。",
			Summary:      "比較公有評估與私有商用部署，涵蓋區域託管、自訂網域與升級規劃。",
			Description:  "私有雲說明 Connect+ 如何從共享評估走向專屬商用部署，將資料所有權、區域位置、自訂網域與升級規劃納入企業採購考量。",
			ImageAlt:     "顯示專屬區域、品牌網域入口與企業控制邊界的私有雲架構。",
			Highlights:   []string{"公有評估與專屬私有商用部署比較", "資料所有權、區域託管與自訂網域", "商用導入、升級路徑與部署支援"},
			Capabilities: []string{"為客戶營運或代管私有區域規劃專屬環境", "反向代理 TLS、網路政策與品牌服務端點", "從評估到生產的發布提升與維護窗口"},
			Outcomes:     []string{"符合企業採購需求", "明確所有權邊界", "建立從 pilot 到 production 的可信路徑"},
		},
		"integrations": {
			Title:        "生態系整合",
			Kicker:       "把產品接入更廣泛的物聯網生態系。",
			Summary:      "Matter Fabric 定位、語音助理、MQTT over TLS、REST API 與 webhook 整合路徑。",
			Description:  "整合頁說明 Realtek Connect+ 如何與智慧家庭生態系與企業後端銜接，包含 Matter、語音助理、安全協定與 webhook 事件交付。",
			ImageAlt:     "顯示 Matter、語音助理、REST API、MQTT over TLS、webhook、App 與企業系統端點的整合中心。",
			Highlights:   []string{"Matter 生態系與 Fabric 部署規劃", "語音助理、REST API、MQTT over TLS 與 webhook", "產品 App、雲端服務與客戶系統的權責邊界"},
			Capabilities: []string{"Matter bridge/controller 規劃與 commissioning touchpoint", "供產品、支援與營運系統使用的安全 REST/MQTT 介面", "把事件交付到 CRM、ticketing 與 analytics 的 webhook"},
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
		"apis":        docZH("API", "透過結構化整合介面開放雲端能力。"),
		"sdks":        docZH("SDK", "記錄打造連網產品體驗所需的開發者介面。"),
		"firmware":    docZH("韌體", "釐清裝置軟體堆疊必須提供的能力。"),
		"cli":         docZH("CLI", "用可重複的命令列流程支援開發者與營運者。"),
		"deployment":  docZH("部署", "說明評估環境如何成熟為商用雲端部署。"),
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
			feature.Highlights = item.Highlights
			feature.Capabilities = item.Capabilities
			feature.Outcomes = item.Outcomes
			feature.Sections = nil
			feature.Table = features.FeatureTable{}
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
		}
		out = append(out, section)
	}
	return out
}

func simplifiedText(input map[string]string) map[string]string {
	out := make(map[string]string, len(input))
	for key, value := range input {
		out[key] = toSimplified(value)
	}
	return out
}

func simplifiedPages(input map[string]PageMeta) map[string]PageMeta {
	out := make(map[string]PageMeta, len(input))
	for key, value := range input {
		value.Title = toSimplified(value.Title)
		value.Description = toSimplified(value.Description)
		out[key] = value
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
		feature.Highlights = simplifiedSlice(feature.Highlights)
		feature.Capabilities = simplifiedSlice(feature.Capabilities)
		feature.Outcomes = simplifiedSlice(feature.Outcomes)
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
		out[index] = section
	}
	return out
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
		"單", "单", "聲明", "声明", "處理", "处理", "聯絡", "联络", "請求", "请求",
		"預", "预", "期", "期", "個月", "个月", "訊息", "讯息", "嵌", "嵌",
		"瀏", "浏",
	)
	return replacer.Replace(value)
}

func ToSimplified(value string) string {
	return toSimplified(value)
}
