package api

import (
	"encoding/json"
	"testing"
)

func TestParseKSIStructure(t *testing.T) {
	client := NewClient()

	// Fetch the actual KSI document
	data, err := client.fetchDocument("FRMR.KSI.key-security-indicators.json")
	if err != nil {
		t.Fatalf("Failed to fetch KSI document: %v", err)
	}

	t.Logf("Fetched %d bytes of KSI data", len(data))

	// Try to parse just the raw structure first
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		t.Fatalf("Failed to unmarshal raw doc: %v", err)
	}

	t.Logf("Top-level keys: %v", getKeys(rawDoc))

	// Check if KSI key exists
	ksiRaw, ok := rawDoc["KSI"]
	if !ok {
		t.Fatal("KSI key not found at top level")
	}

	t.Logf("KSI raw data length: %d bytes", len(ksiRaw))

	// Parse the KSI section
	var ksiMap map[string]json.RawMessage
	if err := json.Unmarshal(ksiRaw, &ksiMap); err != nil {
		t.Fatalf("Failed to unmarshal KSI section: %v", err)
	}

	t.Logf("KSI theme keys: %v", getKeys(ksiMap))

	// Try parsing one theme
	if afrRaw, ok := ksiMap["AFR"]; ok {
		var theme ThemeJSON
		if err := json.Unmarshal(afrRaw, &theme); err != nil {
			t.Fatalf("Failed to unmarshal AFR theme: %v", err)
		}
		t.Logf("AFR theme: id=%s, name=%s, indicators=%d", theme.ID, theme.Name, len(theme.Indicators))
	}

	// Verify we can parse themes correctly
	var ksiThemes map[string]ThemeJSON
	if err := json.Unmarshal(ksiRaw, &ksiThemes); err != nil {
		t.Fatalf("Failed to unmarshal KSI themes: %v", err)
	}

	t.Logf("Parsed %d themes", len(ksiThemes))
	totalIndicators := 0
	for code, theme := range ksiThemes {
		t.Logf("  Theme %s: %s (%d indicators)", code, theme.Name, len(theme.Indicators))
		totalIndicators += len(theme.Indicators)
	}
	t.Logf("Total indicators: %d", totalIndicators)

	if totalIndicators == 0 {
		t.Error("Expected indicators but got 0")
	}
}

func TestParseIndicators(t *testing.T) {
	client := NewClient()

	data, err := client.fetchDocument("FRMR.KSI.key-security-indicators.json")
	if err != nil {
		t.Fatalf("Failed to fetch KSI document: %v", err)
	}

	indicators, err := client.ParseIndicators(data)
	if err != nil {
		t.Fatalf("ParseIndicators failed: %v", err)
	}

	t.Logf("Parsed %d indicators", len(indicators))

	if len(indicators) == 0 {
		t.Error("Expected indicators but got 0")
	}

	// Print first few indicators
	for i, ind := range indicators {
		if i >= 3 {
			break
		}
		t.Logf("Indicator %d: %s - %s", i, ind.ID, ind.Name)
	}
}

func getKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestParseFSIRequirements(t *testing.T) {
	client := NewClient()

	data, err := client.fetchDocument("FRMR.FSI.fedramp-security-inbox.json")
	if err != nil {
		t.Fatalf("Failed to fetch FSI document: %v", err)
	}

	t.Logf("Fetched %d bytes of FSI data", len(data))

	// Parse the raw structure to understand it
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		t.Fatalf("Failed to unmarshal raw doc: %v", err)
	}
	t.Logf("Top-level keys: %v", getKeys(rawDoc))

	// Check FRR section
	frrRaw, ok := rawDoc["FRR"]
	if !ok {
		t.Fatal("FRR key not found")
	}

	var frr map[string]json.RawMessage
	if err := json.Unmarshal(frrRaw, &frr); err != nil {
		t.Fatalf("Failed to unmarshal FRR: %v", err)
	}
	t.Logf("FRR keys: %v", getKeys(frr))

	// Check FSI section within FRR
	fsiRaw, ok := frr["FSI"]
	if !ok {
		t.Fatal("FSI key not found in FRR")
	}
	t.Logf("FSI raw data: %s", string(fsiRaw)[:200])

	// Try parsing as categories
	var categories map[string]RequirementCategory
	if err := json.Unmarshal(fsiRaw, &categories); err != nil {
		t.Logf("Failed to parse as categories: %v", err)
	} else {
		t.Logf("Parsed %d categories", len(categories))
		for name, cat := range categories {
			t.Logf("  Category %s: %d requirements", name, len(cat.Requirements))
		}
	}

	// Now test the actual ParseRequirements function
	reqs, err := client.ParseRequirements(data, "FSI")
	if err != nil {
		t.Fatalf("ParseRequirements failed: %v", err)
	}

	t.Logf("ParseRequirements returned %d requirements", len(reqs))

	if len(reqs) == 0 {
		t.Error("Expected requirements but got 0")
	}

	for i, req := range reqs {
		if i >= 3 {
			break
		}
		t.Logf("Requirement %d: %s - %s (keyword: %s)", i, req.ID, req.Name, req.PrimaryKeyWord)
	}
}
