package test

import (
	"encoding/json"
	"testing"
	"time"
	"timeLedger/app/models"
)

func TestCenterSettings_JSON(t *testing.T) {
	settings := models.CenterSettings{
		AllowPublicRegister: true,
		DefaultLanguage:     "zh-TW",
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled models.CenterSettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.AllowPublicRegister != settings.AllowPublicRegister {
		t.Errorf("Expected AllowPublicRegister %v, got %v", settings.AllowPublicRegister, unmarshaled.AllowPublicRegister)
	}
}

func TestDateRange_JSON(t *testing.T) {
	dateRange := models.DateRange{
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(dateRange)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled models.DateRange
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !unmarshaled.StartDate.Equal(dateRange.StartDate) {
		t.Errorf("Expected StartDate %v, got %v", dateRange.StartDate, unmarshaled.StartDate)
	}
}

func TestRecurrenceRule_JSON(t *testing.T) {
	rule := models.RecurrenceRule{
		Type:     "WEEKLY",
		Interval: 2,
		Weekdays: []int{1, 3, 5},
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled models.RecurrenceRule
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.Type != rule.Type {
		t.Errorf("Expected Type %s, got %s", rule.Type, unmarshaled.Type)
	}

	if len(unmarshaled.Weekdays) != len(rule.Weekdays) {
		t.Errorf("Expected %d weekdays, got %d", len(rule.Weekdays), len(unmarshaled.Weekdays))
	}
}

func TestAuditPayload_JSON(t *testing.T) {
	payload := models.AuditPayload{
		Before: map[string]interface{}{"name": "old name"},
		After:  map[string]interface{}{"name": "new name"},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled models.AuditPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.Before == nil || unmarshaled.After == nil {
		t.Error("Before and After should not be nil")
	}
}

func TestModel_Fields(t *testing.T) {
	teacher := models.Teacher{
		LineUserID:        "TEST_LINE_001",
		Name:              "Test Teacher",
		Email:             "test@example.com",
		Bio:               "Test bio",
		City:              "台北市",
		District:          "大安區",
		PublicContactInfo: "Line: test",
		IsOpenToHiring:    true,
	}

	if teacher.LineUserID == "" {
		t.Error("LineUserID should not be empty")
	}

	if teacher.Name == "" {
		t.Error("Name should not be empty")
	}

	if !teacher.IsOpenToHiring {
		t.Error("IsOpenToHiring should be true")
	}
}
