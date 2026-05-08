package main

import "testing"

// Domain validator regression suite.
//
// /autoplan called for a "scheduler golden-file" test. After inspecting
// company/server/, the codebase has no scheduling algorithm — Staffjoy v2's
// README explicitly says it does not generate or auto-assign shifts; it is
// CRUD over a manual schedule. The closest equivalent worth pinning is the
// pure-function validation logic in helpers.go and timezones.go, which a
// careless refactor (or a Go-version-driven stdlib change) could quietly
// break.

func TestSanitizeDayOfWeek(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"monday", "monday", false},
		{"  Monday ", "monday", false},
		{"FRIDAY", "friday", false},
		{"sunday", "sunday", false},
		{"saturday", "saturday", false},
		{"", "", true},
		{"funday", "", true},
		{"mon", "", true},
	}
	for _, tc := range tests {
		got, err := sanitizeDayOfWeek(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("sanitizeDayOfWeek(%q) error = %v, wantErr %v", tc.in, err, tc.wantErr)
			continue
		}
		if got != tc.want {
			t.Errorf("sanitizeDayOfWeek(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidColor(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{"abc", false},
		{"ABC", false},
		{"a1b2c3", false},
		{"FFFFFF", false},
		{"000000", false},
		{"", true},
		{"abcd", true},      // 4 chars
		{"abcdefg", true},   // 7 chars
		{"#abcdef", true},   // leading #
		{"xyz", true},       // non-hex
		{"abc def", true},   // whitespace
	}
	for _, tc := range tests {
		err := validColor(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("validColor(%q) err=%v, wantErr=%v", tc.in, err, tc.wantErr)
		}
	}
}

func TestValidTimezone(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{"UTC", false},
		{"America/Los_Angeles", false},
		{"Europe/London", false},
		{"Asia/Tokyo", false},
		{"", false}, // time.LoadLocation("") returns UTC, not error
		{"Mars/Olympus_Mons", true},
		{"not-a-tz", true},
	}
	for _, tc := range tests {
		err := validTimezone(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("validTimezone(%q) err=%v, wantErr=%v", tc.in, err, tc.wantErr)
		}
	}
}
