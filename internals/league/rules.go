package league

import (
	"encoding/json"
	"fmt"
	//"regexp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sort"
	"strings"
)

// RuleType defines the type of scoring rule
type RuleType string

const (
	PerUnit    RuleType = "per_unit"    // Points per unit (e.g., yards)
	FixedUnit  RuleType = "fixed_unit"  // Fixed points per occurrence
	RangeBased RuleType = "range_based" // Different points based on ranges
)

// ScoringRule represents a single scoring rule with its value
type ScoringRule struct {
	Type   RuleType           `json:"type"`             // The type of rule
	Value  float64            `json:"value"`            // Point value (used for PerUnit and FixedUnit)
	Ranges map[string]float64 `json:"ranges,omitempty"` // Range-based values (only for RangeBased)
}

// PositionRoster defines how many of each position can be on a roster
type PositionRoster struct {
	QB   int `json:"qb"`   // Quarterbacks
	RB   int `json:"rb"`   // Running Backs
	WR   int `json:"wr"`   // Wide Receivers
	TE   int `json:"te"`   // Tight Ends
	FLEX int `json:"flex"` // FLEX (RB/WR/TE)
	K    int `json:"k"`    // Kickers
	DST  int `json:"dst"`  // Defense/Special Teams
	BN   int `json:"bn"`   // Bench spots
}

// LeagueRules contains all the configuration options for a fantasy league
type LeagueRules struct {
	Name             string                            `json:"name"`
	Description      string                            `json:"description"`
	TeamCount        int                               `json:"team_count"`
	PPR              bool                              `json:"ppr"` // Points Per Reception
	RosterPositions  PositionRoster                    `json:"roster_positions"`
	ScoringRules     map[string]map[string]ScoringRule `json:"scoring_rules"` // Category -> StatType -> ScoringRule
	PlayoffWeekStart int                               `json:"playoff_week_start"`
	PlayoffTeams     int                               `json:"playoff_teams"`
}

// DefaultRules returns the standard fantasy football scoring rules
func DefaultRules() *LeagueRules {
	rules := &LeagueRules{
		Name:             "Standard League",
		Description:      "Default fantasy football league with standard scoring",
		TeamCount:        10,
		PPR:              false,
		PlayoffWeekStart: 15,
		PlayoffTeams:     4,
		RosterPositions: PositionRoster{
			QB:   1,
			RB:   2,
			WR:   2,
			TE:   1,
			FLEX: 1,
			K:    1,
			DST:  1,
			BN:   6,
		},
		ScoringRules: make(map[string]map[string]ScoringRule),
	}

	// Add passing rules
	passingRules := make(map[string]ScoringRule)
	passingRules["passingYards"] = ScoringRule{Type: PerUnit, Value: 0.04}     // 0.04 per yard
	passingRules["passingTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 4} // 4 points per TD
	passingRules["interceptions"] = ScoringRule{Type: FixedUnit, Value: -2}    // -2 per INT
	rules.ScoringRules["passing"] = passingRules

	// Add rushing rules
	rushingRules := make(map[string]ScoringRule)
	rushingRules["rushingYards"] = ScoringRule{Type: PerUnit, Value: 0.1}      // 0.1 per yard
	rushingRules["rushingTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 6} // 6 per TD
	rules.ScoringRules["rushing"] = rushingRules

	// Add receiving rules
	receivingRules := make(map[string]ScoringRule)
	receivingRules["receivingYards"] = ScoringRule{Type: PerUnit, Value: 0.1}      // 0.1 per yard
	receivingRules["receivingTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 6} // 6 per TD
	receivingRules["receptions"] = ScoringRule{Type: FixedUnit, Value: 0}          // 0 by default (PPR = false)
	rules.ScoringRules["receiving"] = receivingRules

	// Add fumble rules
	fumbleRules := make(map[string]ScoringRule)
	fumbleRules["fumblesLost"] = ScoringRule{Type: FixedUnit, Value: -2} // -2 per fumble lost
	rules.ScoringRules["fumbles"] = fumbleRules

	// Add defensive rules
	defensiveRules := make(map[string]ScoringRule)
	defensiveRules["sacks"] = ScoringRule{Type: FixedUnit, Value: 1}               // 1 per sack
	defensiveRules["interceptions"] = ScoringRule{Type: FixedUnit, Value: 2}       // 2 per INT
	defensiveRules["defensiveTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 6} // 6 per TD
	defensiveRules["passesDefended"] = ScoringRule{Type: FixedUnit, Value: 1}      // 1 per pass defended
	rules.ScoringRules["defensive"] = defensiveRules

	// Add kick/punt return rules
	kickReturnRules := make(map[string]ScoringRule)
	kickReturnRules["kickReturnYards"] = ScoringRule{Type: PerUnit, Value: 0.04}     // 0.04 per yard
	kickReturnRules["kickReturnTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 6} // 6 per TD
	rules.ScoringRules["kickReturns"] = kickReturnRules

	puntReturnRules := make(map[string]ScoringRule)
	puntReturnRules["puntReturnYards"] = ScoringRule{Type: PerUnit, Value: 0.04}     // 0.04 per yard
	puntReturnRules["puntReturnTouchdowns"] = ScoringRule{Type: FixedUnit, Value: 6} // 6 per TD
	rules.ScoringRules["puntReturns"] = puntReturnRules

	// Add kicking rules with ranges
	kickingRules := make(map[string]ScoringRule)
	fieldGoalRanges := map[string]float64{
		"0-39":  3, // 3 points for 0-39 yard field goals
		"40-49": 4, // 4 points for 40-49 yard field goals
		"50+":   5, // 5 points for 50+ yard field goals
	}
	kickingRules["fieldGoalsMade"] = ScoringRule{Type: RangeBased, Ranges: fieldGoalRanges}
	kickingRules["extraPointsMade"] = ScoringRule{Type: FixedUnit, Value: 1} // 1 point per extra point
	rules.ScoringRules["kicking"] = kickingRules

	return rules
}

// EnablePPR turns on PPR (Points Per Reception) scoring
func (l *LeagueRules) EnablePPR() {
	l.PPR = true

	// Update the reception points
	if recRules, ok := l.ScoringRules["receiving"]; ok {
		if rule, ok := recRules["receptions"]; ok {
			rule.Value = 1.0
			recRules["receptions"] = rule
		}
	}
}

// DisablePPR turns off PPR (Points Per Reception) scoring
func (l *LeagueRules) DisablePPR() {
	l.PPR = false

	// Update the reception points
	if recRules, ok := l.ScoringRules["receiving"]; ok {
		if rule, ok := recRules["receptions"]; ok {
			rule.Value = 0.0
			recRules["receptions"] = rule
		}
	}
}

// HalfPPR sets up Half PPR (0.5 points per reception) scoring
func (l *LeagueRules) HalfPPR() {
	l.PPR = true

	// Update the reception points
	if recRules, ok := l.ScoringRules["receiving"]; ok {
		if rule, ok := recRules["receptions"]; ok {
			rule.Value = 0.5
			recRules["receptions"] = rule
		}
	}
}

// SetPositionCount updates the roster configuration for a specific position
func (l *LeagueRules) SetPositionCount(position string, count int) error {
	position = strings.ToUpper(position)

	switch position {
	case "QB":
		l.RosterPositions.QB = count
	case "RB":
		l.RosterPositions.RB = count
	case "WR":
		l.RosterPositions.WR = count
	case "TE":
		l.RosterPositions.TE = count
	case "FLEX":
		l.RosterPositions.FLEX = count
	case "K":
		l.RosterPositions.K = count
	case "DST":
		l.RosterPositions.DST = count
	case "BN":
		l.RosterPositions.BN = count
	default:
		return fmt.Errorf("invalid position: %s", position)
	}

	return nil
}

// TotalRosterSize returns the total number of players on a roster
func (l *LeagueRules) TotalRosterSize() int {
	return l.RosterPositions.QB +
		l.RosterPositions.RB +
		l.RosterPositions.WR +
		l.RosterPositions.TE +
		l.RosterPositions.FLEX +
		l.RosterPositions.K +
		l.RosterPositions.DST +
		l.RosterPositions.BN
}

// TotalStartingPlayers returns the number of starting players
func (l *LeagueRules) TotalStartingPlayers() int {
	return l.RosterPositions.QB +
		l.RosterPositions.RB +
		l.RosterPositions.WR +
		l.RosterPositions.TE +
		l.RosterPositions.FLEX +
		l.RosterPositions.K +
		l.RosterPositions.DST
}

// ValidateRules checks if the rules configuration is valid
func (l *LeagueRules) ValidateRules() error {
	// Check team count
	if l.TeamCount < 2 || l.TeamCount > 32 {
		return fmt.Errorf("invalid team count: %d (must be between 2-32)", l.TeamCount)
	}

	// Check playoff teams
	if l.PlayoffTeams < 2 || l.PlayoffTeams > l.TeamCount {
		return fmt.Errorf("invalid playoff teams: %d (must be between 2-%d)", l.PlayoffTeams, l.TeamCount)
	}

	// Check playoff week
	if l.PlayoffWeekStart < 10 || l.PlayoffWeekStart > 17 {
		return fmt.Errorf("invalid playoff start week: %d (must be between 10-17)", l.PlayoffWeekStart)
	}

	// Check roster positions
	if l.RosterPositions.QB < 1 {
		return fmt.Errorf("must have at least 1 QB roster spot")
	}

	if l.TotalRosterSize() < 9 {
		return fmt.Errorf("total roster size must be at least 9 players")
	}

	// Check that starters < total roster size
	if l.TotalStartingPlayers() >= l.TotalRosterSize() {
		return fmt.Errorf("must have at least one bench spot")
	}

	return nil
}

// SetScoringRule sets or updates a specific scoring rule
func (l *LeagueRules) SetScoringRule(category, statType string, value any) error {
	// Make sure the category exists
	if _, ok := l.ScoringRules[category]; !ok {
		l.ScoringRules[category] = make(map[string]ScoringRule)
	}

	// Handle different value types
	switch v := value.(type) {
	case ScoringRule:
		l.ScoringRules[category][statType] = v
	case float64:
		// If just a float is provided, treat it as a fixed-unit point value
		l.ScoringRules[category][statType] = ScoringRule{Type: FixedUnit, Value: v}
	case map[string]float64:
		// If a map is provided, treat it as a range-based rule
		l.ScoringRules[category][statType] = ScoringRule{Type: RangeBased, Ranges: v}
	default:
		return fmt.Errorf("invalid value type for scoring rule: %T", value)
	}

	return nil
}

// GetScoringValue calculates the points for a specific stat
func (l *LeagueRules) GetScoringValue(category, statType string, statValue float64) (float64, error) {
	categoryRules, ok := l.ScoringRules[category]
	if !ok {
		return 0, fmt.Errorf("category not found: %s", category)
	}

	rule, ok := categoryRules[statType]
	if !ok {
		return 0, fmt.Errorf("stat type not found: %s", statType)
	}

	// Calculate points based on rule type
	switch rule.Type {
	case PerUnit:
		return rule.Value * statValue, nil

	case FixedUnit:
		return rule.Value * statValue, nil

	case RangeBased:
		// Handle field goals and other range-based scoring
		if statType == "fieldGoalsMade" && statValue > 0 {
			// Determine the range for field goals
			if statValue >= 50 {
				return rule.Ranges["50+"], nil
			} else if statValue >= 40 {
				return rule.Ranges["40-49"], nil
			} else {
				return rule.Ranges["0-39"], nil
			}
		}
		return 0, fmt.Errorf("range-based scoring not implemented for: %s", statType)

	default:
		return 0, fmt.Errorf("unsupported rule type: %s", rule.Type)
	}
}

// PrintScoringRules displays all current scoring rules in a readable format
func (l *LeagueRules) PrintScoringRules() string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Scoring Rules for %s\n", l.Name))
	output.WriteString(fmt.Sprintf("PPR: %v\n\n", l.PPR))

	// Sort categories for consistent display
	categories := make([]string, 0, len(l.ScoringRules))
	for category := range l.ScoringRules {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	// Create a title-casing object
	caser := cases.Title(language.English)

	for _, category := range categories {
		output.WriteString(fmt.Sprintf("%s:\n", caser.String(category)))
		rules := l.ScoringRules[category]

		// Sort stat types for consistent display
		statTypes := make([]string, 0, len(rules))
		for statType := range rules {
			statTypes = append(statTypes, statType)
		}
		sort.Strings(statTypes)

		for _, statType := range statTypes {
			rule := rules[statType]

			switch rule.Type {
			case PerUnit:
				output.WriteString(fmt.Sprintf("  %s: %.2f points per unit\n", caser.String(statType), rule.Value))
			case FixedUnit:
				output.WriteString(fmt.Sprintf("  %s: %.2f points\n", caser.String(statType), rule.Value))
			case RangeBased:
				output.WriteString(fmt.Sprintf("  %s:\n", caser.String(statType)))

				// Sort ranges for consistent display
				ranges := make([]string, 0, len(rule.Ranges))
				for rng := range rule.Ranges {
					ranges = append(ranges, rng)
				}
				sort.Strings(ranges)

				for _, rng := range ranges {
					output.WriteString(fmt.Sprintf("    %s yards: %.2f points\n", rng, rule.Ranges[rng]))
				}
			}
		}
		output.WriteString("\n")
	}

	return output.String()
}

// GetScoringCategories returns all available scoring categories
func (l *LeagueRules) GetScoringCategories() []string {
	categories := make([]string, 0, len(l.ScoringRules))
	for category := range l.ScoringRules {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// GetStatTypesForCategory returns all stat types for a given category
func (l *LeagueRules) GetStatTypesForCategory(category string) ([]string, error) {
	categoryRules, ok := l.ScoringRules[category]
	if !ok {
		return nil, fmt.Errorf("category not found: %s", category)
	}

	statTypes := make([]string, 0, len(categoryRules))
	for statType := range categoryRules {
		statTypes = append(statTypes, statType)
	}
	sort.Strings(statTypes)
	return statTypes, nil
}

// UpdatePointValue modifies a point value for a specific stat type
func (l *LeagueRules) UpdatePointValue(category, statType string, value float64) error {
	categoryRules, ok := l.ScoringRules[category]
	if !ok {
		return fmt.Errorf("category not found: %s", category)
	}

	rule, ok := categoryRules[statType]
	if !ok {
		return fmt.Errorf("stat type not found: %s", statType)
	}

	if rule.Type == RangeBased {
		return fmt.Errorf("cannot update range-based point value with a single value")
	}

	rule.Value = value
	categoryRules[statType] = rule

	return nil
}

// UpdateRangePointValue modifies a point value for a specific range
func (l *LeagueRules) UpdateRangePointValue(category, statType, rangeKey string, value float64) error {
	categoryRules, ok := l.ScoringRules[category]
	if !ok {
		return fmt.Errorf("category not found: %s", category)
	}

	rule, ok := categoryRules[statType]
	if !ok {
		return fmt.Errorf("stat type not found: %s", statType)
	}

	if rule.Type != RangeBased {
		return fmt.Errorf("stat type is not range-based: %s", statType)
	}

	if _, ok := rule.Ranges[rangeKey]; !ok {
		return fmt.Errorf("range not found: %s", rangeKey)
	}

	rule.Ranges[rangeKey] = value
	categoryRules[statType] = rule

	return nil
}

// ToJSON returns the rules configuration as a JSON string
func (l *LeagueRules) ToJSON() (string, error) {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling rules: %w", err)
	}

	return string(data), nil
}

// FromJSON loads rules from a JSON string
func FromJSON(jsonStr string) (*LeagueRules, error) {
	var rules LeagueRules
	if err := json.Unmarshal([]byte(jsonStr), &rules); err != nil {
		return nil, fmt.Errorf("error unmarshaling rules: %w", err)
	}

	return &rules, nil
}

// CloneRules creates a deep copy of LeagueRules
func (l *LeagueRules) CloneRules() (*LeagueRules, error) {
	// Serialize to JSON
	jsonStr, err := l.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("error marshaling rules for clone: %w", err)
	}

	// Deserialize to create a new instance
	return FromJSON(jsonStr)
}
