package leads

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Lead struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
}

type LeadRecord struct {
	ID        int64
	Name      string
	Company   string
	Email     string
	Interest  string
	Message   string
	CreatedAt time.Time
}

type ListFilter struct {
	Email    string
	Company  string
	Interest string
}

type ListOptions struct {
	Filter ListFilter
	Limit  int
	Offset int
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Init() error {
	if _, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS leads (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  company TEXT,
  email TEXT NOT NULL,
  interest TEXT NOT NULL,
  message TEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`); err != nil {
		return err
	}

	for _, statement := range []string{
		leadValidationTriggerSQL("leads_validate_insert", "INSERT"),
		leadValidationTriggerSQL("leads_validate_update", "UPDATE"),
	} {
		if _, err := r.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Insert(ctx context.Context, lead Lead) error {
	lead = Normalize(lead)
	if validationErrs := Validate(lead); validationErrs != nil {
		return validationErrs
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO leads (name, company, email, interest, message)
VALUES (?, ?, ?, ?, ?)`,
		lead.Name,
		lead.Company,
		lead.Email,
		lead.Interest,
		lead.Message,
	)
	return err
}

func (r *Repository) Count(ctx context.Context, filter ListFilter) (int, error) {
	var count int
	query, args := leadListQuery(`SELECT COUNT(*) FROM leads`, filter, ListOptions{})
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *Repository) List(ctx context.Context, opts ListOptions) ([]LeadRecord, error) {
	opts = normalizeListOptions(opts)
	query, args := leadListQuery(`
SELECT id, name, company, email, interest, message, created_at
FROM leads`, opts.Filter, opts)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []LeadRecord
	for rows.Next() {
		var record LeadRecord
		var createdAt string
		if err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Company,
			&record.Email,
			&record.Interest,
			&record.Message,
			&createdAt,
		); err != nil {
			return nil, err
		}
		record.CreatedAt = parseSQLiteTime(createdAt)
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func normalizeListOptions(opts ListOptions) ListOptions {
	opts.Filter = ListFilter{
		Email:    strings.TrimSpace(opts.Filter.Email),
		Company:  strings.TrimSpace(opts.Filter.Company),
		Interest: strings.TrimSpace(opts.Filter.Interest),
	}
	if opts.Limit < 0 {
		opts.Limit = 100
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}
	return opts
}

func leadListQuery(base string, filter ListFilter, opts ListOptions) (string, []any) {
	clauses := make([]string, 0, 3)
	args := make([]any, 0, 5)

	appendLeadFilterClause(&clauses, &args, "email", filter.Email)
	appendLeadFilterClause(&clauses, &args, "company", filter.Company)
	appendLeadFilterClause(&clauses, &args, "interest", filter.Interest)

	var builder strings.Builder
	builder.WriteString(base)
	if len(clauses) > 0 {
		builder.WriteString("\nWHERE ")
		builder.WriteString(strings.Join(clauses, " AND "))
	}
	if strings.HasPrefix(base, "\nSELECT") || strings.HasPrefix(base, "SELECT") {
		builder.WriteString("\nORDER BY id DESC")
		if opts.Limit > 0 {
			builder.WriteString("\nLIMIT ? OFFSET ?")
			args = append(args, opts.Limit, opts.Offset)
		}
	}

	return builder.String(), args
}

func appendLeadFilterClause(clauses *[]string, args *[]any, column, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	*clauses = append(*clauses, "LOWER(COALESCE("+column+", '')) LIKE ? ESCAPE '\\'")
	*args = append(*args, likePattern(value))
}

func likePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + replacer.Replace(strings.ToLower(strings.TrimSpace(value))) + "%"
}

func parseSQLiteTime(value string) time.Time {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func leadValidationTriggerSQL(name, operation string) string {
	return `
CREATE TRIGGER IF NOT EXISTS ` + name + `
BEFORE ` + operation + ` ON leads
FOR EACH ROW
BEGIN
  SELECT CASE
    WHEN length(trim(COALESCE(NEW.name, ''))) = 0 THEN RAISE(ABORT, 'lead name is required')
    WHEN length(trim(COALESCE(NEW.name, ''))) > 120 THEN RAISE(ABORT, 'lead name exceeds 120 characters')
    WHEN length(trim(COALESCE(NEW.company, ''))) > 160 THEN RAISE(ABORT, 'lead company exceeds 160 characters')
    WHEN length(trim(COALESCE(NEW.email, ''))) = 0 THEN RAISE(ABORT, 'lead email is required')
    WHEN length(trim(COALESCE(NEW.email, ''))) > 254 THEN RAISE(ABORT, 'lead email exceeds 254 characters')
    WHEN length(trim(COALESCE(NEW.interest, ''))) = 0 THEN RAISE(ABORT, 'lead interest is required')
    WHEN length(trim(COALESCE(NEW.interest, ''))) > 120 THEN RAISE(ABORT, 'lead interest exceeds 120 characters')
    WHEN length(trim(COALESCE(NEW.message, ''))) > 2000 THEN RAISE(ABORT, 'lead message exceeds 2000 characters')
  END;
END;`
}
