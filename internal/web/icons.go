package web

import "html/template"

func icon(name string) template.HTML {
	path, ok := iconPaths[name]
	if !ok {
		path = iconPaths["grid"]
	}
	return template.HTML(`<svg class="icon" viewBox="0 0 24 24" aria-hidden="true" focusable="false">` + path + `</svg>`)
}

var iconPaths = map[string]string{
	"api":         `<path d="M7 8h-1.5a3.5 3.5 0 0 0 0 7H7"/><path d="M17 8h1.5a3.5 3.5 0 0 1 0 7H17"/><path d="M8.5 12h7"/><path d="m10 8-2 8"/><path d="m14 8 2 8"/>`,
	"arrow-right": `<path d="M5 12h13"/><path d="m13 6 6 6-6 6"/>`,
	"chart":       `<path d="M4 19V5"/><path d="M4 19h16"/><path d="M8 15v-4"/><path d="M12 15V8"/><path d="M16 15v-6"/><path d="M7 7h2l2 3 2-5 2 8 2-4h2"/>`,
	"cloud-lock":  `<path d="M7 18h10a4 4 0 0 0 .4-7.98A6 6 0 0 0 6.2 8.4 4.8 4.8 0 0 0 7 18Z"/><rect x="9" y="12" width="6" height="5" rx="1"/><path d="M10.5 12v-1.2a1.5 1.5 0 0 1 3 0V12"/>`,
	"device":      `<rect x="7" y="3" width="10" height="18" rx="2"/><path d="M10 17h4"/><path d="M10 7h4"/>`,
	"document":    `<path d="M7 3h7l4 4v14H7z"/><path d="M14 3v5h5"/><path d="M10 12h6"/><path d="M10 16h6"/>`,
	"download":    `<path d="M12 4v10"/><path d="m8 10 4 4 4-4"/><path d="M5 20h14"/>`,
	"filter":      `<path d="M4 6h16"/><path d="M7 12h10"/><path d="M10 18h4"/>`,
	"fleet":       `<rect x="4" y="4" width="5" height="5" rx="1"/><rect x="15" y="4" width="5" height="5" rx="1"/><rect x="4" y="15" width="5" height="5" rx="1"/><rect x="15" y="15" width="5" height="5" rx="1"/><path d="M9 6.5h6"/><path d="M6.5 9v6"/><path d="M17.5 9v6"/><path d="M9 17.5h6"/>`,
	"grid":        `<rect x="4" y="4" width="6" height="6" rx="1"/><rect x="14" y="4" width="6" height="6" rx="1"/><rect x="4" y="14" width="6" height="6" rx="1"/><rect x="14" y="14" width="6" height="6" rx="1"/>`,
	"home":        `<path d="m4 11 8-7 8 7"/><path d="M6 10v10h12V10"/><path d="M10 20v-6h4v6"/><path d="M17 14h2"/>`,
	"mail":        `<rect x="4" y="6" width="16" height="12" rx="2"/><path d="m5 8 7 5 7-5"/>`,
	"nodes":       `<circle cx="6" cy="12" r="2"/><circle cx="18" cy="6" r="2"/><circle cx="18" cy="18" r="2"/><path d="m8 11 8-4"/><path d="m8 13 8 4"/>`,
	"ota":         `<path d="M7 14a5 5 0 0 0 8 2"/><path d="M17 10a5 5 0 0 0-8-2"/><path d="M15 16h3v3"/><path d="M9 8H6V5"/><path d="M12 11v6"/><path d="m9 14 3 3 3-3"/>`,
	"package":     `<path d="m12 3 8 4.5v9L12 21l-8-4.5v-9z"/><path d="m4 7.5 8 4.5 8-4.5"/><path d="M12 12v9"/><path d="m8 5.5 8 4.5"/>`,
	"phone-code":  `<rect x="6" y="3" width="12" height="18" rx="2"/><path d="M10 17h4"/><path d="m10 9-2 2 2 2"/><path d="m14 9 2 2-2 2"/>`,
	"provision":   `<rect x="5" y="6" width="9" height="12" rx="2"/><path d="M8 15h3"/><path d="M17 8.5a5 5 0 0 1 0 7"/><path d="M19.5 6a8 8 0 0 1 0 12"/><path d="M8 4h3"/>`,
	"refresh":     `<path d="M20 12a8 8 0 0 1-13.5 5.8"/><path d="M4 12A8 8 0 0 1 17.5 6.2"/><path d="M17 3v4h4"/><path d="M7 21v-4H3"/>`,
	"shield-user": `<path d="M12 3 5 6v5c0 4.4 2.8 7.8 7 10 4.2-2.2 7-5.6 7-10V6z"/><circle cx="12" cy="10" r="2"/><path d="M8.8 16a3.6 3.6 0 0 1 6.4 0"/>`,
	"terminal":    `<path d="m5 7 5 5-5 5"/><path d="M12 17h7"/>`,
	"upload":      `<path d="M12 20V7"/><path d="m7 12 5-5 5 5"/><path d="M5 20h14"/>`,
	"user-shield": `<path d="M12 3 5 6v5c0 4.4 2.8 7.8 7 10 4.2-2.2 7-5.6 7-10V6z"/><circle cx="12" cy="10" r="2"/><path d="M8.8 16a3.6 3.6 0 0 1 6.4 0"/>`,
	"wifi":        `<path d="M5 10a10 10 0 0 1 14 0"/><path d="M8 13a6 6 0 0 1 8 0"/><path d="M11 16a2 2 0 0 1 2 0"/><path d="M12 19h.01"/>`,
}
