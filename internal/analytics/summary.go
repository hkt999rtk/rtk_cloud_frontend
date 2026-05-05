package analytics

import (
	"context"
	"database/sql"
)

type Summary struct {
	PageViews int
	ClickCTAs int
	Scrolls   int
	Engaged   int
}

type ReferrerOriginCount struct {
	ReferrerOrigin string
	Count          int
}

type ScrollDistribution struct {
	Percent int
	Count   int
}

type CTAByPage struct {
	Page  string
	CTA   string
	Count int
}

func (r *Repository) Summary(ctx context.Context) (Summary, error) {
	if r == nil || r.db == nil {
		return Summary{}, nil
	}

	var summary Summary
	err := r.db.QueryRowContext(ctx, `
SELECT
  COALESCE(SUM(CASE WHEN event = 'page_view' THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN event = 'click_cta' THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN event = 'scroll' THEN 1 ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN event = 'engaged' THEN 1 ELSE 0 END), 0)
FROM analytics_events
`).Scan(&summary.PageViews, &summary.ClickCTAs, &summary.Scrolls, &summary.Engaged)
	if err != nil {
		return Summary{}, err
	}

	return summary, nil
}

func (r *Repository) TopReferrerOrigins(ctx context.Context, limit int) ([]ReferrerOriginCount, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 5
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT COALESCE(referrer_origin, ''), COUNT(*)
FROM analytics_events
WHERE event = 'click_cta'
GROUP BY referrer_origin
ORDER BY COUNT(*) DESC, COALESCE(referrer_origin, '') ASC
LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReferrerOriginCounts(rows)
}

func (r *Repository) ScrollDistribution(ctx context.Context) ([]ScrollDistribution, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT percent, COUNT(*)
FROM analytics_events
WHERE event = 'scroll'
GROUP BY percent
ORDER BY percent ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]ScrollDistribution, 0, 4)
	for rows.Next() {
		var row ScrollDistribution
		if err := rows.Scan(&row.Percent, &row.Count); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) CTAByPage(ctx context.Context, limit int) ([]CTAByPage, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT page, cta, COUNT(*)
FROM analytics_events
WHERE event = 'click_cta'
GROUP BY page, cta
ORDER BY COUNT(*) DESC, page ASC, cta ASC
LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]CTAByPage, 0, limit)
	for rows.Next() {
		var row CTAByPage
		if err := rows.Scan(&row.Page, &row.CTA, &row.Count); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func scanReferrerOriginCounts(rows *sql.Rows) ([]ReferrerOriginCount, error) {
	result := make([]ReferrerOriginCount, 0, 5)
	for rows.Next() {
		var row ReferrerOriginCount
		if err := rows.Scan(&row.ReferrerOrigin, &row.Count); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
