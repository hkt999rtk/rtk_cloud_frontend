package content

import "realtek-connect/internal/features"

func customerZHFeatures() map[string]localizedFeature {
	return map[string]localizedFeature{
		"provision": {
			Title: "Provision 配網", Kicker: "讓裝置快速、安全地完成配網與雲端註冊。",
			Summary:      "建立從首次網路設定到安全雲端註冊的一致導入流程。",
			Description:  "從首次設定到雲端註冊，Realtek Connect+ Provision 協助產品團隊建立一致且安全的裝置導入流程。裝置完成註冊後，即可銜接 OTA 韌體更新、裝置群管理與其他雲端服務；實際配網方式與整合範圍由 Realtek 團隊依產品需求共同規劃。",
			ImageAlt:     "去識別化的 Realtek Connect+ 裝置註冊操作介面",
			Highlights:   []string{"引導使用者完成清楚的首次設定", "在雲端建立可信任的裝置身分", "完成後直接銜接更新與裝置群營運"},
			Capabilities: []string{"依產品需求規劃 Wi-Fi、Bluetooth、QR Code 或其他合適的配網方式", "將裝置註冊與產品需要的雲端服務串連", "讓產品與支援團隊掌握一致的裝置導入狀態"},
			Outcomes:     []string{"降低裝置首次使用的操作阻力", "讓不同產品線重複使用一致的導入經驗", "為每台已註冊裝置建立持續營運的起點"},
			Flows:        []features.FeatureFlow{{Eyebrow: "客戶流程", Title: "從首次設定到雲端服務的三個步驟", Intro: "實際體驗會依產品共同規劃，同時維持簡單清楚的客戶旅程。", Steps: []features.FeatureFlowStep{{Title: "完成裝置首次網路設定", Body: "使用最適合產品的方式，協助裝置連上網路。"}, {Title: "建立裝置身分並註冊至雲端", Body: "建立可信任的裝置紀錄，讓雲端服務能辨識及管理裝置。"}, {Title: "銜接 OTA、裝置群管理與其他服務", Body: "註冊完成後，即可接續產品生命週期所需的雲端能力。"}}}},
		},
		"ota":              zhCustomerFeature("OTA 韌體更新", "讓產品上市後也能安全持續更新。", "規劃韌體版本、分批發布、掌握進度，並處理需要注意的裝置。", "Realtek Connect+ OTA 協助產品與營運團隊有計畫地發布韌體。你可以選擇目標裝置、分階段推送，並在更新過程中掌握整體進度。", "去識別化的韌體版本與 OTA 營運介面", []string{"規劃可控的韌體發布", "依裝置群組或時程分批更新", "掌握進度並處理異常狀況"}),
		"fleet-management": zhCustomerFeature("裝置群管理", "讓每一批裝置都井然有序。", "在同一個營運畫面整理裝置、群組、健康狀態、告警與批次操作。", "Realtek Connect+ 裝置群管理協助營運團隊整理連網產品、找出需要注意的裝置，並針對裝置群執行一致的日常操作。", "去識別化的連網裝置群營運介面", []string{"依產品與營運需求整理裝置群組", "集中查看健康狀態與告警", "對裝置群執行一致的操作"}),
		"smart-home":       zhCustomerFeature("智慧家庭", "讓連網家庭體驗更容易使用。", "簡化首次設定、日常控制、家庭分享與自動化體驗。", "Realtek Connect+ 協助產品團隊打造從首次設定到日常使用都一致的智慧家庭體驗，並保留品牌所需的產品流程與整合彈性。", "智慧家庭設定、控制、分享與自動化體驗", []string{"簡化裝置設定與控制", "支援家庭成員分享", "串連常用的自動化情境"}),
		"user-management":  zhCustomerFeature("使用者管理", "清楚管理使用者與裝置存取。", "隨產品成長整理帳戶、裝置關聯、分享與存取權。", "Realtek Connect+ 使用者管理協助團隊理解客戶帳戶與裝置之間的關係，涵蓋首次擁有、分享及後續支援。", "去識別化的使用者與存取管理介面", []string{"串連客戶帳戶與裝置", "支援分享與存取權設定", "讓支援團隊更容易掌握使用情況"}),
		"app-sdk":          zhCustomerFeature("App SDK", "縮短品牌 App 的整合與上市時間。", "使用 SDK 與參考流程，更快完成配網、控制與帳戶體驗。", "Realtek Connect+ App SDK 協助行動與 Web 團隊整合連網產品的核心體驗，不必從頭建立每一個雲端互動流程。", "去識別化的 SDK 與 App 整合介面", []string{"加速配網與裝置控制體驗", "重複使用一致的帳戶與裝置流程", "更快從評估走向品牌產品"}),
		"video-cloud":      zhCustomerFeature("Video Cloud", "將即時影像融入產品體驗。", "串連裝置、安全雲端觀看與 App 操作流程。", "Realtek Connect+ Video Cloud 協助連網相機產品串連裝置、雲端與客戶 App，建立清楚的即時觀看與營運流程。", "去識別化的即時影像營運介面", []string{"串連裝置與雲端觀看流程", "支援產品的即時影像體驗", "讓團隊掌握影像服務營運狀態"}),
		"insights":         zhCustomerFeature("營運洞察", "把裝置訊號變成可行動的資訊。", "掌握裝置健康、韌體分布與產品營運趨勢。", "Realtek Connect+ 營運洞察協助產品與營運團隊了解連網裝置的表現，並聚焦需要進一步處理的狀況。", "去識別化的裝置報表與營運洞察介面", []string{"觀察裝置健康趨勢", "了解韌體版本分布", "讓團隊優先處理重要問題"}),
		"integrations":     zhCustomerFeature("系統整合", "將雲端連接至既有產品服務。", "透過 API 與產品事件，串連團隊已在使用的系統。", "Realtek Connect+ 提供務實的整合方式，協助裝置營運連接商業系統、應用程式與產品服務；實際整合方式由 Realtek 團隊依產品需求共同規劃。", "連接應用程式與產品服務的雲端整合架構", []string{"連接既有產品服務", "透過 API 支援需要的工作流程", "依產品架構共同規劃整合方式"}),
		"security":         zhCustomerFeature("安全與 PKI", "建立可信任的裝置到雲端安全基礎。", "透過裝置身分、PKI 與存取控制保護連網產品。", "Realtek Connect+ 以可信任的裝置身分、安全憑證與存取控制為基礎，協助團隊保護雲端連線及重要營運操作。", "裝置身分與安全雲端存取架構", []string{"建立可信任的裝置身分", "使用憑證與 PKI 安全基礎", "管理重要產品操作的存取權"}),
		"private-cloud": {
			Title: "雲端方案與額度", Kicker: "優先選擇 Realtek 託管服務，需要時再規劃 Private Cloud。",
			Summary:      "由 Realtek 建置與維運，依實際使用量彈性付費；專屬環境可另行規劃 Private Cloud。",
			Description:  "Realtek 託管服務是建議的開始方式。Realtek 負責平台建置、託管與持續維運，客戶可依實際使用量彈性付費。若產品需要專屬基礎架構、資料位置或治理邊界，Realtek 團隊也可共同規劃 Private Cloud。",
			ImageAlt:     "Realtek 託管服務與 Private Cloud 雲端方案",
			Highlights:   []string{"推薦由 Realtek 建置與維運的託管服務", "依實際使用量彈性付費", "需要時可規劃專屬 Private Cloud"},
			Capabilities: []string{"Realtek 負責平台建置與持續營運", "在 Portal 查看用量、費用、發票與付款", "依專屬需求共同規劃 Private Cloud"},
			Outcomes:     []string{"不用先建立雲端維運團隊即可開始", "讓服務使用與成本更容易掌握", "保留未來導入專屬環境的彈性"},
		},
	}
}

func zhCustomerFeature(title, kicker, summary, description, imageAlt string, highlights []string) localizedFeature {
	return localizedFeature{Title: title, Kicker: kicker, Summary: summary, Description: description, ImageAlt: imageAlt, Highlights: highlights,
		Capabilities: []string{"依產品與客戶需求規劃使用體驗", "將裝置與雲端營運串成一致流程", "由 Realtek 團隊共同規劃適合的整合範圍"},
		Outcomes:     []string{"縮短從產品整合到日常營運的路徑", "讓產品與營運團隊共享清楚的資訊", "保留隨產品持續擴充的彈性"}}
}
