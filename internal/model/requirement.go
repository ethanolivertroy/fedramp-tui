package model

// Impact represents the impact levels for a requirement
type Impact struct {
	Low      bool
	Moderate bool
	High     bool
}

// ImpactString returns a human-readable string of impact levels
func (i Impact) String() string {
	var levels []string
	if i.Low {
		levels = append(levels, "Low")
	}
	if i.Moderate {
		levels = append(levels, "Moderate")
	}
	if i.High {
		levels = append(levels, "High")
	}
	if len(levels) == 0 {
		return "N/A"
	}
	result := levels[0]
	for i := 1; i < len(levels); i++ {
		result += ", " + levels[i]
	}
	return result
}

// Requirement represents a FedRAMP requirement
type Requirement struct {
	ID             string
	DocumentCode   string
	Statement      string
	Name           string
	Impact         Impact
	Affects        []string
	PrimaryKeyWord string
	Note           string
}

// IsMust returns true if this is a MUST requirement
func (r Requirement) IsMust() bool {
	return r.PrimaryKeyWord == "MUST"
}

// IsShould returns true if this is a SHOULD requirement
func (r Requirement) IsShould() bool {
	return r.PrimaryKeyWord == "SHOULD"
}
