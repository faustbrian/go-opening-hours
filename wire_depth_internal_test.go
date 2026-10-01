package openinghours

import (
	"strings"
	"testing"
)

func TestJSONAdmissionHasIndependentDepthBoundary(t *testing.T) {
	// Six leaf edges plus two edges per additional composition level give
	// 36 at the supported 16 levels. Keep the oracle independent of the
	// implementation constant so changing admission does not move the test.
	for _, depth := range []int{36, 37} {
		encoded := strings.Repeat("[", depth) + "0" + strings.Repeat("]", depth)
		err := validateJSON([]byte(encoded))
		if depth == 36 && err != nil {
			t.Fatalf("supported depth 36 rejected: %v", err)
		}
		if depth == 37 && !IsCode(err, CodeLimitExceeded) {
			t.Fatalf("depth 37 = %v, want limit_exceeded", err)
		}
	}
}
