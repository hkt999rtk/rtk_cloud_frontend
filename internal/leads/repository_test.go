package leads

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepositoryInsert(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	if err := repo.Init(); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	if err := repo.Insert(context.Background(), Lead{
		Name:     "Kevin Huang",
		Company:  "Realtek",
		Email:    "kevin@example.com",
		Interest: "OTA",
		Message:  "Interested in rollout control.",
	}); err != nil {
		t.Fatalf("insert lead: %v", err)
	}

	count, err := repo.Count(context.Background())
	if err != nil {
		t.Fatalf("count leads: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	records, err := repo.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("list leads: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].Email != "kevin@example.com" {
		t.Fatalf("email = %q, want kevin@example.com", records[0].Email)
	}
	if records[0].CreatedAt.IsZero() {
		t.Fatal("created_at was not parsed")
	}
}
