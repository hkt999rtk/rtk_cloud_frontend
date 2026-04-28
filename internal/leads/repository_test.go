package leads

import (
	"context"
	"database/sql"
	"strings"
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
		Name:     "  Kevin Huang  ",
		Company:  "  Realtek  ",
		Email:    "  kevin@example.com  ",
		Interest: "  OTA  ",
		Message:  "  Interested in rollout control.  ",
	}); err != nil {
		t.Fatalf("insert lead: %v", err)
	}

	count, err := repo.Count(context.Background(), ListFilter{})
	if err != nil {
		t.Fatalf("count leads: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	records, err := repo.List(context.Background(), ListOptions{Limit: 10})
	if err != nil {
		t.Fatalf("list leads: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].Name != "Kevin Huang" {
		t.Fatalf("name = %q, want Kevin Huang", records[0].Name)
	}
	if records[0].Company != "Realtek" {
		t.Fatalf("company = %q, want Realtek", records[0].Company)
	}
	if records[0].Email != "kevin@example.com" {
		t.Fatalf("email = %q, want kevin@example.com", records[0].Email)
	}
	if records[0].Interest != "OTA" {
		t.Fatalf("interest = %q, want OTA", records[0].Interest)
	}
	if records[0].Message != "Interested in rollout control." {
		t.Fatalf("message = %q, want trimmed value", records[0].Message)
	}
	if records[0].CreatedAt.IsZero() {
		t.Fatal("created_at was not parsed")
	}
}

func TestRepositoryListSupportsFilteringAndPagination(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	if err := repo.Init(); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	for _, lead := range []Lead{
		{
			Name:     "Alpha",
			Company:  "Acme",
			Email:    "alpha@example.com",
			Interest: "Provision",
			Message:  "first",
		},
		{
			Name:     "Beta",
			Company:  "Acme Labs",
			Email:    "beta@example.com",
			Interest: "OTA",
			Message:  "second",
		},
		{
			Name:     "Gamma",
			Company:  "Zenith",
			Email:    "gamma@example.com",
			Interest: "OTA",
			Message:  "third",
		},
	} {
		if err := repo.Insert(context.Background(), lead); err != nil {
			t.Fatalf("insert lead: %v", err)
		}
	}

	count, err := repo.Count(context.Background(), ListFilter{
		Company:  "acme",
		Interest: "ota",
	})
	if err != nil {
		t.Fatalf("count filtered leads: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	records, err := repo.List(context.Background(), ListOptions{
		Filter: ListFilter{
			Interest: "ota",
		},
		Limit:  1,
		Offset: 1,
	})
	if err != nil {
		t.Fatalf("list filtered leads: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].Email != "beta@example.com" {
		t.Fatalf("email = %q, want beta@example.com", records[0].Email)
	}
}

func TestRepositoryInsertRejectsInvalidLead(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	if err := repo.Init(); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	err = repo.Insert(context.Background(), Lead{
		Name:     strings.Repeat("N", NameMaxLength+1),
		Company:  "Realtek",
		Email:    "kevin@example.com",
		Interest: "",
		Message:  strings.Repeat("M", MessageMaxLength+1),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	validationErrs, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("error type = %T, want ValidationErrors", err)
	}
	if validationErrs["name"] != "Name must be 120 characters or fewer." {
		t.Fatalf("name error = %q", validationErrs["name"])
	}
	if validationErrs["interest"] != "Select an area of interest." {
		t.Fatalf("interest error = %q", validationErrs["interest"])
	}
	if validationErrs["message"] != "Message must be 2000 characters or fewer." {
		t.Fatalf("message error = %q", validationErrs["message"])
	}

	count, err := repo.Count(context.Background(), ListFilter{})
	if err != nil {
		t.Fatalf("count leads: %v", err)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
}

func TestRepositoryInitAddsSQLiteValidationGuards(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	if err := repo.Init(); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
INSERT INTO leads (name, company, email, interest, message)
VALUES (?, ?, ?, ?, ?)`,
		"Kevin Huang",
		"Realtek",
		"kevin@example.com",
		"OTA",
		strings.Repeat("M", MessageMaxLength+1),
	)
	if err == nil {
		t.Fatal("expected sqlite validation error")
	}
	if !strings.Contains(err.Error(), "lead message exceeds 2000 characters") {
		t.Fatalf("error = %v, want message length guard", err)
	}

	count, err := repo.Count(context.Background(), ListFilter{})
	if err != nil {
		t.Fatalf("count leads: %v", err)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
}
