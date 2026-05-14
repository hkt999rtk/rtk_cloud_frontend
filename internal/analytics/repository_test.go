package analytics

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"realtek-connect/internal/leads"
)

func TestOpenInitializesSQLiteSchema(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), func() Config {
		cfg := analyticsTestConfig(dbPath)
		return cfg
	}())
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer repository.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("analytics database file was not created: %v", err)
	}

	db := openSQLiteTestDB(t, dbPath)
	defer db.Close()

	for _, name := range []string{
		"analytics_events",
		"idx_analytics_events_ts",
		"idx_analytics_events_event_page",
	} {
		var count int
		if err := db.QueryRowContext(context.Background(), `
SELECT COUNT(*)
FROM sqlite_master
WHERE name = ?`, name).Scan(&count); err != nil {
			t.Fatalf("lookup %s: %v", name, err)
		}
		if count != 1 {
			t.Fatalf("sqlite object %q not initialized", name)
		}
	}

	if got := repository.RetentionDays(); got != DefaultRetentionDays {
		t.Fatalf("retention days = %d, want %d", got, DefaultRetentionDays)
	}
}

func TestOpenDisabledSkipsStorageInitialization(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), Config{
		Enabled:      false,
		DatabasePath: dbPath,
	})
	if err != nil {
		t.Fatalf("open disabled analytics store: %v", err)
	}
	if repository != nil {
		t.Fatal("repository = non-nil, want nil when analytics is disabled")
	}

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatalf("analytics database file exists when disabled: %v", err)
	}
}

func TestAnalyticsStorageIsSeparateFromLeadStorage(t *testing.T) {
	dir := t.TempDir()
	leadPath := filepath.Join(dir, "connectplus.db")
	analyticsPath := filepath.Join(dir, "analytics.db")

	leadDB := openSQLiteTestDB(t, leadPath)
	defer leadDB.Close()

	leadRepository := leads.NewRepository(leadDB)
	if err := leadRepository.Init(); err != nil {
		t.Fatalf("initialize lead schema: %v", err)
	}
	if err := leadRepository.Insert(context.Background(), leads.Lead{
		Name:     "Ada",
		Company:  "Example",
		Email:    "ada@example.com",
		Interest: "evaluation-access",
		Message:  "Please follow up.",
	}); err != nil {
		t.Fatalf("insert lead: %v", err)
	}

	analyticsRepository, err := Open(context.Background(), analyticsTestConfig(analyticsPath))
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer analyticsRepository.Close()

	assertSQLiteObjectExists(t, leadDB, "leads")
	assertSQLiteObjectMissing(t, leadDB, "analytics_events")

	analyticsDB := openSQLiteTestDB(t, analyticsPath)
	defer analyticsDB.Close()

	assertSQLiteObjectExists(t, analyticsDB, "analytics_events")
	assertSQLiteObjectMissing(t, analyticsDB, "leads")
}

func TestInsertEventStoresAnalyticsRow(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), analyticsTestConfig(dbPath))
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer repository.Close()

	percent := 50
	duration := 10
	createdAt := time.Date(2026, 5, 5, 12, 34, 56, 0, time.UTC)
	if err := repository.InsertEvent(context.Background(), Event{
		TS:             1234567890,
		Type:           "click_cta",
		Page:           "home",
		CTA:            "contact_us",
		Percent:        &percent,
		Duration:       &duration,
		Variant:        "",
		ReferrerOrigin: "https://example.com",
		SessionID:      "session-123",
		CreatedAt:      createdAt,
	}); err != nil {
		t.Fatalf("insert event: %v", err)
	}

	db := openSQLiteTestDB(t, dbPath)
	defer db.Close()

	var (
		ts             int64
		eventType      string
		page           string
		cta            sql.NullString
		percentValue   sql.NullInt64
		durationValue  sql.NullInt64
		variant        sql.NullString
		referrerOrigin sql.NullString
		sessionID      string
		createdAtText  string
	)
	if err := db.QueryRowContext(context.Background(), `
SELECT ts, event, page, cta, percent, duration, variant, referrer_origin, session_id, created_at
FROM analytics_events
LIMIT 1`).Scan(&ts, &eventType, &page, &cta, &percentValue, &durationValue, &variant, &referrerOrigin, &sessionID, &createdAtText); err != nil {
		t.Fatalf("query inserted event: %v", err)
	}

	if ts != 1234567890 {
		t.Fatalf("ts = %d, want 1234567890", ts)
	}
	if eventType != "click_cta" || page != "home" {
		t.Fatalf("event/page = %q/%q, want click_cta/home", eventType, page)
	}
	if !cta.Valid || cta.String != "contact_us" {
		t.Fatalf("cta = %+v, want contact_us", cta)
	}
	if !percentValue.Valid || percentValue.Int64 != 50 {
		t.Fatalf("percent = %+v, want 50", percentValue)
	}
	if !durationValue.Valid || durationValue.Int64 != 10 {
		t.Fatalf("duration = %+v, want 10", durationValue)
	}
	if variant.Valid {
		t.Fatalf("variant = %+v, want NULL", variant)
	}
	if !referrerOrigin.Valid || referrerOrigin.String != "https://example.com" {
		t.Fatalf("referrer origin = %+v, want https://example.com", referrerOrigin)
	}
	if sessionID != "session-123" {
		t.Fatalf("session id = %q, want session-123", sessionID)
	}
	if createdAtText != "2026-05-05T12:34:56Z" {
		t.Fatalf("created_at = %q, want 2026-05-05T12:34:56Z", createdAtText)
	}
}

func TestCleanupExpiredEventsUsesDefaultRetention(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db := openSQLiteTestDB(t, dbPath)
	defer db.Close()

	repository := NewRepository(db, 0)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	cutoff := asOf.Add(-DefaultRetentionDays * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, cutoff-1)
	seedAnalyticsEvent(t, db, cutoff)
	seedAnalyticsEvent(t, db, cutoff+1)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted rows = %d, want 1", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, cutoff, cutoff+1)
}

func TestCleanupExpiredEventsUsesCustomRetention(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db := openSQLiteTestDB(t, dbPath)
	defer db.Close()

	repository := NewRepository(db, 14)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	cutoff := asOf.Add(-14 * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, cutoff-86400)
	seedAnalyticsEvent(t, db, cutoff)
	seedAnalyticsEvent(t, db, cutoff+86400)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted rows = %d, want 1", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, cutoff, cutoff+86400)
}

func TestCleanupExpiredEventsPreservesBoundaryTimestamp(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db := openSQLiteTestDB(t, dbPath)
	defer db.Close()

	repository := NewRepository(db, DefaultRetentionDays)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	boundary := asOf.Add(-DefaultRetentionDays * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, boundary)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 0 {
		t.Fatalf("deleted rows = %d, want 0", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, boundary)
}

func TestOpenCleansExpiredEventsOnStartup(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db := openSQLiteTestDB(t, dbPath)

	repository := NewRepository(db, DefaultRetentionDays)
	if err := repository.Init(); err != nil {
		_ = db.Close()
		t.Fatalf("initialize analytics schema: %v", err)
	}

	now := time.Now().UTC()
	seedAnalyticsEvent(t, db, now.Add(-91*24*time.Hour).Unix())
	seedAnalyticsEvent(t, db, now.Add(-89*24*time.Hour).Unix())

	if err := db.Close(); err != nil {
		t.Fatalf("close seeded sqlite: %v", err)
	}

	reopened, err := Open(context.Background(), func() Config {
		cfg := analyticsTestConfig(dbPath)
		cfg.RetentionDays = DefaultRetentionDays
		return cfg
	}())
	if err != nil {
		t.Fatalf("reopen analytics store: %v", err)
	}
	defer reopened.Close()

	verifyDB := openSQLiteTestDB(t, dbPath)
	defer verifyDB.Close()

	if got := countRows(t, verifyDB); got != 1 {
		t.Fatalf("analytics rows = %d, want 1", got)
	}
}

func TestSummaryQueriesAggregateAnalyticsData(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), analyticsTestConfig(dbPath))
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer repository.Close()

	now := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 4; i++ {
		if err := repository.InsertEvent(context.Background(), Event{
			TS:        now.Add(time.Duration(i) * time.Minute).Unix(),
			Type:      "page_view",
			Page:      "home",
			SessionID: "session-page-" + strconv.Itoa(i),
			CreatedAt: now,
		}); err != nil {
			t.Fatalf("seed page_view: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		referrer := "https://example.com"
		page := "home"
		cta := "contact_us"
		if i == 2 {
			referrer = ""
			page = "features"
			cta = "talk_to_sales"
		}
		if err := repository.InsertEvent(context.Background(), Event{
			TS:             now.Add(10 * time.Minute).Add(time.Duration(i) * time.Minute).Unix(),
			Type:           "click_cta",
			Page:           page,
			CTA:            cta,
			ReferrerOrigin: referrer,
			SessionID:      "session-click-" + strconv.Itoa(i),
			CreatedAt:      now,
		}); err != nil {
			t.Fatalf("seed click_cta: %v", err)
		}
	}
	for _, percent := range []int{25, 100} {
		p := percent
		if err := repository.InsertEvent(context.Background(), Event{
			TS:        now.Add(20 * time.Minute).Unix(),
			Type:      "scroll",
			Page:      "home",
			Percent:   &p,
			SessionID: "session-scroll-" + strconv.Itoa(percent),
			CreatedAt: now,
		}); err != nil {
			t.Fatalf("seed scroll: %v", err)
		}
	}
	if err := repository.InsertEvent(context.Background(), Event{
		TS:        now.Add(30 * time.Minute).Unix(),
		Type:      "engaged",
		Page:      "home",
		Duration:  intPtr(10),
		SessionID: "session-engaged",
		CreatedAt: now,
	}); err != nil {
		t.Fatalf("seed engaged: %v", err)
	}

	summary, err := repository.Summary(context.Background())
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.PageViews != 4 || summary.ClickCTAs != 3 || summary.Scrolls != 2 || summary.Engaged != 1 {
		t.Fatalf("summary = %+v", summary)
	}

	referrers, err := repository.TopReferrerOrigins(context.Background(), 5)
	if err != nil {
		t.Fatalf("top referrers: %v", err)
	}
	if len(referrers) != 2 {
		t.Fatalf("referrers = %+v, want 2 rows", referrers)
	}
	if referrers[0].ReferrerOrigin != "https://example.com" || referrers[0].Count != 2 {
		t.Fatalf("first referrer = %+v", referrers[0])
	}
	if referrers[1].ReferrerOrigin != "" || referrers[1].Count != 1 {
		t.Fatalf("second referrer = %+v", referrers[1])
	}

	scrolls, err := repository.ScrollDistribution(context.Background())
	if err != nil {
		t.Fatalf("scroll distribution: %v", err)
	}
	if len(scrolls) != 2 || scrolls[0].Percent != 25 || scrolls[0].Count != 1 || scrolls[1].Percent != 100 || scrolls[1].Count != 1 {
		t.Fatalf("scroll distribution = %+v", scrolls)
	}

	ctaRows, err := repository.CTAByPage(context.Background(), 10)
	if err != nil {
		t.Fatalf("cta by page: %v", err)
	}
	if len(ctaRows) != 2 {
		t.Fatalf("cta rows = %+v, want 2 rows", ctaRows)
	}
	if ctaRows[0].Page != "home" || ctaRows[0].CTA != "contact_us" || ctaRows[0].Count != 2 {
		t.Fatalf("first cta row = %+v", ctaRows[0])
	}
	if ctaRows[1].Page != "features" || ctaRows[1].CTA != "talk_to_sales" || ctaRows[1].Count != 1 {
		t.Fatalf("second cta row = %+v", ctaRows[1])
	}
}

func analyticsTestConfig(dbPath string) Config {
	return Config{
		Enabled:           true,
		DatabasePath:      dbPath,
		RetentionDays:     DefaultRetentionDays,
		UnsafeDisableSync: true,
	}
}

func openSQLiteTestDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if _, err := db.Exec(`PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;`); err != nil {
		_ = db.Close()
		t.Fatalf("configure sqlite test db: %v", err)
	}
	return db
}

func assertSQLiteObjectExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	if countSQLiteObjects(t, db, name) != 1 {
		t.Fatalf("sqlite object %q not initialized", name)
	}
}

func assertSQLiteObjectMissing(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	if countSQLiteObjects(t, db, name) != 0 {
		t.Fatalf("sqlite object %q should not exist", name)
	}
}

func seedAnalyticsEvent(t *testing.T, db *sql.DB, ts int64) {
	t.Helper()

	if _, err := db.ExecContext(context.Background(), `
INSERT INTO analytics_events (
  ts,
  event,
  page,
  created_at
) VALUES (?, ?, ?, ?)`,
		ts,
		"page_view",
		"home",
		time.Unix(ts, 0).UTC().Format(time.RFC3339),
	); err != nil {
		t.Fatalf("seed analytics event: %v", err)
	}
}

func assertAnalyticsEventTimestamps(t *testing.T, db *sql.DB, want ...int64) {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), `
SELECT ts
FROM analytics_events
ORDER BY ts ASC`)
	if err != nil {
		t.Fatalf("query analytics timestamps: %v", err)
	}
	defer rows.Close()

	var got []int64
	for rows.Next() {
		var ts int64
		if err := rows.Scan(&ts); err != nil {
			t.Fatalf("scan analytics timestamp: %v", err)
		}
		got = append(got, ts)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate analytics timestamps: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("timestamps = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("timestamps = %v, want %v", got, want)
		}
	}
}

func countSQLiteObjects(t *testing.T, db *sql.DB, name string) int {
	t.Helper()

	var count int
	if err := db.QueryRowContext(context.Background(), `
SELECT COUNT(*)
FROM sqlite_master
WHERE name = ?`, name).Scan(&count); err != nil {
		t.Fatalf("lookup %s: %v", name, err)
	}
	return count
}

func countRows(t *testing.T, db *sql.DB) int {
	t.Helper()

	var count int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM analytics_events`).Scan(&count); err != nil {
		t.Fatalf("count analytics rows: %v", err)
	}
	return count
}

func intPtr(value int) *int {
	return &value
}
