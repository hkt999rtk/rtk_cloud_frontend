package leads

import (
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	NameMaxLength     = 120
	CompanyMaxLength  = 160
	EmailMaxLength    = 254
	InterestMaxLength = 120
	MessageMaxLength  = 2000
)

type ValidationErrors map[string]string

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	parts := make([]string, 0, len(e))
	for field, message := range e {
		parts = append(parts, field+": "+message)
	}
	sort.Strings(parts)
	return "invalid lead: " + strings.Join(parts, "; ")
}

func Normalize(lead Lead) Lead {
	lead.Name = strings.TrimSpace(lead.Name)
	lead.Company = strings.TrimSpace(lead.Company)
	lead.Email = strings.TrimSpace(lead.Email)
	lead.Interest = strings.TrimSpace(lead.Interest)
	lead.Message = strings.TrimSpace(lead.Message)
	return lead
}

func Validate(lead Lead) ValidationErrors {
	lead = Normalize(lead)

	errors := ValidationErrors{}
	if lead.Name == "" {
		errors["name"] = "Name is required."
	} else if fieldLength(lead.Name) > NameMaxLength {
		errors["name"] = "Name must be 120 characters or fewer."
	}

	if lead.Email == "" {
		errors["email"] = "Email is required."
	} else if fieldLength(lead.Email) > EmailMaxLength {
		errors["email"] = "Email must be 254 characters or fewer."
	}

	if lead.Interest == "" {
		errors["interest"] = "Select an area of interest."
	} else if fieldLength(lead.Interest) > InterestMaxLength {
		errors["interest"] = "Interest must be 120 characters or fewer."
	}

	if fieldLength(lead.Company) > CompanyMaxLength {
		errors["company"] = "Company must be 160 characters or fewer."
	}

	if fieldLength(lead.Message) > MessageMaxLength {
		errors["message"] = "Message must be 2000 characters or fewer."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

func fieldLength(value string) int {
	return utf8.RuneCountInString(value)
}
