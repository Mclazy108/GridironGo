package league

import (
	"testing"
)

func TestDefaultRules(t *testing.T) {
	rules := DefaultRules()

	// Test basic properties
	if rules.Name != "Standard League" {
		t.Errorf("Expected default name to be 'Standard League', got %s", rules.Name)
	}

	if rules.TeamCount != 10 {
		t.Errorf("Expected default team count to be 10, got %d", rules.TeamCount)
	}

	if rules.PPR != false {
		t.Errorf("Expected default PPR to be false, got %v", rules.PPR)
	}

	// Test roster positions
	if rules.RosterPositions.QB != 1 {
		t.Errorf("Expected default QB count to be 1, got %d", rules.RosterPositions.QB)
	}

	if rules.RosterPositions.RB != 2 {
		t.Errorf("Expected default RB count to be 2, got %d", rules.RosterPositions.RB)
	}

	// Test scoring rules exist
	categories := []string{"passing", "rushing", "receiving", "fumbles", "defensive", "kickReturns", "puntReturns", "kicking"}
	for _, category := range categories {
		if _, ok := rules.ScoringRules[category]; !ok {
			t.Errorf("Expected category %s to exist in default rules", category)
		}
	}

	// Test specific scoring rule values
	passingTDRule := rules.ScoringRules["passing"]["passingTouchdowns"]
	if passingTDRule.Type != FixedUnit || passingTDRule.Value != 4 {
		t.Errorf("Expected passing touchdowns to be worth 4 points (FixedUnit)")
	}

	rushingYdsRule := rules.ScoringRules["rushing"]["rushingYards"]
	if rushingYdsRule.Type != PerUnit || rushingYdsRule.Value != 0.1 {
		t.Errorf("Expected rushing yards to be worth 0.1 points per unit")
	}
}

func TestEnableDisablePPR(t *testing.T) {
	rules := DefaultRules()

	// Test enabling PPR
	rules.EnablePPR()
	if !rules.PPR {
		t.Errorf("Expected PPR to be true after enabling")
	}

	// Check if reception value changed
	if recRules, ok := rules.ScoringRules["receiving"]; ok {
		recPoints := recRules["receptions"]
		if recPoints.Value != 1.0 {
			t.Errorf("Expected receptions to be worth 1.0 point after enabling PPR, got %.1f", recPoints.Value)
		}
	} else {
		t.Errorf("Reception rules not found")
	}

	// Test disabling PPR
	rules.DisablePPR()
	if rules.PPR {
		t.Errorf("Expected PPR to be false after disabling")
	}

	// Check if reception value changed back
	if recRules, ok := rules.ScoringRules["receiving"]; ok {
		recPoints := recRules["receptions"]
		if recPoints.Value != 0.0 {
			t.Errorf("Expected receptions to be worth 0.0 points after disabling PPR, got %.1f", recPoints.Value)
		}
	} else {
		t.Errorf("Reception rules not found")
	}

	// Test half PPR
	rules.HalfPPR()
	if !rules.PPR {
		t.Errorf("Expected PPR to be true after enabling half PPR")
	}

	// Check if reception value is 0.5
	if recRules, ok := rules.ScoringRules["receiving"]; ok {
		recPoints := recRules["receptions"]
		if recPoints.Value != 0.5 {
			t.Errorf("Expected receptions to be worth 0.5 points for half PPR, got %.1f", recPoints.Value)
		}
	} else {
		t.Errorf("Reception rules not found")
	}
}

func TestSetPositionCount(t *testing.T) {
	rules := DefaultRules()

	// Test setting position counts
	tests := []struct {
		position string
		count    int
		getter   func() int
	}{
		{"QB", 2, func() int { return rules.RosterPositions.QB }},
		{"RB", 3, func() int { return rules.RosterPositions.RB }},
		{"WR", 4, func() int { return rules.RosterPositions.WR }},
		{"TE", 2, func() int { return rules.RosterPositions.TE }},
		{"FLEX", 2, func() int { return rules.RosterPositions.FLEX }},
		{"K", 0, func() int { return rules.RosterPositions.K }},
		{"DST", 2, func() int { return rules.RosterPositions.DST }},
		{"BN", 8, func() int { return rules.RosterPositions.BN }},
	}

	for _, test := range tests {
		err := rules.SetPositionCount(test.position, test.count)
		if err != nil {
			t.Errorf("Error setting position count for %s: %v", test.position, err)
		}

		if test.getter() != test.count {
			t.Errorf("Expected %s count to be %d, got %d", test.position, test.count, test.getter())
		}
	}

	// Test invalid position
	err := rules.SetPositionCount("INVALID", 1)
	if err == nil {
		t.Errorf("Expected error for invalid position, got nil")
	}
}

func TestTotalRosterSize(t *testing.T) {
	rules := DefaultRules()

	// Calculate expected total manually
	expected := rules.RosterPositions.QB +
		rules.RosterPositions.RB +
		rules.RosterPositions.WR +
		rules.RosterPositions.TE +
		rules.RosterPositions.FLEX +
		rules.RosterPositions.K +
		rules.RosterPositions.DST +
		rules.RosterPositions.BN

	if rules.TotalRosterSize() != expected {
		t.Errorf("Expected total roster size to be %d, got %d", expected, rules.TotalRosterSize())
	}

	// Change a position count and test again
	rules.SetPositionCount("BN", 10)
	expected += 4 // 10 - 6 (original BN count)

	if rules.TotalRosterSize() != expected {
		t.Errorf("Expected total roster size to be %d after change, got %d", expected, rules.TotalRosterSize())
	}
}

func TestValidateRules(t *testing.T) {
	rules := DefaultRules()

	// Test valid rules
	if err := rules.ValidateRules(); err != nil {
		t.Errorf("Expected default rules to be valid, got error: %v", err)
	}

	// Test invalid team count
	rules.TeamCount = 1
	if err := rules.ValidateRules(); err == nil {
		t.Errorf("Expected error for invalid team count")
	}
	rules.TeamCount = 10 // reset

	// Test invalid playoff teams
	rules.PlayoffTeams = 11
	if err := rules.ValidateRules(); err == nil {
		t.Errorf("Expected error for invalid playoff teams")
	}
	rules.PlayoffTeams = 4 // reset

	// Test invalid roster (no QB)
	rules.RosterPositions.QB = 0
	if err := rules.ValidateRules(); err == nil {
		t.Errorf("Expected error for no QB")
	}
	rules.RosterPositions.QB = 1 // reset

	// Test invalid roster (no bench)
	originalBN := rules.RosterPositions.BN
	rules.RosterPositions.BN = 0
	if err := rules.ValidateRules(); err == nil {
		t.Errorf("Expected error for no bench spots")
	}
	rules.RosterPositions.BN = originalBN // reset
}

func TestSetScoringRule(t *testing.T) {
	rules := DefaultRules()

	// Test setting a new rule
	newRule := ScoringRule{Type: FixedUnit, Value: 5.0}
	err := rules.SetScoringRule("newCategory", "newStat", newRule)
	if err != nil {
		t.Errorf("Error setting new scoring rule: %v", err)
	}

	// Verify the rule was set
	rule := rules.ScoringRules["newCategory"]["newStat"]
	if rule.Value != 5.0 || rule.Type != FixedUnit {
		t.Errorf("Expected new rule to be worth 5.0 points (FixedUnit)")
	}

	// Test updating an existing rule
	updatedRule := ScoringRule{Type: FixedUnit, Value: 6.0}
	err = rules.SetScoringRule("passing", "passingTouchdowns", updatedRule)
	if err != nil {
		t.Errorf("Error updating scoring rule: %v", err)
	}

	// Verify the rule was updated
	rule = rules.ScoringRules["passing"]["passingTouchdowns"]
	if rule.Value != 6.0 || rule.Type != FixedUnit {
		t.Errorf("Expected updated rule to be worth 6.0 points (FixedUnit)")
	}

	// Test setting with float64
	err = rules.SetScoringRule("passing", "somethingNew", 3.0)
	if err != nil {
		t.Errorf("Error setting rule with float64: %v", err)
	}

	// Verify it created a ScoringRule
	rule = rules.ScoringRules["passing"]["somethingNew"]
	if rule.Value != 3.0 || rule.Type != FixedUnit {
		t.Errorf("Expected rule from float to have Value 3.0 and Type FixedUnit")
	}

	// Test setting with ranges
	ranges := map[string]float64{
		"0-39":  3,
		"40-49": 4,
		"50+":   5,
	}
	err = rules.SetScoringRule("kicking", "newFieldGoals", ranges)
	if err != nil {
		t.Errorf("Error setting rule with ranges: %v", err)
	}

	// Verify it created a ScoringRule with ranges
	rule = rules.ScoringRules["kicking"]["newFieldGoals"]
	if rule.Type != RangeBased {
		t.Errorf("Expected rule from ranges to have Type RangeBased")
	} else {
		if len(rule.Ranges) != 3 {
			t.Errorf("Expected 3 ranges in rule, got %d", len(rule.Ranges))
		}
		if rule.Ranges["50+"] != 5 {
			t.Errorf("Expected range 50+ to be worth 5 points, got %.1f", rule.Ranges["50+"])
		}
	}
}

func TestGetScoringValue(t *testing.T) {
	rules := DefaultRules()

	// Test per-unit scoring (yards)
	points, err := rules.GetScoringValue("passing", "passingYards", 300)
	if err != nil {
		t.Errorf("Error getting scoring value: %v", err)
	}
	expected := 300 * 0.04 // 300 yards * 0.04 points per yard = 12 points
	if points != expected {
		t.Errorf("Expected %.2f points for 300 passing yards, got %.2f", expected, points)
	}

	// Test fixed scoring (touchdowns)
	points, err = rules.GetScoringValue("rushing", "rushingTouchdowns", 2)
	if err != nil {
		t.Errorf("Error getting scoring value: %v", err)
	}
	expected = 2 * 6 // 2 TDs * 6 points per TD = 12 points
	if points != expected {
		t.Errorf("Expected %.2f points for 2 rushing touchdowns, got %.2f", expected, points)
	}

	// Test range-based scoring (field goals)
	points, err = rules.GetScoringValue("kicking", "fieldGoalsMade", 45)
	if err != nil {
		t.Errorf("Error getting scoring value: %v", err)
	}
	if points != 4 { // 40-49 yard field goal = 4 points
		t.Errorf("Expected 4 points for 45-yard field goal, got %.2f", points)
	}

	// Test non-existent category
	_, err = rules.GetScoringValue("nonexistent", "stat", 1)
	if err == nil {
		t.Errorf("Expected error for non-existent category")
	}

	// Test non-existent stat type
	_, err = rules.GetScoringValue("passing", "nonexistent", 1)
	if err == nil {
		t.Errorf("Expected error for non-existent stat type")
	}
}

func TestPrintScoringRules(t *testing.T) {
	rules := DefaultRules()

	// Just test that it doesn't panic and returns a non-empty string
	output := rules.PrintScoringRules()
	if output == "" {
		t.Errorf("Expected non-empty output from PrintScoringRules")
	}
}

func TestGetScoringCategories(t *testing.T) {
	rules := DefaultRules()

	categories := rules.GetScoringCategories()

	// Test length matches
	if len(categories) != len(rules.ScoringRules) {
		t.Errorf("Expected %d categories, got %d", len(rules.ScoringRules), len(categories))
	}

	// Test all categories are present (order-independent)
	categoryMap := make(map[string]bool)
	for _, category := range categories {
		categoryMap[category] = true
	}

	for category := range rules.ScoringRules {
		if !categoryMap[category] {
			t.Errorf("Expected category %s to be in returned categories", category)
		}
	}
}

func TestGetStatTypesForCategory(t *testing.T) {
	rules := DefaultRules()

	// Test valid category
	statTypes, err := rules.GetStatTypesForCategory("passing")
	if err != nil {
		t.Errorf("Error getting stat types: %v", err)
	}

	// Test length matches
	if len(statTypes) != len(rules.ScoringRules["passing"]) {
		t.Errorf("Expected %d stat types, got %d",
			len(rules.ScoringRules["passing"]), len(statTypes))
	}

	// Test all stat types are present (order-independent)
	statTypeMap := make(map[string]bool)
	for _, statType := range statTypes {
		statTypeMap[statType] = true
	}

	for statType := range rules.ScoringRules["passing"] {
		if !statTypeMap[statType] {
			t.Errorf("Expected stat type %s to be in returned stat types", statType)
		}
	}

	// Test invalid category
	_, err = rules.GetStatTypesForCategory("nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent category")
	}
}

func TestUpdatePointValue(t *testing.T) {
	rules := DefaultRules()

	// Test updating a valid point value
	err := rules.UpdatePointValue("passing", "passingTouchdowns", 6.0)
	if err != nil {
		t.Errorf("Error updating point value: %v", err)
	}

	// Verify the update
	rule := rules.ScoringRules["passing"]["passingTouchdowns"]
	if rule.Value != 6.0 {
		t.Errorf("Expected passing touchdowns to be worth 6.0 points after update")
	}

	// Test updating a range-based rule
	err = rules.UpdatePointValue("kicking", "fieldGoalsMade", 5.0)
	if err == nil {
		t.Errorf("Expected error when updating range-based rule with single value")
	}

	// Test invalid category
	err = rules.UpdatePointValue("nonexistent", "stat", 1.0)
	if err == nil {
		t.Errorf("Expected error for non-existent category")
	}

	// Test invalid stat type
	err = rules.UpdatePointValue("passing", "nonexistent", 1.0)
	if err == nil {
		t.Errorf("Expected error for non-existent stat type")
	}
}

func TestUpdateRangePointValue(t *testing.T) {
	rules := DefaultRules()

	// Test updating a valid range value
	err := rules.UpdateRangePointValue("kicking", "fieldGoalsMade", "50+", 6.0)
	if err != nil {
		t.Errorf("Error updating range point value: %v", err)
	}

	// Verify the update
	rule := rules.ScoringRules["kicking"]["fieldGoalsMade"]
	if rule.Ranges["50+"] != 6.0 {
		t.Errorf("Expected 50+ yard field goals to be worth 6.0 points after update")
	}

	// Test updating a non-range rule
	err = rules.UpdateRangePointValue("passing", "passingTouchdowns", "any", 6.0)
	if err == nil {
		t.Errorf("Expected error when updating non-range rule with range value")
	}

	// Test invalid range
	err = rules.UpdateRangePointValue("kicking", "fieldGoalsMade", "nonexistent", 1.0)
	if err == nil {
		t.Errorf("Expected error for non-existent range")
	}
}

func TestJSONSerialization(t *testing.T) {
	original := DefaultRules()

	// Make some changes
	original.EnablePPR()
	original.SetPositionCount("WR", 3)
	original.UpdatePointValue("passing", "passingTouchdowns", 6.0)

	// Convert to JSON
	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Errorf("Error converting to JSON: %v", err)
	}

	// Convert back from JSON
	restored, err := FromJSON(jsonStr)
	if err != nil {
		t.Errorf("Error converting from JSON: %v", err)
	}

	// Compare some properties
	if restored.PPR != original.PPR {
		t.Errorf("PPR setting not preserved in JSON round-trip")
	}

	if restored.RosterPositions.WR != original.RosterPositions.WR {
		t.Errorf("WR count not preserved in JSON round-trip")
	}

	// Check a rule value
	rule := restored.ScoringRules["passing"]["passingTouchdowns"]
	if rule.Value != 6.0 {
		t.Errorf("Rule value not preserved correctly in JSON round-trip")
	}
}

func TestCloneRules(t *testing.T) {
	original := DefaultRules()

	// Make some changes to the original
	original.EnablePPR()
	original.SetPositionCount("WR", 3)

	// Clone the rules
	clone, err := original.CloneRules()
	if err != nil {
		t.Errorf("Error cloning rules: %v", err)
	}

	// Check that the clone has the same values
	if clone.PPR != original.PPR {
		t.Errorf("PPR setting not preserved in clone")
	}

	if clone.RosterPositions.WR != original.RosterPositions.WR {
		t.Errorf("WR count not preserved in clone")
	}

	// Change the clone and ensure the original is unchanged
	clone.DisablePPR()
	if original.PPR != true {
		t.Errorf("Original PPR changed when clone was modified")
	}

	// Change a rule in the clone
	err = clone.UpdatePointValue("passing", "passingTouchdowns", 7.0)
	if err != nil {
		t.Errorf("Error updating clone's passing TD value: %v", err)
	}

	// Check that original is unchanged
	originalPassingTD := original.ScoringRules["passing"]["passingTouchdowns"].Value
	if originalPassingTD != 4.0 {
		t.Errorf("Original passing TD value changed: expected 4.0, got %.1f", originalPassingTD)
	}

	// Check that clone was updated
	clonePassingTD := clone.ScoringRules["passing"]["passingTouchdowns"].Value
	if clonePassingTD != 7.0 {
		t.Errorf("Clone passing TD value not updated: expected 7.0, got %.1f", clonePassingTD)
	}
}
