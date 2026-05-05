package manual

import "html/template"

type ManualIndex struct {
	Title       string
	Description string
	Sections    []ManualSection
}

type ManualSection struct {
	Slug    string
	Title   string
	Summary string
}

type ManualPage struct {
	Slug        string
	Title       string
	Description string
	BodyHTML    template.HTML
}
