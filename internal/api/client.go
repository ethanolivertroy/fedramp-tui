package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethanolivertroy/fedramp-tui/internal/cache"
	"github.com/ethanolivertroy/fedramp-tui/internal/model"
)

const BaseURL = "https://raw.githubusercontent.com/FedRAMP/docs/main/data"

// DocumentFiles maps document codes to their filenames
var DocumentFiles = map[string]DocumentMetadata{
	"FRD": {Code: "FRD", Name: "FedRAMP Definitions", Description: "Terms and definitions", Filename: "FRMR.FRD.fedramp-definitions.json"},
	"KSI": {Code: "KSI", Name: "Key Security Indicators", Description: "Security indicators with control mappings", Filename: "FRMR.KSI.key-security-indicators.json"},
	"VDR": {Code: "VDR", Name: "Vulnerability Detection & Response", Description: "Vulnerability management requirements", Filename: "FRMR.VDR.vulnerability-detection-and-response.json"},
	"UCM": {Code: "UCM", Name: "Using Cryptographic Modules", Description: "Cryptographic module requirements", Filename: "FRMR.UCM.using-cryptographic-modules.json"},
	"RSC": {Code: "RSC", Name: "Recommended Secure Configuration", Description: "Secure configuration requirements", Filename: "FRMR.RSC.recommended-secure-configuration.json"},
	"ADS": {Code: "ADS", Name: "Authorization Data Sharing", Description: "Data sharing requirements", Filename: "FRMR.ADS.authorization-data-sharing.json"},
	"CCM": {Code: "CCM", Name: "Collaborative Continuous Monitoring", Description: "Continuous monitoring requirements", Filename: "FRMR.CCM.collaborative-continuous-monitoring.json"},
	"FSI": {Code: "FSI", Name: "FedRAMP Security Inbox", Description: "Security inbox procedures", Filename: "FRMR.FSI.fedramp-security-inbox.json"},
	"ICP": {Code: "ICP", Name: "Incident Communications Procedures", Description: "Incident communication requirements", Filename: "FRMR.ICP.incident-communications-procedures.json"},
	"MAS": {Code: "MAS", Name: "Minimum Assessment Scope", Description: "Assessment scope requirements", Filename: "FRMR.MAS.minimum-assessment-scope.json"},
	"PVA": {Code: "PVA", Name: "Persistent Validation & Assessment", Description: "Validation and assessment requirements", Filename: "FRMR.PVA.persistent-validation-and-assessment.json"},
	"SCN": {Code: "SCN", Name: "Significant Change Notifications", Description: "Change notification requirements", Filename: "FRMR.SCN.significant-change-notifications.json"},
}

// DocumentOrder defines the display order of documents
var DocumentOrder = []string{"FRD", "KSI", "VDR", "UCM", "RSC", "ADS", "CCM", "FSI", "ICP", "MAS", "PVA", "SCN"}

// Client is an HTTP client for fetching FedRAMP documents
type Client struct {
	httpClient *http.Client
	baseURL    string
	cache      *cache.Cache
	refresh    bool
}

// ClientOption configures the client
type ClientOption func(*Client)

// WithRefresh forces fresh fetch, ignoring cache
func WithRefresh(refresh bool) ClientOption {
	return func(c *Client) {
		c.refresh = refresh
	}
}

// NewClient creates a new API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		baseURL:    BaseURL,
	}

	// Initialize cache (ignore errors, will just fetch fresh)
	if cache, err := cache.New(); err == nil {
		c.cache = cache
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// FetchResult holds the result of fetching a document
type FetchResult struct {
	Code  string
	Data  []byte
	Error error
}

// FetchAllDocuments fetches all documents in parallel
func (c *Client) FetchAllDocuments() (map[string][]byte, error) {
	results := make(map[string][]byte)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(DocumentFiles))

	for code, meta := range DocumentFiles {
		wg.Add(1)
		go func(code string, meta DocumentMetadata) {
			defer wg.Done()
			data, err := c.fetchDocument(meta.Filename)
			if err != nil {
				errChan <- fmt.Errorf("fetching %s: %w", code, err)
				return
			}
			mu.Lock()
			results[code] = data
			mu.Unlock()
		}(code, meta)
	}

	wg.Wait()
	close(errChan)

	// Collect any errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("errors fetching documents: %v", errs)
	}

	return results, nil
}

func (c *Client) fetchDocument(filename string) ([]byte, error) {
	url := c.baseURL + "/" + filename

	// Check cache first (unless refresh is forced)
	if c.cache != nil && !c.refresh {
		if data, ok := c.cache.Get(url); ok {
			return data, nil
		}
	}

	// Fetch from network
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Save to cache
	if c.cache != nil {
		_ = c.cache.Set(url, data)
	}

	return data, nil
}

// ParseDefinitions parses the FRD document into Definition models
func (c *Client) ParseDefinitions(data []byte) ([]model.Definition, error) {
	var doc DefinitionsDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	definitions := make([]model.Definition, len(doc.FRD.ALL))
	for i, d := range doc.FRD.ALL {
		note := d.Note
		if len(d.Notes) > 0 && note == "" {
			for j, n := range d.Notes {
				if j > 0 {
					note += " "
				}
				note += n
			}
		}

		definitions[i] = model.Definition{
			ID:           d.ID,
			Term:         d.Term,
			Alts:         d.Alts,
			Text:         d.Definition,
			Note:         note,
			Reference:    d.Reference,
			ReferenceURL: d.ReferenceURL,
		}
	}

	return definitions, nil
}

// ParseIndicators parses the KSI document into Indicator models
func (c *Client) ParseIndicators(data []byte) ([]model.Indicator, error) {
	// Parse just the parts we need to avoid type conflicts in other sections
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, fmt.Errorf("parsing raw document: %w", err)
	}

	// Extract the KSI section directly
	ksiRaw, ok := rawDoc["KSI"]
	if !ok {
		return nil, fmt.Errorf("KSI section not found in document")
	}

	var ksiThemes map[string]ThemeJSON
	if err := json.Unmarshal(ksiRaw, &ksiThemes); err != nil {
		return nil, fmt.Errorf("parsing KSI themes: %w", err)
	}

	var indicators []model.Indicator

	// Parse the KSI themes
	for themeCode, theme := range ksiThemes {
		for _, ind := range theme.Indicators {
			controls := make([]model.Control, len(ind.Controls))
			for j, ctrl := range ind.Controls {
				controls[j] = model.Control{
					ControlID: ctrl.ControlID,
					Title:     ctrl.Title,
				}
			}

			indicators = append(indicators, model.Indicator{
				ID:           ind.ID,
				ThemeCode:    themeCode,
				ThemeName:    theme.Name,
				ThemeDesc:    theme.Theme,
				Name:         ind.Name,
				Statement:    ind.Statement,
				Impact: model.Impact{
					Low:      ind.Impact.Low,
					Moderate: ind.Impact.Moderate,
					High:     ind.Impact.High,
				},
				Controls:     controls,
				Reference:    ind.Reference,
				ReferenceURL: ind.ReferenceURL,
				Note:         ind.Note,
				Retired:      ind.Retired,
			})
		}
	}

	return indicators, nil
}

// ParseRequirements parses a requirements document into Requirement models
func (c *Client) ParseRequirements(data []byte, docCode string) ([]model.Requirement, error) {
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	var requirements []model.Requirement

	// Parse FRR section if present
	if frrData, ok := rawDoc["FRR"]; ok {
		var frr map[string]json.RawMessage
		if err := json.Unmarshal(frrData, &frr); err != nil {
			return nil, err
		}

		// Look for the document's section
		if docData, ok := frr[docCode]; ok {
			var categories map[string]RequirementCategory
			if err := json.Unmarshal(docData, &categories); err != nil {
				// Try as a single category
				var singleCat RequirementCategory
				if err2 := json.Unmarshal(docData, &singleCat); err2 == nil {
					requirements = append(requirements, c.extractRequirements(singleCat.Requirements, docCode)...)
				}
			} else {
				for _, cat := range categories {
					requirements = append(requirements, c.extractRequirements(cat.Requirements, docCode)...)
				}
			}
		}
	}

	return requirements, nil
}

func (c *Client) extractRequirements(reqs []RequirementJSON, docCode string) []model.Requirement {
	var requirements []model.Requirement

	for i := range reqs {
		r := &reqs[i]
		// Process the following_information field which may be string or array
		r.UnmarshalFollowingInfo()

		req := model.Requirement{
			ID:             r.ID,
			DocumentCode:   docCode,
			Statement:      r.Statement,
			Name:           r.Name,
			Impact: model.Impact{
				Low:      r.Impact.Low,
				Moderate: r.Impact.Moderate,
				High:     r.Impact.High,
			},
			Affects:        r.Affects,
			PrimaryKeyWord: r.PrimaryKeyWord,
			Note:           r.Note,
		}
		requirements = append(requirements, req)

		// Also extract nested requirements
		if len(r.FollowingInformation) > 0 {
			requirements = append(requirements, c.extractRequirements(r.FollowingInformation, docCode)...)
		}
	}

	return requirements
}

// GetDocumentMetadata returns metadata for all documents
func GetDocumentMetadata() []model.Document {
	docs := make([]model.Document, 0, len(DocumentOrder))
	for _, code := range DocumentOrder {
		meta := DocumentFiles[code]
		docs = append(docs, model.Document{
			Code:        meta.Code,
			Name:        meta.Name,
			Description: meta.Description,
		})
	}
	return docs
}

// ParseDocumentInfo extracts the info section from a document's JSON data
func (c *Client) ParseDocumentInfo(data []byte) (*DocumentInfo, error) {
	var doc struct {
		Info DocumentInfo `json:"info"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc.Info, nil
}

// EnrichDocument populates a Document with info from the JSON
func EnrichDocument(doc *model.Document, info *DocumentInfo) {
	if info == nil {
		return
	}

	// Purpose and expected outcomes
	doc.Purpose = info.FrontMatter.Purpose
	doc.ExpectedOutcomes = info.FrontMatter.ExpectedOutcomes

	// Authority references
	for _, auth := range info.FrontMatter.Authority {
		doc.Authority = append(doc.Authority, model.Authority{
			Reference:    auth.Reference,
			ReferenceURL: auth.ReferenceURL,
			Description:  auth.Description,
		})
	}

	// Releases
	for _, rel := range info.Releases {
		doc.Releases = append(doc.Releases, model.Release{
			ID:            rel.ID,
			PublishedDate: rel.PublishedDate,
			Description:   rel.Description,
		})
	}

	// Effective info (program status)
	if len(info.Effective) > 0 {
		doc.EffectiveInfo = make(map[string]model.EffectiveStatus)
		for version, eff := range info.Effective {
			doc.EffectiveInfo[version] = model.EffectiveStatus{
				Is:            eff.Is,
				CurrentStatus: eff.CurrentStatus,
				StartDate:     eff.StartDate,
				EndDate:       eff.EndDate,
				SignupURL:     eff.SignupURL,
				Comments:      eff.Comments,
			}
		}
	}
}
