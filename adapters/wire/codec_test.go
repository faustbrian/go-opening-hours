//nolint:staticcheck // This compatibility test intentionally exercises deprecated facades.
package openinghourswire_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises deprecated facades.

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
	"testing"

	openinghours "github.com/faustbrian/go-opening-hours"
	openinghourswire "github.com/faustbrian/go-opening-hours/adapters/wire"
	legacy "github.com/faustbrian/go-opening-hours/openinghourswire"
	wire "github.com/faustbrian/go-wire"
)

func TestCodecRoundTripAndNamedIdentity(t *testing.T) {
	if openinghourswire.Format != "opening-hours+json;v=1" || openinghourswire.WireFormat != wire.Format(openinghourswire.Format) {
		t.Fatalf("format constants = %q, %q", openinghourswire.Format, openinghourswire.WireFormat)
	}
	if reflect.TypeOf(openinghourswire.Codec{}) == reflect.TypeOf(legacy.Codec{}) {
		t.Fatal("canonical and legacy Codec reflection identities match")
	}
	want, err := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	codec := openinghourswire.Codec{}
	data, err := codec.Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := want.CanonicalJSON()
	if err != nil || !bytes.Equal(data, canonical) {
		t.Fatalf("Encode() = %q, %v; canonical = %q", data, err, canonical)
	}
	got, err := codec.Decode(data)
	if err != nil || !got.Equal(want) {
		t.Fatalf("Decode() error = %v, equal = %t", err, got.Equal(want))
	}
	for name, input := range map[string][]byte{
		"invalid":        []byte(`{}`),
		"unknown":        []byte(`{"version":1,"timezone":"","weekly":[],"exceptions":[],"metadata":{"label":"","source":"","revision":""},"outside_effective":"closed","unknown":true}`),
		"trailing":       append(append([]byte(nil), data...), []byte(` {}`)...),
		"malformed utf8": {0xff},
		"oversized":      []byte(`{"value":"` + strings.Repeat("x", 1<<20) + `"}`),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := codec.Decode(input); err == nil {
				t.Fatal("invalid bytes accepted")
			}
		})
	}
}

func TestCodecIsDeterministicUnderConcurrentUse(t *testing.T) {
	schedule, _ := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	want, _ := schedule.CanonicalJSON()
	codec := openinghourswire.Codec{}
	var group sync.WaitGroup
	errors := make(chan error, 16)
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			got, err := codec.Encode(schedule)
			if err != nil {
				errors <- err
				return
			}
			if !bytes.Equal(got, want) {
				errors <- errorsMismatch{}
			}
		}()
	}
	group.Wait()
	close(errors)
	for err := range errors {
		t.Fatal(err)
	}
}

type errorsMismatch struct{}

func (errorsMismatch) Error() string { return "concurrent encoding differed" }
