package api

import "testing"

func FuzzParseDefinitions(f *testing.F) {
	// Seed with minimal valid structure
	f.Add([]byte(`{"FRD":{"ALL":[]}}`))
	f.Add([]byte(`{"FRD":{"ALL":[{"id":"test","term":"Test Term","definition":"A test"}]}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseDefinitions(data)
	})
}

func FuzzParseIndicators(f *testing.F) {
	// Seed with minimal valid structure
	f.Add([]byte(`{"KSI":{}}`))
	f.Add([]byte(`{"KSI":{"AFR":{"id":"AFR","name":"Test","indicators":[]}}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseIndicators(data)
	})
}

func FuzzParseRequirements(f *testing.F) {
	// Seed with minimal valid structure and doc codes
	f.Add([]byte(`{"FRR":{"FSI":{}}}`), "FSI")
	f.Add([]byte(`{"FRR":{"VDR":{"requirements":[]}}}`), "VDR")
	f.Add([]byte(`{}`), "FSI")
	f.Add([]byte(`invalid json`), "")

	f.Fuzz(func(t *testing.T, data []byte, docCode string) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseRequirements(data, docCode)
	})
}
