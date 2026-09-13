package main

import (
	"encoding/json"
	"fmt"
	"os"

	"realtek-connect/internal/content"
)

func main() {
	manifest := content.TranslationManifest()
	payload := struct {
		SchemaVersion int                         `json:"schemaVersion"`
		SourceLocale  string                      `json:"sourceLocale"`
		PolicyVersion int                         `json:"policyVersion"`
		Strings       []content.TranslationSource `json:"strings"`
	}{SchemaVersion: 1, SourceLocale: "en", PolicyVersion: 1, Strings: manifest}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
