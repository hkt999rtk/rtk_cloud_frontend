package analytics

import "time"

type Event struct {
	TS             int64
	Type           string
	Page           string
	CTA            string
	Percent        *int
	Duration       *int
	Variant        string
	ReferrerOrigin string
	SessionID      string
	CreatedAt      time.Time
}
