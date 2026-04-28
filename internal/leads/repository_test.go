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
}
