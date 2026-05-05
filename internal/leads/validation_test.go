package leads

import (
	"strings"
	"testing"
)

func TestValidateCountsUnicodeCharacters(t *testing.T) {
	// Interest must be one of the canonical inquiry slugs, so unicode length
	// boundary checks for the interest field are not meaningful here. The
	// interest-not-in-allowed-set case is exercised separately.
	valid := Lead{
		Name:     strings.Repeat("界", NameMaxLength),
		Company:  strings.Repeat("公", CompanyMaxLength),
		Email:    "unicode@example.com",
		Interest: "other",
		Message:  strings.Repeat("文", MessageMaxLength),
	}

	if err := Validate(valid); err != nil {
		t.Fatalf("validate valid unicode lead: %v", err)
	}

	invalid := Lead{
		Name:     strings.Repeat("界", NameMaxLength+1),
		Company:  strings.Repeat("公", CompanyMaxLength+1),
		Email:    strings.Repeat("é", EmailMaxLength+1),
		Interest: "other",
		Message:  strings.Repeat("文", MessageMaxLength+1),
	}

	errs := Validate(invalid)
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if errs["name"] != "Name must be 120 characters or fewer." {
		t.Fatalf("name error = %q", errs["name"])
	}
	if errs["company"] != "Company must be 160 characters or fewer." {
		t.Fatalf("company error = %q", errs["company"])
	}
	if errs["email"] != "Email must be 254 characters or fewer." {
		t.Fatalf("email error = %q", errs["email"])
	}
	if errs["message"] != "Message must be 2000 characters or fewer." {
		t.Fatalf("message error = %q", errs["message"])
	}
}

func TestValidateRejectsInterestOutsideAllowedSet(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		expect string
	}{
		{"empty", "", "Select an area of interest."},
		{"feature-slug", "private-cloud", "Select an area of interest."},
		{"freeform", "I want a free trial", "Select an area of interest."},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			lead := Lead{
				Name:     "Test",
				Email:    "test@example.com",
				Interest: tc.value,
			}
			errs := Validate(lead)
			if errs == nil {
				t.Fatal("expected validation error for invalid interest")
			}
			if errs["interest"] != tc.expect {
				t.Fatalf("interest error = %q, want %q", errs["interest"], tc.expect)
			}
		})
	}

	for _, slug := range AllowedInterests {
		slug := slug
		t.Run("accepts/"+slug, func(t *testing.T) {
			lead := Lead{
				Name:     "Test",
				Email:    "test@example.com",
				Interest: slug,
			}
			if errs := Validate(lead); errs != nil && errs["interest"] != "" {
				t.Fatalf("rejected canonical interest %q: %v", slug, errs)
			}
		})
	}
}
