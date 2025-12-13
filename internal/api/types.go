package api

import "encoding/json"

// DocumentInfo represents the common info structure in all FedRAMP documents
type DocumentInfo struct {
	Name        string                    `json:"name"`
	ShortName   string                    `json:"short_name"`
	Effective   map[string]EffectiveInfo  `json:"effective"`
	Releases    []Release                 `json:"releases"`
	FrontMatter FrontMatter               `json:"front_matter"`
}

// EffectiveInfo represents version-specific applicability
type EffectiveInfo struct {
	Is            string   `json:"is"`
	SignupURL     string   `json:"signup_url"`
	CurrentStatus string   `json:"current_status"`
	StartDate     string   `json:"start_date"`
	EndDate       string   `json:"end_date"`
	Comments      []string `json:"comments"`
	Warnings      []string `json:"warnings"`
}

// RelatedRFC represents a related RFC reference
type RelatedRFC struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	DiscussionURL string `json:"discussion_url"`
	ShortName     string `json:"short_name"`
	FullName      string `json:"full_name"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
}

// Release represents a document release version
type Release struct {
	ID            string       `json:"id"`
	PublishedDate string       `json:"published_date"`
	Description   string       `json:"description"`
	PublicComment bool         `json:"public_comment"`
	RelatedRFCs   []RelatedRFC `json:"related_rfcs"`
}

// FrontMatter contains authority and purpose information
type FrontMatter struct {
	Authority        []Authority `json:"authority"`
	Purpose          string      `json:"purpose"`
	ExpectedOutcomes []string    `json:"expected_outcomes"`
}

// Authority represents a legal authority reference
type Authority struct {
	Reference     string `json:"reference"`
	ReferenceURL  string `json:"reference_url"`
	Description   string `json:"description"`
	Delegation    string `json:"delegation"`
	DelegationURL string `json:"delegation_url"`
}

// RequirementJSON represents a requirement from FRR sections
type RequirementJSON struct {
	ID                   string             `json:"id"`
	Statement            string             `json:"statement"`
	Name                 string             `json:"name"`
	Impact               ImpactJSON         `json:"impact"`
	Affects              []string           `json:"affects"`
	PrimaryKeyWord       string             `json:"primary_key_word"`
	Note                 string             `json:"note"`
	FollowingInformation FollowingInfoField `json:"-"` // Custom unmarshaling
	RawFollowingInfo     json.RawMessage    `json:"following_information"`
}

// FollowingInfoField handles following_information which can be string or []RequirementJSON
type FollowingInfoField []RequirementJSON

// UnmarshalFollowingInfo processes the raw following_information field after initial unmarshal
func (r *RequirementJSON) UnmarshalFollowingInfo() {
	if r.RawFollowingInfo == nil || len(r.RawFollowingInfo) == 0 {
		return
	}
	// Try as array of requirements first
	var reqs []RequirementJSON
	if err := json.Unmarshal(r.RawFollowingInfo, &reqs); err == nil {
		r.FollowingInformation = reqs
		return
	}
	// Otherwise it's a string or other type, ignore
}

// ImpactJSON represents impact levels
type ImpactJSON struct {
	Low      bool `json:"low"`
	Moderate bool `json:"moderate"`
	High     bool `json:"high"`
}

// RequirementCategory represents a category of requirements
type RequirementCategory struct {
	ID           string            `json:"id"`
	Application  string            `json:"application"`
	Name         string            `json:"name"`
	Requirements []RequirementJSON `json:"requirements"`
}

// DefinitionJSON represents a FedRAMP definition
type DefinitionJSON struct {
	ID           string   `json:"id"`
	Term         string   `json:"term"`
	Alts         []string `json:"alts"`
	Definition   string   `json:"definition"`
	Note         string   `json:"note"`
	Notes        []string `json:"notes"`
	Reference    string   `json:"reference"`
	ReferenceURL string   `json:"reference_url"`
}

// IndicatorJSON represents a KSI indicator
type IndicatorJSON struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Statement    string       `json:"statement"`
	Impact       ImpactJSON   `json:"impact"`
	Controls     []ControlJSON `json:"controls"`
	Reference    string       `json:"reference"`
	ReferenceURL string       `json:"reference_url"`
	Note         string       `json:"note"`
	Retired      bool         `json:"retired"`
}

// ControlJSON represents an SP 800-53 control reference
type ControlJSON struct {
	ControlID string `json:"control_id"`
	Title     string `json:"title"`
}

// ThemeJSON represents a KSI theme
type ThemeJSON struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Theme      string          `json:"theme"`
	Indicators []IndicatorJSON `json:"indicators"`
}

// DefinitionsDocument represents the FRD document structure
type DefinitionsDocument struct {
	Schema string       `json:"$schema"`
	ID     string       `json:"$id"`
	Info   DocumentInfo `json:"info"`
	FRD    struct {
		ALL []DefinitionJSON `json:"ALL"`
	} `json:"FRD"`
}

// KSIDocument represents the KSI document structure
type KSIDocument struct {
	Schema string       `json:"$schema"`
	ID     string       `json:"$id"`
	Info   DocumentInfo `json:"info"`
	FRR    struct {
		KSI struct {
			Base RequirementCategory `json:"base"`
		} `json:"KSI"`
	} `json:"FRR"`
	KSI map[string]ThemeJSON `json:"KSI"`
}

// RequirementsDocument represents a generic requirements document (VDR, UCM, RSC, etc.)
type RequirementsDocument struct {
	Schema string       `json:"$schema"`
	ID     string       `json:"$id"`
	Info   DocumentInfo `json:"info"`
	FRR    map[string]map[string]RequirementCategory `json:"FRR"`
}

// DocumentMetadata holds basic info about a FedRAMP document
type DocumentMetadata struct {
	Code        string
	Name        string
	Description string
	Filename    string
}
