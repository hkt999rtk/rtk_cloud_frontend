package leads

import (
	"strings"
	"testing"
)

func TestValidateCountsUnicodeCharacters(t *testing.T) {
	valid := Lead{
		Name:     strings.Repeat("界", NameMaxLength),
		Company:  strings.Repeat("公", CompanyMaxLength),
		Email:    "unicode@example.com",
		Interest: strings.Repeat("类", InterestMaxLength),
		Message:  strings.Repeat("文", MessageMaxLength),
	}

	if err := Validate(valid); err != nil {
		t.Fatalf("validate valid unicode lead: %v", err)
	}

	invalid := Lead{
		Name:     strings.Repeat("界", NameMaxLength+1),
		Company:  strings.Repeat("公", CompanyMaxLength+1),
		Email:    strings.Repeat("é", EmailMaxLength+1),
		Interest: strings.Repeat("类", InterestMaxLength+1),
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
	if errs["interest"] != "Interest must be 120 characters or fewer." {
		t.Fatalf("interest error = %q", errs["interest"])
	}
	if errs["message"] != "Message must be 2000 characters or fewer." {
		t.Fatalf("message error = %q", errs["message"])
	}
}
