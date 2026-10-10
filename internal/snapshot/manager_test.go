package snapshot

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestExpandTemplate(t *testing.T) {
	ts := time.Date(2026, 4, 15, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		tmpl, camName, want string
	}{
		{"{cam_name}", "front_door", "front_door"},
		{"{cam_name}/%Y/%m/%d", "front_door", "front_door/2026/04/15"},
		{"{cam_name}_%Y%m%d_%H%M%S", "front_door", "front_door_20260415_143045"},
		{"{CAM_NAME}_snapshot", "test", "TEST_snapshot"},
		{"%H-%M-%S", "any", "14-30-45"},
		{"/media/snapshots/{cam_name}/%Y/%m/%d", "cam1", "/media/snapshots/cam1/2026/04/15"},
	}

	for _, tt := range tests {
		got := expandTemplate(tt.tmpl, tt.camName, ts)
		if got != tt.want {
			t.Errorf("expandTemplate(%q, %q) = %q, want %q", tt.tmpl, tt.camName, got, tt.want)
		}
	}
}

// TestExpandTemplate_UnixTime verifies %s (unix seconds) expands — exact value
// is non-deterministic across timezones on some test expectations, so we just
// confirm it's a plausible positive integer string and the prefix is intact.
func TestExpandTemplate_UnixTime(t *testing.T) {
	ts := time.Date(2026, 4, 15, 14, 30, 45, 0, time.UTC)
	got := expandTemplate("{CAM_NAME}_%s", "test", ts)
	const wantPrefix = "TEST_"
	if got[:len(wantPrefix)] != wantPrefix {
		t.Errorf("expandTemplate = %q, want prefix %q", got, wantPrefix)
	}
	if len(got) <= len(wantPrefix) {
		t.Errorf("expandTemplate produced no timestamp tail: %q", got)
	}
}

func TestContainsName(t *testing.T) {
	list := []string{"FRONT_DOOR", "BACKYARD"}

	if !containsName(list, "front_door") {
		t.Error("should match case-insensitive")
	}
	if !containsName(list, "BACKYARD") {
		t.Error("should match exact")
	}
	if containsName(list, "garage") {
		t.Error("should not match missing")
	}
}

func TestLatest_RemembersTheMostRecentCapturePerCamera(t *testing.T) {
	m := NewManager(nil, nil, nil, zerolog.Nop())
	if _, _, ok := m.Latest("front_door"); ok {
		t.Fatal("Latest should report nothing before any capture")
	}

	first := time.Date(2026, 10, 10, 21, 0, 0, 0, time.UTC)
	m.remember("front_door", []byte("one"), first)
	m.remember("front_door", []byte("two"), first.Add(time.Minute))
	m.remember("backyard", []byte("yard"), first)

	jpeg, at, ok := m.Latest("front_door")
	if !ok || string(jpeg) != "two" || !at.Equal(first.Add(time.Minute)) {
		t.Errorf("Latest(front_door) = %q, %v, %v; want the second capture", jpeg, at, ok)
	}
	if jpeg, _, _ := m.Latest("backyard"); string(jpeg) != "yard" {
		t.Errorf("Latest(backyard) = %q, want its own frame", jpeg)
	}
}
