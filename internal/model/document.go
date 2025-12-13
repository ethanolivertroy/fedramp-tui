package model

// Document represents a FedRAMP document category
type Document struct {
	Code             string
	Name             string
	Description      string
	RequirementCount int
	// Rich metadata from JSON info section
	Purpose          string
	ExpectedOutcomes []string
	Authority        []Authority
	Releases         []Release
	EffectiveInfo    map[string]EffectiveStatus
}

// Authority represents a legal authority reference
type Authority struct {
	Reference    string
	ReferenceURL string
	Description  string
}

// Release represents a document release version
type Release struct {
	ID            string
	PublishedDate string
	Description   string
}

// EffectiveStatus represents program version status
type EffectiveStatus struct {
	Is            string
	CurrentStatus string
	StartDate     string
	EndDate       string
	SignupURL     string
	Comments      []string
}
