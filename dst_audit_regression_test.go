package openinghours_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v6"
	prepared "github.com/faustbrian/go-opening-hours/v6/compile"
)

type r04OraclePoint struct {
	Instant time.Time `json:"instant"`
	Open    bool      `json:"open"`
}
type r04Oracle struct {
	Transitions []r04OraclePoint `json:"transitions"`
	States      []r04OraclePoint `json:"states"`
}

func r04Fixture(t *testing.T) (openinghours.Schedule, r04Oracle) {
	t.Helper()
	// Set ZONEINFO before process start: time.LoadLocation retains its env choice.
	directory, err := filepath.Abs("testdata/r04/zoneinfo")
	if err != nil || os.Getenv("ZONEINFO") != directory {
		t.Fatal("probe must start with the exact official ZONEINFO directory")
	}
	raw, err := os.ReadFile("testdata/r04/zoneinfo/America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != "d7f2206b3a45989fc9ad63d558922532fa7352280d5f87176bf1db79cb1d1fa9" {
		t.Fatal("timezone input identity mismatch")
	}
	oracleBytes, err := os.ReadFile("testdata/r04/oracle.json")
	if err != nil {
		t.Fatal(err)
	}
	var oracle r04Oracle
	if err := json.Unmarshal(oracleBytes, &oracle); err != nil {
		t.Fatal(err)
	}
	if len(oracle.Transitions) != 4 || len(oracle.States) != 12 {
		t.Fatal("incomplete bounded independent oracle")
	}
	start, err := openinghours.NewLocalTime(1, 30, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	end, err := openinghours.NewLocalTime(2, 30, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	interval, err := openinghours.NewRange(start, end)
	if err != nil {
		t.Fatal(err)
	}
	day, err := openinghours.OpenRanges([]openinghours.Range{interval}, openinghours.RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "America/New_York", Weekly: map[time.Weekday]openinghours.DayRule{time.Sunday: day}})
	if err != nil {
		t.Fatal(err)
	}
	return schedule, oracle
}

func TestR04WallPointMembershipMatchesIndependentOracle(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	for _, point := range oracle.States {
		state, err := schedule.IsOpen(point.Instant)
		if err != nil || state.Open != point.Open {
			t.Errorf("point %s open=%t err=%v; oracle %t", point.Instant, state.Open, err, point.Open)
		}
	}
}

// The original earlier-opening/later-closing envelope covered two hours and
// included a closed repeated period. The selected policy now requires exact
// point membership; retain that historical countercase without requiring it.
func TestR04ExpansionExcludesHistoricallyClosedFoldPeriod(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	start, end := oracle.Transitions[0].Instant, oracle.Transitions[3].Instant
	ranges, err := schedule.EffectiveInstantRanges(start, end)
	if err != nil || len(ranges) != 2 || !ranges[0].Start.Equal(start) || !ranges[1].End.Equal(end) {
		t.Fatalf("point expansion %#v %v", ranges, err)
	}
	fold := oracle.Transitions[1].Instant
	state, err := schedule.IsOpen(fold)
	if err != nil || state.Open {
		t.Fatal("fold wall state should be closed")
	}
	for _, interval := range ranges {
		if !fold.Before(interval.Start) && fold.Before(interval.End) {
			t.Fatal("expansion includes historically closed wall instant")
		}
	}
	duration, err := schedule.OpenDuration(start, end)
	if err != nil || duration != 90*time.Minute {
		t.Fatalf("point expansion duration %s %v", duration, err)
	}
}

func TestR04TransitionsFindFirstActualWallStateChange(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	index, err := prepared.New(schedule)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name  string
		query func(time.Time, time.Duration) (openinghours.Transition, error)
		from  time.Time
		want  r04OraclePoint
	}{
		{"forward-at-offset-jump", schedule.NextTransition, oracle.Transitions[1].Instant.Add(-time.Minute), oracle.Transitions[1]},
		{"forward-second-fold-occurrence", schedule.NextTransition, oracle.Transitions[1].Instant, oracle.Transitions[2]},
		{"backward-at-offset-jump", schedule.PreviousTransition, oracle.Transitions[1].Instant.Add(time.Minute), oracle.Transitions[1]},
		{"backward-second-fold-occurrence", schedule.PreviousTransition, oracle.Transitions[2].Instant.Add(time.Minute), oracle.Transitions[2]},
		{"typed-closing", schedule.NextClosing, oracle.Transitions[1].Instant.Add(-time.Minute), oracle.Transitions[1]},
		{"typed-opening", schedule.NextOpening, oracle.Transitions[1].Instant, oracle.Transitions[2]},
		{"compiled-forward", index.NextTransition, oracle.Transitions[1].Instant.Add(-time.Minute), oracle.Transitions[1]},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transition, err := tc.query(tc.from, 2*time.Hour)
			wantKind := openinghours.TransitionClose
			if tc.want.Open {
				wantKind = openinghours.TransitionOpen
			}
			if err != nil || !transition.Instant.Equal(tc.want.Instant) || transition.Kind != wantKind {
				t.Fatalf("actual first/nearest transition %#v %v; oracle %s kind=%v", transition, err, tc.want.Instant, wantKind)
			}
			before, beforeErr := schedule.IsOpen(transition.Instant.Add(-time.Nanosecond))
			after, afterErr := schedule.IsOpen(transition.Instant)
			if beforeErr != nil || afterErr != nil || before.Open == after.Open || after.Open != tc.want.Open {
				t.Fatal("transition does not identify actual point-state change")
			}
		})
	}
}

// time.LoadLocation reads ZONEINFO once, so isolate the pinned provider input
// before process startup rather than changing the parent's global timezone state.
func r04RunWithPinnedZone(t *testing.T) bool {
	t.Helper()
	if os.Getenv("OPENING_HOURS_R04_TZ_PROBE") == "1" {
		return false
	}
	directory, err := filepath.Abs("testdata/r04/zoneinfo")
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// #nosec G204 -- os.Executable supplies this running test binary; arguments select a package-owned test and never invoke a shell.
	command := exec.CommandContext(ctx, executable, "-test.run=^"+t.Name()+"$", "-test.count=1", "-test.timeout=8s")
	command.Env = append(os.Environ(), "OPENING_HOURS_R04_TZ_PROBE=1", "ZONEINFO="+directory)
	output, err := command.CombinedOutput()
	if err != nil {
		if len(output) > 8192 {
			output = output[:8192]
		}
		t.Fatalf("isolated official-timezone probe: %v: %s", err, output)
	}
	return true
}
